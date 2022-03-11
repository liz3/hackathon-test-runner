package main

import (
	"flag"
	"path/filepath"
	"os"
	"strings"
	"encoding/json"
	"fmt"
)

func main() {
	var action = "read"
	var rootPath = "."
	var file = "out.json"
	flag.StringVar(&rootPath, "path", ".", "The base folder for the Project")
	flag.StringVar(&action, "action", "read", "Set the action, can be execute or read")
	flag.StringVar(&file, "out", "out.json", "Set output file")
	flag.Parse()


	if action == "read" {
		ctx := CreateContext(rootPath)
		err := filepath.Walk(rootPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(path, ".js") && !strings.Contains(path, "node_modules") {
					toRead := ctx.RegisterFile(path)
					Parse(toRead, ctx)
				}
				return nil
			})
		if err != nil {
			panic(err)
		}
		payload, err := json.MarshalIndent(ctx.SerialisedForm(), "", "  ")
		if err != nil {
			panic(err)
		}
		WriteFile(file, payload)
	} else if action == "run" {
		var file = ReadFile(file)
		if file != nil{
			var p FilePayload
			err := json.Unmarshal([]byte(file.Content), &p)
			if err != nil {
				panic(err)
			}
			ctx := CreateRunContext(rootPath)
			RunAll(&p, ctx)
			outData, err := json.MarshalIndent(ctx.Serialise(), "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(outData))
		}
	}
}
