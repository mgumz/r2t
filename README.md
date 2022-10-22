# r2t

**r2t** converts *raw* *to* *text*.

## usage:

    $> r2t < infile > outfile

Environment variables control the output of **r2t**:

### R2T\_ENC

The encoding environment variables defines how the raw bits are converted.
Implemented options:

* b64, 64: use the base64 encoding
* b85, 85: use the base85 encoding
* hex, 16: use hex encoding
* 2: use binary encoding


### R2T\_WRAP

The **R2T_WRAP** environment variables controls at which column the text is
wrapped by a new line.

eg:

    $> env R2T_ENC=64 R2T_WRAP=32 r2t < infile

will wrapp the base64 encoded output at column 32.


## Installation

To install *r2t* to ~/go/bin use:

    $> go install -v github.com/mgumz/r2t@latest

## Building

Build local:

    $> git clone https://github.com/mgumz/r2t
    $> cd r2t
    $> go build .

More advanced build options are available in the Makefile

## License

See the LICENSE file in this repository.

