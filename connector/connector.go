package connector

import (
    "fmt"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    "time"
    "github.com/jmervine/goperf/results"
)

// Connector contains connector data.
type Connector struct {
    waiter *sync.WaitGroup
    tranny chan results.Result

    Path     string
    NumConns int
    Rate     float64
    Verbose  bool
    Results  *results.Results
}

// New generates a new Connector with all the necessaries.
func (conn Connector) New(path string, numconns int) Connector {
    uri, err := url.Parse(path)

    if uri.Scheme == "" {
        uri.Scheme = "http"
    }

    if err != nil {
        panic(err)
    }

    conn.Path = uri.String()
    conn.NumConns = numconns
    conn.waiter = &sync.WaitGroup{}
    conn.tranny = make(chan results.Result)

    conn.Results = &results.Results{
        Took: make([]float64, numconns),
        Code: make([]int, numconns),

        // set to -1 so that it gets the first connection time
        ConnectTime: -1,
    }

    return conn
}

// Run runs the Connector, selecting Parallel or Series based on Rate.
func (conn *Connector) Run() {
    if conn.Rate != 0 {
        conn.Parallel()
    } else {
        conn.Series()
    }
}

// Series runs the Connector serialized.
func (conn *Connector) Series() {
    start := time.Now()

    defer conn.finalize(start)

    for i := 0; i < conn.NumConns; i++ {
        result := conn.Connect()
        result.Index = i
        conn.Results.Add(result)
    }
}

// Parallel runs the Connector parallelized.
func (conn *Connector) Parallel() {
    start := time.Now()

    defer conn.finalize(start)

    for i := 0; i < conn.NumConns; i++ {

        if conn.Rate > 0 && i != 0 {
            time.Sleep(time.Duration((1 / conn.Rate) * float64(time.Second)))
        }

        conn.waiter.Add(1)
        go func(i int) {
            result := conn.Connect()
            result.Index = i
            conn.tranny <- result
            conn.waiter.Done()
        }(i)

        conn.waiter.Add(1)
        go conn.collect()
    }

    conn.waiter.Wait()
}

func (conn *Connector) customDial(network, addr string) (net.Conn, error) {
    start := time.Now()
    c, err := net.Dial(network, addr)

    if conn.Results.ConnectTime == -1 {
        conn.Results.ConnectTime = float64(time.Since(start) / time.Millisecond)
    }

    return c, err
}

// Connect makes a single connection.
func (conn *Connector) Connect() results.Result {

    transport := http.Transport{
        Dial: conn.customDial,
    }

    http.DefaultClient = &http.Client{
        Transport: &transport,
    }

    start := time.Now()
    resp, err := http.Get(conn.Path)
    took := float64(time.Since(start) / time.Millisecond)

    var code int
    var tlen, clen, hlen int64

    if err == nil {
        code = resp.StatusCode
        clen = resp.ContentLength

        if dump, e := httputil.DumpResponse(resp, true); e == nil {
            tlen = int64(len(dump))
            hlen = tlen - clen
        }
    }

    if conn.Verbose {
        if err != nil {
            fmt.Printf(" > Responded with error: %q\n",
                err.(*url.Error).Err.(*net.OpError).Error())
        } else {
            fmt.Printf(" > Responded in %6.2f ms, with code: %d\n", took, code)
        }
    }

    return results.Result{Took: took,
        Code:          code,
        Error:         err,
        TotalLength:   tlen,
        ContentLength: clen,
        HeaderLength:  hlen,
    }
}

/****
 * Private methods
 *****************************************************/

func (conn *Connector) collect() {
    tranny := <-conn.tranny
    conn.Results.Add(tranny)
    conn.waiter.Done()
}

func (conn *Connector) finalize(start time.Time) {
    if conn.Verbose {
        fmt.Println(" > finalizing...\n")
    }

    // Some results data can only be populated if run via Connector.
    conn.Results.Requested = conn.NumConns
    conn.Results.TotalTime = float64(time.Since(start))/float64(time.Second)
    conn.Results.ConnPerSec = float64(conn.NumConns)/conn.Results.TotalTime

    // Finalize results.
    conn.Results.Finalize()
}

