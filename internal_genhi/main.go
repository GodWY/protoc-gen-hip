package internal_genhi

import (
	"fmt"
	"github.com/GodWY/protoc-gen-hip/version"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	ImportPath = "import"
)

// SupportedFeatures reports the set of supported protobuf language features.
var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// GenerateFile generates the contents of a .pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	// filename := file.GeneratedFilenamePrefix + ".pb.go"
	// g := gen.NewGeneratedFile(filename, file.GoImportPath)
	// f := newFileInfo(file)

	// genStandaloneComments(g, f, int32(genid.FileDescriptorProto_Syntax_field_number))
	// genGeneratedHeader(gen, g, f)
	// genStandaloneComments(g, f, int32(genid.FileDescriptorProto_Package_field_number))

	// packageDoc := genPackageKnownComment(f)
	// g.P(packageDoc, "package ", f.GoPackageName)
	// g.P()

	// // Emit a static check that enforces a minimum version of the proto package.
	// if GenerateVersionMarkers {
	// 	g.P("const (")
	// 	g.P("// Verify that this generated code is sufficiently up-to-date.")
	// 	g.P("_ = ", protoimplPackage.Ident("EnforceVersion"), "(", protoimpl.GenVersion, " - ", protoimplPackage.Ident("MinVersion"), ")")
	// 	g.P("// Verify that runtime/protoimpl is sufficiently up-to-date.")
	// 	g.P("_ = ", protoimplPackage.Ident("EnforceVersion"), "(", protoimplPackage.Ident("MaxVersion"), " - ", protoimpl.GenVersion, ")")
	// 	g.P(")")
	// 	g.P()
	// }

	// for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
	// 	genImport(gen, g, f, imps.Get(i))
	// }
	// for _, enum := range f.allEnums {
	// 	genEnum(g, f, enum)
	// }
	// for _, message := range f.allMessages {
	// 	genMessage(g, f, message)
	// }
	// genExtensions(g, f)

	// genReflectFileDescriptor(gen, g, f)

	// ??????http??????

	return generateFile(gen, file, false)
}

func genImport(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, imp protoreflect.FileImport) {
	impFile, ok := gen.FilesByPath[imp.Path()]
	if !ok {
		return
	}
	if impFile.GoImportPath == f.GoImportPath {
		// Don't generate imports or aliases for types in the same Go package.
		return
	}
	// Generate imports for all non-weak dependencies, even if they are not
	// referenced, because other code and tools depend on having the
	// full transitive closure of protocol buffer types in the binary.
	if !imp.IsWeak {
		g.Import(impFile.GoImportPath)
	}
	if !imp.IsPublic {
		return
	}

	// Generate public imports by generating the imported file, parsing it,
	// and extracting every symbol that should receive a forwarding declaration.
	impGen := GenerateFile(gen, impFile)
	impGen.Skip()
	b, err := impGen.Content()
	if err != nil {
		gen.Error(err)
		return
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", b, parser.ParseComments)
	if err != nil {
		gen.Error(err)
		return
	}
	genForward := func(tok token.Token, name string, expr ast.Expr) {
		// Don't import unexported symbols.
		r, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsUpper(r) {
			return
		}
		// Don't import the FileDescriptor.
		if name == impFile.GoDescriptorIdent.GoName {
			return
		}
		// Don't import decls referencing a symbol defined in another package.
		// i.e., don't import decls which are themselves public imports:
		//
		//	type T = somepackage.T
		if _, ok := expr.(*ast.SelectorExpr); ok {
			return
		}
		g.P(tok, " ", name, " = ", impFile.GoImportPath.Ident(name))
	}
	g.P("// Symbols defined in public import of ", imp.Path(), ".")
	g.P()
	for _, decl := range astFile.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					genForward(decl.Tok, spec.Name.Name, spec.Type)
				case *ast.ValueSpec:
					for i, name := range spec.Names {
						var expr ast.Expr
						if i < len(spec.Values) {
							expr = spec.Values[i]
						}
						genForward(decl.Tok, name.Name, expr)
					}
				case *ast.ImportSpec:
				default:
					panic(fmt.Sprintf("can't generate forward for spec type %T", spec))
				}
			}
		}
	}
	g.P()
}

// generateFile generates a _http.pb.go file containing kratos errors definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {

	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	f := newFileInfo(file)
	g.P("// Code generated by protoc-gen-go-hip. DO NOT EDIT.")
	g.P("// versions:", version.String())
	for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
		genImport(gen, g, f, imps.Get(i))
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g, omitempty)
	return g
}

// generateFileContent generates the kratos errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, omitempty bool) {
	if len(file.Services) == 0 {
		return
	}
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the kratos package it is being compiled against.")
	g.P("import (")
	g.P(`"net/http"`)
	g.P(`"github.com/GodWY/gutil"`)
	paths := genImportService(gen, file, g, omitempty)
	for _, path := range paths {
		g.P(`"`, path, `"`)
	}
	g.P()
	g.P(`"github.com/gin-gonic/gin"`)
	//g.P(`"github.com/GodWY/hip/service"`)
	g.P(")")

	// ??????????????????????????????

	// ??????????????????gin???????????????
	for _, service := range file.Services {
		// ???????????????
		serviceComments := string(service.Comments.Leading)
		// ??????????????????middleWare
		var commentsArr = []string{}
		if strings.Contains(serviceComments, "@") {
			// middleWares = parserComment(serviceComments)
			commentsArr = strings.Split(serviceComments, "@")
		}
		si := parserComment(commentsArr)
		if si.ProjectUri == "" {
			si.ProjectUri = "/api"
		}
		g.P("// generated http method")
		g.P("func register", service.GoName, "HttpHandler(srv *gin.Engine, srvs ", service.GoName, "HttpHandler) {")
		if len(si.MiddleWire) > 0 {
			g.P(`   group := srv.Group("`, si.ProjectUri, "/"+strings.ToLower(service.GoName), `" ,`, si.MiddleWire, ")")
		} else {
			g.P(`   group := srv.Group("`, si.ProjectUri, "/"+strings.ToLower(service.GoName), `" )`)
		}
		for _, value := range service.Methods {
			// g.Annotate(value.GoName, value.Location)
			// ????????????
			// leadingComments := appendDeprecationSuffix("",
			// 	value.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated())
			a := string(value.Comments.Leading)
			method, middleWares := hasPathPrefix(a)
			// method := strings.Split(a, "@")
			// valuemiddleWares := parserComment(method[2:])
			methods := strings.Trim(method, "\r\n")
			prefix := value.GoName[0:1]
			last := value.GoName[1:]
			x := strings.ToLower(prefix)
			path := fmt.Sprintf("%v%v", x, last)
			if len(middleWares) > 0 {
				// ??????middle
				var valuemiddleWares string
				for i := 0; i < len(middleWares); i++ {
					m := strings.Trim(middleWares[i], "\r\n")
					valuemiddleWares = valuemiddleWares + m + ","
				}

				g.P("group.", methods, `("/v1/`, path, `"`, ", "+valuemiddleWares, ` srvs.`, value.GoName, ")")
			} else {
				g.P("group.", methods, `("/v1/`, path, `", srvs.`, value.GoName, ")")
			}

		}
		g.P("}")
		g.P()
	}

	// ??????????????????????????????
	for _, service := range file.Services {
		g.P("var T", service.GoName, " ", service.GoName)
		g.P()
		g.P("func Register", service.GoName, "HttpHandler(srv *gin.Engine,", "srvs ", service.GoName, ") {")
		g.P("  tmp := new(", "xxx_", service.GoName, ")")
		g.P("  register", service.GoName, "HttpHandler(srv, tmp)")
		g.P("  T", service.GoName, "=srvs")
		g.P("}")
		g.P()
		// ????????????
		g.P("type ", service.GoName, " interface {")
		genService(gen, file, g, service, omitempty)
		g.P(" }")
	}
	g.P("// generated http handle")

	// ??????????????????
	for _, service := range file.Services {
		g.P("type ", service.GoName, "HttpHandler", " interface {")
		genHttpService(gen, file, g, service, omitempty)
		g.P(" }")
		g.P()
	}

	// ????????????
	// genHttpService implents
	for _, service := range file.Services {
		g.P("type ", "xxx_", service.GoName, " struct {")
		g.P("}")
		g.P()
		genXService(gen, file, g, service, omitempty)
		g.P()
	}

}

// hasMiddleware ????????????????????????middleware
// hasPathPrefix ??????????????????
func hasPathPrefix(comm string) (method string, middlewares []string) {
	xx := strings.Split(comm, "@")
	method = ""
	middlewares = []string{}
	for _, xx := range xx {
		tags := strings.Split(xx, ":")
		if len(tags) > 2 {
			continue
		}
		tag := strings.TrimSpace(tags[0])
		switch tag {
		case "method":
			// ??????gin?????????, ????????????
			method = strings.TrimSpace(strings.ToUpper(tags[1]))
			method = httpMethod(method)
		case "middle":
			middlewares = append(middlewares, tags[1:]...)
		}
	}
	return
}

func genXService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	for _, value := range service.Methods {
		g.P("func (xx *xxx_", service.GoName, ")", value.GoName, "(ctx *gin.Context)", "{")
		g.P("  req := &", value.Input.GoIdent, "{}")
		g.P("  if ok := ctx.Bind(req); ok != nil {")
		g.P(`detail:= "bind request error"`)
		g.P("rt:=gutil.RetFail(10000,detail)")
		g.P(` 	  ctx.JSON(http.StatusOK, rt)`)
		g.P("     return")
		g.P("   }")
		g.P("   rsp, err := ", "T", service.GoName, ".", value.GoName, "(ctx,req)")
		g.P("")
		g.P("   if err != nil {")
		g.P("	  ctx.JSON(http.StatusOK, gutil.RetError(err))")
		g.P("     return")
		g.P("    }")
		g.P(` ctx.JSON(http.StatusOK, gutil.RetSuccess("success",rsp))`)
		g.P("}")
		g.P()
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	for _, value := range service.Methods {
		g.P(value.GoName, "(ctx *gin.Context,in ", "*", value.Input.GoIdent, ")", "(out ", "*", value.Output.GoIdent, ",err error", " )")
	}
}

func genHttpService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	for _, value := range service.Methods {

		g.P(value.GoName, "(ctx ", "*", "gin.Context", ")")
	}
}

// genImportService ???????????????????????????
func genImportService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, omitempty bool) []string {
	// ??????????????????
	importPath := []string{}
	for _, v := range file.Services {
		// ????????????
		serviceComments := string(v.Comments.Leading)
		if serviceComments == "" {
			continue
		}
		var commentsArr = []string{}
		if strings.Contains(serviceComments, "@") {
			// middleWares = parserComment(serviceComments)
			commentsArr = strings.Split(serviceComments, "@")
		}
		for _, c := range commentsArr {
			c = strings.TrimSpace(c)
			c = strings.TrimLeft(c, "\r\n")
			c = strings.TrimRight(c, "\r\n")
			if strings.HasPrefix(c, "import") {
				importPath = append(importPath, strings.Split(c, ":")[1])
				continue
			}
		}
	}
	return importPath
}

const deprecationComment = "// Deprecated: Do not use."

// ServiceInfo ???????????????????????????????????????????????????
type ServiceInfo struct {
	// ?????????
	ProjectUri string
	// ?????????
	MiddleWire string
	// ????????????
	ImportPath []string
}

// parserComment ????????????????????????
func parserComment(comment []string) *ServiceInfo {
	si := &ServiceInfo{}
	var middleware string
	for _, c := range comment {
		c = strings.TrimSpace(c)
		c = strings.TrimLeft(c, "\r\n")
		c = strings.TrimRight(c, "\r\n")
		if strings.HasPrefix(c, "root") {
			si.ProjectUri = strings.Split(c, ":")[1]
			continue
		}
		if strings.HasPrefix(c, "middle") {
			middleware = strings.Split(c, ":")[1]
		}

	}
	si.MiddleWire = middleware
	return si

}

// httpMethod ??????gin???????????????????????????GET???POST
func httpMethod(method string) string {
	newMethod := "GET"
	switch method {
	case "GET", "POST":
		return method
	case "ANY":
		return "Any"
	}
	return newMethod
}
