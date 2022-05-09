package geario

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type nodeBPS struct {
	b    B
	t    int64
	next *nodeBPS
}

type BPS struct {
	put *nodeBPS
	end *nodeBPS
	mut sync.RWMutex
	r   time.Duration
	u   int64
	max int64
}

var pool = &sync.Pool{
	New: func() interface{} {
		return &nodeBPS{}
	},
}

func getNodeBPS(t int64, b B) *nodeBPS {
	n, _ := pool.Get().(*nodeBPS)
	n.t = t
	n.b = b
	n.next = nil
	return n
}

func putNodeBPS(n *nodeBPS) {
	pool.Put(n)
}

func NewBPSAver(r time.Duration) *BPS {
	const min = time.Millisecond
	if r < min {
		r = min
	}
	return &BPS{
		r: r,
		u: int64(r / 10),
	}
}

func (p *BPS) unixNano() int64 {
	return time.Now().UnixNano() / p.u
}

func (p *BPS) Add(b B) {
	now := p.unixNano()

	p.mut.Lock()
	defer p.mut.Unlock()

	if p.put == nil {
		n := getNodeBPS(now, b)
		p.put = n
		p.end = n
	} else if p.put.t == now {
		p.put.b += b
	} else {
		n := getNodeBPS(now, b)
		p.put.next = n
		p.put = n
	}

	p.aver(now)
}

func (p *BPS) clear(now int64) {
	for p.end != nil && now-p.end.t > int64(p.r)/p.u {
		putNodeBPS(p.end)
		p.end = p.end.next
	}
	if p.end == nil {
		p.put = nil
	}
}

func (p *BPS) Next() time.Time {
	p.mut.RLock()
	defer p.mut.RUnlock()

	if p.end != nil {
		return time.Unix(0, p.end.t*p.u).Add(p.r)
	}
	return time.Time{}
}

func (p *BPS) Aver() B {
	p.mut.Lock()
	defer p.mut.Unlock()

	now := p.unixNano()
	return p.aver(now)
}

func (p *BPS) aver(now int64) B {
	p.clear(now)

	d := B(0)
	for i := p.end; i != nil; i = i.next {
		d += i.b
	}
	if int64(d) > p.max {
		p.max = int64(d)
	}
	return d
}

func (p *BPS) MaxAver() B {
	return B(atomic.LoadInt64(&p.max))
}

func (p *BPS) String() string {
	s := ""
	if p.r == time.Second {
		s = "s"
	} else {
		s = p.r.String()
	}
	return fmt.Sprintf("%v/%v", p.Aver(), s)
}
