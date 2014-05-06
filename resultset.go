package perf

import (
    "math"
    "net"
    "net/url"
    "sort"
    "strings"
)

// Performance test results.
type ResultSet struct {
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

// Performance test result transporter.
type ResultTransport struct {
    Index, Code   int
    Took          float64
    Error         error
    TotalLength   int64
    ContentLength int64
    HeaderLength  int64
}

// Add transport data to result set.
func (this *ResultSet) Add(result ResultTransport) {
    this.Took[result.Index] = result.Took
    this.Code[result.Index] = result.Code

    if result.Error != nil {
        this.Errors = append(this.Errors, result.Error)
    }

    if this.TotalLength == 0 {
        this.TotalLength = result.TotalLength
    }

    if this.ContentLength == 0 {
        this.ContentLength = result.ContentLength
    }

    if this.HeaderLength == 0 {
        this.HeaderLength = result.HeaderLength
    }
}

// Finalize results, generating min, max, avg med and percentiles.
func (this *ResultSet) Finalize() {
    this.Replies = len(this.Took)
    this.min()
    this.max()
    this.avg()
    this.med()
    this.pct()

    // Code counts
    for _, code := range this.Code {
        if code < 100 { // ignore
        } else if code < 200 {
            this.Code1xx++
        } else if code < 300 {
            this.Code2xx++
        } else if code < 400 {
            this.Code3xx++
        } else if code < 500 {
            this.Code4xx++
        } else if code < 600 {
            this.Code5xx++
        }
    }

    // Error counts
    this.ErrorsTotal = len(this.Errors)

    for _, err := range this.Errors {
        e := err.(*url.Error).Err.(*net.OpError).Error()
        if strings.Contains(e, "connection refused") {
            this.ErrorsConnRefused++
        } else if strings.Contains(e, "connection reset") {
            this.ErrorsConnReset++
        } else if strings.Contains(e, "connection timed out") {
            this.ErrorsConnTimeout++
        } else if strings.Contains(e, "no free file descriptors") {
            this.ErrorsFdUnavail++
        } else if strings.Contains(e, "no such host") {
            this.ErrorsAddrUnavail++
        } else {
            this.ErrorsOther++
        }
    }
}

// Calculate Percentile from existing Took values.
func (this *ResultSet) CalculatePct(pct int) float64 {
    slice := this.copyTook()

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
func (this *ResultSet) min() {
    slice := this.copyTook()

    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }
    this.TookMin = slice[0]
}

func (this *ResultSet) max() {
    slice := this.copyTook()

    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }
    this.TookMax = slice[len(slice)-1]
}

func (this *ResultSet) avg() {
    slice := this.copyTook()

    var total float64
    for _, n := range slice {
        total += n
    }
    this.TookAvg = total / float64(len(slice))
}

func (this *ResultSet) med() {
    slice := this.copyTook()
    if !sort.Float64sAreSorted(slice) {
        sort.Float64s(slice)
    }

    l := len(slice)
    switch l {
    case 0:
        this.TookMed = float64(0)
    case 1:
        this.TookMed = slice[0]
    case 2:
        this.TookMed = slice[1]
    default:
        if math.Mod(float64(l), 2) == 0 {
            index := int(math.Floor(float64(l)/2) - 1)
            lower := slice[index]
            upper := slice[index+1]
            this.TookMed = (lower + upper) / 2
        } else {
            this.TookMed = slice[l/2]
        }
    }
}

func (this *ResultSet) pct() {
    this.Took85th = this.CalculatePct(85)
    this.Took90th = this.CalculatePct(90)
    this.Took95th = this.CalculatePct(95)
    this.Took99th = this.CalculatePct(99)
}

func (this *ResultSet) copyTook() []float64 {
    slice := make([]float64, len(this.Took))
    copy(slice, this.Took)
    return slice
}
