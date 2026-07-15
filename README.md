# url-cty-funcs

cty types and functions for working with URLs; mainly used in HCL2 templates.

[![CI](https://github.com/tsarna/url-cty-funcs/actions/workflows/ci.yml/badge.svg)](https://github.com/tsarna/url-cty-funcs/actions/workflows/ci.yml)

## Overview

This module provides a [go-cty](https://github.com/zclconf/go-cty) capsule type wrapping Go's `*url.URL`, a companion object type that materializes every URL field as a named attribute, and a set of cty functions for parsing, joining, and encoding URLs in HCL2 expression evaluation contexts.

## Installation

```
go get github.com/tsarna/url-cty-funcs
```

## Usage

```go
import (
    urlcty "github.com/tsarna/url-cty-funcs"
)

// Register all URL functions in an HCL eval context
funcs := urlcty.GetURLFunctions()
// funcs is map[string]function.Function — merge into your eval context
```

The functions are namespaced under `url::` (e.g. `url::parse`, `url::join`). HCL parses
`a::b(x)` natively as a single flat map key, so no special handling is needed — the keys
of the returned map simply contain `::`. This package provides `url::decode` but not
`url::encode`: percent-encoding is cty stdlib's `urlencode`, which a host can also register
under `url::encode` for a symmetric pair.

## Types

### `urlcty.URLCapsuleType`

A cty capsule type wrapping Go's `*url.URL`. Values of this type carry the parsed URL as an opaque handle and can be passed back to URL functions without re-parsing.

### `urlcty.URLObjectType`

A static cty object type with named attributes for every `*url.URL` field, plus a `_capsule` attribute that holds the `URLCapsuleType` value. This is the type returned by `url::parse`, `url::join`, and `url::join_path`, so you can read individual components directly (`u.scheme`, `u.host`, `u.query["key"]`, …) while still passing `u` back to other URL functions.

Attributes:

| Attribute | Type | Notes |
|-----------|------|-------|
| `scheme`, `opaque`, `host`, `hostname`, `port`, `path`, `raw_path`, `raw_query`, `fragment`, `raw_fragment` | `string` | |
| `username`, `password` | `string` | Empty if no userinfo |
| `password_set` | `bool` | Distinguishes `user:@host` from `user@host` |
| `query` | `map(list(string))` | Multi-value parameters preserve order |
| `force_query`, `omit_host` | `bool` | |
| `_capsule` | `URLCapsuleType` | Reusable parsed form |

### Helper functions

```go
urlcty.NewURLCapsule(u *url.URL) cty.Value
urlcty.GetURLFromCapsule(val cty.Value) (*url.URL, error)
urlcty.GetURLFromValue(val cty.Value) (*url.URL, error) // accepts string, capsule, or URL object
urlcty.BuildURLObject(u *url.URL) cty.Value
```

`GetURLFromValue` accepts any of three input shapes — a plain `cty.String`, a `URLCapsuleType` capsule, or a `URLObjectType` object (via its `_capsule` attribute) — which lets URL-consuming functions (`url::join`, `url::join_path`, and user code) take whichever form is most convenient.

### Generic operation support

`*URLWrapper` implements the `Stringable` and `Gettable` interfaces from [rich-cty-types](https://github.com/tsarna/rich-cty-types), enabling:

- `tostring(u)` → canonical URL string
- `get(u, "query_param", key)` → `list(string)` of values for the named query parameter

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `url::parse` | `(string) → url` | Parse a URL string into a URL object |
| `url::join` | `(url, url) → url` | Resolve `ref` against `base` per RFC 3986 |
| `url::join_path` | `(url, string...) → url` | Append path segments to `base`, each segment is path-escaped |
| `url::query_encode` | `(map(string) \| map(list(string))) → string` | Encode a map as a URL-encoded query string |
| `url::query_decode` | `(string) → map(list(string))` | Decode a query string (leading `?` optional) |
| `url::decode` | `(string) → string` | Percent-decode a string (inverse of `urlencode`; `+` → space) |

`url` arguments accept any of: a `string`, a `URLCapsuleType` capsule, or a `URLObjectType` object.

## Signature declarations

Three of these functions have signatures cty cannot express. `url::join`, `url::join_path`,
and `url::query_encode` each take an argument that is a **union** — a URL as a string or a
`url` object, a query map whose values are strings or lists — which cty has no way to name,
so it declares the parameter `dynamic`. And a `dynamic` argument poisons cty's return type
to `dynamic` too, so even the fixed type each returns (`url`, a string) is **invisible** in
reflected metadata: from cty alone, `url::join` reflects with no return type at all.

So `externs.cty` declares those three as [functy](https://github.com/tsarna/functy)
`//functy:extern` declarations, restoring the return and naming the union (as overload
forms, for `url::query_encode`). It is never compiled and declares nothing callable; it
exists so `help()`, generated documentation, and editor tooling can show the real
signatures. The other three functions take a concrete string and are cty-complete, so they
are deliberately not declared. `Externs()` returns the file as opaque bytes — this package
does not import functy:

```go
parser.RegisterType("url", urlcty.URLObjectType)
parser.RegisterExterns(urlcty.Externs(), urlcty.ExternsFilename)
```

Every function and parameter also carries a cty `Description`, so a host that is not a
functy host still has complete metadata.

## Examples

```hcl
# Parse and decompose
u = url::parse("https://user:pass@example.com:8080/v1/users?tag=go&tag=cty#top")
u.scheme     # "https"
u.host       # "example.com:8080"
u.hostname   # "example.com"
u.port       # "8080"
u.path       # "/v1/users"
u.query      # { tag = ["go", "cty"] }

# Query param lookup via get()
get(u, "query_param", "tag")   # ["go", "cty"]

# Join relative / absolute references (RFC 3986)
url::join("https://example.com/a/b", "../c")        # → /c
url::join("https://example.com/base/", "/absolute") # → /absolute

# Append path segments (segments are percent-escaped)
url::join_path("https://api.example.com/v1", "users", "42", "profile")
# → https://api.example.com/v1/users/42/profile

url::join_path("https://example.com/", "hello world")
# tostring(...) → "https://example.com/hello%20world"

# Query encoding / decoding
url::query_encode({ q = "hello world", page = "2" })   # "page=2&q=hello+world"
url::query_encode({ tag = ["go", "cty"] })             # "tag=go&tag=cty"
url::query_decode("a=1&b=2&b=3")                       # { a = ["1"], b = ["2", "3"] }
url::query_decode("?k=v")                              # { k = ["v"] }

# Percent-decoding
url::decode("caf%C3%A9")    # "café"
url::decode("hello+world")  # "hello world"

# Round-trip via tostring
tostring(url::parse("https://example.com/path"))  # "https://example.com/path"
```

## License

BSD 2-Clause — see [LICENSE](LICENSE).
