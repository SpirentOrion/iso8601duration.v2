// Package duration provides an implementation of ISO8601 duration parsing and
// formatting using time.Duration as the underlying type.
package duration

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrBadFormat is returned when parsing fails.
	ErrBadFormat = errors.New("bad format string")

	// ErrNoMonth is returned when a month element is in the format string.
	ErrNoMonth = errors.New("no month elements allowed")

	// ErrNoNegative is returned when a negative Duration is formatted.
	ErrNoNegative = errors.New("cannot format negative duration")

	format = regexp.MustCompile(`^P((?P<year>\d+((\.|,)\d+)?)Y)?((?P<month>\d+((\.|,)\d+)?)M)?((?P<week>\d+((\.|,)\d+)?)W)?((?P<day>\d+((\.|,)\d+)?)D)?(T((?P<hour>\d+((\.|,)\d+)?)H)?((?P<minute>\d+((\.|,)\d+)?)M)?((?P<second>\d+((\.|,)\d+)?)S)?)?$`)
)

const (
	dayTime  = 24 * time.Hour
	weekTime = 7 * 24 * time.Hour
	yearTime = 365 * 24 * time.Hour
)

// Parse parses an ISO8601-formatted duration value and returns a time.Duration.
// Month elements (e.g. "P1M") are not supported.
func Parse(s string) (time.Duration, error) {
	match := format.FindStringSubmatch(strings.TrimSpace(s))
	if match == nil {
		return 0, ErrBadFormat
	}

	var d time.Duration
	var numElems, weekElem, fracElem int

	for i, name := range format.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		whole, frac, hasFrac, err := parseDecimal(part)
		if err != nil {
			return 0, ErrBadFormat
		}

		// Fractional elements must be the last element in the string
		if hasFrac {
			if fracElem > 0 {
				return 0, ErrBadFormat
			}
			fracElem = i
		} else if fracElem > 0 {
			return 0, ErrBadFormat
		}

		switch name {
		case "year":
			d += time.Duration(whole) * yearTime
			if frac != 0 {
				d += time.Duration(frac * float64(yearTime))
			}
		case "month":
			return 0, ErrNoMonth
		case "week":
			d += time.Duration(whole) * weekTime
			if frac != 0 {
				d += time.Duration(frac * float64(weekTime))
			}
			weekElem = i
		case "day":
			d += time.Duration(whole) * dayTime
			if frac != 0 {
				d += time.Duration(frac * float64(dayTime))
			}
		case "hour":
			d += time.Duration(whole) * time.Hour
			if frac != 0 {
				d += time.Duration(frac * float64(time.Hour))
			}
		case "minute":
			d += time.Duration(whole) * time.Minute
			if frac != 0 {
				d += time.Duration(frac * float64(time.Minute))
			}
		case "second":
			d += time.Duration(whole) * time.Second
			if frac != 0 {
				d += time.Duration(frac * float64(time.Second))
			}
		}
		numElems++
	}

	// There must be at least one element in the string
	if numElems == 0 {
		return 0, ErrBadFormat
	}

	// Week elements, when used, must be the only elements in the string
	if weekElem > 0 && numElems > 1 {
		return 0, ErrBadFormat
	}

	return d, nil
}

func parseDecimal(s string) (whole int64, frac float64, hasFrac bool, err error) {
	if sep := strings.IndexAny(s, ".,"); sep != -1 {
		if whole, err = strconv.ParseInt(s[0:sep], 10, 64); err != nil {
			return
		}
		if frac, err = strconv.ParseFloat("."+s[sep+1:], 64); err != nil {
			return
		}
		hasFrac = true
	} else {
		whole, err = strconv.ParseInt(s, 10, 64)
	}
	return
}

// Format returns a string representation of a time.Duration value using ISO8601
// formatting. Negative duration values are not supported.
func Format(d time.Duration) (string, error) {
	if d < 0 {
		return "", ErrNoNegative
	}

	s := bytes.NewBufferString("P")
	if d == 0 {
		s.WriteString("0Y")
		goto done
	}

	if f := d / yearTime; f >= 1 {
		fmt.Fprintf(s, "%dY", f)
		d -= f * yearTime
		if d == 0 {
			goto done
		}
	}

	if f := d / dayTime; f >= 1 {
		fmt.Fprintf(s, "%dD", f)
		d -= f * dayTime
		if d == 0 {
			goto done
		}
	}

	s.WriteString("T")

	if f := d / time.Hour; f >= 1 {
		fmt.Fprintf(s, "%dH", f)
		d -= f * time.Hour
		if d == 0 {
			goto done
		}
	}

	if f := d / time.Minute; f >= 1 {
		fmt.Fprintf(s, "%dM", f)
		d -= f * time.Minute
		if d == 0 {
			goto done
		}
	}

	if d%time.Second == 0 {
		fmt.Fprintf(s, "%dS", d/time.Second)
		goto done
	}

	if d%time.Millisecond == 0 {
		fmt.Fprintf(s, "%.3fS", float64(d)/float64(time.Second))
		goto done
	}

	if d%time.Microsecond == 0 {
		fmt.Fprintf(s, "%.6fS", float64(d)/float64(time.Second))
		goto done
	}

	fmt.Fprintf(s, "%.9fS", float64(d)/float64(time.Second))

done:
	return s.String(), nil
}
