package errors

// AnyErr is helper function that returns first not nil error or nil if there are none.
func AnyErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// AnyFuncErr is helper function that executes funcs one by one, stops if error occurred
// and returns it
func AnyFuncErr(errfuncs ...func() error) error {
	for _, f := range errfuncs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
