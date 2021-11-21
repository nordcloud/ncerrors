package errors

// A collection of all things we know about the error.
type Info struct {
	Message    string   `json:"message,omitempty"`
	Fields     Fields   `json:"fields,omitempty"`
	StackTrace []string `json:"stackTrace,omitempty"`
	FuncName   string   `json:"funcName,omitempty"`
}

// GetInfo retuns a list of infos for each error in chain.
// Even if there is a non-NCError in the chain, GetInfo will try to extract some information.
//
// It returns a simple struct that can easily be marshalled, so it can be directly logged by logger that supports formatting JSON data.
func GetInfo(err error) []Info {
	if err == nil {
		return nil
	}

	var infos []Info

	for err != nil {
		// This one should be type asserted instead of using errors.As,
		// because in case err does not implement Info but instead implements Unwrap
		// we can get an info for the unwrapped error instead
		if infoer, ok := err.(infoer); ok {
			infos = append(infos, infoer.Info())
		} else {
			infos = append(infos, Info{
				Message: err.Error(),
			})
		}

		var unwrapper unwrapper
		if As(err, &unwrapper) {
			err = unwrapper.Unwrap()
		} else {
			err = nil
		}
	}

	return infos
}

type unwrapper interface {
	Unwrap() error
}

type infoer interface {
	Info() Info
}
