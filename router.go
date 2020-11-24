package gateway

import "net/http"

// router is a map with endpoint as key and routerConfig as value.
type router map[string]*routerConfig

// routerConfig holds handlers for different methods.
type routerConfig struct {
	getHandler    ServiceHandler
	postHandler   ServiceHandler
	putHandler    ServiceHandler
	deleteHandler ServiceHandler
}

// ServeHTTP serves HTTP requests.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	config := (*r)[req.URL.RawPath]
	if config == nil {
		// TODO: 404
	}

	switch req.Method {
	case Get:
		// TODO
	case Post:
		// TODO
	case Put:
		// TODO
	case Delete:
		// TODO
	default:
		return
	}
}
