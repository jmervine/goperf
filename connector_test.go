package perf

import (
    "fmt"
    . "github.com/jmervine/check"
)

/***
 * Tests
 *
 * See setup_test.go for Test initialization
 * and helper method definitions.
 ******************************/
func (suite *ConnectorSuite) TestNew(test *C) {
    c := New("http://localhost:9876", 10)

    test.Assert(c.Path, Equals, "http://localhost:9876")
    test.Assert(c.NumConns, Equals, 10)
    test.Assert(c.Verbose, Equals, false)
    test.Assert(c.Rate, Equals, 0)
    test.Assert(c.Results.Took, HasLen, 10)
}

func (suite *ConnectorSuite) TestSeries(test *C) {
    go stubServer()

    c := New("http://localhost:9876", 10)
    c.Series()

    for i := 0; i < 10; i++ {
        test.Assert(c.Results.Took[i], Not(Equals), 0)
        test.Assert(c.Results.Code[i], Equals, 200)
    }

}

func (suite *ConnectorSuite) TestParallel(test *C) {
    go stubServer()

    c := New("http://localhost:9876", 10)
    c.Parallel()

    for i := 0; i < 10; i++ {
        test.Assert(c.Results.Took[i], Not(Equals), 0)
        test.Assert(c.Results.Code[i], Equals, 200)
    }
    test.Assert(c.Results.TookMed, Not(Equals), 0)

}

func (suite *ConnectorSuite) TestRun(test *C) {
    go stubServer()

    c := New("http://localhost:9876", 10)
    c.Run()

    for i := 0; i < 10; i++ {
        test.Assert(c.Results.Took[i], Not(Equals), 0)
        test.Assert(c.Results.Code[i], Equals, 200)
    }
    test.Assert(c.Results.TookMed, Not(Equals), 0)

    c.Rate = -1
    c.Run()

    for i := 0; i < 10; i++ {
        test.Assert(c.Results.Took[i], Not(Equals), 0)
        test.Assert(c.Results.Code[i], Equals, 200)
    }
    test.Assert(c.Results.TookMed, Not(Equals), 0)

    c.Rate = 5
    c.Run()

    for i := 0; i < 10; i++ {
        test.Assert(c.Results.Took[i], Not(Equals), 0)
        test.Assert(c.Results.Code[i], Equals, 200)
    }
}

func (suite *ConnectorSuite) TestConnect(test *C) {
    go stubServer()

    c := New("http://localhost:9876", 10)
    r := c.Connect()

    for i := 0; i < 10; i++ {
        test.Assert(r.Took, Not(Equals), 0)
        test.Assert(r.Code, Equals, 200)
    }
}

/***
 * Examples
 ******************************/
func ExampleConnector() {
    go stubServer()

    c := New("http://localhost:9876", 10)

    /**
     * Note on Rate:
     *
     * If Rate is not zero, Run() will parallelize actions at a Rate (QPS)
     * of the set value.
     *
     * If Rate is zero, Run() will run the connections in a series.
     *
     * Both c.Series() and c.Parallel() can also be called in place of Run().
     *******************************/
    c.Rate = 4 // QPS
    c.Run()

    for i, code := range c.Results.Code {
        fmt.Printf("Code[%d] = %d\n", i, code)
    }

    // Output:
    // Code[0] = 200
    // Code[1] = 200
    // Code[2] = 200
    // Code[3] = 200
    // Code[4] = 200
    // Code[5] = 200
    // Code[6] = 200
    // Code[7] = 200
    // Code[8] = 200
    // Code[9] = 200
}
