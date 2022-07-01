package sqlcommenter

import (
	"context"
	"fmt"
	"sort"
)

type ctxKey struct{}

// attrs is a sorted list of encoded key-value pairs.
type attrs [][2]string

func (a *attrs) insert(k, v string) attrs {
	k = encodeKey(k)
	v = encodeValue(v)

	i := sort.Search(len(*a), func(i int) bool {
		return (*a)[i][0] >= k
	})

	aa := make([][2]string, len(*a)+1)
	copy(aa, (*a)[:i])
	aa[i] = [2]string{k, v}
	copy(aa[i+1:], (*a)[i:])

	return aa
}

// NewContext returns a new Context that carries the sqlcommenter attribute
// key-value pairs in kv.
//
// The returned Context will carry all sqlcommenter attributes that ctx carries
// in addition to those in kv.
//
// Keys are not "overridden". If kv contains a key that ctx already carries,
// then the returned Context will carry multiple attributes with the same key.
//
// NewContext panics if len(kv) is odd.
func NewContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 != 0 {
		panic(fmt.Errorf("NewContext: odd-length kvs: %d", len(kv)))
	}

	attrs, _ := ctx.Value(ctxKey{}).(attrs)
	for i := 0; i < len(kv); i += 2 {
		attrs = attrs.insert(kv[i], kv[i+1])
	}

	return context.WithValue(ctx, ctxKey{}, attrs)
}

func attrsFromContext(ctx context.Context) attrs {
	attrs, _ := ctx.Value(ctxKey{}).(attrs)
	return attrs
}
