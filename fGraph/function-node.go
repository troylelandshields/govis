package fGraph

import (
	"fmt"
)

type functionNode struct {
	name      string
	signature string
	calls     []FunctionNode
}

//FunctionNode represents a function in a graph
type FunctionNode interface {
	Name() string
	Signature() string
	Calls() []FunctionNode
	AddCall(FunctionNode)
	ToString() string
}

func (f *functionNode) Name() string {
	return f.name
}

func (f *functionNode) Signature() string {
	return f.signature
}

func (f *functionNode) Calls() []FunctionNode {
	return f.calls
}

func (f *functionNode) AddCall(calledNode FunctionNode) {
	if calledNode != nil {
		f.calls = append(f.calls, calledNode)
	}
}

func (f *functionNode) ToString() (str string) {
	str += fmt.Sprintf("%s->\n[\n", f.Signature())
	for _, called := range f.calls {
		if called != nil {
			str += fmt.Sprintf("\t %s \n", called.Name())
		}
	}
	str += "\n]\n"
	return str
}

//NewFunctionNode creates a new function node
func NewFunctionNode(name, signature string) FunctionNode {
	functionNode := &functionNode{}
	functionNode.calls = []FunctionNode{}
	functionNode.name = name
	functionNode.signature = signature

	return functionNode
}
