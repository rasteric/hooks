 # Hooks
   -- store global callbacks concurrently

[![GoDoc](https://godoc.org/github.com/rasteric/hooks/go?status.svg)](https://godoc.org/github.com/rasteric/hooks)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasteric/hooks)](https://goreportcard.com/report/github.com/rasteric/hooks)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

Hooks is a package that implements a global repository of callbacks called _hooks_, which are functions that take an array of empty interface values and do not return anything. Every hook has a fixed integer number, which must be attributed by the user of the package. Adding a function returns an integer ID for that function that can be used to remove it later. This might seem redundant and makes the package slower, but it is a requirement for the use it is intended for. The package is thread safe.

## Usage

`Add(hook int, f func(a []interface{})) int`

Add a callback for hook and return a numerical ID for the function that was added.

`Exec(hook int, args ...interface{})`

Execute all functions for the hook if there are any.

`Remove(hook, id int)`

Remove the function with given ID from the hook.

`RemoveAll(hook int)`

Remove all functions for the hook.

`Active(id int) bool`

Returns true if one or more functions for the hook are set, false otherwise. Using this function to check first may be more efficient than calling Exec directly, because of the arguments that Exec takes.
