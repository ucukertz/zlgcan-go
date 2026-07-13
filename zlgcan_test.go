package zlgcan

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

// Change this to match your adapter: ZCAN_USBCAN1, ZCAN_USBCAN2, ZCAN_USBCANFD_100U, etc.
const testDeviceType = ZCAN_USBCAN2

func newTestZCAN(t *testing.T) *ZCAN {
	t.Helper()
	zcanlib, err := NewZCAN(".\\zlgcan_x64\\zlgcan.dll")
	if err != nil {
		t.Fatalf("Failed to load ZCAN DLL: %v", err)
	}
	t.Cleanup(func() { syscall.FreeLibrary(zcanlib.dll) })
	return zcanlib
}

func openTestDevice(t *testing.T, zcanlib *ZCAN) int {
	t.Helper()
	handle := zcanlib.OpenDevice(testDeviceType, 0, 0)
	if handle == INVALID_DEVICE_HANDLE {
		t.Fatalf("OpenDevice(%d) failed", testDeviceType)
	}
	t.Cleanup(func() { zcanlib.CloseDevice(handle) })
	return handle
}

func openAndStartCAN(t *testing.T, zcanlib *ZCAN, channel uint) (chanHandle int) {
	t.Helper()
	handle := openTestDevice(t, zcanlib)

	ip, err := zcanlib.GetIProperty(handle)
	if err != nil {
		t.Fatalf("GetIProperty failed: %v", err)
	}
	ret := zcanlib.SetValue(ip, "0/baud_rate", "500000")
	if ret != ZCAN_STATUS_OK {
		t.Logf("Warning: Set baud_rate returned %d", ret)
	}
	zcanlib.ReleaseIProperty(ip)

	chanHandle = zcanlib.InitCANFD(handle, channel, &ZCAN_CHANNEL_INIT_CONFIG{
		CanType: ZCAN_TYPE_CAN,
		Config: ZCAN_CHANNEL_CONFIG{
			AccCode: 0x00000000,
			AccMask: 0xFFFFFFFF,
			Filter:  0,
			Mode:    0,
		},
	})
	if chanHandle == INVALID_CHANNEL_HANDLE {
		t.Fatalf("InitCANFD failed on channel %d", channel)
	}
	t.Cleanup(func() {
		zcanlib.ResetCAN(chanHandle)
	})

	if zcanlib.StartCAN(chanHandle) != ZCAN_STATUS_OK {
		t.Fatalf("StartCAN failed on channel %d", channel)
	}

	return chanHandle
}

func TestOpenAndCloseDevice(t *testing.T) {
	zcanlib := newTestZCAN(t)
	handle := openTestDevice(t, zcanlib)
	t.Logf("Opened device handle: %d", handle)
}

func TestGetDeviceInfo(t *testing.T) {
	zcanlib := newTestZCAN(t)
	handle := openTestDevice(t, zcanlib)

	info := zcanlib.GetDeviceInf(handle)
	if info == nil {
		t.Fatalf("GetDeviceInf returned nil")
	}
	t.Logf("hw_Version:  %s", info.HwVersion())
	t.Logf("fw_Version:  %s", info.FwVersion())
	t.Logf("dr_Version:  %s", info.DrVersion())
	t.Logf("in_Version:  %s", info.InVersion())
	t.Logf("irq_Num:     %d", info.IrqNum())
	t.Logf("can_Num:     %d", info.CanNum())
	t.Logf("Serial:      %s", info.Serial())
	t.Logf("hw_Type:     %s", info.HwType())
}

func TestDeviceOnLine(t *testing.T) {
	zcanlib := newTestZCAN(t)
	handle := openTestDevice(t, zcanlib)

	ret := zcanlib.IsDeviceOnLine(handle)
	t.Logf("IsDeviceOnLine: %d", ret)
	if ret != ZCAN_STATUS_ONLINE && ret != ZCAN_STATUS_OFFLINE {
		t.Fatalf("IsDeviceOnLine: unexpected value %d", ret)
	}
}

func TestInitCAN(t *testing.T) {
	zcanlib := newTestZCAN(t)
	handle := openTestDevice(t, zcanlib)

	chanHandle := zcanlib.InitCANFD(handle, 0, &ZCAN_CHANNEL_INIT_CONFIG{
		CanType: ZCAN_TYPE_CAN,
		Config:  ZCAN_CHANNEL_CONFIG{Mode: 0},
	})
	if chanHandle == INVALID_CHANNEL_HANDLE {
		t.Fatalf("InitCANFD failed on channel 0")
	}
	t.Logf("Channel handle: %d", chanHandle)
}

func TestGetAndReleaseIProperty(t *testing.T) {
	zcanlib := newTestZCAN(t)
	handle := openTestDevice(t, zcanlib)

	ip, err := zcanlib.GetIProperty(handle)
	if err != nil {
		t.Fatalf("GetIProperty failed: %v", err)
	}
	defer zcanlib.ReleaseIProperty(ip)
	t.Logf("IProperty handle: %+v", ip)
}

func TestOpenAndStartCAN(t *testing.T) {
	zcanlib := newTestZCAN(t)
	chanHandle := openAndStartCAN(t, zcanlib, 0)
	t.Logf("Channel handle: %d", chanHandle)
}

func TestTransmitAndReceive(t *testing.T) {
	zcanlib := newTestZCAN(t)
	chanHandle := openAndStartCAN(t, zcanlib, 0)

	var tx ZCAN_Transmit_Data
	tx.Type = 0
	tx.Frame.GenerateID(0x100, 0, 0, 0)
	tx.Frame.Dlc = 8
	for j := 0; j < 8; j++ {
		tx.Frame.Data[j] = uint8(j)
	}
	sent := zcanlib.Transmit(chanHandle, []ZCAN_Transmit_Data{tx}, 1)
	if sent == 0 {
		t.Fatalf("Transmit returned 0 frames sent")
	}
	t.Logf("Transmitted: sent %d frame(s)", sent)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		n := zcanlib.GetReceiveNum(chanHandle, ZCAN_TYPE_CAN)
		if n > 0 {
			rcv, _ := zcanlib.Receive(chanHandle, n, 0)
			for i := range rcv {
				hex := ""
				for b := range rcv[i].Frame.Data {
					hex += fmt.Sprintf("%02x ", rcv[i].Frame.Data[b])
				}
				t.Logf("CAN reply [%d]: ts:%d id:%x dlc:%d eff:%d rtr:%d err:%d data:%s",
					i, rcv[i].Timestamp, rcv[i].Frame.GetFrameID(), rcv[i].Frame.Dlc,
					rcv[i].Frame.GetFrameEFF(), rcv[i].Frame.GetFrameRTR(), rcv[i].Frame.GetFrameERR(), hex)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Log("No reply within timeout (expected if no remote node)")
}
