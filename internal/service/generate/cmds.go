// Copyright © 2018 by PACE Telematics GmbH. All rights reserved.
// Created at 2018/08/31 by Vincent Landgraf

package generate

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/dave/jennifer/jen"
)

// CommandOptions are applied when generating the different
// microservice commands
type CommandOptions struct {
	DaemonName  string
	ControlName string
}

// NewCommandOptions generate command names using given name
func NewCommandOptions(name string) CommandOptions {
	return CommandOptions{
		DaemonName:  name + "d",
		ControlName: name + "ctl",
	}
}

// Commands generates the microservice commands based of
// the given path
func Commands(path string, options CommandOptions) {
	// Create directories
	dirs := []string{
		filepath.Join(path, "cmd", options.DaemonName),
		filepath.Join(path, "cmd", options.ControlName),
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0770)
		if err != nil {
			log.Fatal(fmt.Printf("Failed to create dir %s: %v", dir, err))
		}
	}

	// Create commands files
	for _, dir := range dirs {
		f, err := os.Create(filepath.Join(dir, "main.go"))
		if err != nil {
			log.Fatal(err)
		}

		code := jen.NewFilePathName("", "main")
		cmdName := filepath.Base(dir)

		if cmdName == options.DaemonName {
			generateDaemonMain(code)
		} else {
			generateControlMain(cmdName, code)
		}
		f.WriteString(copyright())
		f.WriteString(code.GoString())
	}
}

func generateDaemonMain(f *jen.File) {
	httpPkg := "lab.jamit.de/pace/go-microservice/http"
	f.ImportAlias(httpPkg, "pacehttp")
	f.Func().Id("main").Params().BlockFunc(func(g *jen.Group) {
		g.Id("handler").Op(":=").Qual(httpPkg, "Handler").Call()
		g.Id("s").Op(":=").Qual(httpPkg, "Server").Call(jen.Id("handler"))
		g.Qual("log", "Fatal").Call(jen.Id("s").Dot("ListenAndServe").Call())
	})
}

func generateControlMain(cmdName string, f *jen.File) {
	f.Func().Id("main").Params().Block(
		jen.Qual("fmt", "Printf").Call(jen.Lit(cmdName)))
}

// copyright generates copyright statement
func copyright() string {
	stmt := ""
	now := time.Now()
	stmt += fmt.Sprintf("// Copyright © %04d by PACE Telematics GmbH. All rights reserved.\n", now.Year())

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	stmt += fmt.Sprintf("// Created at %04d/%02d/%02d by %s\n\n", now.Year(), now.Month(), now.Day(), u.Name)
	return stmt
}
