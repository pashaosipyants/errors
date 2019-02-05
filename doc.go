/*
Package errors provides error type for detailed stack-traced description and exception(handle-check) style handling

The traditional error handling idiom in Go is roughly akin to

    if err != nil {
            return err
    }

which when applied recursively up the call stack results in error reports
without context or debugging information. The errors package allows
programmers to add context to the failure path in their code in a way
that does not destroy the original value of the error.

Skipstack management

The errors.Wrap function returns a new error that adds context to the
original error by recording a stack trace at the point Wrap is called,
together with the supplied message. For example

Format

Annotations

Skipstack_management
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects
*/
package errors
