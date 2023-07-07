package containers_test

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
	gcpcont "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/build"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/logic"
	gcputil "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type DummyCloudFunction struct {
	service tlogic.IDummyService
	*gcpcont.CloudFunction
}

func NewDummyCloudFunction() *DummyCloudFunction {
	c := DummyCloudFunction{}
	c.CloudFunction = gcpcont.InheritCloudFunctionWithParams(&c, "dummy", "Dummy GCP function")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	c.AddFactory(tbuild.NewDummyFactory())
	return &c
}

func (c *DummyCloudFunction) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.CloudFunction.SetReferences(ctx, references)
	resSrv, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	c.service = resSrv.(tlogic.IDummyService)
}

func (c *DummyCloudFunction) getPageByFilter(res http.ResponseWriter, req *http.Request) {
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

func (c *DummyCloudFunction) getOneById(res http.ResponseWriter, req *http.Request) {
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
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		body["dummy_id"])

	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyCloudFunction) create(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(traceId, req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.service.Create(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummy,
	)

	httpctrl.HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *DummyCloudFunction) update(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)

	dummy, err := c.getDummy(traceId, req)

	if err != nil {
		err := cerr.NewInternalError(traceId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.service.Update(
		cctx.NewContextWithTraceId(req.Context(), c.GetTraceId(req)),
		dummy,
	)
	httpctrl.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyCloudFunction) deleteById(res http.ResponseWriter, req *http.Request) {
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

func (c *DummyCloudFunction) getDummy(traceId string, req *http.Request) (tdata.Dummy, error) {
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

func (c *DummyCloudFunction) Register() {

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
