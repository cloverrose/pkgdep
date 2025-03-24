package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/cloverrose/pkgdep"
)

func main() { unitchecker.Main(pkgdep.Analyzer) }
