package main

import (
	"go-rails/template"
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
)

func main() {
	flag.Parse()

	dir, err := ioutil.ReadDir(flag.Arg(0))
	if err != nil {
		fmt.Println("permission denied", flag.Arg(0))
		os.Exit(-1)
	}
	for _, f := range dir {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html"){
			template.Gen(flag.Arg(0) + "/" + strings.TrimSuffix(f.Name(), ".html"))
		}
	}
}
