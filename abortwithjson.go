package abortwithjson

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `find deprecated AbortWithStatusJSON usage instead of apierror.AbortWithError.
				For more information check https://gl.cmd.su/aura/apierror`

const (
	abortWithStatusJSONMethod = "AbortWithStatusJSON"
	ginImportPath             = `"github.com/gin-gonic/gin"`
)

var Analyzer = &analysis.Analyzer{
	Name:     "abortwithjson",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	importPath := getImport(pass)
	if importPath == nil {
		return nil, nil
	}

	findAbortWithStatusJSON(pass)

	return nil, nil
}

func getImport(pass *analysis.Pass) *string {
	var importPath *string

	inspctr := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}
	inspctr.Preorder(nodeFilter, func(node ast.Node) {
		importSpec := node.(*ast.ImportSpec)
		if importSpec.Path.Value == ginImportPath {
			importPath = &importSpec.Path.Value
			return
		}
	})

	return importPath
}

func findAbortWithStatusJSON(pass *analysis.Pass) {
	inspctr := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspctr.Preorder(nodeFilter, func(node ast.Node) {
		callExpr := node.(*ast.CallExpr)

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if selectorExpr.Sel.Name != abortWithStatusJSONMethod {
			return
		}

		if len(callExpr.Args) != 2 {
			return
		}

		ginCtx := selectorExpr.X.(*ast.Ident).Name
		pass.Report(analysis.Diagnostic{
			Pos:      selectorExpr.Pos(),
			Category: "abortwithjson",
			Message:  fmt.Sprintf("use of deprecated function %s. Use %s instead", ginCtx+"."+selectorExpr.Sel.Name, "apierror.AbortWithError"),
			SuggestedFixes: []analysis.SuggestedFix{
				analysis.SuggestedFix{
					Message: "change gin method to apierror library",
				},
			},
		})
		//pass.Reportf(selectorExpr.Pos(), "use of deprecated function %s. Use %s instead", ginCtx+"."+selectorExpr.Sel.Name, "apierror.AbortWithError")
	})
}
