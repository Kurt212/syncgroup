### Syncgroup

[![Go Report Card](https://goreportcard.com/badge/github.com/kurt212/syncgroup)](https://goreportcard.com/report/github.com/kurt212/syncgroup)
![Build Status](https://github.com/kurt212/syncgroup/actions/workflows/ci.yml/badge.svg)
![GitHub issues](https://img.shields.io/github/issues/kurt212/syncgroup)
![GitHub pull requests](https://img.shields.io/github/issues-pr/kurt212/syncgroup)

[![Go Reference](https://pkg.go.dev/badge/github.com/kurt212/syncgroup.svg)](https://pkg.go.dev/github.com/kurt212/syncgroup)
![Go Version](https://img.shields.io/github/go-mod/go-version/kurt212/syncgroup)

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

Run goroutines in parallel and wait until all of them finish.

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

Collect errors from goroutines.

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

			return fmt.Errorf("error %d", i)
		})
	}

	err := sg.Wait()
	if err != nil {
		/*
			Expected output:
			error 1
			error 3
			error 2
			error 4
			error 6
			error 5
			error 0
			error 9
			error 7
			error 8
		 */
		fmt.Println(err)
	}
}
```

Limit number of concurrent goroutines.

```go
package main

import (
	"fmt"
	"time"

	"github.com/kurt212/syncgroup"
)

func main() {
	sg := syncgroup.New()

	sg.SetLimit(2)

	for i := range 10 {
		sg.Go(func() error {
			fmt.Printf("Go %d\n", i)

			time.Sleep(1 * time.Second)

			return nil
		})
	}

	sg.Wait()
}
```

## Contributing

Feel free to contribute to this project. You can report bugs, suggest features or submit pull requests.

Before submitting a bug report or a feature request, check if there is an existing one and provide as much information as possible.

### Submitting a pull request

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Run checks (`make all`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create a new Pull Request
7. Wait for CI to pass
8. Profit! 🎉
