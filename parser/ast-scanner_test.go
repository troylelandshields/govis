package parser

import (
  "fmt"
	"testing"
  "github.com/troylelandshields/govis/fGraph"
)

func TestParseFile_Simple_Function_Call(t *testing.T) {
  	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
  Exported()
}

func Exported() {
  println("Hello, World!")
}
`
  fGraph := fGraph.NewFunctionGraph()

  parsedCode, dones := ParseFile("src", []byte(src), fGraph)

  for _, done := range dones {
    done <- true
  }

  funcs := parsedCode.Functions()
  
  if len(funcs) != 2 {
    t.Fatal("Expected:", 2, "; Actual:", len(funcs))
  }
  
  mainFunc := funcs[0]
  if mainFunc.Signature() != "func main() " {
    t.Fatalf("Expected: [%s], Actual: [%s]", "func main() ", mainFunc.Signature())
  }
  
  exportedFunc := funcs[1]
  if exportedFunc.Signature() != "func Exported() " {
    t.Fatalf("Expected: [%s], Actual: [%s]", "func Exported() ", exportedFunc.Signature())
  }
  
  mainFuncCalls := mainFunc.Calls()
  if len(mainFuncCalls) != 1 {
    t.Fatalf("Expected: [%d], Actual: [%d]", 1, len(mainFuncCalls))
  }
  
  mainFuncCallExported := mainFuncCalls[0]
  if mainFuncCallExported != exportedFunc {
    t.Fatalf("Expected: [%s], Actual: [%s]", exportedFunc, mainFuncCallExported)
  }
  
  fmt.Println(parsedCode.ToString())
}


func TestParseFile_Simple_Method_Call(t *testing.T) {
  	// src is the input for which we want to print the AST.
	src := `
package main

type example struct {}

func main() {
  e := &example{}
  
  e.Exported()
}

func (e *example) Exported() {
  println("Hello, World!")
}
`
  fGraph := fGraph.NewFunctionGraph()

  parsedCode, dones := ParseFile("src", []byte(src), fGraph)

  for _, done := range dones {
    done <- true
  }

  funcs := parsedCode.Functions()
  
  if len(funcs) != 2 {
    t.Fatal("Expected:", 2, "; Actual:", len(funcs))
  }
  
  mainFunc := funcs[0]
  if mainFunc.Signature() != "func main() " {
    t.Fatalf("Expected: [%s], Actual: [%s]", "func main() ", mainFunc.Signature())
  }
  
  exportedFunc := funcs[1]
  if exportedFunc.Signature() != "func (e *example) Exported() " {
    t.Fatalf("Expected: [%s], Actual: [%s]", "func Exported() ", exportedFunc.Signature())
  }
  
  mainFuncCalls := mainFunc.Calls()
  if len(mainFuncCalls) != 1 {
    t.Fatalf("Expected: [%d], Actual: [%d]", 1, len(mainFuncCalls))
  }
  
  mainFuncCallExported := mainFuncCalls[0]
  if mainFuncCallExported != exportedFunc {
    t.Fatalf("Expected: [%s], Actual: [%s]", exportedFunc, mainFuncCallExported)
  }
  
  fmt.Println(parsedCode.ToString())
}