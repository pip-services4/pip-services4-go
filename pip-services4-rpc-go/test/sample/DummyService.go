package test_sample

import (
	"context"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type DummyService struct {
	commandSet *DummyCommandSet
	entities   []Dummy
}

func NewDummyService() *DummyService {
	dc := DummyService{}
	dc.entities = make([]Dummy, 0)
	return &dc
}

func (c *DummyService) GetCommandSet() *ccomand.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDummyCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *DummyService) GetPageByFilter(ctx context.Context, filter *cquery.FilterParams,
	paging *cquery.PagingParams) (items *cquery.DataPage[Dummy], err error) {

	if filter == nil {
		filter = cquery.NewEmptyFilterParams()
	}
	var key string = filter.GetAsString("key")

	if paging == nil {
		paging = cquery.NewEmptyPagingParams()
	}
	var skip int64 = paging.GetSkip(0)
	var take int64 = paging.GetTake(100)

	var result []Dummy
	for i := 0; i < len(c.entities); i++ {
		var entity Dummy = c.entities[i]
		if key != "" && key != entity.Key {
			continue
		}

		skip--
		if skip >= 0 {
			continue
		}

		take--
		if take < 0 {
			break
		}

		result = append(result, entity)
	}
	var total int64 = (int64)(len(result))
	return cquery.NewDataPage[Dummy](result, int(total)), nil
}

func (c *DummyService) GetOneById(ctx context.Context, id string) (result Dummy, err error) {
	for i := 0; i < len(c.entities); i++ {
		var entity Dummy = c.entities[i]
		if id == entity.Id {
			return entity, nil
		}
	}
	return Dummy{}, nil
}

func (c *DummyService) Create(ctx context.Context, entity Dummy) (result Dummy, err error) {
	if entity.Id == "" {
		entity.Id = keys.IdGenerator.NextLong()
		c.entities = append(c.entities, entity)
	}
	return entity, nil
}

func (c *DummyService) Update(ctx context.Context, newEntity Dummy) (result Dummy, err error) {
	for index := 0; index < len(c.entities); index++ {
		var entity Dummy = c.entities[index]
		if entity.Id == newEntity.Id {
			c.entities[index] = newEntity
			return newEntity, nil

		}
	}
	return Dummy{}, nil
}

func (c *DummyService) DeleteById(ctx context.Context, id string) (result Dummy, err error) {
	var entity Dummy

	for i := 0; i < len(c.entities); {
		entity = c.entities[i]
		if entity.Id == id {
			if i == len(c.entities)-1 {
				c.entities = c.entities[:i]
			} else {
				c.entities = append(c.entities[:i], c.entities[i+1:]...)
			}
		} else {
			i++
		}
		return entity, nil
	}
	return Dummy{}, nil
}

func (c *DummyService) CheckTraceId(ctx context.Context) (result map[string]string, err error) {
	result = map[string]string{"traceId": cctx.GetTraceId(ctx)}
	return result, nil
}

func (c *DummyService) CheckErrorPropagation(ctx context.Context) error {
	return cerr.NewNotFoundError(cctx.GetTraceId(ctx), "NOT_FOUND_TEST", "Not found error")
}

func (c *DummyService) CheckGracefulShutdownContext(ctx context.Context) error {
	defer cctx.DefaultErrorHandlerWithShutdown(ctx)

	panic("called from DummyController.CheckGracefulShutdownContext")
}
