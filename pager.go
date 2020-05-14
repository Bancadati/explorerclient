package explorerclient

import (
	"net/url"
)

const (
	defaultPage = 1
	defaultSize = 100
)

// NewPager returns a pager
func NewPager(page, size int) *Pager {
	return &Pager{p: page, s: size}
}

// Pager for listing
type Pager struct {
	p int
	s int
}

func (p *Pager) apply(v url.Values) {
	if p == nil {
		return
	}

	if p.p < 1 {
		p.p = defaultPage
	}

	if p.s < 1 {
		p.s = defaultSize
	}
}

// Page returns a pager
func Page(page, size int) *Pager {
	return &Pager{p: page, s: size}
}
