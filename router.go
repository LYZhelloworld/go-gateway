package gateway

import "net/http"

// router holds handlers of all endpoints and other configs.
type router struct {
	// config is a map with endpoint as key and routerConfig as value.
	config EndpointConfig
	// errorConfig is a map with status code as key and ServiceHandler as value.
	errorConfig ErrorConfig
}

// ServeHTTP serves HTTP requests.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:        req,
		StatusCode:     http.StatusOK,
		Response:       nil,
		Header:         map[string][]string{},
		responseWriter: w,
	}

	config := r.config[req.URL.RawPath]
	if config == nil {
		r.generalResponse(ctx, http.StatusNotFound)
		ctx.write()
		return
	}

	service, ok := config.getService(req.Method)
	if !ok {
		r.generalResponse(ctx, http.StatusMethodNotAllowed)
		ctx.write()
		return
	}
	ctx.serviceName = service.Name
	if service == nil {
		r.generalResponse(ctx, http.StatusNotFound)
		ctx.write()
		return
	}
	service.Handler(ctx)
	ctx.write()
	return
}

// generalResponse generates error messages depending on the status code.
func (r *router) generalResponse(context *Context, statusCode int) {
	context.StatusCode = statusCode
	if handler, ok := r.errorConfig[statusCode]; ok {
		(*handler)(context)
		return
	}
}
