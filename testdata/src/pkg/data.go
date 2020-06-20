// nolint: deadcode,unused,varcheck
package main

import (
	"fmt"

	"pkg/internal/forbidden" // want `pkg/internal/forbidden should not be used`
	"pkg/internal/helpers"
	context "pkg/internal/old" // want `pkg/internal/old should be replaced with pkg/internal/new`
)

// Forbidden symbols.
var (
	c   = helpers.Constant              // want `pkg/internal/helpers.Constant should not be used`
	v   = helpers.Variable              // want `pkg/internal/helpers.Variable should not be used`
	fun = helpers.FuncFactory()         // want `pkg/internal/helpers.FuncFactory should not be used`
	st  = helpers.StructFactory().Field // want `pkg/internal/helpers.StructFactory should not be used`
	fs  = forbidden.Deprecated
)

func Forbidden(
	_ helpers.StructType, // want `pkg/internal/helpers.StructType should not be used`
	_ helpers.InterfaceType, // want `pkg/internal/helpers.InterfaceType should not be used`
) {
}

// Replaced symbols.
var (
	_ context.Context = context.Background() // want `pkg/internal/old.Context should be replaced with pkg/internal/new.Context`
)

// Permitted symbols.
func main() {
	// Check that variables that shadow package names do not trigger the analysis.
	helpers := struct{ Variable []string }{}
	fmt.Println(helpers.Variable)
}

var (
	// Check that nested fields are not collapsed with the top-level package selector.
	f = helpers.Struct.StructFactory
)
