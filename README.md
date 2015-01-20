ipgrep - ipv4 grep-like tool
============================

`ipgrep` is a simple `IPv4` grep-like tool. It can be used to filter lines in log files matching `IPv4`'s in a specific column.

Features
========

- multiple matching.
- `CIDR` matching (`IPv4` with mask).
- `gzip` io support.

Command line
============

- **-n** `CIDR` network. Defaults to `0.0.0.0`.
- **-t** number of threads. Defaults to `1`.
- **-c** column to filter. Defaults to `1`.
- **-z** enable gzip flag.
- **-i** input file. Defaults to `stdin`.
- **-o** output file. Defaults to `stdout`.

Build
=====

Install `go` compiler and then run:

    go build ipgrep.go

Install
=======

Place the binary in one of yours binary path.

    cp ipgrep /usr/local/bin

Usage
=====

Filter lines with `10.0.0.0/8` at the `1`<sup>st</sup> column in `/var/log/apache/access.log`:

    ipgrep -i /var/log/apache/access.log -c 1 -n 10.0.0.0/8

Filter lines with `10.0.0.0/24` and `10.0.12.0/24` at the `1`<sup>st</sup> column in `/var/log/apache/access.log`:

    ipgrep -i /var/log/apache/access.log -c 1 -n '10.0.0.0/24 10.0.12.0/24'

Filter lines in `/var/log/apache/access.log.gz` to `stdout`:

    ipgrep -z -i /var/log/apache/access.log -c 1 -n 10.0.0.0/8 | gzip -c -d
    
Filter lines in `/var/log/apache/access.log.gz` to `/var/log/apache/intranet_access.log.gz`:

    ipgrep -z -i /var/log/apache/access.log.gz -o /var/log/apache/intranet_access.log.gz -c 1 -n 10.0.0.0/8
    
Version
=======

**0.1.0** according to http://semver.org/.

License
=======

    The MIT License (MIT)
    
    Copyright (c) 2015 Vinícius Tinti
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
