package types

func (e Result) Error() string {
	return e.Message
}
