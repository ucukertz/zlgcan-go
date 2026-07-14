package main

import (
	"fmt"
	"log"
	"time"

	"zlgcan"
)

func main() {
	zcanlib, err := zlgcan.NewZCAN("./zlgcan_x64/zlgcan.dll")
	if err != nil {
		log.Fatal("Failed to load DLL: ", err)
	}

	ch := zcanlib.OpenAndStart(zlgcan.ZCAN_USBCAN2, 0, 0, 500000)
	if ch == zlgcan.INVALID_CHANNEL_HANDLE {
		log.Fatal("OpenAndStart failed")
	}
	defer zcanlib.CloseDevice(ch)

	fmt.Println("CAN started. Ctrl+C to stop.")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	txID := uint32(0x100)
	txExtended := false
	txCount := uint8(0)

	for range ticker.C {
		// Drain all pending RX frames
		for _, msg := range zcanlib.ReceiveAll(ch) {
			eff := "STD"
			if msg.Frame.GetFrameEFF() != 0 {
				eff = "EXT"
			}
			hex := ""
			for b := 0; b < int(msg.Frame.Dlc); b++ {
				hex += fmt.Sprintf("%02X ", msg.Frame.Data[b])
			}
			fmt.Printf("RX [%s] id=0x%03X dlc=%d data=%s\n",
				eff, msg.Frame.GetFrameID(), msg.Frame.Dlc, hex)
		}

		// Send one frame, alternating STD/EXT
		data := [8]byte{txCount, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}
		sent := zcanlib.TransmitOne(ch, txID, data[:], txExtended)

		eff := "STD"
		if txExtended {
			eff = "EXT"
		}
		if sent {
			fmt.Printf("TX [%s] id=0x%03X dlc=8 data=%02X %02X %02X %02X %02X %02X %02X %02X\n",
				eff, txID, data[0], data[1], data[2], data[3],
				data[4], data[5], data[6], data[7])
		} else {
			fmt.Println("TX FAILED")
		}

		txCount++
		txExtended = !txExtended
	}
}
