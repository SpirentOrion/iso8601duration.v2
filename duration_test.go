package duration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseGivenValid(t *testing.T) {
	vecs := []struct {
		in  string
		out time.Duration
	}{
		// Full string
		{"P1Y2DT3H4M5S", yearTime + 2*dayTime + 3*time.Hour + 4*time.Minute + 5*time.Second},

		// Partial strings
		{"P1Y", yearTime},
		{"P2W", 2 * weekTime},
		{"P2D", 2 * dayTime},
		{"PT3H", 3 * time.Hour},
		{"PT4M", 4 * time.Minute},
		{"PT5S", 5 * time.Second},

		// Decimal fractions in smallest parts
		{"P1.5Y", 1.5 * 365 * 24 * time.Hour},
		{"P0.5W", 0.5 * 7 * 24 * time.Hour},
		{"P1Y0.5D", yearTime + 0.5*24*time.Hour},
		{"P1YT0.5H", yearTime + 0.5*60*time.Minute},
		{"P1YT0.5M", yearTime + 0.5*60*time.Second},
		{"P1YT0.5S", yearTime + 500*time.Millisecond},
		{"P1.5D", 1.5 * 24 * time.Hour},
		{"P1DT0.5H", dayTime + 0.5*60*time.Minute},
		{"P1DT0.5M", dayTime + 0.5*60*time.Second},
		{"P1DT0.5S", dayTime + 500*time.Millisecond},
		{"PT1.5H", 1.5 * 60 * time.Minute},
		{"PT1H0.5M", time.Hour + 0.5*60*time.Second},
		{"PT1H0.5S", time.Hour + 500*time.Millisecond},
		{"PT1.5M", 1.5 * 60 * time.Second},
		{"PT1M0.5S", time.Minute + 500*time.Millisecond},
		{"PT0.5S", 500 * time.Millisecond},
	}

	t.Parallel()

	for _, vec := range vecs {
		d, err := Parse(vec.in)
		assert.NoError(t, err, vec.in)
		assert.Equal(t, vec.out, d, vec.in)
	}
}

func TestParseGivenInvalid(t *testing.T) {
	vecs := []struct {
		in  string
		err error
	}{
		// Bad formats
		{"", ErrBadFormat},
		{"asdf", ErrBadFormat},
		{"P", ErrBadFormat},
		{"P1", ErrBadFormat},
		{"P1X", ErrBadFormat},
		{"P1y", ErrBadFormat},
		{"1Y", ErrBadFormat},
		{"P5S1Y", ErrBadFormat},
		{"P1.0Y5S", ErrBadFormat},
		{"P1.0YT5S", ErrBadFormat},
		{"P1.0YT5.0S", ErrBadFormat},
		{"P1Y2W3D4H6M6S", ErrBadFormat},
		{"P1Y1W", ErrBadFormat},
		{"P1S", ErrBadFormat},

		// With month
		{"P0M", ErrNoMonth},
		{"P1M", ErrNoMonth},
		{"P1Y1M", ErrNoMonth},
		{"P0MT1M", ErrNoMonth},
		{"P1MT1M", ErrNoMonth},
	}

	t.Parallel()

	for _, vec := range vecs {
		d, err := Parse(vec.in)
		if assert.Error(t, err, vec.in) {
			assert.Equal(t, vec.err, err, vec.in)
		}
		assert.Equal(t, time.Duration(0), d, vec.in)
	}
}

func TestFormatGivenValid(t *testing.T) {
	t.Parallel()

	vecs := []struct {
		in  time.Duration
		out string
	}{
		// Zero-value
		{time.Duration(0), "P0Y"},

		// Smaller than a second
		{time.Nanosecond, "PT0.000000001S"},
		{time.Microsecond, "PT0.000001S"},
		{time.Millisecond, "PT0.001S"},

		// Smaller than a minute
		{time.Second, "PT1S"},
		{time.Second + time.Millisecond, "PT1.001S"},
		{59 * time.Second, "PT59S"},
		{59*time.Second + time.Millisecond, "PT59.001S"},

		// Smaller than an hour
		{time.Minute, "PT1M"},
		{time.Minute + time.Second, "PT1M1S"},
		{time.Minute + time.Second + time.Millisecond, "PT1M1.001S"},

		// Smaller than a day
		{time.Hour, "PT1H"},
		{time.Hour + time.Minute, "PT1H1M"},
		{time.Hour + time.Minute + time.Second, "PT1H1M1S"},
		{time.Hour + time.Minute + time.Second + time.Millisecond, "PT1H1M1.001S"},

		// Smaller than a week (NB: should not produce week-based values)
		{dayTime + time.Hour, "P1DT1H"},

		// Smaller than a year
		{10*dayTime + time.Hour, "P10DT1H"},
		{10*dayTime + time.Hour + time.Minute, "P10DT1H1M"},
		{10*dayTime + time.Hour + time.Minute + time.Second, "P10DT1H1M1S"},
		{10*dayTime + time.Hour + time.Minute + time.Second + time.Millisecond, "P10DT1H1M1.001S"},

		// Larger than year
		{yearTime + dayTime + time.Hour + time.Minute + time.Second + time.Millisecond, "P1Y1DT1H1M1.001S"},
		{yearTime + 10*dayTime + time.Hour + time.Minute + time.Second + time.Millisecond, "P1Y10DT1H1M1.001S"},
	}

	for _, vec := range vecs {
		s, err := Format(vec.in)
		assert.NoError(t, err, vec.in)
		assert.Equal(t, vec.out, s, vec.in)
	}
}

func TestFormatGivenInvalid(t *testing.T) {
	t.Parallel()

	vecs := []struct {
		in  time.Duration
		err error
	}{
		// Negative durations
		{-1 * time.Millisecond, ErrNoNegative},
		{-1 * time.Second, ErrNoNegative},
	}

	for _, vec := range vecs {
		s, err := Format(vec.in)
		if assert.Error(t, err, vec.in) {
			assert.Equal(t, vec.err, err, vec.in)
		}
		assert.Empty(t, s, vec.in)
	}
}
