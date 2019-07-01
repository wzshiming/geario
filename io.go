package geario

import (
	"io"
	"time"
)

func GearReader(r io.Reader, duration time.Duration, limit B) io.Reader {
	return &gearReader{
		newGear(duration, limit),
		r,
	}
}

func GearWriter(w io.Writer, duration time.Duration, limit B) io.Writer {
	return &gearWriter{
		newGear(duration, limit),
		w,
	}
}

func GearReadWriter(rw io.ReadWriter, duration time.Duration, limit B) io.Writer {
	return &gearReadWriter{
		newGear(duration, limit),
		rw,
	}
}

func newGear(duration time.Duration, limit B) *gear {
	return &gear{
		bps:      NewBPSAver(duration),
		duration: duration,
		limit:    limit,
	}
}

type gear struct {
	bps      *BPS
	duration time.Duration
	limit    B
}

func (g *gear) add(b int) {
	g.bps.Add(B(b))
}

func (g *gear) step() {
	aver := g.bps.Aver()
	if aver < g.limit {
		return
	}
	time.Sleep(g.bps.Next().Sub(time.Now()))
}

type gearReader struct {
	gear   *gear
	reader io.Reader
}

func (g *gearReader) Read(p []byte) (n int, err error) {
	g.gear.step()
	g.gear.add(len(p))
	return g.reader.Read(p)
}

type gearWriter struct {
	gear   *gear
	writer io.Writer
}

func (g *gearWriter) Write(p []byte) (n int, err error) {
	g.gear.step()
	g.gear.add(len(p))
	return g.writer.Write(p)
}

type gearReadWriter struct {
	gear *gear
	rw   io.ReadWriter
}

func (g *gearReadWriter) Read(p []byte) (n int, err error) {
	g.gear.step()
	g.gear.add(len(p))
	return g.rw.Read(p)
}

func (g *gearReadWriter) Write(p []byte) (n int, err error) {
	g.gear.step()
	g.gear.add(len(p))
	return g.rw.Write(p)
}
