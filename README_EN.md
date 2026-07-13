# ZLGCAN Go Wrapper

Go wrapper for the ZLG CAN adapter (`zlgcan.dll`). Provides device control, channel initialization, CAN/CANFD transmit/receive, and property configuration.

## Features

- Device open / close / online detection
- Device info and channel status queries
- CAN / CANFD channel init and start
- CAN / CANFD frame transmit and receive
- Property get/set via `IProperty`
- Auxiliary APIs: device info, channel error, available devices

## Install

```
go get zlgcan
```

Ensure `zlgcan_x64/zlgcan.dll` is in the working directory.

## Usage

```go
package main

import "zlgcan"

func main() {
    z := zlgcan.NewZCAN(".\\zlgcan_x64\\zlgcan.dll")
    dev := z.OpenDevice(zlgcan.ZCAN_USBCAN2, 0, 0)
    if dev == zlgcan.INVALID_DEVICE_HANDLE {
        panic("open failed")
    }
    defer z.CloseDevice(dev)

    ch := can_start(z, dev, 0, 500000) // channel 0 at 500kbps
    if ch == zlgcan.INVALID_CHANNEL_HANDLE {
        panic("channel init failed")
    }
    defer z.CloseDevice(ch)

    // ... transmit / receive ...
}
```

See `expl/main.go` for a full example.

## Testing

```
go test
```

## Notes

- Windows only (tested)
- ZLG device driver must be installed
- Read `USB-CAN-FD-B-API-Manual.md` for API reference

## License

[MIT](https://choosealicense.com/licenses/mit/)
