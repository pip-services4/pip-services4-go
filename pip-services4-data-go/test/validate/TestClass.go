package test_validate

type SubTestClass struct {
	Id         string
	FloatField float32
}

type TestClass struct {
	IntField         int
	StringField1     string
	StringField2     string
	IntArrayField    []int
	StringArrayField []string
	MapField         map[string]any
	SubObjectField   *SubTestClass
	SubArrayField    []*SubTestClass
}
