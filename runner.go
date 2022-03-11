package main

import (
	"encoding/json"
	"os/exec"
	gp "path"
	"fmt"
)

func _run(name, folder, file string) *JasmineResult {
	var executablePath = gp.Join(folder, "node_modules", ".bin", "jasmine")
	var command = executablePath
	var args = []string{file,"--reporter=./reporter.js" }
	if len(name) > 0 {
		args = append(args, "--filter=" + name)
	}
	out, err := exec.Command(command, args...).Output()
	if err != nil && (out == nil || len(out) == 0) {
		return nil
	}
	var d JasmineResult
	err = json.Unmarshal(out, &d)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &d
}

//func (ctx* RunContext) Report(ok bool, duration int, suiteId, file, testId, testName, message *string) {

func Run(tests *FilePayload, ctx *RunContext, id string) {
	for _, test := range tests.Tests {
		if test.Id == id {
			var result = _run(test.Name, ctx.RootPath, test.File)
			if result != nil {
				if result.Status == "passed" {
					ctx.Report(true, result.Duration, test.TestSuites[0], test.File, test.Id, test.Name, "")
				} else {
					for _, entry := range result.Specs {
						if entry.Status == "excluded" {
							continue
						}
						ctx.Report(false, result.Duration, test.TestSuites[0], test.File, test.Id, test.Name, entry.FailedExpectations[0].Message)
						break
					}
				}
			}
		}
	}
}
func searchSuite(p *FilePayload, id string) *SerialisedSuite {
	for _, elem := range p.Suites {
		if elem.Id == id {
			return &elem
		}
	}
	return nil
}
func RunAll(d *FilePayload, ctx *RunContext) {
	var fileMap = make(map[string][]SerialisedTest)
	for _, test := range d.Tests {
		fileMap[test.File] = append(fileMap[test.File], test)
	}
	
	for key, tests := range fileMap {
		var result = _run("", ctx.RootPath, key)
		if result != nil {
			// lets search by exsting tests
			for _, test := range tests {
				var suite = searchSuite(d, test.TestSuites[0])
				if suite == nil {
					continue
				}
				var search = suite.Name + " " + test.Name
				for _, res := range result.Specs {
					if res.FullName != search {
						
						continue
					}
					if res.Status == "passed" {
						ctx.Report(true, res.Duration, suite.Id, test.File, test.Id, test.Name, "")
					} else if res.Status == "failed" {
						ctx.Report(false, res.Duration, test.TestSuites[0], test.File, test.Id, test.Name, res.FailedExpectations[0].Message)
					}
					break
				}
			}
		}
	}
}
