package data

import (
	"encoding/hex"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
)

// IdGenerator Helper class to generate unique object IDs. It supports two types of IDs: long and short.
//	Long IDs are string GUIDs. They are globally unique and 32-character long.
//	ShortIDs are just 9-digit random numbers. They are not guaranteed be unique.
//	Example:
//		IdGenerator.NextLong();      // Possible result: "234ab342c56a2b49c2ab42bf23ff991ac"
//		IdGenerator.NextShort();     // Possible result: "23495247"
var IdGenerator = &_TIdGenerator{}

type _TIdGenerator struct{}

// NextShort generates a random 9-digit random ID (code).
//	Remember: The returned value is not guaranteed to be unique.
//	Returns: string a generated random 9-digit code
func (c *_TIdGenerator) NextShort() string {
	value := 100000000 + rand.Int63n(899999999)
	return strconv.FormatInt(value, 10)
}

// NextLong generates a globally unique 32-digit object ID. The value is a string representation of a GUID value.
//	Returns: string a generated 32-digit object ID
func (c *_TIdGenerator) NextLong() string {
	value := uuid.New()
	return hex.EncodeToString(value[:])
}
