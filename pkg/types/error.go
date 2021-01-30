package types

type FormatError string

func (e FormatError) Error() string {
	return string(e)
}
