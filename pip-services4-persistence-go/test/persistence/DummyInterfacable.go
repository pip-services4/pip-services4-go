package test_persistence

type DummyInterfacable struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

func (d DummyInterfacable) GetId() string {
	return d.Id
}

func (d *DummyInterfacable) SetId(id string) {
	d.Id = id
}

func (d DummyInterfacable) Clone() DummyInterfacable {
	return DummyInterfacable{
		Id:      d.Id,
		Key:     d.Key,
		Content: d.Content,
	}
}
