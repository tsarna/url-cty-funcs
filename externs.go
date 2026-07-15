package urlcty

import _ "embed"

//go:embed externs.cty
var externsCty []byte

// ExternsFilename is the name reported for the embedded declarations in
// diagnostics.
const ExternsFilename = "url-cty-funcs/externs.cty"

// Externs returns the functy `//functy:extern` declarations for the three URL functions
// whose real signature their cty metadata cannot express: their bytes, which this package
// does not itself parse.
//
// url::join, url::join_path, and url::query_encode each take an argument that is a union —
// a URL as a string or a `url` object, a query map whose values are strings or lists —
// which cty has no way to name, so it declares the parameter `dynamic`. And a dynamic
// argument poisons cty's return type to `dynamic` too, so even the fixed type each returns
// (`url`, a string) is invisible in reflected metadata. The declarations restore the return
// and name the union, so that help(), generated documentation, and editor tooling can show
// it.
//
// The other three functions — url::parse, url::query_decode, url::decode — take a concrete
// string and return a fixed type, so their cty metadata is complete and they are
// deliberately not declared here.
//
// The declarations name the `url` object type (see URLObjectType); register it alongside:
//
//	parser.RegisterType("url", urlcty.URLObjectType)
//	parser.RegisterExterns(urlcty.Externs(), urlcty.ExternsFilename)
func Externs() []byte { return externsCty }
