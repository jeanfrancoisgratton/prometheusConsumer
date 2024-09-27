# prometheusConsumer
___


## Why *prometheusConsumer* ?

The Prometheus TS database and monitoring tool manages its server inventory in three modes:
- Static discovery : you manually update its targets list (its inventory) and restart the service
- Dynamic service discovery : you rely on tools such as Hashicorp Consul to perform the discovery and notify Prometheus about it
- File-based service discovery : a specific directory on the Prometheus server is periodically scanned, and changes in that directory is reflected in the inventory

You can combine all three modes.

In my opinion:
- The static mode is the worst, and prone to errors.
- Dynamic (consul/consul-like) mode is very effective but relies on outside tools where communication and configuration can be something of a hurdle
- File-based service discovery is very simple to implement, and effective

*prometheusConsumer* leans on the third option, which is a file-based discovery framework.

The consumer client (named `prometheusSDsendHost`) sends to a corresponding (configured) *prometheusListner*, which will act up the sent command: you can add the current host, remove it, or list all configured hosts on Prometheus.

All you need is to provide a JSON configuration file, and send a command to the configured listener.

This brings up to the ......


## *prometheusConsumer* config

The config file, located in `$HOME/.config/JFG/prometheusConsumer.json` looks like this:

```json
{
  "cacert": "/path/to/at/certificateAuthority/cert.crt",
  "listenerurl": "https://myhost:myport"
}
```

`cacert` being a valid CA certificate, as all communication is https-based
`listenerurl` being the actual URL the listener listens on

**As you see, having a CA certificate is an absolute pre-requisite, as all communications between *prometheusConsumer* and *prometheusListener* are done through TLS 1.2

To create that file, you should run the following: `prometheusSDsendHost -setup`. You will be prompted for the values to enter in the file.

You can also check the current config with `prometheus -printconfig`

## How *prometheusConsumer* works

### Adding the current host

It's as simple as typing: `prometheusSDsendHost -add`

### Removing the current host

It's as simple as typing: `prometheusSDsendHost -rm`


### Listing all hosts that Prometheus knows:

Easy-peasy: `prometheus -ls` . This will return a JSON payload. Hint: the `jq` tool is your friend.

### A note about -add and -rm

Please note that all actions sent to the *prometheusListener* are localhost-based, that is what you add/remove is, the current host (the one that sent the command with `prometheusSDsendHost { -add | -rm } )


## Installing *prometheusConsumer*

### Binary packages 

Packages exist for Alpine, RedHat-based or Debian-based distors ,provided you have the appropriate famillegratton.net repository configured on your host:

- Alpine: `apk add [--no-cache] prometheusConsumer`
- Red-Hat based (RH, CentOS, Rocky, Fedora, OpenSUSE) : `{ dnf | zypper} install prometheusConsumer`
- Debian based (Debian, Ubuntu, some other -untested) : `apt[-get] install prometheusconsumer`

### From source

You need go v1.23.1 to build. You can download it from https://go.dev/dl/

As you can read this, I assume that you already have the repo cloned, somehow

- cd to `src/`
- run: `./build.sh`

Look at the script, there are some default options, and by default the binary is installed in /opt/bin.


