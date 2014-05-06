package perf

import (
    "fmt"
)

type Configurator struct {
    Rate     int
    NumConns int
    Path     string
    Verbose  bool
}

// Setup a new run using a Configurator
func Start(config *Configurator) *ResultSet {
    conn := setup(config)
    conn.Run()
    return conn.Results
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

func Display(r *ResultSet) {
    /**
     * httperf --client=0/1 --server=www.example.com --port=80 --uri=/ --send-buffer=4096 --recv-buffer=16384 --num-conns=10 --num-calls=1
     * httperf: warning: open file limit > FD_SETSIZE; limiting max. # of open files to FD_SETSIZE
     * Maximum connect burst length: 1
     *
     * Total: connections 10 requests 10 replies 10 test-duration 2.019 s
     *
     * Connection rate: 5.0 conn/s (201.9 ms/conn, <=1 concurrent connections)
     * Connection time [ms]: min 174.8 avg 201.9 max 386.2 median 182.5 stddev 64.8
     * Connection time [ms]: connect 89.3
     * Connection length [replies/conn]: 1.000
     *
     * Request rate: 5.0 req/s (201.9 ms/req)
     * Request size [B]: 68.0
     *
     * Reply rate [replies/s]: min 0.0 avg 0.0 max 0.0 stddev 0.0 (0 samples)
     * Reply time [ms]: response 112.5 transfer 0.1
     * Reply size [B]: header 321.0 content 1270.0 footer 0.0 (total 1591.0)
     * Reply status: 1xx=0 2xx=10 3xx=0 4xx=0 5xx=0
     *
     * CPU time [s]: user 0.54 system 1.48 (user 26.7% system 73.1% total 99.9%)
     * Net I/O: 8.0 KB/s (0.1*10^6 bps)
     *
     * Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
     * Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
     *****************************************************/

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
    fmt.Printf("Running: Path=%s NumConns=%d Rate=%v Verbose=%v\n\n",
        config.Path, config.NumConns, config.Rate, config.Verbose)
}
