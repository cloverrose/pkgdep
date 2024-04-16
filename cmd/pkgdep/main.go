package main

import (
	"github.com/cloverrose/pkgdep"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(pkgdep.Analyzer) }
