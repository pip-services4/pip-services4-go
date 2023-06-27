package controllers

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

// HttpRequestDetector helper class that retrieves parameters from HTTP requests.
var HttpRequestDetector = _THttpRequestDetector{}

type _THttpRequestDetector struct {
}

var (
	macOsXRegExp           = regexp.MustCompile(`/CPU( iPhone)? OS ([0-9\._]+) like Mac OS X/`)
	globalUnderscoreRegExp = regexp.MustCompile("/_/g")
	androidRegExp          = regexp.MustCompile(`/Android ([0-9\.]+)[\);]/`)
	webOsRegExp            = regexp.MustCompile(`/webOS\/([0-9\.]+)[\);]\//`)
	intelMacOsXRegExp      = regexp.MustCompile(`/(Intel|PPC) Mac OS X ?([0-9\._]*)[\)\;]/`)
	windowsNtRegExp        = regexp.MustCompile(`/Windows NT ([0-9\._]+)[\);]/`)
)

// DetectPlatform method are detects the platform (using "user-agent")
// from which the given HTTP request was made.
//
//	Parameters:
//		-  req  *http.Request an HTTP request to process.
//	Returns: the detected platform and version. Detectable platforms: "mobile", "iphone",
//		"ipad",  "macosx", "android",  "webos", "mac", "windows". Otherwise - "unknown" will
//		be returned.
func (c *_THttpRequestDetector) DetectPlatform(req *http.Request) string {
	ua := req.Header.Get("user-agent")
	var version string
	var pattern string

	pattern = "/mobile/i"
	match, _ := regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "mobile"
	}

	pattern = "/like Mac OS X/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		result := macOsXRegExp.FindAllStringSubmatch(ua, -1)
		version = globalUnderscoreRegExp.ReplaceAllString(result[0][2], ".")

		pattern = "/iPhone/"
		match, _ = regexp.Match(pattern, ([]byte)(ua))
		if match {
			return "iphone " + version
		}
		pattern = "/iPad/"
		match, _ = regexp.Match(pattern, ([]byte)(ua))

		if match {
			return "ipad " + version
		}
		return "macosx " + version
	}

	pattern = "/Android/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		version = androidRegExp.FindAllStringSubmatch(ua, -1)[0][1]
		return "android " + version
	}

	pattern = `/webOS\//`
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		version = webOsRegExp.FindAllStringSubmatch(ua, -1)[0][1]
		return "webos " + version
	}

	pattern = "/(Intel|PPC) Mac OS X/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		result := intelMacOsXRegExp.FindAllStringSubmatch(ua, -1)
		version = globalUnderscoreRegExp.ReplaceAllString(result[0][2], ".")
		return "mac " + version
	}

	pattern = "/Windows NT/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		version = windowsNtRegExp.FindAllStringSubmatch(ua, -1)[0][1]
		return "windows " + version
	}
	return "unknown"
}

// DetectBrowser detects the browser (using "user-agent") from which the given HTTP request was made.
//
//	Parameters:
//		-  req  *http.Reques an HTTP request to process.
//	Returns: the detected browser. Detectable browsers: "chrome", "msie", "firefox",
//		"safari". Otherwise - "unknown" will be returned.
func (c *_THttpRequestDetector) DetectBrowser(req *http.Request) string {

	ua := req.Header.Get("user-agent")

	var pattern string
	pattern = "/chrome/i"
	match, _ := regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "chrome"
	}

	pattern = "/msie/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "msie"
	}

	pattern = "/firefox/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "firefox"
	}

	pattern = "/safari/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "safari"
	}

	if ua == "" {
		return "unknown"
	}
	return ua
}

// DetectAddress method are detects the IP address from which the given HTTP request was received.
//
//	Parameters:
//		- req *http.Reques an HTTP request to process.
//	Returns the detected IP address (without a port). If no IP is detected -
//		nil will be returned.
func (c *_THttpRequestDetector) DetectAddress(req *http.Request) string {
	var ip string

	if len(req.Header.Get("CF-Connecting-IP")) > 1 {
		ip = req.Header.Get("CF-Connecting-IP")
		ip = string(net.ParseIP(ip))
	} else if len(req.Header.Get("X-Forwarded-For")) > 1 {
		ip = req.Header.Get("X-Forwarded-For")
		ip = string(net.ParseIP(ip))
	} else if len(req.Header.Get("X-Real-IP")) > 1 {
		ip = req.Header.Get("X-Real-IP")
		ip = string(net.ParseIP(ip))
	} else {
		ip = req.RemoteAddr
		if strings.Contains(ip, ":") {
			res := strings.Split(ip, ":")
			if len(res) > 0 {
				ip = net.ParseIP(res[0]).String()
			}
		} else {
			ip = net.ParseIP(ip).String()
		}
	}

	return ip
}

// DetectServerHost method are detects the host name of the request"s destination server.
//
//	Parameters:
//		- req *http.Request  an HTTP request to process.
//	Returns: the destination server"s host name.
func (c *_THttpRequestDetector) DetectServerHost(req *http.Request) string {
	//TODO: Need fix this
	return "" + req.URL.Hostname() // socket.localAddress
}

// DetectServerPort method are detects the request`s destination port number.
//
//	Parameters:
//		- req  *http.Request an HTTP request to process.
//	Returns: the detected port number or <code>80</code> (if none are detected).
func (c *_THttpRequestDetector) DetectServerPort(req *http.Request) string {
	return req.URL.Port() //socket.localPort
}
