package urlcty

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The extern file is scanned here with regexes rather than parsed with functy on purpose:
// this package must not depend on functy (its bytes are opaque to it), and the checks only
// need the names.
var (
	namespaceRE  = regexp.MustCompile(`(?m)^namespace (\w+)\s*$`)
	externDeclRE = regexp.MustCompile(`(?m)^func (\w+)\(`)
)

// noExtern are the functions deliberately left undeclared: they take a concrete string and
// return a fixed type, so their cty metadata already says everything true about them.
var noExtern = map[string]bool{
	"url::parse":        true,
	"url::query_decode": true,
	"url::decode":       true,
}

// declaredExterns returns every declared name, qualified by the file's namespace.
func declaredExterns(t *testing.T) map[string]bool {
	t.Helper()
	src := string(Externs())

	m := namespaceRE.FindStringSubmatch(src)
	require.NotNil(t, m, "externs.cty must declare a namespace")
	ns := m[1] + "::"

	declared := make(map[string]bool)
	for _, d := range externDeclRE.FindAllStringSubmatch(src, -1) {
		declared[ns+d[1]] = true
	}
	return declared
}

// TestExternsCoverTheRightFunctions is the drift guard, in both directions: a function
// whose cty return is hidden by a dynamic argument must be declared to restore it, and a
// function cty describes fully must not be (declaring it would shadow correct, self-
// maintaining metadata with a hand-written copy).
func TestExternsCoverTheRightFunctions(t *testing.T) {
	declared := declaredExterns(t)
	funcs := GetURLFunctions()

	for name := range funcs {
		if noExtern[name] {
			assert.False(t, declared[name],
				"%s() is listed in noExtern but externs.cty declares it: pick one", name)
			continue
		}
		assert.True(t, declared[name],
			"%s() has a dynamic parameter that hides its return type, but no declaration in externs.cty", name)
	}
	for name := range declared {
		assert.Contains(t, funcs, name,
			"externs.cty declares %s(), which GetURLFunctions does not provide", name)
	}
}

// The bytes must declare themselves an extern file, and be namespaced under url::.
func TestExternsAreNamespacedAndDirected(t *testing.T) {
	require.True(t, strings.HasPrefix(string(Externs()), "//functy:extern\n"),
		"externs.cty must begin with the //functy:extern directive")
	for name := range GetURLFunctions() {
		assert.True(t, strings.HasPrefix(name, "url::"), "%s() is not under the url:: namespace", name)
	}
}

// Every function and every parameter carries a cty description. The metadata is the only
// documentation a non-functy cty host can see, and the only thing functy's own doc() reads.
func TestEverythingIsDescribed(t *testing.T) {
	for name, fn := range GetURLFunctions() {
		assert.NotEmpty(t, fn.Description(), "%s() has no cty Description", name)

		for _, p := range fn.Params() {
			assert.NotEmpty(t, p.Description, "%s() parameter %q has no Description", name, p.Name)
		}
		if vp := fn.VarParam(); vp != nil {
			assert.NotEmpty(t, vp.Description, "%s() variadic parameter %q has no Description", name, vp.Name)
		}
	}
}
