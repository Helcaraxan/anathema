package main

import (
	"fmt"

	"external"

	"pkg/internal/forbidden" // want `pkg/internal/forbidden should not be used`
	"pkg/internal/helpers"
	context "pkg/internal/old" // want `pkg/internal/old should be replaced with pkg/internal/new`
)

// Forbidden symbols.
var (
	// Check that direct symbol references are picked up, including in vendored dependencies.
	_ = helpers.Constant  // want `pkg/internal/helpers.Constant should not be used`
	_ = helpers.Variable  // want `pkg/internal/helpers.Variable should not be used`
	_ = external.External // want `external.External should not be used`

	// Check that direct symbol references of functions are picked up in calls and nested selectors.
	_ = helpers.FuncFactory()         // want `pkg/internal/helpers.FuncFactory should not be used`
	_ = helpers.StructFactory().Field // want `pkg/internal/helpers.StructFactory should not be used`
)

func Forbidden(
	// Check that type references are picked up in function definitions.
	_ helpers.StructType, // want `pkg/internal/helpers.StructType should not be used`
	_ helpers.InterfaceType, // want `pkg/internal/helpers.InterfaceType should not be used`
) {
}

type MyStruct struct {
	// Check that type references in embedded structs are picked up.
	helpers.StructType    // want `pkg/internal/helpers.StructType should not be used`
	helpers.InterfaceType // want `pkg/internal/helpers.InterfaceType should not be used`
}

type MyInterface interface {
	// Check that type references in embedded interfaces are picked up.
	helpers.InterfaceType // want `pkg/internal/helpers.InterfaceType should not be used`
}

// Replaced symbols.
var (
	// Check that we get a correct message when a replacement is specified.
	_ context.Context = context.Background() // want `pkg/internal/old.Context should be replaced with pkg/internal/new.Context`
)

// Permitted symbols.
func main() {
	// Check that variables that shadow package names do not trigger the analysis.
	helpers := struct{ Variable []string }{}

	var err error
	// And that selections on universe-scoped symbols (e.g errors.Error()) do not fail.
	fmt.Println(helpers.Variable, err.Error())
}

var (
	// Check that nested fields are not collapsed with the top-level package selector.
	_ = helpers.Struct.StructFactory
	_ = forbidden.Deprecated
)

type StructType struct {
	// Check that we can declare types with names of symbols that are excluded in other packages.
}
