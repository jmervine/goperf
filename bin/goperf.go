package main

import (
    "github.com/jmervine/goperf"
    "flag"
    "os"
    "fmt"
)

var (
    path string
    conns int
    rate float64
    verbose bool
    version bool
)

func init() {
    // config.Path
    //flag.StringVar(&path , "path" , "" , "path")
    flag.StringVar(&path , "u"    , "" , "Target URL.")

    // config.NumConns
    //flag.IntVar(&conns , "num-conns" , 0 , "num-conns")
    flag.IntVar(&conns , "n"         , 0 , "Total number of connections.")

    // config.Rate
    flag.Float64Var(&rate , "r"    , 0 , "Connection rate (per second).")

    // config.Verbose
    //flag.BoolVar(&verbose , "verbose" , false , "verbose")
    flag.BoolVar(&verbose , "v"       , false , "Print verbose messaging.")

    flag.BoolVar(&version , "version", false , "Show version infomration.")

    flag.Parse()

    if version {
        fmt.Printf("goperf version %v\n", perf.Version)
        os.Exit(0)
    }

    if path == "" || conns == 0 {
        flag.Usage()
        os.Exit(0)
    }
}

func main() {
    config := &perf.Configurator{
        Path: path, NumConns: conns, Rate: rate, Verbose: verbose,
    }

    results := perf.Start(config)
    perf.Display(results)
}

