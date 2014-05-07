# goperf
---

Simple HTTPerf clone for performance testing web applications written in Go.

> NOTE: This is the inital commit and shouldn't be considered ready for anyone. That said, it should
> work as outlined below, at least on Linux based systems.

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

## API Usage

```go
import "github.com/jmervine/goperf"
```
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
var Testing bool = false
```

> Flag for disabling certain messaging during test.

```go
var Version = "v0.0.2"
```

> Package version.


### Types

#### Configurator

```go
type Configurator struct {
    Rate     int
    NumConns int
    Path     string
    Verbose  bool
}
```

#### Connector

```go
type Connector struct {
    Path     string
    NumConns int
    Rate     int
    Verbose  bool
    Results  *ResultSet
    // contains filtered or unexported fields
}
```

#### ResultSet

```go
type ResultSet struct {
    Requested  int
    Replies    int
    TotalTime  float64
    ConnPerSec float64

    Took     []float64
    TookMin  float64
    TookMed  float64
    TookAvg  float64
    TookMax  float64
    Took85th float64
    Took90th float64
    Took95th float64
    Took99th float64

    Code    []int
    Code1xx int
    Code2xx int
    Code3xx int
    Code4xx int
    Code5xx int

    Errors            []error
    ErrorsTotal       int
    ErrorsConnTimeout int
    ErrorsConnRefused int
    ErrorsConnReset   int
    ErrorsFdUnavail   int
    ErrorsAddrUnavail int
    ErrorsOther       int

    ContentLength int64
    HeaderLength  int64
    TotalLength   int64
}
```

#### ResultTransport

```go
type ResultTransport struct {
    Index, Code   int
    Took          float64
    Error         error
    TotalLength   int64
    ContentLength int64
    HeaderLength  int64
}
```

### Functions

#### Connect

```go
func (this *Connector) Connect() ResultTransport
```
> A single connection.

#### New

```go
func (connector Connector) New(path string, numconns int) Connector
```
> Generate a new Connector with all the necessaries.

##### Example:
	go stubServer()

	c := Connector{}.New("http://localhost:9876", 10)

	/**
	 * Note on Rate:
	 *
	 * If Rate is not zero, Run() will parallelize actions at a Rate (QPS)
	 * of the set value.
	 *
	 * If Rate is zero, Run() will run the connections in a series.
	 *
	 * Both c.Series() and c.Parallel() can also be called in place of Run().
	 *******************************/
	c.Rate = 4 // QPS
	c.Run()

	for i, code := range c.Results.Code {
	    fmt.Printf("Code[%d] = %d\n", i, code)
	}

#### Parallel

```go
func (this *Connector) Parallel()
```
> Run Connector parallelized based on Rate.

#### Run

```go
func (this *Connector) Run()
```
> Run Connector, selecting Parallel or Series based on Rate.

#### Series

```go
func (this *Connector) Series()
```
> Run Connector serialized.

##### Example:
	// This should typically not be created manaully, but rather by
	// Connector{}.New( ... )

	r := ResultSet{
	    Took: make([]float64, 10),
	    Code: make([]int, 10),
	}

	t := ResultTransport{
	    Index: 0,
	    Took:  300.0,
	    Code:  200,
	}

	r.Add(t)

	// If running this manually, as opposed to via (*Connector{}).Run()
	// you must finalize to create max, min, avg, etc.
	r.Finalize()

	fmt.Printf("Took: %f\n", r.Took[0])
	fmt.Printf("Code: %d\n", r.Code[0])
	fmt.Printf("Max:  %f\n", r.TookMax)

	// Output:
	// Took: 300.000000
	// Code: 200
	// Max:  300.000000

#### Parallel

```go
func Parallel(config *Configurator) *ResultSet
```
> Force Parallel run using a Configurator.

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
func QuickRun(path string, numconns, rate int) *ResultSet
```
> Quickly Run with limited options.

#### Series

```go
func Series(config *Configurator) *ResultSet
```
> Force Series run using a Configurator.

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
func Siege(path string, numconns int) *ResultSet
```
> Force Parallel run, with limited options.

##### Example:
	results := Siege("http://localhost", 100)
	Display(results)

#### Start

```go
func Start(config *Configurator) *ResultSet
```
> Setup a new run using a Configurator

#### Add

```go
func (this *ResultSet) Add(result ResultTransport)
```
> Add transport data to result set.

#### CalculatePct

```go
func (this *ResultSet) CalculatePct(pct int) float64
```
> Calculate Percentile from existing Took values.

#### Finalize

```go
func (this *ResultSet) Finalize()
```
> Finalize results, generating min, max, avg med and percentiles.

#### Connect

```go
func Connect(path string, verbose bool) *ResultTransport
```
> A singled connection, returning a simplified result struct.

##### Example:
	go stubServer()

	results := Connect("http://localhost:9876", false)
	fmt.Printf("Status Code: %v\n", results.Code)

	// Output:
	// Status Code: 200

#### Display

```go
func Display(r *ResultSet)
```
> Display formatted results.

