package fGraph

import (
	"encoding/json"
	"fmt"
)

type functionGraph struct {
	functions map[string]FunctionNode
}

//FunctionGraph contains a graph of functions
type FunctionGraph interface {
	Functions() []FunctionNode
	AddFunctionNode(FunctionNode)
	GetFunctionNode(string) FunctionNode
	ToString() string
	ToJSON() ([]byte, error)
}

//Functions returns function
func (fGraph *functionGraph) Functions() []FunctionNode {
	v := make([]FunctionNode, 0, len(fGraph.functions))

	for _, value := range fGraph.functions {
		v = append(v, value)
	}

	return v
}

//AddFunctionNode adds a function node to the graph
func (fGraph *functionGraph) AddFunctionNode(fNode FunctionNode) {
	fGraph.functions[fNode.Name()] = fNode
}

//GetFunctionNode returns function node for given name
func (fGraph *functionGraph) GetFunctionNode(fName string) FunctionNode {
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

func (fGraph *functionGraph) ToJSON() (str []byte, err error) {

	type fNodeJSON struct {
		ID    string `json:"id"`
		Label string `json:"label"`
    Shape string `json:"shape"`
	}
  
	type fEdgeJSON struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	type fGraphJSON struct {
		Nodes []fNodeJSON `json:"nodes"`
		Edges []fEdgeJSON `json:"edges"`
	}

	fNodesJSON := []fNodeJSON{}
	fEdgesJSON := []fEdgeJSON{}

	for _, fNode := range fGraph.functions {
		for _, call := range fNode.Calls() {

			fEdgesJSON = append(fEdgesJSON, fEdgeJSON{
				From: fNode.Name(),
				To:   call.Name(),
			})

		}

		fNodesJSON = append(fNodesJSON, fNodeJSON{
			ID:    fNode.Name(),
			Label: fNode.Signature(),
      Shape: "square",
		})

	}

	fGJSON := fGraphJSON{
		Nodes: fNodesJSON,
		Edges: fEdgesJSON,
	}

	data, err := json.MarshalIndent(fGJSON, "", "\t")
	return data, err
}
