package test_controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	"github.com/pip-services4/pip-services4-go/pip-services4-http-go/services"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

type DummyRestController struct {
	*services.RestController
	service        tlogic.IDummyService
	numberOfCalls  int
	openApiContent string
	openApiFile    string
}

func NewDummyRestController() *DummyRestController {
	c := &DummyRestController{}
	c.RestController = services.InheritRestController(c)
	c.numberOfCalls = 0
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}

func (c *DummyRestController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	if _val, ok := config.GetAsNullableString("openapi_content"); ok {
		c.openApiContent = _val
	}
	if _val, ok := config.GetAsNullableString("openapi_file"); ok {
		c.openApiFile = _val
	}
	c.RestController.Configure(ctx, config)
}

func (c *DummyRestController) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.RestController.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr == nil && depRes != nil {
		c.service = depRes.(tlogic.IDummyService)
	}
}

func (c *DummyRestController) GetNumberOfCalls() int {
	return c.numberOfCalls
}

func (c *DummyRestController) incrementNumberOfCalls(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	c.numberOfCalls++
	next.ServeHTTP(res, req)
}

func (c *DummyRestController) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)
	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	result, err := c.service.GetPageByFilter(
		ctx,
		cdata.NewFilterParamsFromValue(params), // W! need test
		cdata.NewPagingParamsFromTuples(paginParams),
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) getOneById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}
	result, err := c.service.GetOneById(
		ctx,
		dummyId)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) create(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)
	var dummy tdata.Dummy

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
		utils.ContextHelper.NewContextWithTraceId(req.Context(), traceId),
		dummy,
	)
	c.SendCreatedResult(res, req, result, err)
}

func (c *DummyRestController) update(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	var dummy tdata.Dummy

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
		utils.ContextHelper.NewContextWithTraceId(req.Context(), traceId),
		dummy,
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) deleteById(res http.ResponseWriter, req *http.Request) {
	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	params := req.URL.Query()
	vars := mux.Vars(req)

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}

	result, err := c.service.DeleteById(
		ctx,
		dummyId,
	)
	c.SendDeletedResult(res, req, result, err)
}

func (c *DummyRestController) checkTraceId(res http.ResponseWriter, req *http.Request) {
	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.CheckTraceId(ctx)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestController) checkErrorPropagation(res http.ResponseWriter, req *http.Request) {
	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	err := c.service.CheckErrorPropagation(ctx)
	c.SendError(res, req, err)
}

func (c *DummyRestController) checkGracefulShutdownContext(res http.ResponseWriter, req *http.Request) {
	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	err := c.service.CheckGracefulShutdownContext(ctx)
	c.SendError(res, req, err)
}

func (c *DummyRestController) Register() {
	c.RegisterInterceptor("/dummies$", c.incrementNumberOfCalls)

	c.RegisterRoute(
		http.MethodGet, "/dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("skip", cconv.String).
			WithOptionalProperty("take", cconv.String).
			WithOptionalProperty("total", cconv.String).
			WithOptionalProperty("body", cvalid.NewFilterParamsSchema()).Schema,
		c.getPageByFilter,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/trace_id",
		cvalid.NewObjectSchema().Schema,
		c.checkTraceId,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/error_propagation",
		cvalid.NewObjectSchema().Schema,
		c.checkErrorPropagation,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/graceful_shutdown",
		cvalid.NewObjectSchema().Schema,
		c.checkGracefulShutdownContext,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/{dummy_id}",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.getOneById,
	)

	c.RegisterRoute(
		http.MethodPost, "/dummies",
		cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.create,
	)

	c.RegisterRoute(
		http.MethodPut, "/dummies",
		cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.update,
	)

	c.RegisterRoute(
		http.MethodDelete, "/dummies/{dummy_id}",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.deleteById,
	)

	if c.openApiContent != "" {
		c.RegisterOpenApiSpec(c.openApiContent)
	}

	if c.openApiFile != "" {
		c.RegisterOpenApiSpecFromFile(c.openApiFile)
	}
}
