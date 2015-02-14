package main

import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "github.com/pastebt/gotextcat"
)

// Calculate fingerprint for input file or string, can be used to generate *.lm

func usage() {
        fmt.Println("Usage:", os.Args[0], "-g|-l filename|-")
        fmt.Println("   -g  Generate fingerprint from input, - means stdin")
        fmt.Println("   -l  check language from input, - means stdin")
        os.Exit(1)
}


func main() {
    if len(os.Args) != 3 {
        usage()
    }
    var data []byte
    if os.Args[2] == "-" {
        data, _ = ioutil.ReadAll(os.Stdin)
    } else {
        fin, err := os.Open(os.Args[2])
        if err != nil {
            log.Fatal(err)
        }
        defer fin.Close()
        data, err = ioutil.ReadAll(fin)
    }

    switch os.Args[1] {
    case "-g": gotextcat.PrintFingerPrint(string(data))
    case "-l": {
            gotextcat.Init("../LMI/")
            l1, l2 := gotextcat.GetLanguage(string(data))
            //fmt.Println("%v, %v", l1, l2)
            if l1 != nil {fmt.Println(l1.GetName())}
            if l2 != nil {fmt.Println(l2.GetName())}
        }
    default: usage()
    }
}
