package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (v *vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

func (v *vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	return types.OnPluginStartStatusOK
}

func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &RequestHandler{}
}

type RequestHandler struct {
	// Bring in the callback functions
	types.DefaultHttpContext
}

const (
	XRequestIdHeader = "x-request-id"
)

func (r *RequestHandler) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	err := proxywasm.AddHttpResponseHeader("x-my-response-header", "Some value of the header")
	if err != nil {
		proxywasm.LogCriticalf("failed to add request header: %v", err)
	}
	return types.ActionContinue
}

func (r *RequestHandler) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfof("WASM plugin Handling request")

	// Get the actual request headers from the Envoy Sidecar
	requestHeaders, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		// Allow Envoy Sidecar to forward this request to the upstream service
		return types.ActionContinue
	}

	// Convert the request headers to a map for easier access (more useful in subsequent sections)
	reqHeaderMap := headerArrayToMap(requestHeaders)

	// Get the x-request-id for grouping logs belonging to the same request
	xRequestID := reqHeaderMap[XRequestIdHeader]

	// Now we can take action on this request
	return r.doSomethingWithRequest(reqHeaderMap, xRequestID)
}

// headerArrayToMap is a simple function to convert from array of headers to a Map
func headerArrayToMap(requestHeaders [][2]string) map[string]string {
	headerMap := make(map[string]string)
	for _, header := range requestHeaders {
		headerMap[header[0]] = header[1]
	}
	return headerMap
}

func (r *RequestHandler) doSomethingWithRequest(reqHeaderMap map[string]string, xRequestID string) types.Action {
	// for now, let's just log all the request headers to we get an idea of what we have to work with
	for key, value := range reqHeaderMap {
		proxywasm.LogInfof("  %s: request header --> %s: %s", xRequestID, key, value)
	}

	// Forward request to upstream service, i.e. unblock request
	return types.ActionContinue
}
