package speed

import (
	"fmt"
	"math"
	"strings"
)

type B uint64

const (
	decimal = "KMGTPEZY"
	unit    = "B"
	step    = 1024
)

func (u B) String() string {
	if u == 0 {
		return fmt.Sprintf("%.3f %s", 0., unit)
	}
	u0 := u
	us := []uint64{}
	for u0 != 0 {
		us = append(us, uint64(u0%step))
		u0 /= step
	}

	num := float64(us[len(us)-1])
	if len(us) >= 2 {
		u := us[len(us)-2]
		num += float64(u) / step
	}
	snum := ""
	if len(us) > 1 {
		snum = decimal[len(us)-2 : len(us)-1]
	}

	return fmt.Sprintf("%.3f %s%s", num, snum, unit)
}

func Parse(p string) B {
	f := float64(0)
	u := ""

	fmt.Sscanf(strings.ToUpper(p), "%f%s", &f, &u)
	switch len(u) {
	case 0:
		return B(f)
	case 1, 2:
		for k, v := range decimal {
			if v == rune(u[0]) {
				return B(f * math.Pow(step, float64(k+1)))
			}
		}
	}
	return B(f)
}
