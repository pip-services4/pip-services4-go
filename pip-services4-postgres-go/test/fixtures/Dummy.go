package fixtures

type Dummy struct {
	Id      string `bson:"_id" json:"id"`
	Key     string `bson:"key" json:"key"`
	Content string `bson:"content" json:"content"`
}

func (d *Dummy) SetId(id string) {
	d.Id = id
}

func (d Dummy) GetId() string {
	return d.Id
}

func (d Dummy) Clone() Dummy {
	return Dummy{
		Id:      d.Id,
		Key:     d.Key,
		Content: d.Content,
	}
}
