package services_test

import (
	"context"
	"encoding/json"
	"net/http"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	gcpctrl "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/controllers"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/logic"
	gcputil "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type DummyCloudFunctionController struct {
	*gcpctrl.CloudFunctionController

	service tlogic.IDummyService
}

func NewDummyCloudFunctionService() *DummyCloudFunctionController {
	c := DummyCloudFunctionController{}

	c.CloudFunctionController = gcpctrl.InheritCloudFunctionController(&c, "dummies")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))

	return &c
}

func (c *DummyCloudFunctionController) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.CloudFunctionController.SetReferences(ctx, references)

	depRes, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr == nil && depRes != nil {
		c.service = depRes.(tlogic.IDummyService)
	} else {
		panic("Can't find Service reference")
	}
}

func (c *DummyCloudFunctionController) getPageByFilter(res http.ResponseWriter, req *http.Request) {
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
		cquery.NewFilterParamsFromValue(params),
		cquery.NewPagingParamsFromTuples(paginParams),
	)

	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyCloudFunctionController) getOneById(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)
	var body map[string]string

	err := gcputil.CloudFunctionRequestHelper.DecodeBody(req, &body)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	defer req.Body.Close()

	result, err := c.service.GetOneById(
		cctx.NewContextWithTraceId(req.Context(), traceId),
		body["dummy_id"])

	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyCloudFunctionController) create(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(traceId, req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.service.Create(
		cctx.NewContextWithTraceId(req.Context(), traceId),
		dummy,
	)

	httpctrl.HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *DummyCloudFunctionController) update(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(traceId, req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.service.Update(
		cctx.NewContextWithTraceId(req.Context(), traceId),
		dummy,
	)
	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyCloudFunctionController) deleteById(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	var body map[string]string

	err := gcputil.CloudFunctionRequestHelper.DecodeBody(req, &body)
	defer req.Body.Close()

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	dummyId := body["dummy_id"]

	result, err := c.service.DeleteById(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummyId,
	)
	httpctrl.HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *DummyCloudFunctionController) getDummy(traceId string, req *http.Request) (tdata.Dummy, error) {
	var body map[string]any
	var dummy tdata.Dummy

	err := gcputil.CloudFunctionRequestHelper.DecodeBody(req, &body)
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

func (c *DummyCloudFunctionController) Register() {

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
