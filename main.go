/*
Package perf is a simple HTTPerf clone for performance testing web applications written in Go.

This is designed to be run as a command line too, however, can be hooked in to as an API as well.

CLI Usage:

    $ ./goperf-v0.0.1 -help
    Usage of ./goperf-v0.0.1:
      -n=0: Total number of connections.
      -r=0: Connection rate (per second).
      -u="": Target URL.
      -v=false: Print verbose messaging.
      -version=false: Show version infomration.

*/
package perf

import (
    "fmt"
    "github.com/jmervine/goperf/connector"
    "github.com/jmervine/goperf/results"
)

// Version is package version.
var Version = "v0.0.3"

// Testing is a flag for disabling certain messaging during test.
var Testing = false

// Configurator is a basic data struct for configuring runs.
type Configurator struct {
    Rate     int
    NumConns int
    Path     string
    Verbose  bool
}

// QuickRun limited options.
func QuickRun(path string, numconns, rate int) *results.Results {
    config := &Configurator{
        Path:     path,
        NumConns: numconns,
        Rate:     rate,
    }

    return Start(config)
}

// Siege forces Parallel run, with limited options.
func Siege(path string, numconns int) *results.Results {
    config := &Configurator{
        Path:     path,
        NumConns: numconns,
    }

    return Parallel(config)
}

// Start a new run using a Configurator
func Start(config *Configurator) *results.Results {
    conn := setup(config)
    conn.Run()
    return conn.Results
}

// Parallel forces a parallel run using a Configurator.
func Parallel(config *Configurator) *results.Results {
    conn := setup(config)
    conn.Parallel()
    return conn.Results
}

// Series forces a run using a Configurator, running request in series.
func Series(config *Configurator) *results.Results {
    conn := setup(config)
    conn.Series()
    return conn.Results
}

// Connect makes a singled connection, returning a simplified result struct.
func Connect(path string, verbose bool) *results.Result {
    conn := connector.Connector{}.New(path, 0)
    conn.Verbose = verbose

    result := conn.Connect()
    return &result
}

// Display formatted results.
func Display(r *results.Results) {
    fmt.Printf("Total: requested %d replies %d test-duration %6.2fs\n",
        r.Requested, len(r.Took), r.TotalTime)
    fmt.Println()

    fmt.Printf("Connection rate: %6.2f conn/s\n", r.ConnPerSec)
    fmt.Printf("Connection time [ms]: min %6.2f avg %6.2f max %6.2f med %6.2f\n",
        r.TookMin, r.TookAvg, r.TookMax, r.TookMed)
    fmt.Printf("Connection time [ms]: 85th %6.2f 90th %6.2f 95th %6.2f 99th %6.2f\n",
        r.Took85th, r.Took90th, r.Took95th, r.Took99th)
    fmt.Printf("Connection time [ms]: connect %6.2f\n", r.ConnectTime)
    fmt.Println()

    fmt.Printf("Reply size [B]: content %v header/footer %v (total %v)\n",
        r.ContentLength, r.HeaderLength, r.TotalLength)
    fmt.Printf("Reply status: 1xx=%d 2xx=%d 3xx=%d 4xx=%d 5xx=%d\n",
        r.Code1xx, r.Code2xx, r.Code3xx, r.Code4xx, r.Code5xx)
    fmt.Println()

    fmt.Printf("Errors: total %d conn-timeout %d conn-refused %d conn-reset %d\n",
        r.ErrorsTotal, r.ErrorsConnTimeout, r.ErrorsConnRefused, r.ErrorsConnReset)
    fmt.Printf("Errors: fd-unavail %d addr-unavail %d other %d\n",
        r.ErrorsFdUnavail, r.ErrorsAddrUnavail, r.ErrorsOther)
    fmt.Println()
}

/****
 * Private methods
 *****************************************************/

// Setup Connector via Configurator
func setup(config *Configurator) *connector.Connector {
    validate(config)
    header(config)
    conn := connector.Connector{}.New(config.Path, config.NumConns)
    conn.Rate = config.Rate
    conn.Verbose = config.Verbose
    return &conn
}

func validate(config *Configurator) {
    if config.Path == "" {
        panic("Path is required.")
    }

    if config.NumConns == 0 {
        panic("NumConns is required and cannot be zero.")
    }
}

func header(config *Configurator) {
    // Hide header when testing.
    if !Testing {
        fmt.Printf("Running: Path=%s NumConns=%d Rate=%v Verbose=%v\n\n",
            config.Path, config.NumConns, config.Rate, config.Verbose)
    }
}
