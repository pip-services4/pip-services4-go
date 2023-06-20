package connect

import (
	"net/url"
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

var ConnectionUtils = _TConnectionUtils{}

// _TConnectionUtils a set of utility functions to process connection parameters
type _TConnectionUtils struct{}

// Concat concatenates two options by combining duplicated properties into comma-separated list
//
//	Parameters:
//		- options1 first options to merge
//		- options2 second options to merge
//	Returns: keys when define it limits only to specific keys
func (c *_TConnectionUtils) Concat(options1 *config.ConfigParams, options2 *config.ConfigParams, keys ...string) *config.ConfigParams {
	options := config.NewConfigParamsFromValue(options1)
	for _, key := range options2.Keys() {
		value1 := options1.GetAsString(key)
		value2 := options2.GetAsString(key)

		if value1 != "" && value2 != "" {
			if len(keys) == 0 || indexOf(keys, key) >= 0 {
				options.SetAsObject(key, value1+","+value2)
			}
		} else if value1 != "" {
			options.SetAsObject(key, value1)
		} else if value2 != "" {
			options.SetAsObject(key, value2)
		}
	}
	return options
}

func (c *_TConnectionUtils) concatValues(value1 string, value2 string) string {
	if value1 == "" {
		return value2
	}
	if value2 == "" {
		return value1
	}
	return value1 + "," + value2
}

// ParseUri parses URI into config parameters.
//
//	The URI shall be in the following form:
//		protocol://username@password@host1:port1,host2:port2,...?param1=abc&param2=xyz&...
//	Parameters:
//		- uri the URI to be parsed
//		- defaultProtocol a default protocol
//		- defaultPort a default port
//	Returns: a configuration parameters with URI elements
func (c *_TConnectionUtils) ParseUri(uri string, defaultProtocol string, defaultPort int) *config.ConfigParams {
	options := config.NewEmptyConfigParams()

	if uri == "" {
		return options
	}

	uri = strings.TrimSpace(uri)

	// Process parameters
	pos := strings.Index(uri, "?")
	if pos > 0 {
		params := uri[pos+1:]
		uri = uri[:pos]

		paramsList := strings.Split(params, "&")
		for _, param := range paramsList {
			pos := strings.Index(param, "=")
			if pos >= 0 {
				key, _ := url.QueryUnescape(param[:pos])
				value, _ := url.QueryUnescape(param[pos+1:])
				options.SetAsObject(key, value)
			} else {
				param, _ = url.QueryUnescape(param)
				options.SetAsObject(param, "")
			}
		}
	}

	// Process protocol
	pos = strings.Index(uri, "://")
	if pos > 0 {
		protocol := uri[:pos]
		uri = uri[pos+3:]
		options.SetAsObject("protocol", protocol)
	} else {
		options.SetAsObject("protocol", defaultProtocol)
	}

	// Process user and password
	pos = strings.Index(uri, "@")
	if pos > 0 {
		userAndPass := uri[:pos]
		uri = uri[pos+1:]

		pos = strings.Index(userAndPass, ":")
		if pos > 0 {
			options.SetAsObject("username", userAndPass[:pos])
			options.SetAsObject("password", userAndPass[pos+1:])
		} else {
			options.SetAsObject("username", userAndPass)
		}
	}

	// Process host and ports
	// options.setAsObject("servers", c.concatValues(options.getAsString("servers"), uri));
	servers := strings.Split(uri, ",")
	for _, server := range servers {
		pos = strings.Index(server, ":")
		if pos > 0 {
			options.SetAsObject("servers", c.concatValues(options.GetAsString("servers"), server))
			options.SetAsObject("host", c.concatValues(options.GetAsString("host"), server[:pos]))
			options.SetAsObject("port", c.concatValues(options.GetAsString("port"), server[pos+1:]))
		} else {
			options.SetAsObject("servers", c.concatValues(options.GetAsString("servers"), server+":"+convert.StringConverter.ToString(defaultPort)))
			options.SetAsObject("host", c.concatValues(options.GetAsString("host"), server))
			options.SetAsObject("port", c.concatValues(options.GetAsString("port"), convert.StringConverter.ToString(defaultPort)))
		}
	}

	return options
}

// ComposeUri composes URI from config parameters.
//
//	The result URI will be in the following form:
//		protocol://username@password@host1:port1,host2:port2,...?param1=abc&param2=xyz&...
//	Parameters:
//		- options configuration parameters
//		- defaultProtocol a default protocol
//		- defaultPort a default port
//	Returns: a composed URI
func (c *_TConnectionUtils) ComposeUri(options *config.ConfigParams, defaultProtocol string, defaultPort int) string {
	builder := ""

	protocol := options.GetAsStringWithDefault("protocol", defaultProtocol)
	if protocol != "" {
		builder = protocol + "://" + builder
	}

	if username, ok := options.GetAsNullableString("username"); ok && username != "" {
		builder += username
		if password, ok := options.GetAsNullableString("password"); ok && password != "" {
			builder += ":" + password
		}
		builder += "@"
	}

	servers := ""
	defaultPortStr := ""
	if defaultPort > 0 {
		defaultPortStr = convert.StringConverter.ToString(defaultPort)
	}
	hosts := strings.Split(options.GetAsStringWithDefault("host", "???"), ",")
	ports := strings.Split(options.GetAsStringWithDefault("port", defaultPortStr), ",")
	for index := range hosts {
		if len(servers) > 0 {
			servers += ","
		}

		host := hosts[index]
		servers += host

		port := defaultPortStr
		if len(ports) > index && ports[index] != "" {
			port = ports[index]
		}

		if port != "" {
			servers += ":" + port
		}
	}
	builder += servers

	params := ""
	reservedKeys := []string{"protocol", "host", "port", "username", "password", "servers"}
	for _, key := range options.Keys() {
		if indexOf(reservedKeys, key) >= 0 {
			continue
		}

		if len(params) > 0 {
			params += "&"
		}
		params += url.QueryEscape(key)

		if value, ok := options.GetAsNullableString(key); ok && value != "" {
			params += "=" + url.QueryEscape(value)
		}
	}

	if len(params) > 0 {
		builder += "?" + params
	}

	return builder
}

// Include includes specified keys from the config parameters.
//
//	Parameters:
//		- options configuration parameters to be processed.
//		- keys a list of keys to be included.
//	Returns: a processed config parameters.
func (c *_TConnectionUtils) Include(options *config.ConfigParams, keys ...string) *config.ConfigParams {
	if len(keys) == 0 {
		return options
	}

	result := config.NewEmptyConfigParams()

	for _, key := range options.Keys() {
		if indexOf(keys, key) >= 0 {
			result.SetAsObject(key, options.GetAsString(key))
		}
	}

	return result
}

// Exclude specified keys from the config parameters.
//
//	Parameters:
//		- options configuration parameters to be processed.
//		- keys a list of keys to be excluded.
//	Returns: a processed config parameters.
func (c *_TConnectionUtils) Exclude(options *config.ConfigParams, keys ...string) *config.ConfigParams {
	if len(keys) == 0 {
		return options
	}

	if options == nil {
		return nil
	}

	result := config.NewConfigParamsFromValue(options)
	for _, key := range keys {
		result.Remove(key)
	}
	return result
}
