package date_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/server/model/date"
)

func TestDate_String(t *testing.T) {
	d, err := date.NewFromString("2010-03-02")
	require.NoError(t, err)
	assert.Equal(t, "2010-03-02", d.String())
}
