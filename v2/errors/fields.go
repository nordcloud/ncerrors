package errors

type Fields map[string]interface{}

func (f Fields) Add(key string, val interface{}) Fields {
	newFields := f.copy()
	newFields[key] = val

	return newFields
}

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
