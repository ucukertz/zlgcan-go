package zlgcan

/*
#cgo windows LDFLAGS: -lws2_32
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// ZCAN device type constants. Used with OpenDevice to specify the hardware model.
// See ZLG USB-CAN-FD-B API Manual §3.1.
const (
	ZCAN_PCI5121              = 0x1
	ZCAN_PCI9810              = 0x2
	ZCAN_USBCAN1              = 0x3
	ZCAN_USBCAN2              = 0x4
	ZCAN_PCI9820              = 0x5
	ZCAN_CAN232               = 0x6
	ZCAN_PCI5110              = 0x7
	ZCAN_CANLITE              = 0x8
	ZCAN_ISA9620              = 0x9
	ZCAN_ISA5420              = 0xa
	ZCAN_PC104CAN             = 0xb
	ZCAN_CANETUDP             = 0xc
	ZCAN_CANETE               = 0xc
	ZCAN_DNP9810              = 0xd
	ZCAN_PCI9840              = 0xe
	ZCAN_PC104CAN2            = 0xf
	ZCAN_PCI9820I             = 0x10
	ZCAN_CANETTCP             = 0x11
	ZCAN_PCIE_9220            = 0x12
	ZCAN_PCI5010U             = 0x13
	ZCAN_USBCAN_E_U           = 0x14
	ZCAN_USBCAN_2E_U          = 0x15
	ZCAN_PCI5020U             = 0x16
	ZCAN_EG20T_CAN            = 0x17
	ZCAN_PCIE9221             = 0x18
	ZCAN_WIFICAN_TCP          = 0x19
	ZCAN_WIFICAN_UDP          = 0x1a
	ZCAN_PCIe9120             = 0x1b
	ZCAN_PCIe9110             = 0x1c
	ZCAN_PCIe9140             = 0x1d
	ZCAN_USBCAN_4E_U          = 0x1f
	ZCAN_CANDTU_200UR         = 0x20
	ZCAN_CANDTU_MINI          = 0x21
	ZCAN_USBCAN_8E_U          = 0x22
	ZCAN_CANREPLAY            = 0x23
	ZCAN_CANDTU_NET           = 0x24
	ZCAN_CANDTU_100UR         = 0x25
	ZCAN_PCIE_CANFD_100U      = 0x26
	ZCAN_PCIE_CANFD_200U      = 0x27
	ZCAN_PCIE_CANFD_400U      = 0x28
	ZCAN_USBCANFD_200U        = 0x29
	ZCAN_USBCANFD_100U        = 0x2a
	ZCAN_USBCANFD_MINI        = 0x2b
	ZCAN_CANFDCOM_100IE       = 0x2c
	ZCAN_CANSCOPE             = 0x2d
	ZCAN_CLOUD                = 0x2e
	ZCAN_CANDTU_NET_400       = 0x2f
	ZCAN_CANFDNET_TCP         = 0x30
	ZCAN_CANFDNET_200U_TCP    = 0x30
	ZCAN_CANFDNET_UDP         = 0x31
	ZCAN_CANFDNET_200U_UDP    = 0x31
	ZCAN_CANFDWIFI_TCP        = 0x32
	ZCAN_CANFDWIFI_100U_TCP   = 0x32
	ZCAN_CANFDWIFI_UDP        = 0x33
	ZCAN_CANFDWIFI_100U_UDP   = 0x33
	ZCAN_CANFDNET_400U_TCP    = 0x34
	ZCAN_CANFDNET_400U_UDP    = 0x35
	ZCAN_CANFDBLUE_200U       = 0x36
	ZCAN_CANFDNET_100U_TCP    = 0x37
	ZCAN_CANFDNET_100U_UDP    = 0x38
	ZCAN_CANFDNET_800U_TCP    = 0x39
	ZCAN_CANFDNET_800U_UDP    = 0x3a
	ZCAN_USBCANFD_800U        = 0x3b
	ZCAN_PCIE_CANFD_100U_EX   = 0x3c
	ZCAN_PCIE_CANFD_400U_EX   = 0x3d
	ZCAN_PCIE_CANFD_200U_MINI = 0x3e
	ZCAN_PCIE_CANFD_200U_M2   = 0x3f
	ZCAN_CANFDDTU_400_TCP     = 0x40
	ZCAN_CANFDDTU_400_UDP     = 0x41
	ZCAN_CANFDWIFI_200U_TCP   = 0x42
	ZCAN_CANFDWIFI_200U_UDP   = 0x43
	ZCAN_OFFLINE_DEVICE       = 0x62
	ZCAN_VIRTUAL_DEVICE       = 0x63
)

// Status codes returned by most ZCAN API functions.
// See ZLG USB-CAN-FD-B API Manual §3.2, §3.6.
const (
	ZCAN_STATUS_ERR         = 0
	ZCAN_STATUS_OK          = 1
	ZCAN_STATUS_ONLINE      = 2
	ZCAN_STATUS_OFFLINE     = 3
	ZCAN_STATUS_UNSUPPORTED = 4
)

// Channel type identifiers for InitCAN and GetReceiveNum.
// See ZLG USB-CAN-FD-B API Manual §2.2, §3.11.
const (
	ZCAN_TYPE_CAN   = 0x0
	ZCAN_TYPE_CANFD = 0x1
)

// Sentinel values indicating an invalid device or channel handle.
// See ZLG USB-CAN-FD-B API Manual §3.1, §3.5.
const (
	INVALID_DEVICE_HANDLE  = 0
	INVALID_CHANNEL_HANDLE = 0
)

// ZCAN_CHANNEL_CONFIG defines the inner configuration for a CAN or CANFD channel.
// For CAN: use AccCode, AccMask, Filter, Mode.
// For CANFD: additionally use AbitTiming, DbitTiming, Brp, Pad.
// See ZLG USB-CAN-FD-B API Manual §2.2.
type ZCAN_CHANNEL_CONFIG struct {
	AccCode    uint32 // Acceptance code. Set to 0 to accept all. §2.2 CAN.
	AccMask    uint32 // Acceptance mask. 0xFFFFFFFF = receive all. §2.2 CAN.
	AbitTiming uint32 // Arbitration domain timing. Do not set directly; use IProperty. §2.2 CANFD.
	DbitTiming uint32 // Data domain timing. Do not set directly; use IProperty. §2.2 CANFD.
	Brp        uint32 // Baud prescaler for CANFD. Set to 0. §2.2 CANFD.
	Filter     uint8  // Filtering mode: 0 = double, 1 = single. §2.2 CAN.
	Mode       uint8  // Working mode: 0 = normal, 1 = listen-only. §2.2 CAN.
	Pad        uint16 // Data alignment for CANFD. Do not set. §2.2 CANFD.
	Reserved   uint32 // Reserved.
}

// ZCAN_CHANNEL_INIT_CONFIG is the initialization structure for a CAN channel.
// Pass to InitCAN or InitCANFD before starting the channel.
// See ZLG USB-CAN-FD-B API Manual §2.2, §3.5.
type ZCAN_CHANNEL_INIT_CONFIG struct {
	CanType uint32              // 0 = TYPE_CAN, 1 = TYPE_CANFD. §2.2.
	Config  ZCAN_CHANNEL_CONFIG // Channel-specific configuration. §2.2.
}

// Type aliases for clarity.
type ZCAN_CANFD_CHANNEL_INIT_CONFIG = ZCAN_CHANNEL_INIT_CONFIG
type ZCAN_NORMAL_CHANNEL_INIT_CONFIG = ZCAN_CHANNEL_INIT_CONFIG

// ZCAN_DEVICE_INFO holds basic device information returned by GetDeviceInf.
// See ZLG USB-CAN-FD-B API Manual §2.1, §3.3.
type ZCAN_DEVICE_INFO struct {
	hw_Version     uint16   // Hardware version (hex). §2.1.
	fw_Version     uint16   // Firmware version (hex). §2.1.
	dr_Version     uint16   // Driver version (hex). §2.1.
	in_Version     uint16   // Interface library version (hex). §2.1.
	irq_Num        uint16   // Interrupt number used by the board. §2.1.
	can_Num        uint8    // Number of CAN channels. §2.1.
	str_Serial_Num [20]uint8 // Serial number string (null-terminated). §2.1.
	str_hw_Type    [40]uint8 // Hardware type string (null-terminated). §2.1.
	reserved       [8]uint16 // Reserved.
}

// HwVersion returns the hardware version as "V<major>.<minor>".
func (i *ZCAN_DEVICE_INFO) HwVersion() string {
	return fmt.Sprintf("V%d.%02X", i.hw_Version>>8, i.hw_Version&0xFF)
}

// FwVersion returns the firmware version as "V<major>.<minor>".
func (i *ZCAN_DEVICE_INFO) FwVersion() string {
	return fmt.Sprintf("V%d.%02X", i.fw_Version>>8, i.fw_Version&0xFF)
}

// DrVersion returns the driver version as "V<major>.<minor>".
func (i *ZCAN_DEVICE_INFO) DrVersion() string {
	return fmt.Sprintf("V%d.%02X", i.dr_Version>>8, i.dr_Version&0xFF)
}

// InVersion returns the interface library version as "V<major>.<minor>".
func (i *ZCAN_DEVICE_INFO) InVersion() string {
	return fmt.Sprintf("V%d.%02X", i.in_Version>>8, i.in_Version&0xFF)
}

// IrqNum returns the interrupt number used by the board.
func (i *ZCAN_DEVICE_INFO) IrqNum() uint16 {
	return i.irq_Num
}

// CanNum returns the number of CAN channels on the device.
func (i *ZCAN_DEVICE_INFO) CanNum() uint8 {
	return i.can_Num
}

// Serial returns the device serial number string.
func (i *ZCAN_DEVICE_INFO) Serial() string {
	return strings.TrimRight(string(i.str_Serial_Num[:]), "\x00 ")
}

// HwType returns the hardware type string (e.g. "USBCANFD0002").
func (i *ZCAN_DEVICE_INFO) HwType() string {
	return strings.TrimRight(string(i.str_hw_Type[:]), "\x00 ")
}

// ZCAN_CAN_FRAME represents a standard CAN 2.0 frame (8-byte payload).
// The Id field encodes both the frame ID and flags (EFF/RTR/ERR) in the upper bits.
// See ZLG USB-CAN-FD-B API Manual §2.3.
type ZCAN_CAN_FRAME struct {
	Id     uint32   // Frame ID with flags in upper 3 bits. §2.3.
	Dlc    uint8    // Data length (0..8). §2.3.
	__pad  uint8    // Padding (ignore).
	__res0 uint8    // Reserved.
	__res1 uint8    // Reserved.
	Data   [8]uint8 // Frame payload. §2.3.
}

// GenerateID sets the frame ID and applies flag bits.
// eff: extended frame flag (0=standard 11-bit, 1=extended 29-bit).
// rtr: remote transmission request (0=data, 1=remote).
// err: error frame flag (must be 0).
// See ZLG USB-CAN-FD-B API Manual §2.3.
func (f *ZCAN_CAN_FRAME) GenerateID(id uint32, eff, rtr, err uint8) {
	f.Id = id
	if eff != 0 {
		f.Id |= 0x80000000
	}
	if rtr != 0 {
		f.Id |= 0x40000000
	}
	if err != 0 {
		f.Id |= 0x20000000
	}
}

// GetFrameID extracts the actual frame ID (lower 29 bits) from the encoded Id.
func (f *ZCAN_CAN_FRAME) GetFrameID() uint32 {
	return f.Id & 0x1FFFFFFF
}

// GetFrameEFF returns 1 if this is an extended frame (29-bit ID), 0 for standard (11-bit).
func (f *ZCAN_CAN_FRAME) GetFrameEFF() uint8 {
	if (f.Id & 0x80000000) != 0 {
		return 1
	}
	return 0
}

// GetFrameRTR returns 1 if this is a remote frame, 0 for data frame.
func (f *ZCAN_CAN_FRAME) GetFrameRTR() uint8 {
	if (f.Id & 0x40000000) != 0 {
		return 1
	}
	return 0
}

// GetFrameERR returns 1 if this is an error frame, 0 otherwise.
func (f *ZCAN_CAN_FRAME) GetFrameERR() uint8 {
	if (f.Id & 0x20000000) != 0 {
		return 1
	}
	return 0
}

// ZCAN_CANFD_FRAME represents a CAN FD frame (up to 64-byte payload).
// See ZLG USB-CAN-FD-B API Manual §2.4.
type ZCAN_CANFD_FRAME struct {
	Id     uint32    // Frame ID with flags in upper 3 bits. §2.4.
	Len    uint8     // Data length (0..64). §2.4.
	Flags  uint8     // CANFD flags: bit 0 = BRS, bit 1 = ESI. §2.4.
	__res0 uint8     // Reserved.
	__res1 uint8     // Reserved.
	Data   [64]uint8 // Frame payload. §2.4.
}

// GenerateID sets the frame ID and applies flag bits (same as ZCAN_CAN_FRAME).
func (f *ZCAN_CANFD_FRAME) GenerateID(id uint32, eff, rtr, err uint8) {
	f.Id = id
	if eff != 0 {
		f.Id |= 0x80000000
	}
	if rtr != 0 {
		f.Id |= 0x40000000
	}
	if err != 0 {
		f.Id |= 0x20000000
	}
}

// GetFrameID extracts the actual frame ID (lower 29 bits).
func (f *ZCAN_CANFD_FRAME) GetFrameID() uint32 {
	return f.Id & 0x1FFFFFFF
}

// GetFrameEFF returns 1 for extended frame, 0 for standard.
func (f *ZCAN_CANFD_FRAME) GetFrameEFF() uint8 {
	if (f.Id & 0x80000000) != 0 {
		return 1
	}
	return 0
}

// GetFrameRTR returns 1 for remote frame, 0 for data frame.
func (f *ZCAN_CANFD_FRAME) GetFrameRTR() uint8 {
	if (f.Id & 0x40000000) != 0 {
		return 1
	}
	return 0
}

// GetFrameERR returns 1 for error frame, 0 otherwise.
func (f *ZCAN_CANFD_FRAME) GetFrameERR() uint8 {
	if (f.Id & 0x20000000) != 0 {
		return 1
	}
	return 0
}

// GenerateFlags sets CANFD-specific flags.
// brs: Bit Rate Switch (0 = fixed, 1 = switch to data rate).
// esi: Error State Indicator (0 = error active, 1 = error passive).
// See ZLG USB-CAN-FD-B API Manual §2.4.
func (f *ZCAN_CANFD_FRAME) GenerateFlags(brs, esi, __res uint8) {
	f.Flags = 0
	if brs != 0 {
		f.Flags |= 0x01
	}
	if esi != 0 {
		f.Flags |= 0x02
	}
}

// GetFrameBRS returns 1 if Bit Rate Switch is enabled.
func (f *ZCAN_CANFD_FRAME) GetFrameBRS() uint8 {
	if (f.Flags & 0x01) != 0 {
		return 1
	}
	return 0
}

// GetFrameESI returns 1 if Error State Indicator is set.
func (f *ZCAN_CANFD_FRAME) GetFrameESI() uint8 {
	if (f.Flags & 0x02) != 0 {
		return 1
	}
	return 0
}

// ZCAN_Transmit_Data is the transmit structure for CAN 2.0 frames.
// Used with Transmit. Type field: 0=normal, 1=single, 2=self-reception.
// See ZLG USB-CAN-FD-B API Manual §2.5, §3.9.
type ZCAN_Transmit_Data struct {
	Frame ZCAN_CAN_FRAME // Frame to send. §2.5.
	Type  uint32         // Send type: 0=normal (auto-retry), 1=single, 2=self-reception. §2.5.
}

// ZCAN_Receive_Data is the receive structure for CAN 2.0 frames.
// Timestamp is in microseconds from device startup.
// See ZLG USB-CAN-FD-B API Manual §2.7, §3.12.
type ZCAN_Receive_Data struct {
	Frame     ZCAN_CAN_FRAME // Received frame. §2.7.
	Timestamp uint64         // Timestamp in microseconds. §2.7.
}

// ZCAN_TransmitFD_Data is the transmit structure for CAN FD frames.
// See ZLG USB-CAN-FD-B API Manual §2.6, §3.10.
type ZCAN_TransmitFD_Data struct {
	Frame ZCAN_CANFD_FRAME // CANFD frame to send. §2.6.
	Type  uint32           // Send type (same as ZCAN_Transmit_Data). §2.6.
}

// ZCAN_ReceiveFD_Data is the receive structure for CAN FD frames.
// Timestamp is in microseconds from device startup.
// See ZLG USB-CAN-FD-B API Manual §2.8, §3.13.
type ZCAN_ReceiveFD_Data struct {
	Frame     ZCAN_CANFD_FRAME // Received CANFD frame. §2.8.
	Timestamp uint64           // Timestamp in microseconds. §2.8.
}

// ZCAN_IProperty is the property configuration interface for getting/setting
// device parameters (baud rate, filters, etc.) via path strings.
// Obtained with GetIProperty, released with ReleaseIProperty.
// See ZLG USB-CAN-FD-B API Manual §2.9, §3.14, §3.15.
type ZCAN_IProperty struct {
	SetValue     *[0]byte // Set attribute value. §2.9.
	GetValue     *[0]byte // Get attribute value. §2.9.
	GetPropertys *[0]byte // Return all attributes. §2.9.
}

// ZCAN is the main wrapper around the ZLG zlgcan.dll.
// Create with NewZCAN, then use methods to control CAN hardware.
type ZCAN struct {
	dll syscall.Handle // Handle to the loaded zlgcan.dll.
}

// NewZCAN loads the ZLG CAN DLL and returns a ZCAN instance.
// dllPath should point to zlgcan.dll (e.g. ".\\zlgcan_x64\\zlgcan.dll").
// Returns error if not on Windows or if the DLL cannot be loaded.
func NewZCAN(dllPath string) (*ZCAN, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("unsupported OS")
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load DLL: %w", err)
	}
	return &ZCAN{dll: dll}, nil
}

// OpenDevice opens the CAN device. A device can only be opened once.
// deviceType: hardware model constant (e.g. ZCAN_USBCAN2).
// deviceIndex: device index (0 for first device).
// Returns INVALID_DEVICE_HANDLE on failure.
// See ZLG USB-CAN-FD-B API Manual §3.1.
func (zc *ZCAN) OpenDevice(deviceType int, deviceIndex int, reserved int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_OpenDevice")
	if err != nil {
		fmt.Println("Failed to get ZCAN_OpenDevice:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceType), uintptr(deviceIndex), uintptr(reserved))
	return int(ret)
}

// CloseDevice closes the device. Each OpenDevice must be paired with a CloseDevice.
// See ZLG USB-CAN-FD-B API Manual §3.2.
func (zc *ZCAN) CloseDevice(deviceHandle int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_CloseDevice")
	if err != nil {
		fmt.Println("Failed to get ZCAN_CloseDevice:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle))
	return int(ret)
}

// GetDeviceInf retrieves device information (versions, serial, channel count).
// Returns nil on failure.
// See ZLG USB-CAN-FD-B API Manual §3.3.
func (zc *ZCAN) GetDeviceInf(deviceHandle int) *ZCAN_DEVICE_INFO {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetDeviceInf")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetDeviceInf:", err)
		return nil
	}
	info := ZCAN_DEVICE_INFO{}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle), uintptr(unsafe.Pointer(&info)))
	if ret == ZCAN_STATUS_OK {
		return &info
	}
	fmt.Println("ZCAN_GetDeviceInf failed with code:", ret)
	return nil
}

// IsDeviceOnLine checks whether the device is online.
// Returns ZCAN_STATUS_ONLINE or ZCAN_STATUS_OFFLINE.
// See ZLG USB-CAN-FD-B API Manual §3.4.
func (zc *ZCAN) IsDeviceOnLine(deviceHandle int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_IsDeviceOnLine")
	if err != nil {
		fmt.Println("Failed to get ZCAN_IsDeviceOnLine:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle))
	return int(ret)
}

// InitCAN initializes a CAN channel. Returns a channel handle for use with
// StartCAN, Transmit, Receive, etc. Must call SetCANBaudRate first.
// See ZLG USB-CAN-FD-B API Manual §3.5.
func (zc *ZCAN) InitCAN(deviceHandle int, canIndex uint, initConfig *ZCAN_CHANNEL_INIT_CONFIG) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_InitCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_InitCAN:", err)
		return INVALID_CHANNEL_HANDLE
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle), uintptr(canIndex), uintptr(unsafe.Pointer(initConfig)))
	return int(ret)
}

// InitCANFD is an alias for InitCAN. Both set baud rate via IProperty.
// See ZLG USB-CAN-FD-B API Manual §3.5.
func (zc *ZCAN) InitCANFD(deviceHandle int, canIndex uint, initConfig *ZCAN_CHANNEL_INIT_CONFIG) int {
	return zc.InitCAN(deviceHandle, canIndex, initConfig)
}

// StartCAN starts the CAN channel. Must be called after InitCAN.
// See ZLG USB-CAN-FD-B API Manual §3.6.
func (zc *ZCAN) StartCAN(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_StartCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_StartCAN:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

// ResetCAN resets the CAN channel. Call StartCAN to recover.
// See ZLG USB-CAN-FD-B API Manual §3.7.
func (zc *ZCAN) ResetCAN(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ResetCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ResetCAN:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

// SetCANBaudRate sets the baud rate for a CAN channel via IProperty.
// canIndex: channel index (0 = channel 1, 1 = channel 2).
// baudrate: desired baud rate (e.g. 500000, 1000000).
// Must be called before InitCAN.
// See ZLG USB-CAN-FD-B API Manual §4 (Attribute List).
func (zc *ZCAN) SetCANBaudRate(deviceHandle int, canIndex int, baudrate int) uint {
	ip, err := zc.GetIProperty(deviceHandle)
	if err != nil || ip == nil {
		fmt.Println("Failed to get IProperty:", err)
		return ZCAN_STATUS_ERR
	}
	defer zc.ReleaseIProperty(ip)

	path := fmt.Sprintf("%d/baud_rate", canIndex)
	value := fmt.Sprintf("%d", baudrate)

	return zc.SetValue(ip, path, value)
}

// ClearBuffer clears the library receive buffer for the specified channel.
// See ZLG USB-CAN-FD-B API Manual §3.8.
func (zc *ZCAN) ClearBuffer(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ClearBuffer")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ClearBuffer:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

// Transmit sends CAN 2.0 frames. Returns the number of frames successfully sent.
// stdMsg: slice of ZCAN_Transmit_Data to send.
// length: number of frames to send (typically len(stdMsg)).
// See ZLG USB-CAN-FD-B API Manual §3.9.
func (zc *ZCAN) Transmit(channelHandle int, stdMsg []ZCAN_Transmit_Data, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_Transmit")
	if err != nil {
		fmt.Println("Failed to get ZCAN_Transmit:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&stdMsg[0])), uintptr(length))
	return uint(ret)
}

// TransmitFD sends CAN FD frames. Returns the number of frames successfully sent.
// fdMsg: slice of ZCAN_TransmitFD_Data to send.
// length: number of frames to send (typically len(fdMsg)).
// See ZLG USB-CAN-FD-B API Manual §3.10.
func (zc *ZCAN) TransmitFD(channelHandle int, fdMsg []ZCAN_TransmitFD_Data, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_TransmitFD")
	if err != nil {
		fmt.Println("Failed to get ZCAN_TransmitFD:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&fdMsg[0])), uintptr(length))
	return uint(ret)
}

// Receive reads CAN 2.0 frames from the buffer.
// rcvNum: max frames to read (use GetReceiveNum to query available).
// waitTime: block timeout in ms. -1 = wait forever, 0 = non-blocking.
// Returns the received frames and actual count.
// See ZLG USB-CAN-FD-B API Manual §3.12.
func (zc *ZCAN) Receive(channelHandle int, rcvNum uint, waitTime int) ([]ZCAN_Receive_Data, uint) {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_Receive")
	if err != nil {
		fmt.Println("Failed to get ZCAN_Receive:", err)
		return nil, 0
	}
	msgs := make([]ZCAN_Receive_Data, rcvNum)
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&msgs[0])), uintptr(rcvNum), uintptr(waitTime))
	return msgs, uint(ret)
}

// ReceiveFD reads CAN FD frames from the buffer.
// rcvNum: max frames to read (use GetReceiveNum with CANFD type).
// waitTime: block timeout in ms. -1 = wait forever, 0 = non-blocking.
// Returns the received frames and actual count.
// See ZLG USB-CAN-FD-B API Manual §3.13.
func (zc *ZCAN) ReceiveFD(channelHandle int, rcvNum uint, waitTime int) ([]ZCAN_ReceiveFD_Data, uint) {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ReceiveFD")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ReceiveFD:", err)
		return nil, 0
	}
	if waitTime == 0 {
		waitTime = -1
	}
	msgs := make([]ZCAN_ReceiveFD_Data, rcvNum)
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&msgs[0])), uintptr(rcvNum), uintptr(waitTime))
	return msgs, uint(ret)
}

// GetReceiveNum returns the number of pending CAN/CANFD frames in the buffer.
// canType: 0 = CAN, 1 = CANFD.
// See ZLG USB-CAN-FD-B API Manual §3.11.
func (zc *ZCAN) GetReceiveNum(channelHandle int, canType uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetReceiveNum")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetReceiveNum:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(canType))
	return uint(ret)
}

// TransmitData sends raw data via the ZCAN_TransmitData DLL export.
// Use with caution — prefer Transmit/TransmitFD for typed access.
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) TransmitData(channelHandle int, data unsafe.Pointer, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_TransmitData")
	if err != nil {
		fmt.Println("Failed to get ZCAN_TransmitData:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(data), uintptr(length))
	return uint(ret)
}

// ReceiveData receives raw data via the ZCAN_ReceiveData DLL export.
// Use with caution — prefer Receive/ReceiveFD for typed access.
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) ReceiveData(channelHandle int, data unsafe.Pointer, length uint, waitTime int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ReceiveData")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ReceiveData:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(data), uintptr(length), uintptr(waitTime))
	return uint(ret)
}

// GetIProperty returns the property configuration interface for the device.
// Use SetValue/GetValue to configure baud rate, filters, etc.
// Must call ReleaseIProperty when done.
// See ZLG USB-CAN-FD-B API Manual §3.14.
func (zc *ZCAN) GetIProperty(deviceHandle int) (*ZCAN_IProperty, error) {
	proc, err := syscall.GetProcAddress(zc.dll, "GetIProperty")
	if err != nil {
		return nil, fmt.Errorf("failed to get GetIProperty: %w", err)
	}
	ret, _, callErr := syscall.SyscallN(proc, uintptr(deviceHandle))
	if callErr != 0 {
		return nil, fmt.Errorf("error calling GetIProperty: %w", callErr)
	}
	return (*ZCAN_IProperty)(unsafe.Pointer(ret)), nil
}

// SetValue sets a device attribute via IProperty.
// path: attribute path (e.g. "0/baud_rate" for channel 0 baud rate).
// value: attribute value as string (e.g. "500000").
// Returns ZCAN_STATUS_OK or ZCAN_STATUS_ERR.
// See ZLG USB-CAN-FD-B API Manual §4 (Attribute List).
func (zc *ZCAN) SetValue(ip *ZCAN_IProperty, path, value string) uint {
	cPath := C.CString(path)
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cPath))
	defer C.free(unsafe.Pointer(cValue))
	ret, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(ip.SetValue)), uintptr(unsafe.Pointer(cPath)), uintptr(unsafe.Pointer(cValue)))
	return uint(ret)
}

// GetValue gets a device attribute via IProperty.
// path: attribute path (e.g. "0/baud_rate").
// Returns the value as a string, or empty string on failure.
// See ZLG USB-CAN-FD-B API Manual §4 (Attribute List).
func (zc *ZCAN) GetValue(ip *ZCAN_IProperty, path string) string {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	ret, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(ip.GetValue)), uintptr(unsafe.Pointer(cPath)))
	return C.GoString((*C.char)(unsafe.Pointer(ret)))
}

// ReleaseIProperty releases the property interface obtained with GetIProperty.
// Must be called when done with IProperty operations.
// See ZLG USB-CAN-FD-B API Manual §3.15.
func (zc *ZCAN) ReleaseIProperty(ip *ZCAN_IProperty) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ReleaseIProperty")
	if err != nil {
		fmt.Println("Failed to get ReleaseIProperty:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(unsafe.Pointer(ip)))
	return uint(ret)
}

// ZCAN_SetReference sets device attributes via the ZCAN API.
// RefType=0: baud rate, pData = DWORD timing value (e.g. 0x060007 = 500kbps).
// See ZLG USB-CAN-FD-B API Manual §6.2.11.
func (zc *ZCAN) ZCAN_SetReference(channelHandle int, refType uint, pData unsafe.Pointer) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_SetReference")
	if err != nil {
		fmt.Println("Failed to get ZCAN_SetReference:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(refType), uintptr(pData))
	return uint(ret)
}

// ZCAN_GetReference retrieves device attributes via the ZCAN API.
// See ZLG USB-CAN-FD-B API Manual §6.2.11.
func (zc *ZCAN) ZCAN_GetReference(channelHandle int, refType uint, pData unsafe.Pointer) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetReference")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetReference:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(refType), uintptr(pData))
	return uint(ret)
}

// ReadChannelErrInfo reads the error information for the specified channel.
// Returns raw error data and status code.
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) ReadChannelErrInfo(channelHandle int) ([]byte, uint) {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ReadChannelErrInfo")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ReadChannelErrInfo:", err)
		return nil, 0
	}
	buf := make([]byte, 64)
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return buf, uint(ret)
}

// ZCAN_CHANNEL_STATUS holds the current CAN channel status including
// error counters and bus state. Read with ReadChannelStatus.
// See ZLG USB-CAN-FD-B API Manual.
type ZCAN_CHANNEL_STATUS struct {
	ATYPE           uint8    // Adapter type.
	ErrorCode       uint8    // Error code.
	Passive_Err_Tx  uint8    // Passive error TX count.
	Passive_Err_Rx  uint8    // Passive error RX count.
	EPL_Cnt         uint8    // Error passive limit count.
	Status          uint8    // Bus status flags.
	Warn_Cnt        uint8    // Warning count.
	Err_Cnt         [8]uint8 // Error counters per component.
}

// ReadChannelStatus reads the CAN channel status (bus state, error counters).
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) ReadChannelStatus(channelHandle int) (*ZCAN_CHANNEL_STATUS, uint) {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ReadChannelStatus")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ReadChannelStatus:", err)
		return nil, 0
	}
	status := ZCAN_CHANNEL_STATUS{}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&status)))
	return &status, uint(ret)
}

// GetAvailableDevices returns the number of ZLG CAN devices currently connected.
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) GetAvailableDevices() int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetAvailableDevices")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetAvailableDevices:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc)
	return int(ret)
}

// ZCAN_DEVICE_INFO_EX holds extended device information including version strings.
// See ZLG USB-CAN-FD-B API Manual.
type ZCAN_DEVICE_INFO_EX struct {
	Index        uint32    // Device index.
	DeviceType   uint32    // Device type constant.
	ChannelMask  uint32    // Bitmask of available channels.
	ChannelNum   uint32    // Number of channels.
	Serial       [32]uint8 // Serial number string.
	HardwareVers [48]uint8 // Hardware version string.
	FirmwareVers [48]uint8 // Firmware version string.
	DriverVers   [48]uint8 // Driver version string.
	KernelVers   [48]uint8 // Kernel version string.
	Name         [64]uint8 // Device name string.
}

// GetDeviceInfoEx retrieves extended device information for the given device index.
// Returns the info struct and status code (1 = success).
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) GetDeviceInfoEx(deviceIndex int) (*ZCAN_DEVICE_INFO_EX, int) {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetDeviceInfoEx")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetDeviceInfoEx:", err)
		return nil, 0
	}
	info := ZCAN_DEVICE_INFO_EX{}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceIndex), uintptr(unsafe.Pointer(&info)))
	return &info, int(ret)
}

// SetDeviceChangeCallback registers a callback for device hot-plug events.
// callback: function pointer to the callback.
// See ZLG USB-CAN-FD-B API Manual.
func (zc *ZCAN) SetDeviceChangeCallback(callback uintptr) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_SetDeviceChangeCallback")
	if err != nil {
		fmt.Println("Failed to get ZCAN_SetDeviceChangeCallback:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, callback)
	return uint(ret)
}

// can_start is the internal helper that performs the full channel setup sequence:
// SetCANBaudRate → InitCANFD → StartCAN.
// Used by the convenience methods below.
func can_start(zcan *ZCAN, handle int, ch int) int {
	zcan.SetCANBaudRate(handle, ch, 500000)
	initCfg := ZCAN_CHANNEL_INIT_CONFIG{
		CanType: ZCAN_TYPE_CAN,
	}
	initCfg.Config.AccCode = 0x00000000
	initCfg.Config.AccMask = 0xFFFFFFFF
	initCfg.Config.Mode = 0
	canHandle := zcan.InitCANFD(handle, uint(ch), &initCfg)
	if canHandle != INVALID_CHANNEL_HANDLE {
		if zcan.StartCAN(canHandle) != ZCAN_STATUS_OK {
			return INVALID_CHANNEL_HANDLE
		}
	}
	return canHandle
}

// ---------------------------------------------------------------------------
// Convenience methods — wrap common multi-step flows into single calls.
// ---------------------------------------------------------------------------

// OpenAndStart opens the device, sets baud rate, initializes and starts the CAN channel.
// Returns the channel handle, or INVALID_CHANNEL_HANDLE on failure.
// This combines: OpenDevice → SetCANBaudRate → InitCAN → StartCAN.
//
//	z, _ := zlgcan.NewZCAN(".\\zlgcan_x64\\zlgcan.dll")
//	ch := z.OpenAndStart(zlgcan.ZCAN_USBCAN2, 0, 0, 500000)
//	if ch == zlgcan.INVALID_CHANNEL_HANDLE { ... }
//	defer z.CloseDevice(ch)
func (zc *ZCAN) OpenAndStart(deviceType, deviceIndex int, canIndex uint, baudrate int) int {
	handle := zc.OpenDevice(deviceType, deviceIndex, 0)
	if handle == INVALID_DEVICE_HANDLE {
		return INVALID_CHANNEL_HANDLE
	}

	zc.SetCANBaudRate(handle, int(canIndex), baudrate)

	initCfg := ZCAN_CHANNEL_INIT_CONFIG{
		CanType: ZCAN_TYPE_CAN,
		Config: ZCAN_CHANNEL_CONFIG{
			AccCode: 0x00000000,
			AccMask: 0xFFFFFFFF,
			Mode:    0,
		},
	}
	ch := zc.InitCANFD(handle, canIndex, &initCfg)
	if ch == INVALID_CHANNEL_HANDLE {
		return INVALID_CHANNEL_HANDLE
	}
	if zc.StartCAN(ch) != ZCAN_STATUS_OK {
		return INVALID_CHANNEL_HANDLE
	}
	return ch
}

// OpenAndStartFD opens the device, sets CANFD baud rates, and starts the CANFD channel.
// arbBaudrate: arbitration domain baud rate (e.g. 500000).
// dataBaudrate: data domain baud rate (e.g. 2000000).
// Returns the channel handle, or INVALID_CHANNEL_HANDLE on failure.
func (zc *ZCAN) OpenAndStartFD(deviceType, deviceIndex int, canIndex uint, arbBaudrate, dataBaudrate int) int {
	handle := zc.OpenDevice(deviceType, deviceIndex, 0)
	if handle == INVALID_DEVICE_HANDLE {
		return INVALID_CHANNEL_HANDLE
	}

	ip, err := zc.GetIProperty(handle)
	if err != nil || ip == nil {
		return INVALID_CHANNEL_HANDLE
	}
 basePath := fmt.Sprintf("%d/", canIndex)
	zc.SetValue(ip, basePath+"canfd_abit_baud_rate", fmt.Sprintf("%d", arbBaudrate))
	zc.SetValue(ip, basePath+"canfd_dbit_baud_rate", fmt.Sprintf("%d", dataBaudrate))
	zc.ReleaseIProperty(ip)

	initCfg := ZCAN_CHANNEL_INIT_CONFIG{
		CanType: ZCAN_TYPE_CANFD,
		Config: ZCAN_CHANNEL_CONFIG{
			AccCode: 0x00000000,
			AccMask: 0xFFFFFFFF,
			Mode:    0,
		},
	}
	ch := zc.InitCANFD(handle, canIndex, &initCfg)
	if ch == INVALID_CHANNEL_HANDLE {
		return INVALID_CHANNEL_HANDLE
	}
	if zc.StartCAN(ch) != ZCAN_STATUS_OK {
		return INVALID_CHANNEL_HANDLE
	}
	return ch
}

// TransmitOne sends a single standard CAN frame.
// id: frame ID (11-bit for standard, 29-bit for extended).
// data: payload bytes (max 8).
// extended: true for 29-bit extended frame.
// Returns true if sent successfully.
func (zc *ZCAN) TransmitOne(channelHandle int, id uint32, data []byte, extended bool) bool {
	eff := uint8(0)
	if extended {
		eff = 1
	}

	var tx ZCAN_Transmit_Data
	tx.Frame.GenerateID(id, eff, 0, 0)
	tx.Frame.Dlc = uint8(len(data))
	copy(tx.Frame.Data[:], data)
	tx.Type = 0 // normal send

	sent := zc.Transmit(channelHandle, []ZCAN_Transmit_Data{tx}, 1)
	return sent > 0
}

// TransmitOneFD sends a single CANFD frame.
// id: frame ID.
// data: payload bytes (max 64).
// extended: true for 29-bit extended frame.
// brs: true to enable Bit Rate Switch.
// Returns true if sent successfully.
func (zc *ZCAN) TransmitOneFD(channelHandle int, id uint32, data []byte, extended, brs bool) bool {
	eff := uint8(0)
	if extended {
		eff = 1
	}

	var tx ZCAN_TransmitFD_Data
	tx.Frame.GenerateID(id, eff, 0, 0)
	tx.Frame.Len = uint8(len(data))
	copy(tx.Frame.Data[:], data)
	if brs {
		tx.Frame.GenerateFlags(1, 0, 0)
	}
	tx.Type = 0

	sent := zc.TransmitFD(channelHandle, []ZCAN_TransmitFD_Data{tx}, 1)
	return sent > 0
}

// ReceiveAll drains all pending CAN frames from the buffer (non-blocking).
// Returns the received frames.
func (zc *ZCAN) ReceiveAll(channelHandle int) []ZCAN_Receive_Data {
	n := zc.GetReceiveNum(channelHandle, ZCAN_TYPE_CAN)
	if n == 0 {
		return nil
	}
	msgs, count := zc.Receive(channelHandle, n, 0)
	return msgs[:count]
}

// ReceiveAllFD drains all pending CANFD frames from the buffer (non-blocking).
// Returns the received frames.
func (zc *ZCAN) ReceiveAllFD(channelHandle int) []ZCAN_ReceiveFD_Data {
	n := zc.GetReceiveNum(channelHandle, ZCAN_TYPE_CANFD)
	if n == 0 {
		return nil
	}
	msgs, count := zc.ReceiveFD(channelHandle, n, 0)
	return msgs[:count]
}

// TransmitMulti sends multiple standard CAN frames at once.
// msgs: frames to send (Type defaults to 0 = normal).
// Returns the number of frames successfully sent.
func (zc *ZCAN) TransmitMulti(channelHandle int, msgs []ZCAN_Transmit_Data) uint {
	return zc.Transmit(channelHandle, msgs, uint(len(msgs)))
}

// TransmitMultiFD sends multiple CANFD frames at once.
// Returns the number of frames successfully sent.
func (zc *ZCAN) TransmitMultiFD(channelHandle int, msgs []ZCAN_TransmitFD_Data) uint {
	return zc.TransmitFD(channelHandle, msgs, uint(len(msgs)))
}
