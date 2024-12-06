package main

type Errors struct {
	errors []error
}

func NewErrors() *Errors {
	return &Errors{
		errors: []error{},
	}
}

func (e *Errors) Append(errs ...error) {
	for _, err := range errs {
		if err != nil {
			e.errors = append(e.errors, err)
		}
	}
}

func (e *Errors) Error() string {
	s := ""
	for _, err := range e.errors {
		s += err.Error() + "\n"
	}
	return s
}
