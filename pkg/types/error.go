package types

type ScriptError struct {
	Message  string
	ExitCode int
}

func (e *ScriptError) Error() string {
	return e.Message
}
