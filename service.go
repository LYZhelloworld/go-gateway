package gateway

import "strings"

const (
	baseServiceHandler = "*"
)

// Service is the configuration of one service.
type Service struct {
	// Name is the identifier of the service, separated by dots, with parent Services before sub Services.
	// For example: "foo.bar.baz".
	//
	// A service can be handled by a more generic service name (the request of which can be forwarded to other services).
	// For example: "foo.bar" can handle "foo.bar.baz" requests.
	// But "foo.bar.baz" cannot handle "foo.bar".
	//
	// An asterisk (*) means a service handler for all services, if there is no other services that are more specific.
	Name ServiceName
	// Handler is the handler function of a service.
	Handler ServiceHandler
}

// ServiceHandler is a function that handles the service.
type ServiceHandler func(context *Context)

type ServiceName string

func (n *ServiceName) split() []string {
	return strings.Split(string(*n), ".")
}

// match matches a name with the current service.
// Distance is the distance between serviceName and Service.Name. 0 means the service name is exact the same.
// The value of distance is meaningful only if the ok value is true.
func (n *Service) match(name ServiceName) (ok bool, distance int) {
	thisName := n.Name
	if thisName == baseServiceHandler {
		thisName = ""
	}
	this := thisName.split()
	other := name.split()
	lenThis := len(this)
	lenOther := len(other)
	distance = lenOther - lenThis
	if lenOther < lenThis {
		// the service name is too generic
		ok = false
		return
	}
	for i := 0; i < lenThis; i++ {
		if this[i] != other[i] {
			ok = false
			return
		}
	}
	ok = true
	return
}
