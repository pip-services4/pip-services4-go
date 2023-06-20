package auth

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// CredentialParams contains credentials to authenticate against external services.
// They are used together with connection parameters, but usually stored
// in a separate store, protected from unauthorized access.
//
//	Configuration parameters
//
//		- store_key: key to retrieve parameters from credential store
//		- username: user name
//		- user: alternative to username
//		- password: user password
//		- pass: alternative to password
//		- access_id: application access id
//		- client_id: alternative to access_id
//		- access_key: application secret key
//		- client_key: alternative to access_key
//		- secret_key: alternative to access_key
//
// In addition to standard parameters CredentialParams may contain any number of custom parameters
// see config.ConfigParams
// see ConnectionParams
// see CredentialResolver
// see ICredentialStore
//
//	Example:
//		credential := NewCredentialParamsFromTuples(
//			"user", "jdoe",
//			"pass", "pass123",
//			"pin", "321"
//		);
//		username := credential.Username();  // Result: "jdoe"
//		password := credential.Password();  // Result: "pass123"
type CredentialParams struct {
	*config.ConfigParams
}

const (
	CredentialsParamSectionKey string = "credentials"
	CredentialParamSectionKey  string = "credential"
	CredentialParamStoreKey    string = "store_key"
	CredentialParamUsername    string = "username"
	CredentialParamUser        string = "user"
	CredentialParamPassword    string = "password"
	CredentialParamPass        string = "pass"
	CredentialParamAccessId    string = "access_id"
	CredentialParamClientId    string = "client_id"
	CredentialParamAccessKey   string = "access_key"
	CredentialParamClientKey   string = "client_key"
	CredentialParamSecretKey   string = "secret_key"
)

// NewEmptyCredentialParams creates a new credential parameters and fills it with values.
//
//	Returns: *CredentialParams
func NewEmptyCredentialParams() *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewEmptyConfigParams(),
	}
}

// NewCredentialParams creates a new credential parameters and fills it with values.
//
//	Parameters:
//		- values map[string]string an object to be converted
//		into key-value pairs to initialize these credentials.
//	Returns: *CredentialParams
func NewCredentialParams(values map[string]string) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParams(values),
	}
}

// NewCredentialParamsFromValue method that creates a ConfigParams object
// based on the values that are stored in the 'value' object's properties.
//
//	Parameters:
//		- value any configuration parameters in the form of an object with properties.
//	Returns: *ConfigParams generated ConfigParams.
func NewCredentialParamsFromValue(value any) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParamsFromValue(value),
	}
}

// NewCredentialParamsFromTuples creates a new CredentialParams object filled with
// provided key-value pairs called tuples.
// Tuples parameters contain a sequence of key1, value1, key2, value2, ... pairs.
//
//	Parameters:
//		- tuples ...any the tuples to fill a new CredentialParams object.
//	Returns: *CredentialParams a new CredentialParams object.
func NewCredentialParamsFromTuples(tuples ...any) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParamsFromTuplesArray(tuples),
	}
}

// NewCredentialParamsFromTuplesArray static method for creating a CredentialParams from an array of tuples.
//
//	Parameters:
//		- tuples []any the key-value tuples array to initialize the new StringValueMap with.
//	Returns: CredentialParams the CredentialParams created and filled by the 'tuples' array provided.
func NewCredentialParamsFromTuplesArray(tuples []any) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParamsFromTuplesArray(tuples),
	}
}

// NewCredentialParamsFromString creates a new CredentialParams object filled with key-value pairs serialized as a string.
//
//	Parameters:
//	- line string a string with serialized key-value
//		pairs as "key1=value1;key2=value2;..."
//		Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//	Returns: *CredentialParams a new CredentialParams object.
func NewCredentialParamsFromString(line string) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParamsFromString(line),
	}
}

// NewCredentialParamsFromMaps static method for creating a CredentialParams using the maps passed as parameters.
//
//	Parameters:
//		- maps ...map[string]string the maps passed to this method to create a StringValueMap with.
//	Returns: *CredentialParams the CredentialParams created.
func NewCredentialParamsFromMaps(maps ...map[string]string) *CredentialParams {
	return &CredentialParams{
		ConfigParams: config.NewConfigParamsFromMaps(maps...),
	}
}

// NewManyCredentialParamsFromConfig retrieves all CredentialParams from configuration parameters
// from "credentials" section. If "credential" section is present instead,
// then it returns a list with only one CredentialParams.
//
//	Parameters:
//		- config *config.ConfigParams a configuration parameters to retrieve credentials
//	Returns: []*CredentialParams a list of retrieved CredentialParams
func NewManyCredentialParamsFromConfig(config *config.ConfigParams) []*CredentialParams {
	result := make([]*CredentialParams, 0)

	credentials := config.GetSection(CredentialsParamSectionKey)

	if credentials.Len() > 0 {
		for _, section := range credentials.GetSectionNames() {
			credential := credentials.GetSection(section)
			result = append(result, NewCredentialParams(credential.Value()))
		}
	} else {
		credential := config.GetSection(CredentialParamSectionKey)
		if credential.Len() > 0 {
			result = append(result, NewCredentialParams(credential.Value()))
		}
	}

	return result
}

// NewCredentialParamsFromConfig кetrieves a single CredentialParams from
// configuration parameters from "credential" section. If "credentials"
// section is present instead, then is returns only the first credential element.
//
//	Parameters:
//		- config *config.ConfigParams, containing a section named "credential(s)".
//	Returns []*CredentialParams the generated CredentialParams object.
func NewCredentialParamsFromConfig(config *config.ConfigParams) *CredentialParams {
	credentials := NewManyCredentialParamsFromConfig(config)
	if len(credentials) > 0 {
		return credentials[0]
	}
	return nil
}

// UseCredentialStore сhecks if these credential parameters shall be
// retrieved from CredentialStore. The credential parameters are
// redirected to CredentialStore when store_key parameter is set.
//
//	Returns: bool true if credentials shall be retrieved from CredentialStore
func (c *CredentialParams) UseCredentialStore() bool {
	if _, ok := c.GetAsNullableString(CredentialParamStoreKey); ok {
		return true
	}
	return false
}

// StoreKey gets the key to retrieve these credentials from CredentialStore.
// If this key is null, then all parameters are already present.
//
//	Returns: string the store key to retrieve credentials.
func (c *CredentialParams) StoreKey() string {
	return c.GetAsString(CredentialParamStoreKey)
}

// SetStoreKey sets the key to retrieve these parameters from CredentialStore.
//
//	Parameters: value string a new key to retrieve credentials.
func (c *CredentialParams) SetStoreKey(value string) {
	c.Put(CredentialParamStoreKey, value)
}

// Username gets the user name. The value can be stored in parameters "username" or "user".
//
//	Returns: string the user name.
func (c *CredentialParams) Username() string {
	if username, ok := c.GetAsNullableString(CredentialParamUsername); ok {
		return username
	}
	return c.GetAsString(CredentialParamUser)
}

// SetUsername sets the user name.
//
//	Parameters: value string a new user name.
func (c *CredentialParams) SetUsername(value string) {
	c.Put(CredentialParamUsername, value)
}

// Password get the user password. The value can be stored in parameters "password" or "pass".
//
//	Returns: string the user password.
func (c *CredentialParams) Password() string {
	if password, ok := c.GetAsNullableString(CredentialParamPassword); ok {
		return password
	}
	return c.GetAsString(CredentialParamPass)
}

// SetPassword sets the user password.
//
//	Parameters: value string a new user password.
func (c *CredentialParams) SetPassword(value string) {
	c.Put(CredentialParamPassword, value)
}

// AccessId gets the application access id. The value can be stored in parameters "access_id" pr "client_id"
//
//	Returns: string the application access id.
func (c *CredentialParams) AccessId() string {
	if accessId, ok := c.GetAsNullableString(CredentialParamAccessId); ok {
		return accessId
	}
	return c.GetAsString(CredentialParamClientId)
}

// SetAccessId sets the application access id.
//
//	Parameters: value: string a new application access id.
func (c *CredentialParams) SetAccessId(value string) {
	c.Put(CredentialParamAccessId, value)
}

// AccessKey the application secret key.
// The value can be stored in parameters "access_key", "client_key" or "secret_key".
//
//	ReturnsЖ string the application secret key.
func (c *CredentialParams) AccessKey() string {
	if accessKey, ok := c.GetAsNullableString(CredentialParamAccessKey); ok {
		return accessKey
	}
	return c.GetAsString(CredentialParamClientKey)
}

// SetAccessKey sets the application secret key.
//
//	Parameters: value string a new application secret key.
func (c *CredentialParams) SetAccessKey(value string) {
	c.Put(CredentialParamAccessKey, value)
}
