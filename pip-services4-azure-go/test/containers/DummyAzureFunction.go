package containers_test

import (
	"context"
	"encoding/json"
	"net/http"

	azurecont "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/build"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/data"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/logic"
	azureutil "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/utils"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type DummyAzureFunction struct {
	service tlogic.IDummyService
	*azurecont.AzureFunction
}

func NewDummyAzureFunction() *DummyAzureFunction {
	c := DummyAzureFunction{}
	c.AzureFunction = azurecont.InheritAzureFunctionWithParams(&c, "dummy", "Dummy Azure function")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	c.AddFactory(tbuild.NewDummyFactory())
	return &c
}

func (c *DummyAzureFunction) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.AzureFunction.SetReferences(ctx, references)
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	c.service = resCtrl.(tlogic.IDummyService)
}

func (c *DummyAzureFunction) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.GetPageByFilter(
		ctx,
		cdata.NewFilterParamsFromValue(params),
		cdata.NewPagingParamsFromTuples(paginParams),
	)

	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) getOneById(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)
	var body map[string]string

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	defer req.Body.Close()

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.GetOneById(ctx, body["dummy_id"])

	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) create(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.Create(ctx, dummy)

	httpctrl.HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *DummyAzureFunction) update(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.Update(ctx, dummy)
	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) deleteById(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	var body map[string]string

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)
	defer req.Body.Close()

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	dummyId := body["dummy_id"]

	ctx := utils.ContextHelper.NewContextWithTraceId(req.Context(), c.GetTraceId(req))
	result, err := c.service.DeleteById(ctx, dummyId)
	httpctrl.HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *DummyAzureFunction) getDummy(req *http.Request) (tdata.Dummy, error) {
	var body map[string]any
	var dummy tdata.Dummy

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)
	defer req.Body.Close()

	if err != nil {
		return tdata.Dummy{}, err
	}

	dummyBytes, err := json.Marshal(body["dummy"])

	if err != nil {
		return tdata.Dummy{}, err
	}

	err = json.Unmarshal(dummyBytes, &dummy)

	if err != nil {
		return tdata.Dummy{}, err
	}

	return dummy, nil
}

func (c *DummyAzureFunction) Register() {

	c.RegisterAction(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty(
			"body", cvalid.NewObjectSchema().WithOptionalProperty(
				"filter", cvalid.NewFilterParamsSchema())).WithOptionalProperty(
			"paging", cvalid.NewPagingParamsSchema()).Schema,
		c.getPageByFilter,
	)

	c.RegisterAction(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String)).Schema,
		c.getOneById,
	)

	c.RegisterAction(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema())).Schema,
		c.create,
	)

	c.RegisterAction(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema())).Schema,
		c.update,
	)

	c.RegisterAction(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String)).Schema,
		c.deleteById,
	)
}
