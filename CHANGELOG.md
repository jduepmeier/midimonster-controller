## changelog

## 0.3.1 (2024-06-29)

### Fix

- **deps**: update module github.com/gorilla/websocket to v1.5.3

## 0.3.0 (2024-06-29)

### Feat

- **websocket**: send status on status change without polling
- **websocket**: send logs in batches
- **websocket**: send status on cmd without run loop
- **websocket**: stream logs instead of polling
- add websocket and toastify

### Fix

- **deps**: update module github.com/stretchr/testify to v1.9.0
- **websocket**: fix websocket hangs on logs and status
- **websocket**: use getStatus constant
- **websocket**: use location for websocket address instead of hardcoding it
- **deps**: update module github.com/rs/zerolog to v1.33.0
- **deps**: update module github.com/jessevdk/go-flags to v1.6.1
- **deps**: update module github.com/rs/zerolog to v1.32.0
- **deps**: update module github.com/gorilla/mux to v1.8.1
- **deps**: update module github.com/rs/zerolog to v1.31.0
- **deps**: update module github.com/rs/zerolog to v1.30.0
- **deps**: update module github.com/stretchr/testify to v1.8.4
- **deps**: update module github.com/stretchr/testify to v1.8.3
- **deps**: update module github.com/rs/zerolog to v1.29.1
- **midimonster**: remove deprecated io/ioutil package

### Refactor

- rename unused function parameters

## 0.2.0 (2023-03-19)

### Feat

- **deps**: update golang => 1.20

### Fix

- **deps**: update module gopkg.in/yaml.v2 to v3
- **deps**: update module github.com/rs/zerolog to v1.29.0
- **deps**: update module github.com/google/go-cmp to v0.5.9
- **deps**: update module github.com/coreos/go-systemd/v22 to v22.5.0

## 0.1.0 (2021-11-18)

* initial version