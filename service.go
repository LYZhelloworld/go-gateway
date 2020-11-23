package gateway

import "strings"

// Service is the configuration of one service.
type Service struct {
	// Name is the identifier of the service, separated by dots, with parent Services before sub Services.
	// For example: "grand_parent_service.parent_service.service".
	Name ServiceName
	// Handler is the handler function of a service
	Handler ServiceHandler
}

// ServiceHandler is a function that handles the service
type ServiceHandler func(c *Context)

type ServiceName string

func (n *ServiceName) split() []string {
	return strings.Split(string(*n), ".")
}

// match matches a name with the current service.
// Distance is the distance between serviceName and Service.Name. 0 means the service name is exact the same.
// The value of distance is meaningful only if the ok value is true.
func (n *Service) match(name ServiceName) (ok bool, distance int) {
	this := n.Name.split()
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
