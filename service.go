package gateway

// Service is the configuration of one service.
// The key is the identifier of the service, separated by dots, with parent service before sub service.
// For example: "foo.bar.baz".
//
// The value is the handler function of a service.
//
// A service can be handled by a more generic service name (the request of which can be forwarded to other service).
// For example: "foo.bar" can handle "foo.bar.baz" requests.
// But "foo.bar.baz" cannot handle "foo.bar".
//
// An asterisk (*) means a service handler for all service, if there is no other service that are more specific.
type Service map[string]Handler

// Handler is a function that handles the service.
type Handler func(context *Context)
