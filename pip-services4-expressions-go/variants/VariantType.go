package variants

// VariantType represents types of variant
type VariantType int

// Defines supported types of variant values.
const (
	Null     VariantType = iota
	Integer  VariantType = iota
	Long     VariantType = iota
	Float    VariantType = iota
	Double   VariantType = iota
	String   VariantType = iota
	Boolean  VariantType = iota
	DateTime VariantType = iota
	TimeSpan VariantType = iota
	Object   VariantType = iota
	Array    VariantType = iota
)
