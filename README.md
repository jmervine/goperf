# goperf

Simple HTTPerf clone for performance testing web applications written in Go.

> NOTE: This is the inital commit and shouldn't be considered ready for anyone. That said, it should
> work as outlined below, at least on Linux based systems.

#### Supports: Go 1.1+

## Install

```
# manual
$ git clone https://github.com/jmervine/goperf.git $GOPATH/src/github.com/jmervine/goperf
$ cd $GOPATH/src/github.com/jmervine/goperf
$ make build/exe
$ cp pkg/goperf-VERSION $BIN/goperf

# $BIN is a directory of your choosing in your $PATH
```

## Usage

```
$ ./goperf-v0.0.1 -help
Usage of ./goperf-v0.0.1:
  -n=0: Total number of connections.
  -r=0: Connection rate (per second).
  -u="": Target URL.
  -v=false: Print verbose messaging.
  -version=false: Show version infomration.
```

## [API Documentation](http://godoc.org/github.com/jmervine/goperf)

```go
import "github.com/jmervine/goperf"
```
Package perf is a simple HTTPerf clone for performance testing web applications
written in Go.

This is designed to be run as a command line too, however, can be hooked in to
as an API as well.

CLI Usage:

    $ ./goperf-v0.0.1 -help
    Usage of ./goperf-v0.0.1:
      -n=0: Total number of connections.
      -r=0: Connection rate (per second).
      -u="": Target URL.
      -v=false: Print verbose messaging.
      -version=false: Show version infomration.

##### Example:
	// Start()
	config := &Configurator{
	    Path:     "http://localhost",
	    NumConns: 100,
	    Rate:     10,
	    Verbose:  true,
	}

	results := Start(config)
	Display(results)

	// QuickRun()
	quick := QuickRun("http://localhost", 100, 10)
	Display(quick)

### Variables

```go
var Testing = false
```

> Testing is a flag for disabling certain messaging during test.

```go
var Version = "v0.0.4"
```

> Version is package version.


### Types

#### Configurator

```go
type Configurator struct {
    Rate     float64
    NumConns int
    Path     string
    Verbose  bool
}
```




#### Connect

```go
func Connect(path string, verbose bool) *results.Result
```
> Connect makes a singled connection, returning a simplified result struct.

##### Example:
	go stubServer()

	results := Connect("http://localhost:9876", false)
	fmt.Printf("Status Code: %v\n", results.Code)

	// Output:
	// Status Code: 200

#### Display

```go
func Display(r *results.Results)
```
> Display formatted results.


#### Parallel

```go
func Parallel(config *Configurator) *results.Results
```
> Parallel forces a parallel run using a Configurator.

##### Example:
	config := &Configurator{
	    Path:     "http://localhost",
	    NumConns: 100,
	    Rate:     10,
	    Verbose:  true,
	}

	results := Parallel(config)
	Display(results)

#### QuickRun

```go
func QuickRun(path string, numconns int, rate float64) *results.Results
```
> QuickRun limited options.


#### Series

```go
func Series(config *Configurator) *results.Results
```
> Series forces a run using a Configurator, running request in series.

##### Example:
	config := &Configurator{
	    Path:     "http://localhost",
	    NumConns: 100,
	    Rate:     10,
	    Verbose:  true,
	}

	results := Parallel(config)
	Display(results)

#### Siege

```go
func Siege(path string, numconns int) *results.Results
```
> Siege forces Parallel run, with limited options.

##### Example:
	results := Siege("http://localhost", 100)
	Display(results)

#### Start

```go
func Start(config *Configurator) *results.Results
```
> Start a new run using a Configurator



