package errors

/*
	retrieve contexts from an error
*/

// Gets value by key saved in err. Returns nil if not found.
func ValueE(err error, key interface{}) interface{} {
	var errVal *_errorValue
	for errIteration := err; AsE(errIteration, &errVal); errIteration = errVal.error {
		if errVal.key == key {
			return errVal.value
		}
	}
	return nil
}

func SuppressedE(err error) (supps []error) {
	var errSupp *_errorSuppressed
	for errIteration := err; AsE(errIteration, &errSupp); errIteration = errSupp.error {
		supps = append(supps, errSupp.suppressed)
	}
	return
}
