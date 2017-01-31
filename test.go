package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	// "template"

	"log"
	// "os"
)

// var (
// 	generatedTemplate = template.Must(template.New("render").Parse(`// generated by joiner -- DO NOT EDIT
// package {{.Package}}
// import (
// 	"fmt"
// 	"strings"
// )
// {{range .Types}}
// func (t {{.Name}}) Map( ) string {
// 	return fmt.Sprintf("%#v", t)
// }
// {{end}}
// type Join{{.Name}} []{{.Name}}
// func (j Join{{.Name}}) With(sep string) string {
// 	all := make([]string, 0, len(j))
// 	for _, s := range j {
// 		all = append(all, s.String())
// 	}
// 	return strings.Join(all, sep)
// }
// {{end}}`))
// )

/*
CollectionType contains all the information about an annotated collection type
required to generate Map, Reduce, and Filter.

Ex.
type Users []*User
CollectionType{
  Ident: "Users",
  TypeExpr: "[]*User"
}

We parse everything into strings because the problem is simple enough that we
can use the template package to generate the code that we want instead of
constructing manually constructing an AST tree.
*/
type CollectionType struct {
	Ident    string
	TypeExpr string
}

func parseFile(inputPath string) []*CollectionType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Could not parse file: %s", err)
	}

	// ast.Print(fset, f)
	packageName := identifyPackage(f)
	if packageName == "" {
		log.Fatalf("Could not determine package name of %s", inputPath)
	}

	typeSpecs := getTypeSpecsFromAST(f)
	collectionTypes := []*CollectionType{}
	for _, typeSpec := range typeSpecs {
		collectionTypes = append(
			collectionTypes,
			getCollectionTypeFromTypeSpec(typeSpec, fset),
		)
	}

	return collectionTypes
}

func getTypeSpecsFromAST(tree *ast.File) []*ast.TypeSpec {
	typeSpecs := []*ast.TypeSpec{}

	for _, decl := range tree.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		isAnnotated := isAnnotatedTypeDecl(genDecl)
		if !isAnnotated {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			hasName := typeSpec.Name != nil
			_, isArrayType := typeSpec.Type.(*ast.ArrayType)
			if !hasName || !isArrayType {
				continue
			}

			typeSpecs = append(typeSpecs, typeSpec)
		}
	}
	return typeSpecs
}

func isAnnotatedTypeDecl(decl *ast.GenDecl) bool {
	// Not a type declaration
	if decl.Tok != token.TYPE {
		return false
	}

	// Doesn't have an annotation
	if decl.Doc == nil {
		return false
	}

	for _, comment := range decl.Doc.List {
		if strings.Contains(comment.Text, "slice++") {
			return true
		}
	}

	return false
}

func getCollectionTypeFromTypeSpec(
	typeSpec *ast.TypeSpec,
	fset *token.FileSet,
) *CollectionType {
	typeIdent := typeSpec.Name.Name
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, typeSpec.Type); err != nil {
		panic(err)
	}
	typeExpr := buf.String()
	return &CollectionType{
		Ident:    typeIdent,
		TypeExpr: typeExpr,
	}
}

func identifyPackage(f *ast.File) string {
	if f.Name == nil {
		return ""
	}
	return f.Name.Name
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("joiner: ")

	collectionTypes := parseFile("test1.go")
	for _, collectionType := range collectionTypes {
		log.Printf("%#v\n", collectionType)
	}
}