# Syncgroup

[![Go Report Card](https://goreportcard.com/badge/github.com/kurt212/syncgroup)](https://goreportcard.com/report/github.com/kurt212/syncgroup)
![Build Status](https://github.com/kurt212/syncgroup/actions/workflows/ci.yml/badge.svg)
![GitHub issues](https://img.shields.io/github/issues/kurt212/syncgroup)
![GitHub pull requests](https://img.shields.io/github/issues-pr/kurt212/syncgroup)

[![Go Reference](https://pkg.go.dev/badge/github.com/kurt212/syncgroup.svg)](https://pkg.go.dev/github.com/kurt212/syncgroup)
![Go Version](https://img.shields.io/github/go-mod/go-version/kurt212/syncgroup)

## Introduction

Syncgroup is a Go package that provides an abstract synchronization mechanism, allowing you to run independent tasks in separate goroutines and collect all occurred errors. It is similar to `errgroup` but does not cancel the context of the goroutines if any of them returns an error.

## Key Features

- **Convenient API**: Run goroutines and wait for their completion with `sg.Go(func() error)` and `sg.Wait()`.
- **Error Handling**: Collects all errors from goroutines and returns them as a single error, wrapped according to [Go 1.13 errors wrapping rules](https://go.dev/blog/go1.13-errors).
- **Panic Recovery**: Recovers panics in goroutines and returns them as errors.
- **Concurrency Limiting**: Set a limit on the number of concurrent goroutines.

## Differences from Industry Standards

**SyncGroup** offers several enhancements over `sync.WaitGroup` and `errgroup`:

- **Better API**: Simplified and more convenient.
- **Error Handling**: Automatically collects and returns all errors.
- **Panic Recovery**: Recovers from panics and returns them as errors.
- **Concurrency Limiting**: Allows setting a limit on concurrent goroutines.
- **Comprehensive Error Collection**: Unlike `errgroup`, `SyncGroup` does not cancel the context of the goroutines if any of them returns an error. 
`errgroup` is designed to obtain the result only if all jobs are successful. When multiple errors occur, `errgroup` only returns the first one and ignores the rest.

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

Limit the number of concurrent goroutines.

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
