# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-07-15

### Changed

- **BREAKING: the functions are namespaced under `url::`.** All six already carried a `url`
  prefix — a poor-man's namespace — so this makes it a real one: the leaf names drop the
  prefix and sort together. HCL parses `a::b(x)` natively as a single flat map key, so this
  is a naming change, not a structural one. **Existing `.vcl`/`.cty` files must be updated.**

  | was | is |
  | --- | --- |
  | `urlparse` | `url::parse` |
  | `urljoin` | `url::join` |
  | `urljoinpath` | `url::join_path` |
  | `urlqueryencode` | `url::query_encode` |
  | `urlquerydecode` | `url::query_decode` |
  | `urldecode` | `url::decode` |

  The `url` object type keeps its name — it is a type, not a function. This package provides
  `url::decode` but not `url::encode`: percent-encoding is cty stdlib's `urlencode`, which a
  host can register under `url::encode` for a symmetric pair.

### Added

- **`Externs()` — the real signatures of the three functions cty cannot describe.**
  `externs.cty` (embedded; exposed as opaque bytes via `Externs()` and `ExternsFilename`)
  declares `url::join`, `url::join_path`, and `url::query_encode` as
  [functy](https://github.com/tsarna/functy) `//functy:extern` declarations. Each takes a
  union argument — a URL as a string or a `url` object, a query map of strings or lists —
  that cty declares `dynamic`, which then poisons its return type to `dynamic` and hides it.
  From cty alone, `url::join` reflects with no return type at all. The declarations restore
  the return and name the union. The other three functions take a concrete string and are
  cty-complete, so they are deliberately not declared. This package does not import functy.

  ```go
  parser.RegisterType("url", urlcty.URLObjectType)
  parser.RegisterExterns(urlcty.Externs(), urlcty.ExternsFilename)
  ```

## [0.1.0] - earlier

- Initial release: `urlparse`, `urljoin`, `urljoinpath`, `urlqueryencode`,
  `urlquerydecode`, `urldecode`, and the `url` capsule/object types.
