package geario

import (
	"fmt"
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
	const min = time.Millisecond
	if r < min {
		r = min
	}
	return &BPS{
		r: r,
		pool: &sync.Pool{
			New: func() interface{} {
				return &nodeBPS{}
			},
		},
		u: int64(min),
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

	for p.End != nil && now-p.End.t > int64(p.r)/p.u {
		p.pool.Put(p.End)
		p.End = p.End.Next
	}
	if p.End == nil {
		p.Put = nil
	}
}

func (p *BPS) Next() time.Time {
	p.mut.RLock()
	defer p.mut.RUnlock()

	if p.End != nil {
		return time.Unix(0, p.End.t*p.u).Add(p.r)
	}
	return time.Time{}
}

func (p *BPS) Aver() B {
	p.mut.RLock()
	defer p.mut.RUnlock()

	d := B(0)
	for i := p.End; i != nil; i = i.Next {
		d += i.b
	}
	return d
}

func (p *BPS) String() string {
	return fmt.Sprintf("%v/%v", p.Aver(), p.r)
}
