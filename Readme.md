 # Hooks
-- store global callbacks concurrently

[![GoDoc](https://godoc.org/github.com/rasteric/hooks/go?status.svg)](https://godoc.org/github.com/rasteric/hooks)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasteric/hooks)](https://goreportcard.com/report/github.com/rasteric/hooks)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

Hooks is a package that implements a global repository of callbacks called _hooks_, which are functions that take an arbitrary number of arguments and return a value and an error. Every hook has a fixed integer ID, which must be attributed by the user of the package. The package is thread safe.

## Usage

`type HookFunc func(a []interface{}) (interface{}, error)`

The function signature of a hook.

`Add(id int, f HookFunc)`

Add the hook with given ID. An existing callback will be overwritten.

`Exec(id int, args ...interface{})`

Execute the hook with given ID and arguments.

`Remove(id int)`

Remove the hook with given ID.

