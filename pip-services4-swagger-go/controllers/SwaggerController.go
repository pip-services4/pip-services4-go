package controllers

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cctrls "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	"github.com/rakyll/statik/fs"
	"goji.io/pattern"

	_ "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/resources"
)

type SwaggerController struct {
	*cctrls.RestController
	routes map[string]string
	fs     http.FileSystem
}

func NewSwaggerController() *SwaggerController {
	c := SwaggerController{}
	c.RestController = cctrls.InheritRestController(&c)
	c.BaseRoute = "swagger"
	c.routes = map[string]string{}

	sfs, err := fs.NewWithNamespace("swagger")
	if err != nil {
		panic(err)
	}
	c.fs = sfs

	return &c
}

func (c *SwaggerController) calculateContentType(fileName string) string {
	if strings.HasSuffix(fileName, ".html") {
		return "text/html"
	}
	if strings.HasSuffix(fileName, ".css") {
		return "text/css"
	}
	if strings.HasSuffix(fileName, ".js") {
		return "application/javascript"
	}
	if strings.HasSuffix(fileName, ".png") {
		return "image/png"
	}
	return "text/plain"
}

func (c *SwaggerController) getSwaggerFile(res http.ResponseWriter, req *http.Request) {
	var vars map[string]any = make(map[string]any, 0)
	if reqVars, ok := req.Context().Value(pattern.AllVariables).(map[pattern.Variable]any); ok {
		for k, v := range reqVars {
			vars[string(k)] = v
		}
	}

	fileName := ""
	if fname, ok := vars["file_name"]; ok {
		if val, ok := fname.(string); ok {
			fileName = strings.ToLower(val)
		}
	}

	r, err := c.fs.Open("/" + fileName)
	if err != nil {
		res.WriteHeader(404)
		io.WriteString(res, err.Error())
		return
	}
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		res.WriteHeader(500)
		io.WriteString(res, err.Error())
		return
	}

	res.Header().Set("Content-Length", cconv.StringConverter.ToString(len(content)))
	res.Header().Set("Content-Type", c.calculateContentType(fileName))
	res.WriteHeader(200)
	res.Write(content)
}

func (c *SwaggerController) getIndex(res http.ResponseWriter, req *http.Request) {
	r, err := c.fs.Open("/index.html")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	contentBytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	content := string(contentBytes)

	builder := strings.Builder{}
	builder.WriteString("[")
	for k, v := range c.routes {
		if builder.Len() > 1 {
			builder.WriteString(",")
		}
		builder.WriteString("{name:\"")
		name := strings.ReplaceAll(k, "\"", "\\\"")
		builder.WriteString(name)
		builder.WriteString("\",url:\"")
		url := strings.ReplaceAll(v, "\"", "\\\"")
		builder.WriteString(url)
		builder.WriteString("\"}")
	}
	builder.WriteString("]")

	content = strings.ReplaceAll(content, "[/*urls*/]", builder.String())

	res.Header().Add("Content-Type", "text/html")
	res.Header().Add("Content-Length", cconv.StringConverter.ToString(len(content)))
	res.WriteHeader(200)
	io.WriteString(res, content)
}

func (c *SwaggerController) redirectToIndex(res http.ResponseWriter, req *http.Request) {
	url := req.RequestURI
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	http.Redirect(res, req, url+"index.html", http.StatusSeeOther)
}

func (c *SwaggerController) composeSwaggerRoute(baseRoute string, route string) string {
	if baseRoute != "" {
		if route == "" {
			route = "/"
		}
		if !strings.HasPrefix(route, "/") {
			route = "/" + route
		}
		if !strings.HasPrefix(baseRoute, "/") {
			baseRoute = "/" + baseRoute
		}
		route = baseRoute + route
	}

	return route
}

func (c *SwaggerController) RegisterOpenApiSpec(baseRoute string, swaggerRoute string) {
	route := c.composeSwaggerRoute(baseRoute, swaggerRoute)
	if baseRoute == "" {
		baseRoute = "default"
	}
	c.routes[baseRoute] = route
}

func (c *SwaggerController) Register() {
	// A hack to redirect default base route
	baseRoute := c.BaseRoute
	c.BaseRoute = ""
	c.RegisterRoute(
		"get", baseRoute, nil, c.redirectToIndex,
	)
	c.BaseRoute = baseRoute

	c.RegisterRoute(
		"get", "/", nil, c.redirectToIndex,
	)

	c.RegisterRoute(
		"get", "/index.html", nil, c.getIndex,
	)

	c.RegisterRoute(
		"get", "/:file_name", nil, c.getSwaggerFile,
	)
}
