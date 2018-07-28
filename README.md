=======
Veterani
========

There is nothing really useful for the "general public." In `src/veterani2013/iof` there is a partial typedefinition for the annotated structs that can be used to parse the IOF 2.x format used at http://oris.orientacnisporty.cz/. Othervise it is pretty boring text processing. It would be probably better done in someting like Perl or Python. The code is heavy on sorting and that is pretty painful with Go sort package.

The code uses SQLite. I had fun times choosing from the three available bindings for Go, at the end I went with what seemed the most popular.


        for i in 2282 2328 2283 2329; do wget http://oris.orientacnisporty.cz/ExportVysledkuCSOS?id=${i}; done
        
        for f in *; do iconv -f WINDOWS-1250 -t UTF-8 $f > $f.txt; done

Usage:

        export GOPATH=`pwd`

        go run src/veterani2013/zpracuj_data.go -clubs clubs2018.txt -results 2018/ -suffix .txt
        go run src/veterani2013/vypis_vysledky.go
        go run src/veterani2013/vypis_v_kategorie.go

        go run src/veterani2013/csv2fw.go 10,25,7,1,6 x.csv >
        go run src/veterani2013/dvoudenniveteraniada.go clubs.txt klasifikace.txt hlavni.txt >