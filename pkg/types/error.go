package types

type FormatError string

func (e FormatError) Error() string {
	return string(e)
}

func (e Result) Error() string {
	return e.Message
}
