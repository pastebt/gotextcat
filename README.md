gotextcat
=========

libtextcat golang implementation


*.lm filename has format like  lang_name.id.lm

Install
=======

cd $GOROOT/src/pkg
mkdir -p github.com/pastebt
cd github.com/pastebt
git clone https://github.com/pastebt/gotextcat

copy LMI/*.lm data into /usr/share/gotextcat/data/LMI/ ,

go test github.com/pastebt/gotextcat


DEMO
====
* cd demo
* go build fp.go
got generate utility fp, which can be used to generate *.lm file
