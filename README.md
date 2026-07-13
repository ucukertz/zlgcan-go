# ZLGCAN Go Wrapper

ZLG CAN 适配器（`zlgcan.dll`）的 Go 语言封装。提供设备控制、通道初始化、CAN/CANFD 收发及属性配置接口。

## 功能

- 设备打开 / 关闭 / 在线检测
- 获取设备信息与通道状态
- CAN / CANFD 通道初始化与启动
- CAN / CANFD 帧发送与接收
- 通过 `IProperty` 获取和设置设备属性
- 参考设备信息、通道错误、可用设备查询等辅助接口

## 安装

```
go get zlgcan
```

确保 `zlgcan_x64/zlgcan.dll` 在工作目录中。

## 使用

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

    ch := zcan_open(z, dev, 0, 500000) // opens channel 0 at 500kbps
    if ch == zlgcan.INVALID_CHANNEL_HANDLE {
        panic("channel init failed")
    }
    defer z.CloseDevice(ch)

    // ... send / receive ...
}
```

更多示例见 `expl/main.go`。

## 测试

```
go test
```

## 注意事项

- 仅 Windows 环境测试
- 需正确安装 ZLG 设备驱动
- 使用前请阅读 `USB-CAN-FD-B-API-Manual.md`

## 许可证

[MIT](https://choosealicense.com/licenses/mit/)
