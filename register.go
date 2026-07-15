package urlcty

import "github.com/zclconf/go-cty/cty/function"

// GetURLFunctions returns all URL-related cty functions for registration
// in an HCL2 eval context.
//
// The names are namespaced under `url::`. All six already carried a `url` prefix — a
// poor-man's namespace — so this makes it a real one: the leaf names drop the prefix and
// sort together. HCL parses `a::b(x)` natively as a single flat map key, so this is a
// naming change, not a structural one.
//
// This package provides `url::decode` but not `url::encode`: percent-encoding is cty
// stdlib's `urlencode`. A host that wants the pair symmetric can register that function
// under `url::encode` as well.
func GetURLFunctions() map[string]function.Function {
	return map[string]function.Function{
		"url::parse":        makeURLParseFunc(),
		"url::join":         makeURLJoinFunc(),
		"url::join_path":    makeURLJoinPathFunc(),
		"url::query_encode": makeURLQueryEncodeFunc(),
		"url::query_decode": makeURLQueryDecodeFunc(),
		"url::decode":       makeURLDecodeFunc(),
	}
}
