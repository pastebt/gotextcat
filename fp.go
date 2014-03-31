package main

import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "fgdwcfgo/lang"
)


func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "filename|-")
        os.Exit(1)
    }
    var data []byte
    if os.Args[1] == "-" {
        data, _ = ioutil.ReadAll(os.Stdin)
    } else {
        fin, err := os.Open(os.Args[1])
        if err != nil {
            log.Fatal(err)
        }
        defer fin.Close()
        data, err = ioutil.ReadAll(fin)
    }
    lang.PrintFingerPrint(string(data))
}
