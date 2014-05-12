package results

import (
    "fmt"
    "testing"
    . "github.com/jmervine/GoT"
)

//
// Tests
//
// See setup_Go(T).go for Test initialization
// and helper method definitions.
//

func TestResults(T *testing.T) {
    r := newRS(2)
    Go(T).AssertEqual(r.Took[0], 0.0, "")
    Go(T).AssertEqual(r.Code[0], 0, "")
}

func TestResult(T *testing.T) {
    t := newRT(0, 300.0, 200)
    Go(T).AssertEqual(t.Took, 300.0, "")
    Go(T).AssertEqual(t.Code, 200, "")
}

func TestAdd(T *testing.T) {
    r := newRS(10)
    t := newRT(0, 300.0, 200)
    r.Add(t)
    Go(T).AssertEqual(r.Took[0], 300.0, "")
    Go(T).AssertEqual(r.Code[0], 200, "")

    t = newRT(9, 400.0, 500)
    Go(T).AssertEqual(t.Took, 400.0, "")
    Go(T).AssertEqual(t.Code, 500, "")

    r.Add(t)
    Go(T).AssertEqual(r.Took[9], 400.0, "")
    Go(T).AssertEqual(r.Code[9], 500, "")
}

func TestMin(T *testing.T) {
    r := populatedRS(5)

    r.min()
    Go(T).AssertEqual(r.TookMin, 100.0, "")
}

func TestMax(T *testing.T) {
    r := populatedRS(5)

    r.max()
    Go(T).AssertEqual(r.TookMax, 300.0, "")
}

func TestAvg(T *testing.T) {
    r := populatedRS(5)

    r.avg()
    Go(T).AssertEqual(r.TookAvg, 200.0, "")
}

func TestMed(T *testing.T) {
    r := populatedRS(5)

    r.med()
    Go(T).AssertEqual(r.TookMed, 200.0, "")
}

func TestCalculatePct(T *testing.T) {
    r := populatedRS(20)

    Go(T).AssertEqual(r.CalculatePct(80), 850.0, "")
    Go(T).AssertEqual(r.CalculatePct(75), 800.0, "")
}

func TestPct(T *testing.T) {
    r := populatedRS(20)

    r.pct()
    Go(T).AssertEqual(r.Took99th, 1050.0, "")
    Go(T).AssertEqual(r.Took95th, 1000.0, "")
    Go(T).AssertEqual(r.Took90th, 950.0, "")
    Go(T).AssertEqual(r.Took85th, 900.0, "")
}

func TestFinalize(T *testing.T) {
    r := populatedRS(5)

    r.Finalize()
    Go(T).AssertEqual(r.TookMin, 100.0, "")
    Go(T).AssertEqual(r.TookAvg, 200.0, "")
    Go(T).AssertEqual(r.TookMed, 200.0, "")
    Go(T).AssertEqual(r.TookMax, 300.0, "")
    Go(T).AssertEqual(r.Took99th, 300.0, "")
}

/***
 * Examples
 ******************************/

func ExampleResults() {
    // This should typically not be created manaully, but rather by
    // Connector{}.New( ... )

    r := Results{
        Took: make([]float64, 10),
        Code: make([]int, 10),
    }

    t := Result{
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

/***
 * Helpers
 ******************************/

func newRS(l int) Results {
    return Results{
        Took: make([]float64, l),
        Code: make([]int, l),
    }
}

func newRT(i int, t float64, c int) Result {
    return Result{
        Index: i,
        Took:  t,
        Code:  c,
    }
}

func populatedRS(l int) Results {
    r := newRS(l)
    n, p := 100.0, 50.0
    for i := 0; i < l; i++ {
        t := newRT(i, n, 200)
        r.Add(t)
        n += p
    }
    return r
}

