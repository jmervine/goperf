package perf

import (
    "fmt"
    . "github.com/jmervine/GoT"
    "io/ioutil"
    "net/http"
    "strings"
    "testing"
    "time"
)

var StubServerRunning bool = false

func init() {
    Testing = true
}

func TestVersion(T *testing.T) {
    content, err := ioutil.ReadFile("VERSION")
    if err != nil {
        panic(err)
    }

    Go(T).AssertEqual(Version, strings.TrimSpace(string(content)))
}

func TestQuickRun(T *testing.T) {
    go stubServer()

    rs := QuickRun("http://localhost:9876", 5, 5)

    Go(T).AssertLength(rs.Took, 5)
    Go(T).AssertLength(rs.Code, 5)
    Go(T).AssertLength(rs.Errors, 0)
}

func TestSiege(T *testing.T) {
    go stubServer()

    rs := Siege("http://localhost:9876", 5)
    Go(T).AssertLength(rs.Took, 5)
    Go(T).AssertLength(rs.Code, 5)
    Go(T).AssertLength(rs.Errors, 0)
}

func TestStart(T *testing.T) {
    go stubServer()

    rs := Start(newConf())

    Go(T).AssertLength(rs.Took, 5)
    Go(T).AssertLength(rs.Code, 5)
    Go(T).AssertLength(rs.Errors, 0)
}

func TestParallel(T *testing.T) {
    go stubServer()

    rs := Parallel(newConf())

    Go(T).AssertLength(rs.Took, 5)
    Go(T).AssertLength(rs.Code, 5)
    Go(T).AssertLength(rs.Errors, 0)
}

func TestSeries(T *testing.T) {
    go stubServer()

    rs := Series(newConf())

    Go(T).AssertLength(rs.Took, 5)
    Go(T).AssertLength(rs.Code, 5)
    Go(T).AssertLength(rs.Errors, 0)
}

func TestConnect(T *testing.T) {
    go stubServer()

    r := Connect("http://localhost:9876", false)

    Go(T).AssertEqual(r.Code, 200)
}

/***
 * Examples
 ******************************/

func Example() {
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
}

func ExampleSiege() {
    results := Siege("http://localhost", 100)
    Display(results)
}

func ExampleParallel() {
    config := &Configurator{
        Path:     "http://localhost",
        NumConns: 100,
        Rate:     10,
        Verbose:  true,
    }

    results := Parallel(config)
    Display(results)
}

func ExampleSeries() {
    config := &Configurator{
        Path:     "http://localhost",
        NumConns: 100,
        Rate:     10,
        Verbose:  true,
    }

    results := Parallel(config)
    Display(results)
}

func ExampleConnect() {
    go stubServer()

    results := Connect("http://localhost:9876", false)
    fmt.Printf("Status Code: %v\n", results.Code)

    // Output:
    // Status Code: 200
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

func newConf() *Configurator {
    return &Configurator{
        Path:     "http://localhost:9876",
        NumConns: 5,
        Rate:     5,
    }
}
