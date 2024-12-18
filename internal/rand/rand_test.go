package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringBytes(t *testing.T) {
    type args struct {
         n int
    }
    tests := []struct {
        name string
        args args
    }{
        {
			name: "positive test #1",
			args: args{
				n: 2,
			},
		},
		{
			name: "positive test #2",
			args: args{
				n: 7,
			},
		},
		{
			name: "positive test #3",
			args: args{
				n: 12,
			},
		},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
			got := RandStringBytes(tt.args.n)
			assert.Equal(t, len(got), tt.args.n)
         })
     }
}