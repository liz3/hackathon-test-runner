package main

import (
	"strconv"
	gp "path"
)

type FileEntry struct {
	RelativePath string
	AbsolutePath string
}

type JasmineResult struct {
  Id string `json:"id"`
  Description string `json:"description"`
  FullName string `json:"fullName"`
  FailedExpectations []interface{} `json:"failedExpectations"`
  DeprecationWarnings []interface{} `json:"deprecationWarnings"`
  Duration int `json:"duration"`
  Properties interface{} `json:"properties"`
  Status string `json:"status"`
  Specs []struct {
      Id string `json:"id"`
      Description string `json:"description"`
      FullName string `json:"fullName"`
      FailedExpectations []struct {
          MatcherName string `json:"matcherName"`
          Message string `json:"message"`
          Stack string `json:"stack"`
          Passed bool `json:"passed"`
          Expected interface{} `json:"expected"`
          Actual interface{} `json:"actual"`
        } `json:"failedExpectations"`
      PassedExpectations []interface{} `json:"passedExpectations"`
      DeprecationWarnings []interface{} `json:"deprecationWarnings"`
      PendingReason string `json:"pendingReason"`
      Duration int `json:"duration"`
      Properties interface{} `json:"properties"`
      DebugLogs interface{} `json:"debugLogs"`
      Status string `json:"status"`
    } `json:"specs"`
}

type Test struct {
	Id string
	Name string
}
type Suite struct {
	Id string
	Tests []*Test
	Name string
	File string
}

type SerialisedSuite struct {
	Id string `json:"id"`
	Name string `json:"label"`
	File string `json:"file"`
}
type SerialisedTest struct {
	Id string `json:"id"`
	Name string `json:"label"`
	File string `json:"file"`
	TestSuites []string `json:"test-suites"` // lol??
}

type FilePayload struct {
	Tests []SerialisedTest `json:"testCases"`
	Suites []SerialisedSuite `json:"testSuites"`
}

type Context struct {
	id int
	Suits []*Suite
	Files []FileEntry
	RootPath string
}

func CreateContext(path string) *Context {
	var ctx = Context{
		RootPath: path,
		Suits: make([]*Suite, 0),
		id: 0,
	}
	return &ctx
}
func (ctx *Context) GetID() string {
	 ctx.id += 1
	return "t-" + strconv.Itoa(ctx.id)
}
func (ctx *Context) RegisterFile(path string) string {
	e := FileEntry{
		RelativePath: path,
		AbsolutePath: gp.Join(ctx.RootPath, path),
	}
	ctx.Files = append(ctx.Files, e)
	return e.RelativePath
}

func (ctx *Context) CreateSuite(name string) *Suite {
	ctx.id += 1
	var id = "s-" + strconv.Itoa(ctx.id)
	var s = Suite{
		Id: id,
		Name: name,
		File: ctx.Files[len(ctx.Files)-1].RelativePath,

	}
	ctx.Suits = append(ctx.Suits, &s)

	return &s
}
func (ctx* Context) SerialisedForm() FilePayload {
	tests := make([]SerialisedTest, 0)
	suites := make([]SerialisedSuite, 0)

	for _, suite := range ctx.Suits {
		for _, test := range suite.Tests {
			tests = append(tests, SerialisedTest{
				Id: test.Id, 
				Name: test.Name,
				File: suite.File,
				TestSuites: []string{suite.Id},
			})
		}
		suites = append(suites, SerialisedSuite{
			Id: suite.Id,
			Name: suite.Name,
			File: suite.File,
		})
	}
	return FilePayload{
		Tests: tests,
		Suites: suites,
	}
}
func (ctx *Context) CreateTest(suite *Suite, name string) *Test {
	var t = Test{
		Id: ctx.GetID(),
		Name: name,
	}
	suite.Tests = append(suite.Tests, &t)
	return &t
}