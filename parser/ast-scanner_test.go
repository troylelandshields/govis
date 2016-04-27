package parser

import (
	"fmt"
	"testing"

	"github.com/troylelandshields/govis/fGraph"
)

type testData struct {
	description string
	src         string
}

func TestMainCallsExported(t *testing.T) {
	tDatas := []*testData{
		&testData{description: "Simple Function Call", src: `
package main
func main() {
  Exported()
}

func Exported() {
  println("Hello, World!")
}
`},
		&testData{description: "Simple Method Call", src: `
package main

type example struct {}

func main() {
  e := &example{}
  
  e.Exported()
}

func (e *example) Exported() {
  println("Hello, World!")
}
`},
		&testData{description: "Simple Function Call Assigned to Variable", src: `
package main
func main() {
  e := Exported()
  println(e)
}

func Exported() string{
  return "Hello, World!"
}
`},
		&testData{description: "Simple Function Call Within Function Call", src: `
package main
func main() {
  println(Exported())
}

func Exported() string{
  return "Hello, World!"
}
`},
	}

  fmt.Printf("Running [%d] tests\n\n", len(tDatas))
	for _, singleTest := range tDatas {
    fmt.Printf("Running test: [%s]\n", singleTest.description)
		verifyMainCallsExported(t, singleTest)
	}
}

func verifyMainCallsExported(t *testing.T, tData *testData) {

	fGraph := fGraph.NewFunctionGraph()

	parsedCode, dones := ParseFile("src", []byte(tData.src), fGraph)

	for _, done := range dones {
		done <- true
	}

	expectFoundFuncsCount(t, parsedCode, 2)

	mainFunc := expectFunctionWithName(t, parsedCode, "main")

	exportedFunc := expectFunctionWithName(t, parsedCode, "Exported")

	expectFunctionCallCount(t, mainFunc, 1)

	expectFunctionNodesToEqual(t, mainFunc.Calls()[0], exportedFunc)
}

func expectFoundFuncsCount(t *testing.T, fG fGraph.FunctionGraph, expected int) {
	if len(fG.Functions()) != expected {
		t.Fatalf("Expected to find [%d] functions, but actually found [%d]\n", expected, len(fG.Functions()))
	}
}

func expectFunctionWithName(t *testing.T, fG fGraph.FunctionGraph, expected string) fGraph.FunctionNode {

	fNode := fG.GetFunctionNode(expected)

	if fNode == nil {
		t.Fatalf("Expected a function with signature [%s], actual signature is [%s]\n", expected, fNode.Signature())
	}

	return fNode
}

func expectFunctionNodesToEqual(t *testing.T, expectedFNode, actualFNode fGraph.FunctionNode) {
	if expectedFNode != actualFNode {
		t.Fatalf("Expected function node [%s] but actually have [%s]\n", expectedFNode, actualFNode)
	}
}

func expectFunctionCallCount(t *testing.T, fNode fGraph.FunctionNode, expected int) {
	if len(fNode.Calls()) != expected {
		t.Fatalf("Expected [%s] to call [%d] functions, actually called [%d]\n", fNode.Signature(), expected, len(fNode.Calls()))
	}
}
