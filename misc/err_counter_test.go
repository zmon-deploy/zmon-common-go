package misc

import (
	"github.com/zmon-deploy/zmon-common-go/timeutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	c := NewErrorCounter(time.Hour, 3)
	err := errors.New("some error")

	c.Put(err)
	require.False(t, c.IsOverThreshold())
	c.Put(err)
	require.False(t, c.IsOverThreshold())
	c.Put(err)
	require.True(t, c.IsOverThreshold())
}

func TestExpire(t *testing.T) {
	c := NewErrorCounter(time.Hour, 3)
	err := errors.New("some error")

	c.clock = timeutil.NewFixedClock(time.Now().Add(-2 * time.Hour))
	c.Put(err)
	require.Equal(t, 1, c.count)

	c.clock = timeutil.NewClock()
	c.Put(err)
	require.Equal(t, 1, c.count) // 2 시간 전에 들어간 값은 drop 되어 count 가 다시 1이 되어야 함
	c.Put(err)
	require.Equal(t, 2, c.count)
}

func TestPullErr(t *testing.T) {
	c := NewErrorCounter(time.Hour, 3)
	err := errors.New("some error")

	require.Nil(t, c.Err())

	c.Put(err)
	c.Put(err)
	require.NotNil(t, c.Err())

	splited := strings.Split(c.Err().Error(), ";")
	require.Equal(t, 2, len(splited))
}