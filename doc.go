/*
Package errors provides error type for detailed stack-traced description and exception(handle-check) style handling

The traditional error handling idiom in Go is roughly akin to

    if err != nil {
            return err
    }

Format

Annotations

Skipstack_management
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects
*/
package errors
