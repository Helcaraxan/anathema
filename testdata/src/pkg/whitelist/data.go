package main

import (
	"fmt" // want `fmt should not be used`

	"pkg/internal/helpers"
	new_context "pkg/internal/new"
	old_context "pkg/internal/old" // want `pkg/internal/old should not be used`
)

// Forbidden symbols.
var (
	_ = helpers.Variable         // want `pkg/internal/helpers.Variable should not be used`
	_ = old_context.Background() // want `pkg/internal/old.Background should not be used`
)

func main() {
	fmt.Print("foo") // want `fmt.Print should not be used`
}

// Permitted symbols.
var (
	_ = new_context.Background()
)
