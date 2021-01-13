package errors

// AnyE is a helper function that returns first not nil error or nil if there are none.
func AnyE(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// AnyFuncE is a helper function that executes funcs one by one, stops if error occurred
// and returns it
func AnyFuncE(errfuncs ...func() error) error {
	for _, f := range errfuncs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
