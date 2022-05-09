package geario

import (
	"io"
	"net"
	"sync/atomic"
	"time"
)

// GearReader is limit the speed of reading from io.Reader
func GearReader(r io.Reader, duration time.Duration, limit B) io.Reader {
	g := NewGear(duration, limit)
	return g.Reader(r)
}

// GearWriter is limit the speed of writing to io.Writer
func GearWriter(w io.Writer, duration time.Duration, limit B) io.Writer {
	g := NewGear(duration, limit)
	return g.Writer(w)
}

// GearReadWriter is limit the speed of reading and writing from io.ReadWriter
func GearReadWriter(rw io.ReadWriter, duration time.Duration, limit B) io.ReadWriter {
	g := NewGear(duration, limit)
	return g.ReadWriter(rw)
}

// GearConn is limit the speed of reading and writing from net.Conn
func GearConn(conn net.Conn, duration time.Duration, limit B) net.Conn {
	g := NewGear(duration, limit)
	return g.Conn(conn)
}

func NewGear(duration time.Duration, limit B) *Gear {
	return &Gear{
		bps:   NewBPSAver(duration),
		limit: int64(limit),
	}
}

type Gear struct {
	bps   *BPS
	limit int64
	aver  int64
	total uint64
}

func (g *Gear) add(b int) {
	atomic.AddUint64(&g.total, uint64(b))
	g.bps.Add(B(b))
}

func (g *Gear) step() int64 {
	limit := int64(g.Limit())
	aver := int64(g.bps.Aver())
	atomic.StoreInt64(&g.aver, aver)
	if limit < 0 {
		return limit
	}
	if aver < limit {
		return limit - aver
	}
	next := g.bps.Next()
	wait := next.Sub(time.Now())
	if wait > 0 {
		time.Sleep(wait)
	}
	return g.step()
}

func (g *Gear) SetLimit(b B) {
	atomic.StoreInt64(&g.limit, int64(b))
}

func (g *Gear) Limit() B {
	return B(atomic.LoadInt64(&g.limit))
}

func (g *Gear) MaxAver() B {
	return g.bps.MaxAver()
}

func (g *Gear) BPS() *BPS {
	return g.bps
}

func (g *Gear) Aver() B {
	return g.bps.Aver()
}

func (g *Gear) LastAver() B {
	return B(atomic.LoadInt64(&g.aver))
}

func (g *Gear) Total() B {
	return B(atomic.LoadUint64(&g.total))
}

func (g *Gear) Reader(r io.Reader) io.Reader {
	return &gearReader{g, r}
}

func (g *Gear) Writer(w io.Writer) io.Writer {
	return &gearWriter{g, w}
}

func (g *Gear) ReadWriter(rw io.ReadWriter) io.ReadWriter {
	return ReadWriterGear(rw, g, g)
}

func ReadWriterGear(rw io.ReadWriter, r, w *Gear) io.ReadWriter {
	return struct {
		io.Reader
		io.Writer
	}{
		r.Reader(rw),
		w.Writer(rw),
	}
}

func (g *Gear) Conn(rw net.Conn) net.Conn {
	return ConnGear(rw, g, g)
}

func ConnGear(rw net.Conn, r, w *Gear) net.Conn {
	type raw interface {
		Close() error
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
		SetDeadline(t time.Time) error
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error
	}
	return struct {
		raw
		io.Reader
		io.Writer
	}{
		rw,
		r.Reader(rw),
		w.Writer(rw),
	}
}

type gearReader struct {
	gear   *Gear
	reader io.Reader
}

func (g *gearReader) Read(p []byte) (n int, err error) {
	limit := int(g.gear.step())
	if limit > 0 && len(p) > limit {
		p = p[:limit]
	}
	n, err = g.reader.Read(p)
	g.gear.add(n)
	return n, err
}

type gearWriter struct {
	gear   *Gear
	writer io.Writer
}

func (g *gearWriter) Write(p []byte) (n int, err error) {
	limit := int(g.gear.step())
	if limit > 0 {
		for limit < len(p) {
			i, err := g.write(p[:limit])
			n += i
			if err != nil {
				return n, err
			}
			p = p[limit:]
			limit = int(g.gear.step())
		}
	}
	i, err := g.write(p)
	n += i
	if err != nil {
		return n, err
	}
	return n, nil
}

func (g *gearWriter) write(p []byte) (n int, err error) {
	n, err = g.writer.Write(p)
	g.gear.add(n)
	return n, err
}
