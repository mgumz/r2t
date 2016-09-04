# r2t

**r2t** converts *raw* *to* *text*.

## usage:

    $> r2t < infile > outfile

environment variables control the output of **r2t**:

### ENC

the encoding environment variables defines how the raw bits are converted.
implemented options:

* b64, 64: use the base64 encoding
* b85, 85: use the base85 encoding


### WRAP

the **WRAP** environment variables controls at which column the text is
wrapped by a new line.

eg:

    $> env ENC=64 WRAP=40 r2t < infile

will wrapp the base64 encoded output at column 40.


## building

    $> export GOPATH=`pwd`
    $> go get -v github.com/mgumz/r2t
    $> cp bin/r2t /usr/local/bin

