Dense
=======

File compression utility

Dependencies
------------
* Golang standard library (>= 1.8 should work)

Installation
------------
* ``go get github.com/lk16/dense``
* ``go install github.com/lk16/dense``

Example with filename parameters
-----------
``$ echo "this is some content" > testfile
$ dense -i testfile -o testfile.dense
$ dense -d -i testfile.dense -o testfile.out
$ cat testfile.out 
$ this is some content``

Example with stdin/stdout
-----------
``$ dense <testfile >testfile.dense
dense -d <testfile.dense >testfile.out``

