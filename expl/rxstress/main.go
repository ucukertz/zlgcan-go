package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zlgcan"
)

func main() {
	zcanlib, err := zlgcan.NewZCAN("./zlgcan_x64/zlgcan.dll")
	if err != nil {
		log.Fatal("Failed to load DLL: ", err)
	}

	handle := zcanlib.OpenDevice(zlgcan.ZCAN_USBCAN2, 0, 0)
	if handle == zlgcan.INVALID_DEVICE_HANDLE {
		log.Fatal("OpenDevice failed")
	}
	defer zcanlib.CloseDevice(handle)

	ip, err := zcanlib.GetIProperty(handle)
	if err != nil || ip == nil {
		log.Fatal("GetIProperty failed: ", err)
	}
	zcanlib.SetValue(ip, "0/baud_rate", "500000")
	zcanlib.ReleaseIProperty(ip)

	channel := zcanlib.InitCANFD(handle, 0, &zlgcan.ZCAN_CHANNEL_INIT_CONFIG{
		CanType: zlgcan.ZCAN_TYPE_CAN,
		Config: zlgcan.ZCAN_CHANNEL_CONFIG{
			AccCode: 0x00000000,
			AccMask: 0xFFFFFFFF,
			Filter:  0,
			Mode:    0,
		},
	})
	if channel == zlgcan.INVALID_CHANNEL_HANDLE {
		log.Fatal("InitCANFD failed")
	}
	if zcanlib.StartCAN(channel) != zlgcan.ZCAN_STATUS_OK {
		log.Fatal("StartCAN failed")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("CAN RX stress test running. Ctrl+C to stop.")

	var total uint64
	var batch uint64
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		<-sig
		ticker.Stop()
	}()

loop:
	for {
		select {
		case <-ticker.C:
			fmt.Printf("total_rx=%d last_1s=%d\n", total, batch)
			batch = 0
		case <-sig:
			break loop
		default:
		}

		n := zcanlib.GetReceiveNum(channel, zlgcan.ZCAN_TYPE_CAN)
		if n == 0 {
			continue
		}

		_, num := zcanlib.Receive(channel, n, 0)
		total += uint64(num)
		batch += uint64(num)
	}

	fmt.Printf("\nDone. total_rx=%d\n", total)
}
