# kodama

[![CI](https://github.com/winebarrel/kodama/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/kodama/actions/workflows/ci.yml)

kodama is a DNS server that returns any private IP address.

## Usage

```
Usage: kodama <domain> [flags]

Arguments:
  <domain>    Accepted domain ($KODAMA_DOMAIN).

Flags:
  -h, --help                Show help.
      --addr=":53"          Listening address ($KODAMA_ADDR).
      --ttl=300             Dynamic Record TTL ($KODAMA_TTL).
      --zone-data=STRING    Zone file data ($KODAMA_ZONE_DATA).
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

### Use Zone Data

```
$ cat zonefile
www.example.com. IN A 127.0.0.2
www2.example.com. IN A 127.0.0.3
$ export KODAMA_ZONE_DATA=$(cat zonefile)
$ kodama example.com &
$ dig @127.0.0.1 +short www.example.com
127.0.0.2
$ dig @127.0.0.1 +short www2.example.com
127.0.0.3
```

### Run with Docker

```
$ docker run --rm -p 53:53 -p 53:53/udp ghcr.io/winebarrel/kodama example.com
```
