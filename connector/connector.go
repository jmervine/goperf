package connector

import (
    "fmt"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strconv"
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
    Rate     int
    Verbose  bool
    Results  *results.Results
}

// New generates a new Connector with all the necessaries.
func (connector Connector) New(path string, numconns int) Connector {
    //connector := Connector{}
    connector.Path = path
    connector.NumConns = numconns
    connector.waiter = &sync.WaitGroup{}
    connector.tranny = make(chan results.Result)

    connector.Results = &results.Results{
        Took: make([]float64, numconns),
        Code: make([]int, numconns),
    }

    return connector
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
            time.Sleep(time.Second / time.Duration(conn.Rate))
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

// Connect makes a single connection.
func (conn *Connector) Connect() results.Result {
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
            fmt.Printf(" > Responded in %4v ms, with code: %d\n", trimFloat(took, 2), code)
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
    conn.Results.TotalTime = trimFloat(float64(time.Since(start))/float64(time.Second), 3)
    conn.Results.ConnPerSec = trimFloat(float64(conn.NumConns)/conn.Results.TotalTime, 3)

    // Finalize results.
    conn.Results.Finalize()
}

func trimFloat(float float64, points int) float64 {
    ff, err := strconv.ParseFloat(strconv.FormatFloat(float, byte('f'), points, 64), 64)
    if err != nil {
        panic(err)
    }
    return ff
}
