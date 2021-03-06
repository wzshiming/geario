package geario

import (
	"fmt"
	"math"
	"strings"
)

const (
	decimal = "KMGTPEZY"
	step    = 1024
)

const (
	KB  B = 1e3
	MB  B = 1e6
	GB  B = 1e9
	TB  B = 1e12
	PB  B = 1e15
	EB  B = 1e18
	ZB  B = 1e21
	YB  B = 1e24
	KiB B = step
	MiB B = step * step
	GiB B = step * step * step
	TiB B = step * step * step * step
	PiB B = step * step * step * step * step
	EiB B = step * step * step * step * step * step
	ZiB B = step * step * step * step * step * step * step
	YiB B = step * step * step * step * step * step * step * step
)

// B is byte
type B float64

func (u B) String() string {
	u0 := math.Abs(float64(u))

	steps := 0
	last := u0
	for u0 >= step {
		u0 /= step
		last = u0
		steps++
	}

	snum := ""
	if steps != 0 {
		snum = decimal[steps-1 : steps]
	}
	if snum == "" {
		return fmt.Sprintf("%gB", last)
	}
	return fmt.Sprintf("%g%siB", last, snum)
}

// Parse text to byte type
func Parse(p string) (B, error) {
	f := 0.
	u := ""

	_, err := fmt.Sscanf(strings.ToUpper(p), "%f%s", &f, &u)
	if err != nil {
		return 0, err
	}
	u = strings.ToUpper(u)
	u = strings.TrimSuffix(u, "B")
	switch u {
	case "":
		return B(f), nil
	case "K":
		return B(f) * KB, nil
	case "M":
		return B(f) * MB, nil
	case "G":
		return B(f) * GB, nil
	case "T":
		return B(f) * TB, nil
	case "P":
		return B(f) * PB, nil
	case "E":
		return B(f) * EB, nil
	case "Z":
		return B(f) * ZB, nil
	case "Y":
		return B(f) * YB, nil
	case "KI":
		return B(f) * KiB, nil
	case "MI":
		return B(f) * MiB, nil
	case "GI":
		return B(f) * GiB, nil
	case "TI":
		return B(f) * TiB, nil
	case "PI":
		return B(f) * PiB, nil
	case "EI":
		return B(f) * EiB, nil
	case "ZI":
		return B(f) * ZiB, nil
	case "YI":
		return B(f) * YiB, nil
	}
	return 0, fmt.Errorf("parse failure %q", p)
}
