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
* hex, 16: use hex encoding
* 2: use binary encoding


### WRAP

the **WRAP** environment variables controls at which column the text is
wrapped by a new line.

eg:

    $> env ENC=64 WRAP=32 r2t < infile

will wrapp the base64 encoded output at column 32.


## installation

to install *r2t* to ~/go/bin use:

    $> go install -v github.com/mgumz/r2t@latest

## building

build local:

    $> git clone https://github.com/mgumz/r2t
    $> cd r2t
    $> go build .

