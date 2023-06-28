package connect

import (
	"context"
	"strings"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Contains connection parameters to authenticate against Amazon Web Services (AWS)
// and connect to specific AWS resource.
//
// The struct is able to compose and parse AWS resource ARNs.
//
// Configuration parameters
//
//   - access_id:     application access id
//   - client_id:     alternative to access_id
//   - access_key:    application secret key
//   - client_key:    alternative to access_key
//   - secret_key:    alternative to access_key
//
// # In addition to standard parameters CredentialParams may contain any number of custom parameters
//
// # See AwsConnectionResolver
//
// Example
//
//	connection := NewAwsConnectionParamsFromTuples(
//	    "region", "us-east-1",
//	    "access_id", "XXXXXXXXXXXXXXX",
//	    "secret_key", "XXXXXXXXXXXXXXX",
//	    "service", "s3",
//	    "bucket", "mybucket"
//	)
//
//	region := connection.Region()                     // Result: "us-east-1"
//	accessId := connection.AccessId()                 // Result: "XXXXXXXXXXXXXXX"
//	secretKey := connection.AccessKey()               // Result: "XXXXXXXXXXXXXXX"
//	pin := connection.GetAsNullableString("bucket")      // Result: "mybucket"
type AwsConnectionParams struct {
	*cconf.ConfigParams
}

// NewAwsConnectionParams creates an new instance of the connection parameters.
//   - values 	(optional) an object to be converted into key-value pairs to initialize this connection.
func NewAwsConnectionParams(values map[string]string) *AwsConnectionParams {
	c := AwsConnectionParams{}
	c.ConfigParams = cconf.NewConfigParams(values)
	return &c
}

// NewEmptyAwsConnectionParams creates an new instance of the connection parameters.
func NewEmptyAwsConnectionParams() *AwsConnectionParams {
	c := AwsConnectionParams{}
	c.ConfigParams = cconf.NewEmptyConfigParams()
	return &c
}

// GetPartition Gets the AWS partition name.
// Returns the AWS partition name.
func (c *AwsConnectionParams) GetPartition() string {
	res, ok := c.GetAsNullableString("partition")
	if ok {
		return res
	}
	return "aws"
}

// SetPartition Sets the AWS partition name.
//   - value a new AWS partition name.
func (c *AwsConnectionParams) SetPartition(value string) {
	c.Put("partition", value)
}

// GetService gets the AWS service name.
// Returns the AWS service name.
func (c *AwsConnectionParams) GetService() string {
	res, ok := c.GetAsNullableString("service")
	if ok {
		return res
	}
	res, ok = c.GetAsNullableString("protocol")
	if ok {
		return res
	}
	return ""
}

// SetService sets the AWS service name.
//   - value a new AWS service name.
func (c *AwsConnectionParams) SetService(value string) {
	c.Put("service", value)
}

// GetRegion gets the AWS region.
// Returns the AWS region.
func (c *AwsConnectionParams) GetRegion() string {
	res, ok := c.GetAsNullableString("region")
	if ok {
		return res
	}
	return ""
}

// SetRegion Sets the AWS region.
//   - value a new AWS region.
func (c *AwsConnectionParams) SetRegion(value string) {
	c.Put("region", value)
}

// GetAccount Gets the AWS account id.
// Returns the AWS account id.
func (c *AwsConnectionParams) GetAccount() string {
	res, ok := c.GetAsNullableString("account")
	if ok {
		return res
	}
	return ""
}

// SetAccount Sets the AWS account id.
// - value the AWS account id.
func (c *AwsConnectionParams) SetAccount(value string) {
	c.Put("account", value)
}

// GetResourceType gets the AWS resource type.
// Returns the AWS resource type.
func (c *AwsConnectionParams) GetResourceType() string {
	res, ok := c.GetAsNullableString("resource_type")
	if ok {
		return res
	}
	return ""
}

// SetResourceType sets the AWS resource type.
//   - value a new AWS resource type.
func (c *AwsConnectionParams) SetResourceType(value string) {
	c.Put("resource_type", value)
}

// GetResource gets the AWS resource id.
//
//	Returns the AWS resource id.
func (c *AwsConnectionParams) GetResource() string {
	res, ok := c.GetAsNullableString("resource")
	if ok {
		return res
	}
	return ""
}

// SetResource sets the AWS resource id.
//   - value a new AWS resource id.
func (c *AwsConnectionParams) SetResource(value string) {
	c.Put("resource", value)
}

// GetArn gets the AWS resource ARN.
// If the ARN is not defined it automatically generates it from other properties.
// Returns the AWS resource ARN.
func (c *AwsConnectionParams) GetArn() string {
	res, ok := c.GetAsNullableString("arn")
	if ok {
		return res
	}

	arn := "arn"
	partition := c.GetPartition()
	arn += ":" + partition
	service := c.GetService()
	arn += ":" + service
	region := c.GetRegion()
	arn += ":" + region
	account := c.GetAccount()
	arn += ":" + account
	resourceType := c.GetResourceType()
	if resourceType != "" {
		arn += ":" + resourceType
	}
	resource := c.GetResource()
	arn += ":" + resource

	return arn
}

// SetArn sets the AWS resource ARN.
// When it sets the value, it automatically parses the ARN
// and sets individual parameters.
//   - value a new AWS resource ARN.
func (c *AwsConnectionParams) SetArn(value string) {

	c.Put("arn", value)

	if value != "" {
		tokens := strings.Split(value, ":")
		c.SetPartition(tokens[1])
		c.SetService(tokens[2])
		c.SetRegion(tokens[3])
		c.SetAccount(tokens[4])
		if len(tokens) > 6 {
			c.SetResourceType(tokens[5])
			c.SetResource(tokens[6])
		} else {
			temp := tokens[5]
			pos := strings.Index(temp, "/")
			if pos > 0 {
				c.SetResourceType(temp[:pos])
				c.SetResource(temp[pos+1:])
			} else {
				c.SetResourceType("")
				c.SetResource(temp)
			}
		}
	}
}

// GetAccessId gets the AWS access id.
// Returns the AWS access id.
func (c *AwsConnectionParams) GetAccessId() string {
	res, ok := c.GetAsNullableString("access_id")
	if ok {
		return res
	}
	res, ok = c.GetAsNullableString("client_id")
	if ok {
		return res
	}
	return ""
}

// SetAccessId sets the AWS access id.
//   - value the AWS access id.
func (c *AwsConnectionParams) SetAccessId(value string) {
	c.Put("access_id", value)
}

// GetAccessKey gets the AWS client key.
// Returns the AWS client key.
func (c *AwsConnectionParams) GetAccessKey() string {
	res, ok := c.GetAsNullableString("access_key")
	if ok {
		return res
	}
	res, ok = c.GetAsNullableString("client_key")
	if ok {
		return res
	}
	return ""
}

// SetAccessKey sets the AWS client key.
//   - value a new AWS client key.
func (c *AwsConnectionParams) SetAccessKey(value string) {
	c.Put("access_key", value)
}

// NewAwsConnectionParamsFromString creates a new AwsConnectionParams object filled with key-value pairs serialized as a string.
//   - line 	a string with serialized key-value pairs as "key1=value1;key2=value2;..."
//     Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//
// Returns			a new AwsConnectionParams object.
func NewAwsConnectionParamsFromString(line string) *AwsConnectionParams {
	strinMap := cdata.NewStringValueMapFromString(line)
	return NewAwsConnectionParams(strinMap.Value())
}

// Validates this connection parameters
//   - ctx context.Context    execution context to trace execution through call chain.
//
// Returns a ConfigException or null if validation passed successfully.
func (c *AwsConnectionParams) Validate(ctx context.Context) *cerr.ApplicationError { //ConfigException\
	traceId := cctx.GetTraceId(ctx)
	arn := c.GetArn()
	if arn == "arn:aws::::" {
		return cerr.NewConfigError(
			traceId,
			"NO_AWS_CONNECTION",
			"AWS connection is not set")
	}

	if c.GetAccessId() == "" {
		return cerr.NewConfigError(
			traceId,
			"NO_ACCESS_ID",
			"No access_id is configured in AWS credential")
	}

	if c.GetAccessKey() == "" {
		return cerr.NewConfigError(
			traceId,
			"NO_ACCESS_KEY",
			"No access_key is configured in AWS credential")
	}
	return nil
}

// Retrieves AwsConnectionParams from configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
//   - config   configuration parameters
//
// Returns			the generated AwsConnectionParams object.
// See NewAwsConnectionParamsMergeConfigs
func NewAwsConnectionParamsFromConfig(config *cconf.ConfigParams) *AwsConnectionParams {
	result := NewEmptyAwsConnectionParams()

	credentials := cauth.NewManyCredentialParamsFromConfig(config)
	for _, credential := range credentials {
		result.Append(credential.Value())
	}

	connections := cconn.NewManyConnectionParamsFromConfig(config)
	for _, connection := range connections {
		result.Append(connection.Value())
	}

	return result
}

// NewAwsConnectionParamsMergeConfigs retrieves AwsConnectionParams from multiple configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
//   - configs   a list with configuration parameters
//
// Returns the generated AwsConnectionParams object.
// See NewAwsConnectionParamsFromConfig
func NewAwsConnectionParamsMergeConfigs(configs []*cconf.ConfigParams) *AwsConnectionParams {

	var maps []map[string]string
	maps = make([]map[string]string, 0)

	for _, conf := range configs {
		maps = append(maps, conf.Value())
	}

	config := cconf.NewConfigParamsFromValue(maps)

	return NewAwsConnectionParams(config.Value())
}
