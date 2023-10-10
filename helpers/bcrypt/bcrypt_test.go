package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndVerifyPassword(t *testing.T) {
	t.Run("Hash Valid OK", func(t *testing.T) {
		newHash, err := HashPassword("123")
		assert.Equal(t, err, nil)

		err = VerifyPassword(newHash, "123")
		assert.Equal(t, err, nil)
		
	})

	t.Run("Hash Invalid OK", func(t *testing.T) {
		newHash, err := HashPassword("123")
		assert.Equal(t, err, nil)

		err = VerifyPassword(newHash, "12345")
		assert.NotEqual(t, err, nil)
		
	})
}