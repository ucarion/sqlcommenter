package sqlcommenter

import (
	"context"
	"reflect"
	"strconv"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx1 := context.Background()
	ctx2 := NewContext(ctx1, "k1", "v1")
	ctx3 := NewContext(ctx2, "k2", "v2")
	ctx4 := NewContext(ctx1, "a", "a")
	ctx5 := NewContext(ctx4, "z", "z")
	ctx6 := NewContext(ctx5, "b", "b")
	ctx7 := NewContext(ctx4, "a", "aa")

	testCases := []struct {
		ctx   context.Context
		attrs attrs
	}{
		{
			ctx:   ctx1,
			attrs: nil,
		},
		{
			ctx:   ctx2,
			attrs: attrs{{"k1", "'v1'"}},
		},
		{
			ctx:   ctx3,
			attrs: attrs{{"k1", "'v1'"}, {"k2", "'v2'"}},
		},
		{
			ctx:   ctx4,
			attrs: attrs{{"a", "'a'"}},
		},
		{
			ctx:   ctx5,
			attrs: attrs{{"a", "'a'"}, {"z", "'z'"}},
		},
		{
			ctx:   ctx6,
			attrs: attrs{{"a", "'a'"}, {"b", "'b'"}, {"z", "'z'"}},
		},
		{
			ctx:   ctx7,
			attrs: attrs{{"a", "'aa'"}, {"a", "'a'"}},
		},
	}

	for i, tt := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !reflect.DeepEqual(attrsFromContext(tt.ctx), tt.attrs) {
				t.Fatalf("bad attrsFromContext, want: %v, got: %v", tt.attrs, attrsFromContext(tt.ctx))
			}
		})
	}
}
