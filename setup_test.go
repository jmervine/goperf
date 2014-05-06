package perf

import (
    "fmt"
    . "github.com/jmervine/check"
    "net/http"
    "testing"
    "time"
)

var StubServerRunning bool = false

func init() {
    Testing = true
}

/***
 * Setup
 ******************************/
func Test(t *testing.T) {
    TestingT(t)
}

/***
 * ResultSet
 ******************************/
type MainSuite struct{}

var _ = Suite(&MainSuite{})

/***
 * ResultSet
 ******************************/
type ResultSetSuite struct{}

var _ = Suite(&ResultSetSuite{})

/***
 * Connector
 ******************************/
type ConnectorSuite struct{}

var _ = Suite(&ConnectorSuite{})

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

func stubRS() ResultSet {
    r := ResultSet{}
    took := 200.0
    for i := 0; i < 10; i++ {
        r.Add(ResultTransport{
            Index: i,
            Took:  took,
            Code:  200,
        })
        took += 10
    }
    return r
}

func stubRT(index int) ResultTransport {
    return ResultTransport{
        Index: index,
        Took:  300,
        Code:  200,
    }
}

func newRS(l int) ResultSet {
    return ResultSet{
        Took: make([]float64, l),
        Code: make([]int, l),
    }
}

func populatedRS(l int) ResultSet {
    r := newRS(l)
    n, p := 100.0, 50.0
    for i := 0; i < l; i++ {
        t := newRT(i, n, 200)
        r.Add(t)
        n += p
    }
    return r
}

func newRT(i int, t float64, c int) ResultTransport {
    return ResultTransport{
        Index: i,
        Took:  t,
        Code:  c,
    }
}

func newConf() *Configurator {
    return &Configurator{
        Path:     "http://localhost:9876",
        NumConns: 5,
        Rate:     5,
    }
}
