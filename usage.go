package main

const usage = `r2t - convert raw into text
usage: r2t < in.raw

environment variables:

R2T_ENC - encoding
          b85, 85, a85 - ascii85 encoding
          b64, 64      - base64 encoding
          16, hex      - hex encoding
          2            - binary encoding
          <unset>      - passthrough

R2T_WRAP - wrap at column c
`
