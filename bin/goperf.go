package main

import (
    . "github.com/jmervine/goperf"
    "flag"
    "os"
)

var (
    path string
    conns, rate int
    verbose bool
)

func init() {
    // config.Path
    //flag.StringVar(&path , "path" , "" , "path")
    flag.StringVar(&path , "p"    , "" , "path")

    // config.NumConns
    //flag.IntVar(&conns , "num-conns" , 0 , "num-conns")
    flag.IntVar(&conns , "n"         , 0 , "num-conns")

    // config.Rate
    //flag.IntVar(&rate , "rate" , 0 , "rate")
    flag.IntVar(&rate , "r"    , 0 , "rate")

    // config.Verbose
    //flag.BoolVar(&verbose , "verbose" , false , "verbose")
    flag.BoolVar(&verbose , "v"       , false , "verbose")

    flag.Parse()

    if path == "" || conns == 0 {
        flag.Usage()
        os.Exit(0)
    }
}

func main() {
    config := &Configurator{
        Path: path, NumConns: conns, Rate: rate, Verbose: verbose,
    }

    results := Start(config)
    Display(results)
}
