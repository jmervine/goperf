package perf

import (
    "fmt"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strconv"
    "sync"
    "time"
)

type Connector struct {
    waiter *sync.WaitGroup
    tranny chan ResultTransport

    Path     string
    NumConns int
    Rate     int
    Verbose  bool
    Results  *ResultSet
}

func New(path string, numconns int) Connector {
    connector := Connector{}
    connector.Path = path
    connector.NumConns = numconns
    connector.waiter = &sync.WaitGroup{}
    connector.tranny = make(chan ResultTransport)

    connector.Results = &ResultSet{
        Took: make([]float64, numconns),
        Code: make([]int, numconns),
    }

    return connector
}

func (this *Connector) Run() {
    if this.Rate != 0 {
        this.Parallel()
    } else {
        this.Series()
    }
}

func (this *Connector) Series() {
    start := time.Now()

    defer this.finalize(start)

    for i := 0; i < this.NumConns; i++ {
        result := this.Connect()
        result.Index = i
        this.Results.Add(result)
    }
}

func (this *Connector) Parallel() {
    start := time.Now()

    defer this.finalize(start)

    for i := 0; i < this.NumConns; i++ {

        if this.Rate > 0 && i != 0 {
            time.Sleep(time.Second / time.Duration(this.Rate))
        }

        this.waiter.Add(1)
        go func(i int) {
            result := this.Connect()
            result.Index = i
            this.tranny <- result
            this.waiter.Done()
        }(i)

        this.waiter.Add(1)
        go this.collect()
    }

    this.waiter.Wait()
}

func (this *Connector) Connect() ResultTransport {
    start := time.Now()
    resp, err := http.Get(this.Path)
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

    if this.Verbose {
        if err != nil {
            fmt.Printf(" > Responded with error: %q\n",
                err.(*url.Error).Err.(*net.OpError).Error())
        } else {
            fmt.Printf(" > Responded in %4v ms, with code: %d\n", trimFloat(took, 2), code)
        }
    }

    return ResultTransport{Took: took,
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

func (this *Connector) collect() {
    tranny := <-this.tranny
    this.Results.Add(tranny)
    this.waiter.Done()
}

func (this *Connector) finalize(start time.Time) {
    if this.Verbose {
        fmt.Println(" > finalizing...\n")
    }

    // Some results data can only be populated if run via Connector.
    this.Results.Requested = this.NumConns
    this.Results.TotalTime = trimFloat(float64(time.Since(start))/float64(time.Second), 3)
    this.Results.ConnPerSec = trimFloat(float64(this.NumConns)/this.Results.TotalTime, 3)

    // Finalize results.
    this.Results.Finalize()
}

func trimFloat(float float64, points int) float64 {
    ff, err := strconv.ParseFloat(strconv.FormatFloat(float, byte('f'), points, 64), 64)
    if err != nil {
        panic(err)
    }
    return ff
}
