#gotextcat

libtextcat golang implementation

*.lm filename has format like  lang_name.id.lm


##Install
```bash
$ cd $GOROOT/src/pkg
$ mkdir -p github.com/pastebt
$ cd github.com/pastebt
$ git clone https://github.com/pastebt/gotextcat
```
##setup data

copy LMI/*.lm data into /usr/share/gotextcat/data/LMI/ ,

##test
```bash
$ go test github.com/pastebt/gotextcat
```


#DEMO
```bash
$ cd $GOROOT/src/pkg/github.com/pastebt/gotextcat/demo
$ go build fp.go
```
got generate utility fp, which can be used to generate *.lm file
