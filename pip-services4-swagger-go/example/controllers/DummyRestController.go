package example_services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	cservices "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	data "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/data"
	logic "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/logic"
	"github.com/rakyll/statik/fs"

	_ "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/resources"
)

type DummyRestController struct {
	*cservices.RestController
	service logic.IDummyService
}

func NewDummyRestController() *DummyRestController {
	c := DummyRestController{}
	c.RestController = cservices.InheritRestController(&c)
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return &c
}

func (c *DummyRestController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)
}

func (c *DummyRestController) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.RestController.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr == nil && depRes != nil {
		c.service = depRes.(logic.IDummyService)
	} else {
		panic(depErr.Error())
	}
}

func (c *DummyRestController) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	result, err := c.service.GetPageByFilter(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		cquery.NewFilterParamsFromValue(params), // W! need test
		cquery.NewPagingParamsFromTuples(paginParams),
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) getOneById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)
	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}
	result, err := c.service.GetOneById(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummyId)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) create(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)
	var dummy data.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}

	result, err := c.service.Create(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummy,
	)
	c.SendCreatedResult(res, req, result, err)
}

func (c *DummyRestController) update(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	var dummy data.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}
	result, err := c.service.Update(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummy,
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) deleteById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}

	result, err := c.service.DeleteById(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummyId,
	)
	c.SendDeletedResult(res, req, result, err)
}

func (c *DummyRestController) Register() {
	statikFS, err := fs.NewWithNamespace("example")
	if err != nil {
		panic(err)
	}
	r, err := statikFS.Open("/dummies.yml")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	c.SwaggerRoute = "/dummies/swagger"
	c.RegisterOpenApiSpec(string(content))

	c.RegisterRoute(
		"get", "/dummies",
		cvalid.NewObjectSchema().
			WithOptionalProperty("key", cconv.String).
			WithOptionalProperty("skip", cconv.Long).
			WithOptionalProperty("take", cconv.Long).
			WithOptionalProperty("total", cconv.String).Schema,
		c.getPageByFilter,
	)

	c.RegisterRoute(
		"get", "/dummies/{dummy_id}",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.getOneById,
	)

	c.RegisterRoute(
		"post", "/dummies",
		cvalid.NewObjectSchema().
			WithRequiredProperty("body", data.NewDummySchema()).Schema,
		c.create,
	)

	c.RegisterRoute(
		"put", "/dummies",
		cvalid.NewObjectSchema().
			WithRequiredProperty("body", data.NewDummySchema()).Schema,
		c.update,
	)

	c.RegisterRoute(
		"delete", "/dummies/{dummy_id}",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.deleteById,
	)
}
