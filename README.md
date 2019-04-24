# mackerel-plugin-pinging

ICMP Ping RTT custom mackerel plugin


## usage

```
Usage:
  mackerel-plugin-pinging [OPTIONS]

Application Options:
      --host=       Hostname to ping
      --max-rtt=    wait time. Max RTT(ms) (default: 200)
      --count=      Count Sending ping (default: 10)
      --key-prefix= Metric key prefix
  -v, --version     Show version

Help Options:
  -h, --help        Show this help message
  ```

## sample

```
pinging.googledns_rtt_count.success     10.000000       1556095818
pinging.googledns_rtt_count.error       0.000000        1556095818
pinging.googledns_rtt_ms.average        6.684400        1556095818
pinging.googledns_rtt_ms.max    4.464000        1556095818
pinging.googledns_rtt_ms.90_percentile  5.506000        1556095818
```
