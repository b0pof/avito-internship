package main

import (
	"github.com/b0pof/avito-internship/internal/app"
)

func main() {
	a := app.MustInit()
	a.Run()
}
