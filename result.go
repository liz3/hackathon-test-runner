package main

const failed string = ("One of more test-cases failed!")

type TestRun struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Message string `json:"failureMsg,omitempty"`
	Duration int `json:"duration"`
}
type SuiteRun struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Duration int `json:"duration"`
	TestCount int `json:"numTests"`
	File string `json:"file"`
	Message string `json:"failureMsg,omitempty"`
}

type RunPair struct {
	Suite *SuiteRun
	Runs []*TestRun
}
type RunResultSerialise struct {
	Runs []*TestRun `json:"testCases"`
	Suites []*SuiteRun `json:"testSuites"`
}
type RunContext struct {
	Suites map[string]*RunPair
	RootPath string
}

func CreateRunContext(p string) *RunContext {
	return &RunContext{RootPath:p,Suites: make(map[string]*RunPair),}
}
func (ctx* RunContext) Serialise() *RunResultSerialise {
	var res = RunResultSerialise{}
	for _,v := range ctx.Suites {
		for _, t := range v.Runs {
		res.Runs = append(res.Runs, t)
		}
		res.Suites = append(res.Suites, v.Suite)
	}
	return &res
}
func (ctx* RunContext) EnsureSuite(suiteId, file string) *RunPair {
	if val, ok := ctx.Suites[suiteId]; ok {
    	return val
    }
    ctx.Suites[suiteId] = &RunPair{
    		Suite: &SuiteRun{
    			Id: suiteId,
    			Status: "passed",
    			Duration: 0,
    			TestCount: 0,
    			Message: "",
    			File: file,
    		},
    		Runs: make([]*TestRun, 0),
    	}	
   	return ctx.Suites[suiteId]
}
func (ctx* RunContext) Report(ok bool, duration int, suiteId, file, testId, testName string, message string) {
	var suite = ctx.EnsureSuite(suiteId, file)
	var status = "passed"
	suite.Suite.TestCount++;
	if !ok {
		suite.Suite.Message = failed
		suite.Suite.Status = "failed"
		status = "failed"
	}
	suite.Suite.Duration += duration

	suite.Runs = append(suite.Runs, &TestRun{
		Id: testId,
		Duration: duration,
		Message: message,
		Status: status,
	})
}