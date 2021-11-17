package errors

// A collection of all things we know about the error.
type Info struct {
	Message    string   `json:"message,omitempty"`
	Fields     Fields   `json:"fields,omitempty"`
	StackTrace []string `json:"stackTrace,omitempty"`
	FuncName   string   `json:"funcName,omitempty"`
}

func GetInfo(err error) []Info {
	if err == nil {
		return nil
	}

	var infos []Info

	for err != nil {
		// This one should be type asserted instead of using errors.As,
		// because in case err does not implement Info but instead implements Unwrap
		// we can get an info for the unwrapped error instead
		if infoer, ok := err.(Infoer); ok {
			infos = append(infos, infoer.Info())
		} else {
			infos = append(infos, Info{
				Message: err.Error(),
			})
		}

		var unwrapper Unwrapper
		if As(err, &unwrapper) {
			err = unwrapper.Unwrap()
		} else {
			err = nil
		}
	}

	return infos
}
