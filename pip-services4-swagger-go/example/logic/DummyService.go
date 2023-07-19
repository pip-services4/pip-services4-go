package example_logic

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	data "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/data"
)

type DummyService struct {
	commandSet *DummyCommandSet
	entities   []data.Dummy
}

func NewDummyService() *DummyService {
	dc := DummyService{}
	dc.entities = make([]data.Dummy, 0)
	return &dc
}

func (c *DummyService) GetCommandSet() *ccomand.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDummyCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *DummyService) GetPageByFilter(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (items *cquery.DataPage[data.Dummy], err error) {

	if filter == nil {
		filter = cquery.NewEmptyFilterParams()
	}
	var key string = filter.GetAsString("key")

	if paging == nil {
		paging = cquery.NewEmptyPagingParams()
	}
	var skip int64 = paging.GetSkip(0)
	var take int64 = paging.GetTake(100)

	var result []data.Dummy
	for i := 0; i < len(c.entities); i++ {
		var entity data.Dummy = c.entities[i]
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
	return cquery.NewDataPage(result, len(result)), nil
}

func (c *DummyService) GetOneById(ctx context.Context, id string) (result *data.Dummy, err error) {
	for i := 0; i < len(c.entities); i++ {
		var entity data.Dummy = c.entities[i]
		if id == entity.Id {
			return &entity, nil
		}
	}
	return nil, nil
}

func (c *DummyService) Create(ctx context.Context, entity data.Dummy) (result *data.Dummy, err error) {
	if entity.Id == "" {
		entity.Id = keys.IdGenerator.NextLong()
	}
	c.entities = append(c.entities, entity)
	return &entity, nil
}

func (c *DummyService) Update(ctx context.Context, newEntity data.Dummy) (result *data.Dummy, err error) {
	for index := 0; index < len(c.entities); index++ {
		var entity data.Dummy = c.entities[index]
		if entity.Id == newEntity.Id {
			c.entities[index] = newEntity
			return &newEntity, nil

		}
	}
	return nil, nil
}

func (c *DummyService) DeleteById(ctx context.Context, id string) (result *data.Dummy, err error) {
	var entity data.Dummy

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
		return &entity, nil
	}
	return nil, nil
}
