### Syncgroup

This is a package that contains an implementation of an abstract
synchronisation mechanism - synchronisation group.
The main idea is to have an ability to run independent tasks in separate goroutines which way return errors.
A user can wait until all goroutines finish running and collect all occurred errors.

The design is similar to errgroup (https://godoc.org/golang.org/x/sync/errgroup),
but it does not cancel the context of the goroutines if any of them returns an error.

## Documentation
See more on [godoc site](https://godoc.org/github.com/Kurt212/syncgroup)
