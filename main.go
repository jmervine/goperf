package perf

import (
    "fmt"
)

var Version = "v0.0.2"
var Testing bool = false

type Configurator struct {
    Rate     int
    NumConns int
    Path     string
    Verbose  bool
}

// Quickly Run with limited options.
func QuickRun(path string, numconns, rate int) *ResultSet {
    config := &Configurator{
        Path:     path,
        NumConns: numconns,
        Rate:     rate,
    }

    return Start(config)
}

// Force Parallel run, with limited options.
func Siege(path string, numconns int) *ResultSet {
    config := &Configurator{
        Path:     path,
        NumConns: numconns,
    }

    return Parallel(config)
}

// Setup a new run using a Configurator
func Start(config *Configurator) *ResultSet {
    conn := setup(config)
    conn.Run()
    return conn.Results
}

// Force Parallel run using a Configurator.
func Parallel(config *Configurator) *ResultSet {
    conn := setup(config)
    conn.Parallel()
    return conn.Results
}

// Force Series run using a Configurator.
func Series(config *Configurator) *ResultSet {
    conn := setup(config)
    conn.Series()
    return conn.Results
}

// A singled connection, returning a simplified result struct.
func Connect(path string, verbose bool) *ResultTransport {
    conn := New(path, 0)
    conn.Verbose = verbose

    result := conn.Connect()
    return &result
}

// Display formatted results.
func Display(r *ResultSet) {
    fmt.Printf("Total: requested %d replies %d test-duration %vs\n",
        r.Requested, len(r.Took), r.TotalTime)
    fmt.Println()

    fmt.Printf("Connection rate: %v conn/s\n", r.ConnPerSec)
    fmt.Printf("Connection time [ms]: min %v avg %v max %v med %v\n",
        r.TookMin, r.TookAvg, r.TookMax, r.TookMed)
    fmt.Printf("Connection time [ms]: 85th %v 90th %v 95th %v 99th %v\n",
        r.Took85th, r.Took90th, r.Took95th, r.Took99th)
    fmt.Println()

    fmt.Printf("Reply size [B]: content %v header/footer %v (total %v)\n",
        r.ContentLength, r.HeaderLength, r.TotalLength)
    fmt.Printf("Reply status: 1xx=%v 2xx=%v 3xx=%v 4xx=%v 5xx=%v\n",
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
func setup(config *Configurator) *Connector {
    validate(config)
    header(config)
    conn := New(config.Path, config.NumConns)
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
