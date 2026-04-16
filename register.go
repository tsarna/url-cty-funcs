package urlcty

import "github.com/zclconf/go-cty/cty/function"

// GetURLFunctions returns all URL-related cty functions for registration
// in an HCL2 eval context.
func GetURLFunctions() map[string]function.Function {
	return map[string]function.Function{
		"urlparse":       makeURLParseFunc(),
		"urljoin":        makeURLJoinFunc(),
		"urljoinpath":    makeURLJoinPathFunc(),
		"urlqueryencode": makeURLQueryEncodeFunc(),
		"urlquerydecode": makeURLQueryDecodeFunc(),
		"urldecode":      makeURLDecodeFunc(),
	}
}
