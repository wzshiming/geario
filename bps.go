package speed

import (
	"sync"
	"time"
)

type nodeBPS struct {
	b    B
	t    int64
	Next *nodeBPS
}

type BPS struct {
	Put  *nodeBPS
	End  *nodeBPS
	mut  sync.RWMutex
	pool *sync.Pool
	r    time.Duration
	u    int64
}

func NewBPSAver(r time.Duration) *BPS {
	if r < time.Second {
		r = time.Second
	}
	return &BPS{
		r: r,
		pool: &sync.Pool{
			New: func() interface{} {
				return &nodeBPS{}
			},
		},
		u: int64(time.Millisecond),
	}
}

func (p *BPS) unixNano() int64 {
	return time.Now().UnixNano() / p.u
}

func (p *BPS) Add(b B) {
	p.mut.Lock()
	defer p.mut.Unlock()

	now := p.unixNano()

	if p.Put != nil {
		if p.Put.t == now {
			p.Put.b += b
		} else {
			n, _ := p.pool.Get().(*nodeBPS)
			n.t = now
			n.b = b
			n.Next = nil
			p.Put.Next = n
			p.Put = p.Put.Next
		}
	} else {
		n, _ := p.pool.Get().(*nodeBPS)
		n.t = now
		n.b = b
		n.Next = nil
		p.Put = n
		p.End = p.Put
	}
}

func (p *BPS) limit(d time.Duration) {
	p.mut.Lock()
	defer p.mut.Unlock()

	now := p.unixNano()
	for p.End != nil && now-p.End.t > int64(d)/p.u {
		p.pool.Put(p.End)
		p.End = p.End.Next
	}
	if p.End == nil {
		p.Put = nil
	}
}

func (p *BPS) Number() B {
	p.limit(p.r)
	p.mut.RLock()
	defer p.mut.RUnlock()

	d := B(0)
	for i := p.End; i != nil; i = i.Next {
		d += i.b
	}

	return d / B(p.r/time.Second)
}

func (p *BPS) String() string {
	return p.Number().String() + "/S"
}
