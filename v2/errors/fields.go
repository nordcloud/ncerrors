package errors

// Fields is used to store additional information about error.
// To avoid duplication of information contained within Fields, you should not add any data that is already available
// to the caller of the function where the error happend.
// Information that should be contained within fields:
// - data created by current function
// - data received from external sources
// - a specific element of a collection for which the error has happend
type Fields map[string]interface{}

// Add allows to add a single field to existing Fields.
func (f Fields) Add(key string, val interface{}) Fields {
	newFields := f.copy()
	newFields[key] = val

	return newFields
}

// Extend allows to extend existing fields with a new collection of Fields.
// If given key already exists in the Fields it will be overriden by the new Fields value.
func (f Fields) Extend(extFields Fields) Fields {
	newFields := f.copy()
	for k, v := range extFields {
		newFields[k] = v
	}

	return newFields
}

func (f Fields) copy() Fields {
	newFields := make(map[string]interface{}, len(f))
	for k, v := range f {
		newFields[k] = v
	}

	return newFields
}
