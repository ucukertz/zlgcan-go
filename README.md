# ZLGCAN Go Wrapper

[中文文档](README_CN.md)

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
go get github.com/ucukertz/zlgcan-go
```

## Usage

```go
package main

import (
    "log"
    "github.com/ucukertz/zlgcan-go"
)

func main() {
    z, err := zlgcan.NewZCAN("")
    if err != nil {
        log.Fatal(err)
    }

    ch := z.OpenAndStart(zlgcan.ZCAN_USBCAN2, 0, 0, 500000)
    if ch == zlgcan.INVALID_CHANNEL_HANDLE {
        log.Fatal("OpenAndStart failed")
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
