package sqlcommenter

import (
	"context"
	"strconv"
	"testing"
)

func BenchmarkNewContext(b *testing.B) {
	b.Run("1 1", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 1, 1) })
	b.Run("1 10", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 1, 10) })
	b.Run("1 100", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 1, 100) })
	b.Run("10 1", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 10, 1) })
	b.Run("10 10", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 10, 10) })
	b.Run("10 100", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 10, 100) })
	b.Run("100 1", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 100, 1) })
	b.Run("100 10", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 100, 10) })
	b.Run("100 100", func(b *testing.B) { benchmarkNewContextWithCardinality(b, 100, 100) })
}

var result string // prevent compiler from optimizing away appendComment calls

func benchmarkNewContextWithCardinality(b *testing.B, numAttrs, numOps int) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		for i := 0; i < numAttrs; i++ {
			ctx = NewContext(ctx, strconv.Itoa(i), strconv.Itoa(i))
		}

		for i := 0; i < numOps; i++ {
			result = appendComment(attrsFromContext(ctx), "select 1")
		}
	}
}
