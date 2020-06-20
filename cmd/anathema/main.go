package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Helcaraxan/anathema"
)

func main() {
	singlechecker.Main(anathema.Analysis(nil))
}
