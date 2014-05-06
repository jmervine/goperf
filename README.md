# goperf

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
  -n=0: num-conns
  -p="": path
  -r=0: rate
  -v=false: verbose
```

## API Usage

```
import "github.com/jmervine/goperf"
```

### Documentation

```
PACKAGE DOCUMENTATION

package perf
    import "."



FUNCTIONS


func Display(r *ResultSet)


TYPES

type Configurator struct {
    Rate     int
    NumConns int
    Path     string
    Verbose  bool
}



type Connector struct {
    Path     string
    NumConns int
    Rate     int
    Verbose  bool
    Results  *ResultSet
    // contains filtered or unexported fields
}

    Example:
    go stubServer()
    
    c := New("http://localhost:9876", 10)
    
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
    
    // Output:
    // Code[0] = 200
    // Code[1] = 200
    // Code[2] = 200
    // Code[3] = 200
    // Code[4] = 200
    // Code[5] = 200
    // Code[6] = 200
    // Code[7] = 200
    // Code[8] = 200
    // Code[9] = 200


func New(path string, numconns int) Connector



func (this *Connector) Connect() ResultTransport


func (this *Connector) Parallel()


func (this *Connector) Run()


func (this *Connector) Series()


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
    Performance test results.

    Example:
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


func Parallel(config *Configurator) *ResultSet
    Force Parallel run using a Configurator.



func QuickRun(path string, numconns, rate int) *ResultSet
    Quickly Run with limited options.



func Series(config *Configurator) *ResultSet
    Force Series run using a Configurator.



func Siege(path string, numconns int) *ResultSet
    Force Parallel run, with limited options.



func Start(config *Configurator) *ResultSet
    Setup a new run using a Configurator



func (this *ResultSet) Add(result ResultTransport)
    Add transport data to result set.


func (this *ResultSet) CalculatePct(pct int) float64
    Calculate Percentile from existing Took values.


func (this *ResultSet) Finalize()
    FFinalize results, generating min, max, avg med and percentiles.


type ResultTransport struct {
    Index, Code   int
    Took          float64
    Error         error
    TotalLength   int64
    ContentLength int64
    HeaderLength  int64
}
    Performance test result transporter.



func Connect(path string, verbose bool) *ResultTransport
    A singled connection, returning a simplified result struct.




SUBDIRECTORIES

    bin
    pkg

```


