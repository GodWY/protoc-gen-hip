package internal_genhi

import (
	"fmt"
	"github.com/GodWY/protoc-gen-hip/template_hip"
	"github.com/GodWY/protoc-gen-hip/version"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func genWithTemplate(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	f := newFileInfo(file)
	for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
		genImport(gen, g, f, imps.Get(i))
	}
	//generateFileContent(gen, file, g, omitempty)
	g.P(buildTemplate(gen, file, g))
	return g
}

func buildTemplate(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) string {
	tl := template_hip.NewHttp()
	imports := []string{
		`"net/http"`,
		`"github.com/GodWY/gutil"`,

		`"github.com/gin-gonic/gin"`,
	}

	tl.AddPackageName(string(file.GoPackageName))
	tl.AddVersion(version.String())
	for _, service := range file.Services {
		comments := string(service.Comments.Leading)
		serInfo := new(ServerInfo)
		serInfo.SpiltServiceComments(comments)
		if serInfo.Path == "" {
			serInfo.Path = "api/" + strings.ToLower(service.GoName)
		}
		sg := &template_hip.ServiceGroups{
			// 服务的名字
			ServiceName: service.GoName,
			Routers:     nil,
			GroupPath:   serInfo.Path,
			//MiddleWare:  serInfo.Middles[0],
			Leading: serInfo.Docs,
		}
		if len(serInfo.Middles) > 0 {
			sg.MiddleWare = serInfo.Middles[0]
		}

		routs := make([]*template_hip.Routers, 0, len(service.Methods))
		// 遍历路由
		for _, m := range service.Methods {
			mInfo := &MethodInfo{}
			mInfo.SpiltMethodComments(m.Comments.Leading.String())

			prefix := m.GoName[0:1]
			x := strings.ToLower(prefix)
			path := fmt.Sprintf("/%v%v", x, m.GoName[1:])
			rout := &template_hip.Routers{
				Path:     path,
				FuncName: m.GoName,
				In:       m.Input.GoIdent.GoName,
				Out:      m.Output.GoIdent.GoName,
				Method:   mInfo.Method,
				Leading:  mInfo.Docs,
			}
			if len(mInfo.Middles) > 0 {
				rout.MiddleWare = mInfo.Middles[0]
			}
			routs = append(routs, rout)
		}
		sg.Routers = routs
		tl.AddGroups(sg)
		imports = append(imports, serInfo.Imports...)

	}
	tl.AddImports(imports...)
	return tl.Execute()
}

// SpiltServiceComments 分割服务元数据
// @middle: 组路由
// @root: 自定义路由组。默认使用服务名称 /api/{定义的路由组}
// @import: 自定义导入包，此元组用于导入服务自定义的middle
func (s *ServerInfo) SpiltServiceComments(leading string) []string {
	var commentsArr []string
	// 分割元数据
	if strings.Contains(leading, "@") {
		// middleWares = parserComment(serviceComments)
		commentsArr = strings.Split(leading, "@")
	}
	var meta []Meta
	// 遍历数组
	for _, c := range commentsArr {
		c = strings.TrimSpace(c)
		c = strings.TrimLeft(c, "\r\n")
		c = strings.TrimRight(c, "\r\n")
		if strings.HasPrefix(c, "root") {
			s.Path = strings.Split(c, ":")[1]
		}

		if strings.HasPrefix(c, "middle") {
			middleware := strings.Split(c, ":")[1]
			s.Middles = append(s.Middles, middleware)
		}

		if strings.HasPrefix(c, "import") {
			imports := strings.Split(c, ":")[1]
			s.Imports = append(s.Imports, imports)
		}

		if strings.HasPrefix(c, "doc") {
			s.Docs = strings.Split(c, ":")[1]
		}
		s.Meta = meta
	}
	if len(s.Middles) > 0 {
		finalMiddle := strings.Join(s.Middles, ",")
		s.Middles = []string{finalMiddle}
	}
	return nil
}

// SpiltMethodComments 分割服务元数据
// @middle: 组路由
// @root: 自定义路由组。默认使用服务名称 /api/{定义的路由组}
// @import: 自定义导入包，此元组用于导入服务自定义的middle
func (m *MethodInfo) SpiltMethodComments(leading string) []string {
	var commentsArr []string
	// 分割元数据
	if strings.Contains(leading, "@") {
		// middleWares = parserComment(serviceComments)
		commentsArr = strings.Split(leading, "@")
	}
	var meta []Meta
	// 遍历数组
	for _, c := range commentsArr {
		c = strings.TrimSpace(c)
		c = strings.TrimLeft(c, "\r\n")
		c = strings.TrimRight(c, "\r\n")
		if strings.HasPrefix(c, "middle") {
			middleware := strings.Split(c, ":")[1]
			m.Middles = append(m.Middles, middleware)
		}

		if strings.HasPrefix(c, "import") {
			imports := strings.Split(c, ":")[1]
			m.Imports = append(m.Imports, imports)
		}
		if strings.HasPrefix(c, "method") {
			method := strings.Split(c, ":")[1]
			m.Method = httpMethod(method)
		}
		if strings.HasPrefix(c, "doc") {
			m.Docs = strings.Split(c, ":")[1]
		}
		m.Meta = meta
	}
	if len(m.Middles) > 0 {
		finalMiddle := strings.Join(m.Middles, ",")
		m.Middles = []string{finalMiddle}
	}
	return nil
}

type MethodInfo struct {
	// Meta 元数据
	Meta []Meta
	// 服务名字
	Name string
	// Path 路径
	Path    string
	Middles []string
	Imports []string
	Method  string
	Docs    string
}
type ServerInfo struct {
	// Meta 元数据
	Meta []Meta
	// 服务名字
	Name string
	// Path 路径
	Path    string
	Middles []string
	Imports []string
	Docs    string
}

type Meta struct {
	Key   string
	Value []string
}
