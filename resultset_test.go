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
func (suite *ResultSetSuite) TestResultSet(test *C) {
    r := newRS(2)
    test.Assert(r.Took[0], Equals, 0.0)
    test.Assert(r.Code[0], Equals, 0)
}

func (suite *ResultSetSuite) TestResultTransport(test *C) {
    t := newRT(0, 300.0, 200)
    test.Assert(t.Took, Equals, 300.0)
    test.Assert(t.Code, Equals, 200)
}

func (suite *ResultSetSuite) TestAdd(test *C) {
    r := newRS(10)
    t := newRT(0, 300.0, 200)
    r.Add(t)
    test.Assert(r.Took[0], Equals, 300.0)
    test.Assert(r.Code[0], Equals, 200)

    t = newRT(9, 400.0, 500)
    test.Assert(t.Took, Equals, 400.0)
    test.Assert(t.Code, Equals, 500)

    r.Add(t)
    test.Assert(r.Took[9], Equals, 400.0)
    test.Assert(r.Code[9], Equals, 500)
}

func (suite *ResultSetSuite) TestMin(test *C) {
    r := populatedRS(5)

    r.min()
    test.Assert(r.TookMin, Equals, 100.0)
}

func (suite *ResultSetSuite) TestMax(test *C) {
    r := populatedRS(5)

    r.max()
    test.Assert(r.TookMax, Equals, 300.0)
}

func (suite *ResultSetSuite) TestAvg(test *C) {
    r := populatedRS(5)

    r.avg()
    test.Assert(r.TookAvg, Equals, 200.0)
}

func (suite *ResultSetSuite) TestMed(test *C) {
    r := populatedRS(5)

    r.med()
    test.Assert(r.TookMed, Equals, 200.0)
}

func (suite *ResultSetSuite) TestCalculatePct(test *C) {
    r := populatedRS(20)

    test.Assert(r.CalculatePct(80), Equals, 850.0)
    test.Assert(r.CalculatePct(75), Equals, 800.0)
}

func (suite *ResultSetSuite) TestPct(test *C) {
    r := populatedRS(20)

    r.pct()
    test.Assert(r.Took99th, Equals, 1050.0)
    test.Assert(r.Took95th, Equals, 1000.0)
    test.Assert(r.Took90th, Equals, 950.0)
    test.Assert(r.Took85th, Equals, 900.0)
}

func (suite *ResultSetSuite) TestFinalize(test *C) {
    r := populatedRS(5)

    r.Finalize()
    test.Assert(r.TookMin, Equals, 100.0)
    test.Assert(r.TookAvg, Equals, 200.0)
    test.Assert(r.TookMed, Equals, 200.0)
    test.Assert(r.TookMax, Equals, 300.0)
    test.Assert(r.Took99th, Equals, 300.0)
}

/***
 * Examples
 ******************************/
func ExampleResultSet() {
    // This should typically not be created manaully, but rather by
    // Connector{}.New( ... )

    r := ResultSet{
        Took: make([]float64, 10),
        Code: make([]int, 10),
    }

    t := ResultTransport{
        Index: 0,
        Took:  300.0,
        Code:  200,
    }

    r.Add(t)

    // If running this manually, as opposed to via (*Connector{}).Run()
    // you must finalize to create max, min, avg, etc.
    r.Finalize()

    fmt.Printf("Took: %f\n", r.Took[0])
    fmt.Printf("Code: %d\n", r.Code[0])
    fmt.Printf("Max:  %f\n", r.TookMax)

    // Output:
    // Took: 300.000000
    // Code: 200
    // Max:  300.000000
}
