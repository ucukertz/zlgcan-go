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
	zcanlib, err := zlgcan.NewZCAN("")
	if err != nil {
		log.Fatal("Failed to load DLL: ", err)
	}

	ch := zcanlib.OpenAndStart(zlgcan.ZCAN_USBCAN2, 0, 0, 500000)
	if ch == zlgcan.INVALID_CHANNEL_HANDLE {
		log.Fatal("OpenAndStart failed")
	}
	defer zcanlib.CloseDevice(ch)

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

		n := zcanlib.GetReceiveNum(ch, zlgcan.ZCAN_TYPE_CAN)
		if n == 0 {
			continue
		}

		_, num := zcanlib.Receive(ch, n, 0)
		total += uint64(num)
		batch += uint64(num)
	}

	fmt.Printf("\nDone. total_rx=%d\n", total)
}
