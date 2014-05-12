package connector

import (
    "fmt"
    "testing"
    "time"
    "net/http"
    . "github.com/jmervine/GoT"
)

var StubServerRunning bool = false

func TestNew(T *testing.T) {
    c := Connector{}.New("http://localhost:9876", 10)

    Go(T).AssertEqual(c.Path, "http://localhost:9876")
    Go(T).AssertEqual(c.NumConns, 10)
    Go(T).AssertEqual(c.Verbose, false)
    Go(T).AssertEqual(c.Rate, 0)
    Go(T).AssertLength(c.Results.Took, 10)
}

func TestSeries(T *testing.T) {
    go stubServer()

    c := Connector{}.New("http://localhost:9876", 10)
    c.Series()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(c.Results.Took[i], 0)
        Go(T).AssertEqual(c.Results.Code[i], 200)
    }

}

func TestParallel(T *testing.T) {
    go stubServer()

    c := Connector{}.New("http://localhost:9876", 10)
    c.Parallel()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(c.Results.Took[i], 0)
        Go(T).AssertEqual(c.Results.Code[i], 200)
    }
    Go(T).RefuteEqual(c.Results.TookMed, 0)

}

func TestRun(T *testing.T) {
    go stubServer()

    c := Connector{}.New("http://localhost:9876", 10)
    c.Run()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(c.Results.Took[i], 0)
        Go(T).AssertEqual(c.Results.Code[i], 200)
    }
    Go(T).RefuteEqual(c.Results.TookMed, 0)

    c.Rate = -1
    c.Run()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(c.Results.Took[i], 0)
        Go(T).AssertEqual(c.Results.Code[i], 200)
    }
    Go(T).RefuteEqual(c.Results.TookMed, 0)

    c.Rate = 5
    c.Run()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(c.Results.Took[i], 0)
        Go(T).AssertEqual(c.Results.Code[i], 200)
    }
}

func TestConnect(T *testing.T) {
    go stubServer()

    c := Connector{}.New("http://localhost:9876", 10)
    r := c.Connect()

    for i := 0; i < 10; i++ {
        Go(T).RefuteEqual(r.Took, 0)
        Go(T).AssertEqual(r.Code, 200)
    }
}

/***
 * Examples
 ******************************/

func ExampleConnector_New() {
    go stubServer()

    c := Connector{}.New("http://localhost:9876", 10)

    //
    // Note on Rate:
    //
    // If Rate is not zero, Run() will parallelize actions at a Rate (QPS)
    // of the set value.
    //
    // If Rate is zero, Run() will run the connections in a series.
    //
    // Both c.Series() and c.Parallel() can also be called in place of Run().
    //
    c.Rate = 4 // QPS
    c.Run()

    for i, code := range c.Results.Code {
        fmt.Printf("Code[%d] = %d\n", i, code)
    }
}

/***
 * Helpers
 ******************************/
func stubServer() {
    if StubServerRunning {
        return
    }

    StubServerRunning = true
    defer func() { StubServerRunning = false }()

    // Starting a stub server on :9876 to handle incoming requests
    // for example.
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(5 * time.Millisecond)
        fmt.Fprintln(w, "hello web")
    })
    http.ListenAndServe(":9876", nil)
}

