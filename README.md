# Hackathon test runner
I did particiipate in a hackathon ig

https://github.com/LambdaTest/hackathon/blob/main/README.md

## How i did it
This is a test runner for the Javascript jasmine test framework in golang.

It uses [otto](https://github.com/robertkrimen/otto), the go runtime for . but ONLY its parser, it doesnt execute anything through it.  
i included the folder in vendor because i edited otto!
I made the parser support const, let and arrow function to the point, that the parser is happy with it and gives me a valid AST, which was the most fun part in this project,
by default otto only supports es5, see the test cases which now can use arrow functons!.  

Then i search the AST for all test cases and save them in the file as required, running then happens via executing jasmine through a sub process, it then prints the result to console

Generate test cases for all JS files in test(excluding node_modules):
```
go run . -path tests
```
this will generate out.json

then to run the test cases use
```
go run . -action run -path tests
```
