package mathutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NCk(t *testing.T) {
	tests := []struct {
		n    int
		k    int
		want float64
	}{
		{0, 0, 1},
		{1, 0, 1},
		{5, 0, 1},
		{0, 1, 0},
		{0, 5, 0},
		{2, 1, 2},
		{4, 2, 6},
		{33, 7, 4272048},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("nCk(%d, %d)", tt.n, tt.k), func(t *testing.T) {
			if got := NCk(tt.n, tt.k); got != tt.want {
				t.Errorf("nCk() =N%v, want %v", got, tt.want)
			}
		})
	}

	n := 49
	k := 12
	assert.Equal(t, NCk(n, k), NCk(n, n-k), "mirror identity")
	assert.Equal(t, NCk(n, k), NCk(n-1, k-1)+NCk(n-1, k), "pascal's triangle identity")

}
