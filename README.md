# midimonster-controller

This app controls [midimonster](https://midimonster.net) via web interface.


## Prerequisites

* golang (https://golang.org) for building the binary
* a midimonster installation

## Quick Start


## Setup

First create a config file (copy the [config.sample.yaml](config.sample.yaml)).
The controller supports two ways of controlling midimonster.

### Process

The simplest is the process control type. It spawns the midimonster process as a subprocess.
To enable the mode open the config file and configure the following settings:

```yaml
controlType: process
process:
  binPath: <path to midimonster executable> (e.g. /usr/bin/midimonster)
```


### Systemd

This mode uses systemd for management of the midimonster daemon.
It uses dbus to communicate to systemd. Make sure dbus is running.

```yaml
controlType: systemd
systemd:
    unitName: midimonster.service
```


## Start

Start the controller with

```
./bin/midimonster-controller --config <path to config.yaml>
```

After this open the web browser and open the configured url (http://<bind>:<port>/)