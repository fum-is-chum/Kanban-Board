package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedResponse(t *testing.T) {
	t.Run("Failed Response OK", func(t *testing.T) {
		res := FailedResponse("Some message")
		assert.Equal(t, res["status"], "failed")
		assert.Equal(t, res["message"], "Some message")
	})
}

func TestSuccessResponse(t *testing.T) {
	t.Run("Sucess Response OK", func(t *testing.T) {
		res := SuccessResponse("Some message")
		assert.Equal(t, res["status"], "success")
		assert.Equal(t, res["message"], "Some message")
	})
}

func TestSuccessWithDataResponse(t *testing.T) {
	t.Run("Sucess With Data Response OK", func(t *testing.T) {
		res := SuccessWithDataResponse("Some message", "some data")
		assert.Equal(t, res["status"], "success")
		assert.Equal(t, res["message"], "Some message")
		assert.Equal(t, res["data"], "some data")
	})
}
