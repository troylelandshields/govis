package fGraph

import "fmt"

type functionGraph struct {
	functions   map[string]FunctionNode
}

//FunctionGraph contains a graph of functions
type FunctionGraph interface {
  Functions() []FunctionNode
  AddFunctionNode(FunctionNode)
  GetFunctionNode(string) FunctionNode
  ToString() string
}

//Functions returns function
func (fGraph *functionGraph) Functions() []FunctionNode {
  v := make([]FunctionNode, 0, len(fGraph.functions))
  
	for  _, value := range fGraph.functions {
    v = append(v, value)
  }
  
  return v
}

//AddFunctionNode adds a function node to the graph
func (fGraph *functionGraph) AddFunctionNode(fNode FunctionNode) {
  fGraph.functions[fNode.Name()] = fNode
}

//GetFunctionNode returns function node for given name
func (fGraph *functionGraph) GetFunctionNode(fName string) FunctionNode{
  return fGraph.functions[fName]
}

//NewFunctionGraph returns a new Function Graph
func NewFunctionGraph() FunctionGraph {
  fGraph := &functionGraph{}
  
  fGraph.functions = make(map[string]FunctionNode)
  
  return fGraph
}

//ToString prints this damn thing
func (fGraph *functionGraph) ToString() (str string) {
	str += fmt.Sprintf("functions:\n")
	for _, fun := range fGraph.functions {
    if fun != nil {
		  str += fun.ToString()
    }
	}
	str += "\n"
	return str
}
