package payload_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/server/model/payload"
)

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := payload.Payload{
			Scheme:  payload.SchemePlainProto,
			Version: "VER#",
			Blob:    []byte("content"),
		}

		p, err := payload.New(payload.SchemePlainProto, "VER#", []byte("content"))
		require.NoError(t, err)
		assert.Equal(t, expected, p)
	})

	t.Run("InvalidScheme", func(t *testing.T) {
		p, err := payload.New(payload.Scheme(300), "VER#", []byte("content"))
		require.ErrorIs(t, err, payload.ErrInvalidScheme)
		assert.Equal(t, payload.Payload{}, p)
	})

	t.Run("InvalidVersion", func(t *testing.T) {
		p, err := payload.New(payload.SchemePlainProto, "", []byte("content"))
		require.ErrorIs(t, err, payload.ErrEmptyVersion)
		assert.Equal(t, payload.Payload{}, p)
	})
}
