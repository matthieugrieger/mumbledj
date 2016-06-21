package duration_test

import (
	"testing"
	"time"

	duration "github.com/channelmeter/iso8601duration"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	t.Parallel()

	// test with bad format
	_, err := duration.FromString("asdf")
	assert.Equal(t, err, duration.ErrBadFormat)

	// test with month
	_, err = duration.FromString("P1M")
	assert.Equal(t, err, duration.ErrNoMonth)

	// test with good full string
	dur, err := duration.FromString("P1Y2DT3H4M5S")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Years)
	assert.Equal(t, 2, dur.Days)
	assert.Equal(t, 3, dur.Hours)
	assert.Equal(t, 4, dur.Minutes)
	assert.Equal(t, 5, dur.Seconds)

	// test with good week string
	dur, err = duration.FromString("P1W")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Weeks)
}

func TestString(t *testing.T) {
	t.Parallel()

	// test empty
	d := duration.Duration{}
	assert.Equal(t, d.String(), "P")

	// test only larger-than-day
	d = duration.Duration{Years: 1, Days: 2}
	assert.Equal(t, d.String(), "P1Y2D")

	// test only smaller-than-day
	d = duration.Duration{Hours: 1, Minutes: 2, Seconds: 3}
	assert.Equal(t, d.String(), "PT1H2M3S")

	// test full format
	d = duration.Duration{Years: 1, Days: 2, Hours: 3, Minutes: 4, Seconds: 5}
	assert.Equal(t, d.String(), "P1Y2DT3H4M5S")

	// test week format
	d = duration.Duration{Weeks: 1}
	assert.Equal(t, d.String(), "P1W")
}

func TestToDuration(t *testing.T) {
	t.Parallel()

	d := duration.Duration{Years: 1}
	assert.Equal(t, d.ToDuration(), time.Hour*24*365)

	d = duration.Duration{Weeks: 1}
	assert.Equal(t, d.ToDuration(), time.Hour*24*7)

	d = duration.Duration{Days: 1}
	assert.Equal(t, d.ToDuration(), time.Hour*24)

	d = duration.Duration{Hours: 1}
	assert.Equal(t, d.ToDuration(), time.Hour)

	d = duration.Duration{Minutes: 1}
	assert.Equal(t, d.ToDuration(), time.Minute)

	d = duration.Duration{Seconds: 1}
	assert.Equal(t, d.ToDuration(), time.Second)
}
