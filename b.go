package geario

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type B float64

func (b B) String() string {
	return BytesSize(b)
}

// See: http://en.wikipedia.org/wiki/Binary_prefix
const (
	// Decimal

	KB B = 1000
	MB B = 1000 * KB
	GB B = 1000 * MB
	TB B = 1000 * GB
	PB B = 1000 * TB
	EB B = 1000 * PB
	ZB B = 1000 * EB
	YB B = 1000 * ZB

	// Binary

	KiB B = 1024
	MiB B = 1024 * KiB
	GiB B = 1024 * MiB
	TiB B = 1024 * GiB
	PiB B = 1024 * TiB
	EiB B = 1024 * PiB
	ZiB B = 1024 * EiB
	YiB B = 1024 * ZiB
)

type unitMap map[string]B

var (
	decimalMap = unitMap{"k": KB, "m": MB, "g": GB, "t": TB, "p": PB, "e": EB, "z": ZB, "y": YB}
	binaryMap  = unitMap{"k": KiB, "m": MiB, "g": GiB, "t": TiB, "p": PiB, "e": EiB, "z": ZiB, "y": YiB}
	sizeRegex  = regexp.MustCompile(`^(\d+(\.\d+)*) ?([kKmMgGtTpPeEzZyY])?[iI]?[bB]?$`)
)

var decimapAbbrs = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var binaryAbbrs = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

func getSizeAndUnit(size B, base B, _map []string) (B, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}

// CustomSize returns a human-readable approximation of a size
// using custom format.
func CustomSize(format string, size B, base B, _map []string) string {
	size, unit := getSizeAndUnit(size, base, _map)
	return fmt.Sprintf(format, size, unit)
}

// HumanSizeWithPrecision allows the size to be in any precision
func HumanSizeWithPrecision(size B, precision int) string {
	size, unit := getSizeAndUnit(size, 1000.0, decimapAbbrs)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

// BytesSize returns a human-readable size in bytes, kibibytes,
// mebibytes, gibibytes, or tebibytes (eg. "44kiB", "17MiB").
func BytesSize(size B) string {
	return CustomSize("%.4g%s", size, 1024.0, binaryAbbrs)
}

// FromHumanSize returns an integer from a human-readable specification of a
// size using SI standard (eg. "44kB", "17MB").
func FromHumanSize(size string) (B, error) {
	return parseSize(size, decimalMap)
}

// FromBytesSize returns an integer from a human-readable specification of a
// size using binary standard (eg. "44kiB", "17MiB").
func FromBytesSize(size string) (B, error) {
	return parseSize(size, binaryMap)
}

// Parses the human-readable size string into the amount it represents.
func parseSize(sizeStr string, uMap unitMap) (B, error) {
	matches := sizeRegex.FindStringSubmatch(sizeStr)
	if len(matches) != 4 {
		return -1, fmt.Errorf("invalid size: '%s'", sizeStr)
	}

	size, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return -1, err
	}

	unitPrefix := strings.ToLower(matches[3])
	if mul, ok := uMap[unitPrefix]; ok {
		size *= float64(mul)
	}

	return B(size), nil
}
