### Syncgroup

[![Go Reference](https://pkg.go.dev/badge/github.com/kurt212/syncgroup.svg)](https://godoc.org/github.com/kurt212/syncgroup)

This is a package that contains an implementation of an abstract
synchronisation mechanism - synchronisation group.
The main idea is to have an ability to run independent tasks in separate goroutines which way return errors.
A user can wait until all goroutines finish running and collect all occurred errors.

The design is similar to errgroup (https://godoc.org/golang.org/x/sync/errgroup),
but it does not cancel the context of the goroutines if any of them returns an error.

## Documentation
See more on [godoc site](https://godoc.org/github.com/kurt212/syncgroup)

## Usage

### Installation

```shell
go get github.com/kurt212/syncgroup
```

### Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/kurt212/syncgroup"
)

func main() {
	sg := syncgroup.New()

	for i := range 10 {
		sg.Go(func() error {
			time.Sleep(1 * time.Second)

			fmt.Printf("Hello from %d\n", i)

			return nil
		})
	}

	sg.Wait()
}
```
