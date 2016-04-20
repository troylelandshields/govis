package parser

import (
	"fmt"
	"go/ast"
	"go/format"
	goParser "go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/troylelandshields/govis/fGraph"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func validSyntax(file []byte) bool {
	_, err := format.Source(file)

	if err == nil {
		return true
	}

	return false
}

func goParseFile(path string, src interface{}) (ast *ast.File) {
	fset := token.NewFileSet()
	ast, err := goParser.ParseFile(fset, path, src, 0)
	check(err)
	return
}

func parseImportDeclarations(importSpecs []ast.Spec) []string {
	imports := make([]string, len(importSpecs))
	for i, spec := range importSpecs {
		importSpec := spec.(*ast.ImportSpec)
		importString, _ := strconv.Unquote(importSpec.Path.Value)
		imports[i] = strings.ToLower(importString)
	}

	return imports
}

//ParseFile a go file into something that we can analyze
func ParseFile(path string, src []byte, fG fGraph.FunctionGraph) (fGraph.FunctionGraph, []chan bool) {

	//convert file to byte array.
	if len(src) == 0 && path != "" {
		fileInBytes, err := ioutil.ReadFile(path)
		check(err)

		if validSyntax(fileInBytes) == false {
			fmt.Printf("File %s has invalid go syntax", path)
		}

		src = fileInBytes
	} else if len(src) == 0 && path == "" {
		panic("Huh??")
	}

	syntaxTree := goParseFile(path, src)

	dones := []chan bool{}
	for _, declaration := range syntaxTree.Decls {
		switch declType := declaration.(type) {
		case *ast.GenDecl:
			genericDeclaration := (*ast.GenDecl)(declType)
			switch genericDeclaration.Tok {
			// case token.IMPORT:
			// 	// packageInfo.addImports(parseImportDeclarations(genericDeclaration.Specs))
			// case token.TYPE:
			// 	// packageInfo.addTypes(parseTypeDeclarations(genericDeclaration.Specs))
			// case token.VAR:
			// 	// packageInfo.addVars(parseVariableDeclarations(genericDeclaration.Specs, "var"))
			// case token.CONST:
			// packageInfo.addVars(parseVariableDeclarations(genericDeclaration.Specs, "const"))
			default:
				//fmt.Println("unknown generic declaration")
			}
		case *ast.FuncDecl:
			funcDecl := (*ast.FuncDecl)(declType)

			fNode, done := createFuncNode(funcDecl, src, fG)
			fG.AddFunctionNode(fNode)

			dones = append(dones, done)
		default:
			//fmt.Println("unknown declaration")
		}
	}

	return fG, dones
}

func createFuncNode(funcDecl *ast.FuncDecl, fileInBytes []byte, fG fGraph.FunctionGraph) (fGraph.FunctionNode, chan bool) {

	name := funcDecl.Name.Name
	signature := string(fileInBytes[funcDecl.Pos()-1 : funcDecl.Body.Lbrace-1])

	functionNode := fGraph.NewFunctionNode(name, signature)

	done := make(chan bool)

	go func() {
		funcCallNames := getFuncCallNames(funcDecl, fileInBytes)
		select {
		case <-done:
			for _, fName := range funcCallNames {
        fmt.Printf("[%s] calls [%s]\n", functionNode.Name(), fName)
				calledFNode := fG.GetFunctionNode(fName)
				functionNode.AddCall(calledFNode)
			}
		}
	}()

	return functionNode, done
}

func getFuncCallNames(funcDecl *ast.FuncDecl, src []byte) []string {
	callNames := []string{}

	for _, stmt := range funcDecl.Body.List {
		switch stmtType := stmt.(type) {
		case *ast.ExprStmt:
			exprStmt := (*ast.ExprStmt)(stmtType)
			callWithSignature := string(src[exprStmt.X.Pos()-1 : exprStmt.X.End()-1])

			r, _ := regexp.Compile("(.*)\\(.*\\)")

			matches := r.FindAllStringSubmatch(callWithSignature, -1)

			if len(matches) != 1 {
				panic("Uhhh, what?")
			}

			callWithoutSignature := matches[0][1]
			callNames = append(callNames, callWithoutSignature)

		// case *ast.SendStmt:
		// 	/fmt.Println("Send stmt", stmtType)
		// case *ast.LabeledStmt:
		// 	fmt.Println("Labeled stmt", stmtType)
		// case *ast.BlockStmt:
		// 	fmt.Println("block stmt", stmtType)
		// default:
		// 	fmt.Println("Unknown stmt type", stmtType)
		}
	}

	return callNames
}
