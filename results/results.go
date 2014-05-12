package results

import (
    "math"
    "net"
    "net/url"
    "sort"
    "strings"
)

// Results is a container for the performance test results.
type Results struct {
    Requested  int
    Replies    int
    TotalTime  float64
    ConnPerSec float64

    Took     []float64
    TookMin  float64
    TookMed  float64
    TookAvg  float64
    TookMax  float64
    Took85th float64
    Took90th float64
    Took95th float64
    Took99th float64

    Code    []int
    Code1xx int
    Code2xx int
    Code3xx int
    Code4xx int
    Code5xx int

    Errors            []error
    ErrorsTotal       int
    ErrorsConnTimeout int
    ErrorsConnRefused int
    ErrorsConnReset   int
    ErrorsFdUnavail   int
    ErrorsAddrUnavail int
    ErrorsOther       int

    ContentLength int64
    HeaderLength  int64
    TotalLength   int64
}

/**
 * Public Methods
 ******************************************/

// Result is the performance test result transporter.
type Result struct {
    Index, Code   int
    Took          float64
    Error         error
    TotalLength   int64
    ContentLength int64
    HeaderLength  int64
}

// Add adds Result data to Results.
func (res *Results) Add(result Result) {
    res.Took[result.Index] = result.Took
    res.Code[result.Index] = result.Code

    if result.Error != nil {
        res.Errors = append(res.Errors, result.Error)
    }

    if res.TotalLength == 0 {
        res.TotalLength = result.TotalLength
    }

    if res.ContentLength == 0 {
        res.ContentLength = result.ContentLength
    }

    if res.HeaderLength == 0 {
        res.HeaderLength = result.HeaderLength
    }
}

// Finalize finalizes results, generating min, max, avg med and percentiles.
func (res *Results) Finalize() {
    res.Replies = len(res.Took)
    res.min()
    res.max()
    res.avg()
    res.med()
    res.pct()

    // Code counts
    for _, code := range res.Code {
        if code < 100 { // ignore
        } else if code < 200 {
            res.Code1xx++
        } else if code < 300 {
            res.Code2xx++
        } else if code < 400 {
            res.Code3xx++
        } else if code < 500 {
            res.Code4xx++
        } else if code < 600 {
            res.Code5xx++
        }
    }

    // Error counts
    res.ErrorsTotal = len(res.Errors)

    for _, err := range res.Errors {
        e := err.(*url.Error).Err.(*net.OpError).Error()
        if strings.Contains(e, "connection refused") {
            res.ErrorsConnRefused++
        } else if strings.Contains(e, "connection reset") {
            res.ErrorsConnReset++
        } else if strings.Contains(e, "connection timed out") {
            res.ErrorsConnTimeout++
        } else if strings.Contains(e, "no free file descriptors") {
            res.ErrorsFdUnavail++
        } else if strings.Contains(e, "no such host") {
            res.ErrorsAddrUnavail++
        } else {
            res.ErrorsOther++
        }
    }
}

// CalculatePct calculates percentiles from existing Took values.
func (res *Results) CalculatePct(pct int) float64 {
    slice := res.copyTook()

    l := len(slice)
    switch l {
    case 0:
        return float64(0)
    case 1:
        return slice[0]
    case 2:
        return slice[1]
    }

    index := int(math.Floor(((float64(l)/100)*float64(pct))+0.5) - 1)
    return slice[index]
}

/**
 * Private Methods
 ******************************************/

func (res *Results) min() {
    slice := res.copyTook()

    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }
    res.TookMin = slice[0]
}

func (res *Results) max() {
    slice := res.copyTook()

    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }
    res.TookMax = slice[len(slice)-1]
}

func (res *Results) avg() {
    slice := res.copyTook()

    var total float64
    for _, n := range slice {
        total += n
    }
    res.TookAvg = total / float64(len(slice))
}

func (res *Results) med() {
    slice := res.copyTook()
    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }

    l := len(slice)
    switch l {
    case 0:
        res.TookMed = float64(0)
    case 1:
        res.TookMed = slice[0]
    case 2:
        res.TookMed = slice[1]
    default:
        if math.Mod(float64(l), 2) == 0 {
            index := int(math.Floor(float64(l)/2) - 1)
            lower := slice[index]
            upper := slice[index+1]
            res.TookMed = (lower + upper) / 2
        } else {
            res.TookMed = slice[l/2]
        }
    }
}

func (res *Results) pct() {
    res.Took85th = res.CalculatePct(85)
    res.Took90th = res.CalculatePct(90)
    res.Took95th = res.CalculatePct(95)
    res.Took99th = res.CalculatePct(99)
}

func (res *Results) copyTook() []float64 {
    slice := make([]float64, len(res.Took))
    copy(slice, res.Took)
    return slice
}
