/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uerror

// Must takes a variadic number of interface{} arguments and panics if the last argument is a non-nil error.
// This allows Must function to handle functions returning multiple values, where the last value is an error.
//
// Examples:
//
//   - With a function returning only an error:
//     Must(doSomethingRisky())
//
//   - With a function returning a result and an error:
//     data, err := os.ReadFile("./myfile.txt")
//     Must(data, err) // Panics if err is non-nil
//
//   - Direct usage with a function that returns an error:
//     Must(os.Remove("./myfile.txt")) // Panics if Remove returns an error
func Must(args ...any) {
	if len(args) == 0 {
		return
	}

	lastArg := args[len(args)-1]
	if err, ok := lastArg.(error); ok && err != nil {
		panic(err)
	}
}
