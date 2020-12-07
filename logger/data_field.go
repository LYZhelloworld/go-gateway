package logger

// field is a data field that may contain a key/value pair, a lazy-loaded field, or error.
type field interface {
	// getKey gets key.
	getKey() string
	// getValue gets value.
	getValue() interface{}
}

// normalField is the normal key/value pair.
type normalField struct {
	key   string
	value interface{}
}

// getKey returns key.
func (n normalField) getKey() string {
	return n.key
}

// getValue returns value.
func (n normalField) getValue() interface{} {
	return n.value
}

// errorField is the error field.
type errorField struct {
	err error
}

// getKey returns key.
func (e errorField) getKey() string {
	return "err"
}

// getValue returns value.
func (e errorField) getValue() interface{} {
	return e.err.Error()
}

// lazyField is the lazy-loaded key/value pair.
type lazyField struct {
	key   string
	value func() interface{}
}

// getKey returns key.
func (l lazyField) getKey() string {
	return l.key
}

// getValue returns value.
func (l lazyField) getValue() interface{} {
	return l.value()
}
