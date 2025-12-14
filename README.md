# kodama

[![CI](https://github.com/winebarrel/kodama/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/kodama/actions/workflows/ci.yml)

kodama is a DNS server that returns any private IP address.

## Usage

```
Usage: kodama <domain> [flags]

Arguments:
  <domain>    Accepted domain ($KODAMA_DOMAIN).

Flags:
  -h, --help          Show help.
  -a, --addr=":53"    Listening address ($KODAMA_ADDR).
      --version
```

```
$ kodama example.com &
[1] 25270
$ dig @127.0.0.1 +short 127-0-0-1.example.com
127.0.0.1
$ dig @127.0.0.1 +short web-192-168-10-1.example.com
192.168.10.1
```

### Run with Docker

```
$ docker run --rm -p 53:53 -p 53:53/udp ghcr.io/winebarrel/kodama example.com
```
