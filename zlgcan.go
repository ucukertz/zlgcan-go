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

// ZCAN_TYPES
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

// Constants (device types, status codes, etc.)
const (
	ZCAN_STATUS_ERR         = 0
	ZCAN_STATUS_OK          = 1
	ZCAN_STATUS_ONLINE      = 2
	ZCAN_STATUS_OFFLINE     = 3
	ZCAN_STATUS_UNSUPPORTED = 4

	ZCAN_TYPE_CAN   = 0x0
	ZCAN_TYPE_CANFD = 0x1

	INVALID_DEVICE_HANDLE  = 0
	INVALID_CHANNEL_HANDLE = 0
)

type ZCAN_CHANNEL_CONFIG struct {
	AccCode    uint32
	AccMask    uint32
	AbitTiming uint32
	DbitTiming uint32
	Brp        uint32
	Filter     uint8
	Mode       uint8
	Pad        uint16
	Reserved   uint32
}

type ZCAN_CHANNEL_INIT_CONFIG struct {
	CanType uint32
	Config  ZCAN_CHANNEL_CONFIG
}

type ZCAN_CANFD_CHANNEL_INIT_CONFIG = ZCAN_CHANNEL_INIT_CONFIG
type ZCAN_NORMAL_CHANNEL_INIT_CONFIG = ZCAN_CHANNEL_INIT_CONFIG

// Device Info
type ZCAN_DEVICE_INFO struct {
	hw_Version     uint16
	fw_Version     uint16
	dr_Version     uint16
	in_Version     uint16
	irq_Num        uint16
	can_Num        uint8
	str_Serial_Num [20]uint8
	str_hw_Type    [40]uint8
	reserved       [8]uint16
}

func (i *ZCAN_DEVICE_INFO) HwVersion() string {
	return fmt.Sprintf("V%d.%02X", i.hw_Version>>8, i.hw_Version&0xFF)
}

func (i *ZCAN_DEVICE_INFO) FwVersion() string {
	return fmt.Sprintf("V%d.%02X", i.fw_Version>>8, i.fw_Version&0xFF)
}

func (i *ZCAN_DEVICE_INFO) DrVersion() string {
	return fmt.Sprintf("V%d.%02X", i.dr_Version>>8, i.dr_Version&0xFF)
}

func (i *ZCAN_DEVICE_INFO) InVersion() string {
	return fmt.Sprintf("V%d.%02X", i.in_Version>>8, i.in_Version&0xFF)
}

func (i *ZCAN_DEVICE_INFO) IrqNum() uint16 {
	return i.irq_Num
}

func (i *ZCAN_DEVICE_INFO) CanNum() uint8 {
	return i.can_Num
}

func (i *ZCAN_DEVICE_INFO) Serial() string {
	return strings.TrimRight(string(i.str_Serial_Num[:]), "\x00 ")
}

func (i *ZCAN_DEVICE_INFO) HwType() string {
	return strings.TrimRight(string(i.str_hw_Type[:]), "\x00 ")
}

// CAN Frame
type ZCAN_CAN_FRAME struct {
	Id     uint32
	Dlc    uint8
	__pad  uint8
	__res0 uint8
	__res1 uint8
	Data   [8]uint8
}

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

func (f *ZCAN_CAN_FRAME) GetFrameID() uint32 {
	return f.Id & 0x1FFFFFFF
}

func (f *ZCAN_CAN_FRAME) GetFrameEFF() uint8 {
	if (f.Id & 0x80000000) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CAN_FRAME) GetFrameRTR() uint8 {
	if (f.Id & 0x40000000) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CAN_FRAME) GetFrameERR() uint8 {
	if (f.Id & 0x20000000) != 0 {
		return 1
	}
	return 0
}

// CANFD Frame
type ZCAN_CANFD_FRAME struct {
	Id     uint32
	Len    uint8
	Flags  uint8
	__res0 uint8
	__res1 uint8
	Data   [64]uint8
}

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

func (f *ZCAN_CANFD_FRAME) GetFrameID() uint32 {
	return f.Id & 0x1FFFFFFF
}

func (f *ZCAN_CANFD_FRAME) GetFrameEFF() uint8 {
	if (f.Id & 0x80000000) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CANFD_FRAME) GetFrameRTR() uint8 {
	if (f.Id & 0x40000000) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CANFD_FRAME) GetFrameERR() uint8 {
	if (f.Id & 0x20000000) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CANFD_FRAME) GenerateFlags(brs, esi, __res uint8) {
	f.Flags = 0
	if brs != 0 {
		f.Flags |= 0x01
	}
	if esi != 0 {
		f.Flags |= 0x02
	}
}

func (f *ZCAN_CANFD_FRAME) GetFrameBRS() uint8 {
	if (f.Flags & 0x01) != 0 {
		return 1
	}
	return 0
}

func (f *ZCAN_CANFD_FRAME) GetFrameESI() uint8 {
	if (f.Flags & 0x02) != 0 {
		return 1
	}
	return 0
}

// Transmit/Receive Data
type ZCAN_Transmit_Data struct {
	Frame ZCAN_CAN_FRAME
	Type  uint32
}

type ZCAN_Receive_Data struct {
	Frame     ZCAN_CAN_FRAME
	Timestamp uint64
}

type ZCAN_TransmitFD_Data struct {
	Frame ZCAN_CANFD_FRAME
	Type  uint32
}

type ZCAN_ReceiveFD_Data struct {
	Frame     ZCAN_CANFD_FRAME
	Timestamp uint64
}

// IProperty
type ZCAN_IProperty struct {
	SetValue     *[0]byte
	GetValue     *[0]byte
	GetPropertys *[0]byte
}

// ZCAN wrapper
type ZCAN struct {
	dll syscall.Handle
}

// Constructor
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

func (zc *ZCAN) OpenDevice(deviceType int, deviceIndex int, reserved int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_OpenDevice")
	if err != nil {
		fmt.Println("Failed to get ZCAN_OpenDevice:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceType), uintptr(deviceIndex), uintptr(reserved))
	return int(ret)
}

func (zc *ZCAN) CloseDevice(deviceHandle int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_CloseDevice")
	if err != nil {
		fmt.Println("Failed to get ZCAN_CloseDevice:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle))
	return int(ret)
}

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

func (zc *ZCAN) IsDeviceOnLine(deviceHandle int) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_IsDeviceOnLine")
	if err != nil {
		fmt.Println("Failed to get ZCAN_IsDeviceOnLine:", err)
		return -1
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle))
	return int(ret)
}

func (zc *ZCAN) InitCAN(deviceHandle int, canIndex uint, initConfig *ZCAN_CHANNEL_INIT_CONFIG) int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_InitCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_InitCAN:", err)
		return INVALID_CHANNEL_HANDLE
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(deviceHandle), uintptr(canIndex), uintptr(unsafe.Pointer(initConfig)))
	return int(ret)
}

func (zc *ZCAN) InitCANFD(deviceHandle int, canIndex uint, initConfig *ZCAN_CHANNEL_INIT_CONFIG) int {
	return zc.InitCAN(deviceHandle, canIndex, initConfig)
}

func (zc *ZCAN) StartCAN(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_StartCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_StartCAN:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

func (zc *ZCAN) ResetCAN(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ResetCAN")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ResetCAN:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

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

func (zc *ZCAN) ClearBuffer(channelHandle int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ClearBuffer")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ClearBuffer:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle))
	return uint(ret)
}

func (zc *ZCAN) Transmit(channelHandle int, stdMsg []ZCAN_Transmit_Data, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_Transmit")
	if err != nil {
		fmt.Println("Failed to get ZCAN_Transmit:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&stdMsg[0])), uintptr(length))
	return uint(ret)
}

func (zc *ZCAN) TransmitFD(channelHandle int, fdMsg []ZCAN_TransmitFD_Data, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_TransmitFD")
	if err != nil {
		fmt.Println("Failed to get ZCAN_TransmitFD:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(unsafe.Pointer(&fdMsg[0])), uintptr(length))
	return uint(ret)
}

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

func (zc *ZCAN) GetReceiveNum(channelHandle int, canType uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetReceiveNum")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetReceiveNum:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(canType))
	return uint(ret)
}

func (zc *ZCAN) TransmitData(channelHandle int, data unsafe.Pointer, length uint) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_TransmitData")
	if err != nil {
		fmt.Println("Failed to get ZCAN_TransmitData:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(data), uintptr(length))
	return uint(ret)
}

func (zc *ZCAN) ReceiveData(channelHandle int, data unsafe.Pointer, length uint, waitTime int) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_ReceiveData")
	if err != nil {
		fmt.Println("Failed to get ZCAN_ReceiveData:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(data), uintptr(length), uintptr(waitTime))
	return uint(ret)
}

// IProperty functions
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

func (zc *ZCAN) SetValue(ip *ZCAN_IProperty, path, value string) uint {
	cPath := C.CString(path)
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cPath))
	defer C.free(unsafe.Pointer(cValue))
	ret, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(ip.SetValue)), uintptr(unsafe.Pointer(cPath)), uintptr(unsafe.Pointer(cValue)))
	return uint(ret)
}

func (zc *ZCAN) GetValue(ip *ZCAN_IProperty, path string) string {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	ret, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(ip.GetValue)), uintptr(unsafe.Pointer(cPath)))
	return C.GoString((*C.char)(unsafe.Pointer(ret)))
}

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
// RefType=0: baud rate, pData = DWORD timing value (0x060007 = 500kbps).
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
func (zc *ZCAN) ZCAN_GetReference(channelHandle int, refType uint, pData unsafe.Pointer) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetReference")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetReference:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, uintptr(channelHandle), uintptr(refType), uintptr(pData))
	return uint(ret)
}

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

type ZCAN_CHANNEL_STATUS struct {
	ATYPE           uint8
	ErrorCode       uint8
	Passive_Err_Tx  uint8
	Passive_Err_Rx  uint8
	EPL_Cnt         uint8
	Status          uint8
	Warn_Cnt        uint8
	Err_Cnt         [8]uint8
}

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

func (zc *ZCAN) GetAvailableDevices() int {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_GetAvailableDevices")
	if err != nil {
		fmt.Println("Failed to get ZCAN_GetAvailableDevices:", err)
		return 0
	}
	ret, _, _ := syscall.SyscallN(proc)
	return int(ret)
}

type ZCAN_DEVICE_INFO_EX struct {
	Index         uint32
	DeviceType    uint32
	ChannelMask   uint32
	ChannelNum    uint32
	Serial        [32]uint8
	HardwareVers  [48]uint8
	FirmwareVers  [48]uint8
	DriverVers    [48]uint8
	KernelVers    [48]uint8
	Name          [64]uint8
}

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

func (zc *ZCAN) SetDeviceChangeCallback(callback uintptr) uint {
	proc, err := syscall.GetProcAddress(zc.dll, "ZCAN_SetDeviceChangeCallback")
	if err != nil {
		fmt.Println("Failed to get ZCAN_SetDeviceChangeCallback:", err)
		return ZCAN_STATUS_ERR
	}
	ret, _, _ := syscall.SyscallN(proc, callback)
	return uint(ret)
}

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
