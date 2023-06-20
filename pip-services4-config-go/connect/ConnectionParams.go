package connect

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// ConnectionParams contains connection parameters to connect to external services.
// They are used together with credential parameters, but usually stored separately
// from more protected sensitive values.
//
//	Configuration parameters:
//		- discovery_key: key to retrieve parameters from discovery service
//		- protocol: connection protocol like http, https, tcp, udp
//		- host: host name or IP address
//		- port: port number
//		- uri: resource URI or connection string with all parameters in it
//
// In addition to standard parameters ConnectionParams may contain any number of custom parameters
//
//	see ConfigParams
//	see CredentialParams
//	see ConnectionResolver
//	see IDiscovery
//	Example ConnectionParams object usage:
//		connection := NewConnectionParamsFromTuples(
//			ConnectionParamProtocol, "http",
//			ConnectionParamHost, "10.1.1.100",
//			ConnectionParamPort, "8080",
//			ConnectionParamCluster, "mycluster"
//		);
//		host := connection.Host();                             				// Result: "10.1.1.100"
//		port := connection.Port();                             				// Result: 8080
//		cluster := connection.GetAsNullableString(ConnectionParamCluster);  // Result: "mycluster"
type ConnectionParams struct {
	*config.ConfigParams
}

const (
	SectionNameConnections      = "connections"
	SectionNameConnection       = "connection"
	ConnectionParamDiscoveryKey = "discovery_key"
	ConnectionParamProtocol     = "protocol"
	ConnectionParamHost         = "host"
	ConnectionParamIp           = "ip"
	ConnectionParamPort         = "port"
	ConnectionParamURI          = "uri"
	ConnectionParamCluster      = "cluster"
)

// NewEmptyConnectionParams creates a new connection parameters and fills it with values.
//
//	Returns: *ConnectionParams
func NewEmptyConnectionParams() *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewEmptyConfigParams(),
	}
}

// NewConnectionParams creates a new connection parameters and fills it with values.
//
//	Parameters:
//		- values map[string]string an object to be
//		converted into key-value pairs to initialize this connection.
//	Returns: *ConnectionParams
func NewConnectionParams(values map[string]string) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParams(values),
	}
}

// NewConnectionParamsFromValue method that creates a ConfigParams object based on the values that
// are stored in the 'value' object's properties.
//
//	see RecursiveObjectReader.GetProperties
//	Parameters:
//		- value any configuration parameters in the form of an object with properties.
//	Returns: ConnectionParams generated ConnectionParams.
func NewConnectionParamsFromValue(value any) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParamsFromValue(value),
	}
}

// NewConnectionParamsFromTuples creates a new ConnectionParams object filled with provided key-value
// pairs called tuples. Tuples parameters contain a sequence of
// key1, value1, key2, value2, ... pairs.
//
//	Parameters:
//		- tuples ...any the tuples to fill a new ConnectionParams object.
//	Returns: *ConnectionParams a new ConnectionParams object.
func NewConnectionParamsFromTuples(tuples ...any) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParamsFromTuplesArray(tuples),
	}
}

// NewConnectionParamsFromTuplesArray method for creating a StringValueMap from an array of tuples.
//
//	Parameters:
//		- tuples []any the key-value tuples array to initialize the new StringValueMap with.
//	Returns: *ConnectionParams the ConnectionParams created and filled by the 'tuples' array provided.
func NewConnectionParamsFromTuplesArray(tuples []any) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParamsFromTuplesArray(tuples),
	}
}

// NewConnectionParamsFromString creates a new ConnectionParams object filled with
// key-value pairs serialized as a string.
//
//	Parameters:
//		- line string a string with serialized key-value pairs as
//		"key1=value1;key2=value2;..." Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//	Returns: *ConnectionParams a new ConnectionParams object.
func NewConnectionParamsFromString(line string) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParamsFromString(line),
	}
}

// NewConnectionParamsFromMaps static method for creating a StringValueMap using the maps passed as parameters.
//
//	Parameters:
//		- maps ...map[string]string the maps passed to this method to create a StringValueMap with.
//	Returns: ConnectionParams the ConnectionParams created.
func NewConnectionParamsFromMaps(maps ...map[string]string) *ConnectionParams {
	return &ConnectionParams{
		ConfigParams: config.NewConfigParamsFromMaps(maps...),
	}
}

// NewManyConnectionParamsFromConfig retrieves all ConnectionParams from configuration parameters
// from "connections" section. If "connection" section is present instead, than it returns a
// list with only one ConnectionParams.
//
//	Parameters:
//		- config *config.ConfigParams a configuration parameters to retrieve connections
//	Returns: []*ConnectionParams a list of retrieved ConnectionParams
func NewManyConnectionParamsFromConfig(config *config.ConfigParams) []*ConnectionParams {
	result := make([]*ConnectionParams, 0)

	connections := config.GetSection(SectionNameConnections)

	if connections.Len() > 0 {
		for _, section := range connections.GetSectionNames() {
			connection := connections.GetSection(section)
			result = append(result, NewConnectionParams(connection.Value()))
		}
	} else {
		connection := config.GetSection(SectionNameConnection)
		if connection.Len() > 0 {
			result = append(result, NewConnectionParams(connection.Value()))
		}
	}

	return result
}

// NewConnectionParamsFromConfig retrieves a single ConnectionParams from configuration parameters
// from "connection" section. If "connections" section is present instead,
// then is returns only the first connection element.
//
//	Parameters:
//		- config *config.ConfigParams ConnectionParams, containing a section named "connection(s)".
//	Returns: *ConnectionParams the generated ConnectionParams object.
func NewConnectionParamsFromConfig(config *config.ConfigParams) *ConnectionParams {
	connections := NewManyConnectionParamsFromConfig(config)
	if len(connections) > 0 {
		return connections[0]
	}
	return nil
}

// UseDiscovery checks if these connection parameters shall be retrieved from DiscoveryService.
// The connection parameters are redirected to DiscoveryService when discovery_key parameter is set.
//
//	Returns: bool true if connection shall be retrieved from DiscoveryService
func (c *ConnectionParams) UseDiscovery() bool {
	return c.GetAsString(ConnectionParamDiscoveryKey) != ""
}

// DiscoveryKey gets the key to retrieve this connection from DiscoveryService.
// If this key is null, then all parameters are already present.
//
//	see UseDiscovery
//	Returns: string the discovery key to retrieve connection.
func (c *ConnectionParams) DiscoveryKey() string {
	return c.GetAsString(ConnectionParamDiscoveryKey)
}

// SetDiscoveryKey sets the key to retrieve these parameters from DiscoveryService.
//
//	Parameters:
//		- value string a new key to retrieve connection.
func (c *ConnectionParams) SetDiscoveryKey(value string) {
	c.Put(ConnectionParamDiscoveryKey, value)
}

// Protocol gets the connection protocol.
//
//	Returns: string the connection protocol or the default value if it's not set.
func (c *ConnectionParams) Protocol() string {
	return c.GetAsString(ConnectionParamProtocol)
}

// ProtocolWithDefault gets the connection protocol.
//
//	Parameters:
//		- defaultValue string the default protocol
//	Returns: string the connection protocol or the default value if it's not set.
func (c *ConnectionParams) ProtocolWithDefault(defaultValue string) string {
	return c.GetAsStringWithDefault(ConnectionParamProtocol, defaultValue)
}

// SetProtocol sets the connection protocol.
//
//	Parameters:
//		- value string a new connection protocol.
func (c *ConnectionParams) SetProtocol(value string) {
	c.Put(ConnectionParamProtocol, value)
}

// Host gets the host name or IP address.
//
//	Returns: string the host name or IP address.
func (c *ConnectionParams) Host() string {
	host := c.GetAsString(ConnectionParamHost)
	if host != "" {
		return host
	}
	return c.GetAsString(ConnectionParamIp)
}

// SetHost sets the host name or IP address.
//
//	Parameters:
//		- value string a new host name or IP address.
func (c *ConnectionParams) SetHost(value string) {
	c.Put(ConnectionParamHost, value)
}

// Port gets the port number.
//
//	Returns int the port number.
func (c *ConnectionParams) Port() int {
	return c.GetAsInteger(ConnectionParamPort)
}

// PortWithDefault gets the port number.
//
//	Parameters:
//		- defaultValue int default port number
//	Returns: int the port number.
func (c *ConnectionParams) PortWithDefault(defaultValue int) int {
	return c.GetAsIntegerWithDefault(ConnectionParamPort, defaultValue)
}

// SetPort sets the port number.
//
//	see Host
//	Parameters: value int a new port number.
func (c *ConnectionParams) SetPort(value int) {
	c.Put(ConnectionParamPort, value)
}

// Uri gets the resource URI or connection string. Usually it includes all connection parameters in it.
//
//	Returns: string the resource URI or connection string.
func (c *ConnectionParams) Uri() string {
	return c.GetAsString(ConnectionParamURI)
}

// SetUri sets the resource URI or connection string.
//
//	Parameters:
//		- value string a new resource URI or connection string.
func (c *ConnectionParams) SetUri(value string) {
	c.Put(ConnectionParamURI, value)
}
