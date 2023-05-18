package repo_err

type EntityNotFoundError struct {
	Msg string
}

func (e *EntityNotFoundError) Error() string {
	return e.Msg
}
