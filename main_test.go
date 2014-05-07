package perf

import (
    "fmt"
    . "github.com/jmervine/check"
    "io/ioutil"
    "strings"
)

/***
 * Tests
 *
 * See setup_test.go for Test initialization
 * and helper method definitions.
 ******************************/
func (suite *MainSuite) TestVersion(test *C) {
    content, err := ioutil.ReadFile("VERSION")
    if err != nil {
        panic(err)
    }

    test.Assert(Version, Equals, strings.TrimSpace(string(content)))
}

func (suite *MainSuite) TestQuickRun(test *C) {
    go stubServer()

    rs := QuickRun("http://localhost:9876", 5, 5)

    test.Assert(rs.Took, HasLen, 5)
    test.Assert(rs.Code, HasLen, 5)
    test.Assert(rs.Errors, HasLen, 0)
}

func (suite *MainSuite) TestSiege(test *C) {
    go stubServer()

    rs := Siege("http://localhost:9876", 5)
    test.Assert(rs.Took, HasLen, 5)
    test.Assert(rs.Code, HasLen, 5)
    test.Assert(rs.Errors, HasLen, 0)
}

func (suite *MainSuite) TestStart(test *C) {
    go stubServer()

    rs := Start(newConf())

    test.Assert(rs.Took, HasLen, 5)
    test.Assert(rs.Code, HasLen, 5)
    test.Assert(rs.Errors, HasLen, 0)
}

func (suite *MainSuite) TestParallel(test *C) {
    go stubServer()

    rs := Parallel(newConf())

    test.Assert(rs.Took, HasLen, 5)
    test.Assert(rs.Code, HasLen, 5)
    test.Assert(rs.Errors, HasLen, 0)
}

func (suite *MainSuite) TestSeries(test *C) {
    go stubServer()

    rs := Series(newConf())

    test.Assert(rs.Took, HasLen, 5)
    test.Assert(rs.Code, HasLen, 5)
    test.Assert(rs.Errors, HasLen, 0)
}

func (suite *MainSuite) TestConnect(test *C) {
    go stubServer()

    r := Connect("http://localhost:9876", false)

    test.Assert(r.Code, Equals, 200)
}

func (suite *MainSuite) TestHeavyLoad(test *C) {
    go stubServer()

    rs := Siege("http://localhost:9876", 10000)
    test.Assert(rs.Replies, Equals, 10000)

    qr := QuickRun("http://localhost:9876", 50000, 5000)
    test.Assert(qr.Replies, Equals, 50000)
}

/***
 * Examples
 ******************************/
func ExampleDisplay() {
    results := QuickRun("http://localhost", 100, 10)
    Display(results)
}

func ExampleQuickRun() {
    results := QuickRun("http://localhost", 100, 10)
    Display(results)
}

func ExampleSiege() {
    results := Siege("http://localhost", 100)
    Display(results)
}

func ExampleStart() {
    config := &Configurator{
        Path:     "http://localhost",
        NumConns: 100,
        Rate:     10,
        Verbose:  true,
    }

    results := Start(config)
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
