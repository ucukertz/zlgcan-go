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

	handle := zcanlib.OpenDevice(zlgcan.ZCAN_USBCAN2, 0, 0)
	if handle == zlgcan.INVALID_DEVICE_HANDLE {
		log.Fatal("OpenDevice failed")
	}

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

	defer zcanlib.CloseDevice(handle)

	fmt.Println("CAN started. Ctrl+C to stop.")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	txID := uint32(0x100)
	txExtended := false
	txCount := uint8(0)

	for range ticker.C {
		for {
			n := zcanlib.GetReceiveNum(channel, zlgcan.ZCAN_TYPE_CAN)
			if n == 0 {
				break
			}
			rcv, num := zcanlib.Receive(channel, n, 0)
			for i := uint(0); i < num; i++ {
				frame := &rcv[i].Frame
				eff := "STD"
				if frame.GetFrameEFF() != 0 {
					eff = "EXT"
				}
				rtr := ""
				if frame.GetFrameRTR() != 0 {
					rtr = " RTR"
				}
				hex := ""
				for b := 0; b < int(frame.Dlc); b++ {
					hex += fmt.Sprintf("%02X ", frame.Data[b])
				}
				fmt.Printf("RX [%s] id=0x%03X dlc=%d%s data=%s\n",
					eff, frame.GetFrameID(), frame.Dlc, rtr, hex)
			}
		}

		var tx zlgcan.ZCAN_Transmit_Data
		tx.Type = 0
		if txExtended {
			tx.Frame.GenerateID(0x100000+uint32(txCount), 1, 0, 0)
		} else {
			tx.Frame.GenerateID(txID, 0, 0, 0)
		}
		tx.Frame.Dlc = 8
		data := [8]byte{txCount, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}
		tx.Frame.Data = data

		sent := zcanlib.Transmit(channel, []zlgcan.ZCAN_Transmit_Data{tx}, 1)
		eff := "STD"
		if txExtended {
			eff = "EXT"
		}
		if sent > 0 {
			fmt.Printf("TX [%s] id=0x%03X dlc=8 data=%02X %02X %02X %02X %02X %02X %02X %02X\n",
				eff, tx.Frame.GetFrameID(), data[0], data[1], data[2], data[3],
				data[4], data[5], data[6], data[7])
		} else {
			fmt.Println("TX FAILED")
		}

		txCount++
		txExtended = !txExtended
	}
}
