package err_test

import (
	"testing"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := errors.NewID()
		assert.Len(t, id, 16)

		for _, r := range id {
			assert.Contains(t, errors.Alphabet, string(r))
		}
	}
}
