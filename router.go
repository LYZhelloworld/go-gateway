package gateway

import "net/http"

// router holds handlers of all endpoints and other configs.
type router struct {
	// endpointConfig is a map with endpoint as key and routerConfig as value.
	endpointConfig map[string]*routerConfig

	// errorConfig is a map with status code as key and ServiceHandler as value.
	errorConfig map[int]*ServiceHandler // TODO
}

// routerConfig holds handlers for different methods.
type routerConfig struct {
	getHandler    ServiceHandler
	postHandler   ServiceHandler
	putHandler    ServiceHandler
	deleteHandler ServiceHandler
}

// ServeHTTP serves HTTP requests.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request: req,
	}

	config := r.endpointConfig[req.URL.RawPath]
	if config == nil {
		r.generalResponse(ctx, http.StatusNotFound)
		_ = r.write(w, ctx)
		return
	}

	var handler ServiceHandler
	switch req.Method {
	case http.MethodGet:
		if config.getHandler != nil {
			handler = config.getHandler
		}
	case http.MethodPost:
		if config.postHandler != nil {
			handler = config.postHandler
		}
	case http.MethodPut:
		if config.putHandler != nil {
			handler = config.putHandler
		}
	case http.MethodDelete:
		if config.deleteHandler != nil {
			handler = config.deleteHandler
		}
	default:
		r.generalResponse(ctx, http.StatusMethodNotAllowed)
		_ = r.write(w, ctx)
		return
	}
	if handler == nil {
		r.generalResponse(ctx, http.StatusNotFound)
		_ = r.write(w, ctx)
		return
	}
	handler(ctx)
	_ = r.write(w, ctx)
	return
}

// write writes response to the http.ResponseWriter.
func (r *router) write(w http.ResponseWriter, context *Context) (err error) {
	w.WriteHeader(context.statusCode)
	_, err = w.Write(context.response)
	return
}

// generalResponse generates error messages depending on the status code.
func (r *router) generalResponse(context *Context, statusCode int) {
	context.SetStatusCode(statusCode)
	if handler, ok := r.errorConfig[statusCode]; ok {
		(*handler)(context)
		return
	}
}
