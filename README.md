# cisco_exporter
Exporter for metrics from devices running Cisco (NX-OS/IOS XE/IOS) (via SSH) https://prometheus.io/

This is a fork of https://github.com/lwlcom/cisco_exporter that seems to no longer be maintained.

# flags
Name     | Description | Default
---------|-------------|---------
version | Print version information. |
web.listen-address | Address on which to expose metrics and web interface. | :9362
web.telemetry-path | Path under which to expose metrics. | /metrics
ssh.targets | Comma seperated list of hosts to scrape |
ssh.user | Username to use for SSH connection | cisco_exporter
ssh.keyfile | Key file to use for SSH connection | cisco_exporter
ssh.timeout | Timeout in seconds to use for SSH connection | 5
debug | Show verbose debug output | false
legacy.ciphers | Allow insecure legacy ciphers: aes128-cbc 3des-cbc aes192-cbc aes256-cbc | false
config.file | Path to config file |
dynamic-interface-labels | Parse interface and BGP descriptions to get labels dynamically | true
interface-description-regex | Give a regex to retrieve the interface description labels | `\[([^=\]]+)(=[^\]]+)?\]`

If `-config-file` is set all settings are read from the file and command line flags
are ignored.

# metrics

The following metric collectors are supported To enable or disable a collector pass a flag `--<name>.enabled=false`, where `<name>` is the name of the collector. Ot set it under features in config file.

Name     | Description | OS | Default
---------|-------------|----|--------
bgp | BGP (message count, prefix counts per peer, session state) | IOS XE/NX-OS | enabled
environment | Environment (temperatures, state of power supply) | NX-OS/IOS XE/IOS | enabled
facts | System information (OS Version, memory: total/used/free, cpu: 5s/1m/5m/interrupts) | IOS XE/IOS | enabled
interfaces | Interfaces (transmitted/received: bytes/errors/drops, admin/oper state) | NX-OS (*_drops is always 0)/IOS XE/IOS | enabled
optics | Optical signals (tx/rx) & temp | NX-OS/IOS XE/IOS | enabled
neighbors | Count of ARP & IPv6 ND entries | IOS XE/IOS | disabled
inventory | S/N & other info for liecards transceivers and other FRU | IOS XE | disabled

## Install
```bash
go get -u github.com/matejv/cisco_exporter
```

## Usage

### Binary
```bash
./cisco_exporter -ssh.targets="host1.example.com,host2.example.com:2233,172.16.0.1" -ssh.keyfile=cisco_exporter
```

```bash
./cisco_exporter -config.file=config.yml
```

## Config file
The exporter can be configured with a YAML based config file:

```yaml
---
debug: false
legacy_ciphers: false
# default values
timeout: 5
batch_size: 10000
username: default-username
password: default-password
key_file: /path/to/key
dynamic_labels: true

devices:
  - host: host1.example.com
    key_file: /path/to/key
    timeout: 5
    batch_size: 10000
    features: # enable/disable per host
      bgp: false
  - host: host2.example.com:2233
    username: exporter
    password: secret
  - host: router.*.example.com
    # Tell the exporter that this hostname should be used as a pattern when loading
    # device-specific configurations. This example would match against a hostname
    # like "router1.example.com".
    host_pattern: true
    username: exporter
    password: secret


features:
  bgp: true
  environment: true
  facts: true
  interfaces: true
  neighbors: true
  optics: true
  inventory: false

```

## Dynamic Labels

Dynamic labels can be parsed from interface descriptions. Supports key/value pairs or flags.

The feature can be disabled via command line with `-dynamic-interface-labels=false` or via config with `dynamic_labels: false`. You cannot enable or disable the feature per-host.

Default regex for parsing is `\[([^=\]]+)(=[^\]]+)?\]`. Example:

```
Description: example-r1 [prod] [cust=shop]

label1: prod
value1: 1

label2: cust
value2: shop

# or inf prometheus format:
my_metric{prod="1", cust="shop", description="example-r1 [prod] [cust=shop]"...}
```

You can modify the regex used for parsing via `-description-regex` or in config file via `description_regex` key, either globally or per-host.

This feature was ported from [junos_exporter](https://github.com/czerwonk/junos_exporter).

## Third Party Components
This software uses components of the following projects
* Prometheus Go client library (https://github.com/prometheus/client_golang)

## License
(c) Martin Poppen, 2018; Matej Vadnjal, 2023. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/
