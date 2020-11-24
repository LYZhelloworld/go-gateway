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

	var handler ServiceHandler
	switch req.Method {
	case http.MethodGet:
		if config.getHandler != nil {
			ctx.serviceName = config.getHandler.Name
			handler = config.getHandler.Handler
		}
	case http.MethodPost:
		if config.postHandler != nil {
			ctx.serviceName = config.postHandler.Name
			handler = config.postHandler.Handler
		}
	case http.MethodPut:
		if config.putHandler != nil {
			ctx.serviceName = config.putHandler.Name
			handler = config.putHandler.Handler
		}
	case http.MethodDelete:
		if config.deleteHandler != nil {
			ctx.serviceName = config.deleteHandler.Name
			handler = config.deleteHandler.Handler
		}
	default:
		r.generalResponse(ctx, http.StatusMethodNotAllowed)
		ctx.write()
		return
	}
	if handler == nil {
		r.generalResponse(ctx, http.StatusNotFound)
		ctx.write()
		return
	}
	handler(ctx)
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
