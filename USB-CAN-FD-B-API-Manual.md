# USB-CAN-FD-B API Library User Manual

> **Version:** V1.3  
> **Date:** 2022.08.08

## Revision History

| Version | Date       | Description |
|---------|------------|-------------|
| V1.0    | 2022.01.12 | First Edition |
| V1.1    | 2022.03.23 | Newly added library function (chapter 3.16~3.20) |
| V1.2    | 2022.04.03 | Improve filtering Configuration; Add APIs (chapter 3.21~3.25) |
| V1.3    | 2022.08.08 | Improve Functionality, Compatible with CANtest/CANPro; Add Chapter 6 |

## Contents

1. [Overview](#1-overview)
2. [Data Structure Definition](#2-data-structure-definition)
3. [APIs Description](#3-apis-description)
4. [Attribute List](#4-attribute-list)
5. [Flow of Using API](#5-flow-of-using-api)
6. [Compatible with ZLG ControlCAN.dll API Library Manual](#6-compatible-with-zlg-controlcandll-api-library-manual)

---

## 1 Overview

If users only use USBCANFD devices for CAN/CANFD bus debugging, you can directly use the provided CANtest software to test data transmission and reception.

If users plan to write software programs for their own products, please carefully read the following instructions and refer to the demo we provided.

**Development Library Files:** `ControlCANFD.lib`, `ControlCANFD.dll`  
**VC Platform Function Declaration File:** `ControlCANFD.h`, `config.h`

> **Note 1:** `ControlCANFD.lib` / `ControlCANFD.dll` rely on the VC2008 runtime, which is typically included in most systems but may need to be installed on very few lean systems.

> **Note 2:** The secondary development interface functions and data structures supported by this device are compatible with ZLG's interface and data structures.

---

## 2 Data Structure Definition

### 2.1 ZCAN_DEVICE_INFO

This structure contains basic information about the device, populated by `ZCAN_GetDeviceInf`.

```c
typedef struct tagZCAN_DEVICE_INFO {
    USHORT hw_Version;
    USHORT fw_Version;
    USHORT dr_Version;
    USHORT in_Version;
    USHORT irq_Num;
    BYTE   can_Num;
    UCHAR  str_Serial_Num[20];
    UCHAR  str_hw_Type[40];
    USHORT reserved[4];
} ZCAN_DEVICE_INFO;
```

| Member | Description |
|--------|-------------|
| `hw_Version` | Hardware version number, in hexadecimal. For example, `0x0100` represents V1.00. |
| `fw_Version` | Firmware version number, hexadecimal. |
| `dr_Version` | Driver version number, hexadecimal. |
| `in_Version` | Interface library version number, hexadecimal. |
| `irq_Num` | The interrupt number used by the board. |
| `can_Num` | Number of channels. |
| `str_Serial_Num` | Serial number of the board, e.g. `"USBCANFD0002"` (including string terminator `'\0'`). |
| `str_hw_Type` | Hardware type. |
| `reserved` | Reserved only, not set. |

### 2.2 ZCAN_CHANNEL_INIT_CONFIG

This structure defines the parameters for channel initialization configuration and must be initialized before calling `ZCAN_InitCAN`.

```c
typedef struct tagZCAN_CHANNEL_INIT_CONFIG {
    UINT can_type; // 0 = TYPE_CAN, 1 = TYPE_CANFD
    union {
        struct {
            UINT acc_code;
            UINT acc_mask;
            UINT reserved;
            BYTE filter;
            BYTE timing0;
            BYTE timing1;
            BYTE mode;
        } can;
        struct {
            UINT   acc_code;
            UINT   acc_mask;
            UINT   abit_timing;
            UINT   dbit_timing;
            UINT   brp;
            BYTE   filter;
            BYTE   mode;
            USHORT pad;
            UINT   reserved;
        } canfd;
    };
} ZCAN_CHANNEL_INIT_CONFIG;
```

**`can_type`**: Device type, `0` = CAN, `1` = CANFD.

#### CAN Context

| Member | Description |
|--------|-------------|
| `acc_code` | Acceptance code. Recommended setting: `0`. |
| `acc_mask` | Mask code. Bits `0` = "relevant bits", bits `1` = "irrelevant bits". Recommended: `0xFFFFFFFF` (receive all). |
| `reserved` | Reserved only, not set. |
| `filter` | Filtering method: `1` = single filtering, `0` = double filtering. |
| `timing0` | Ignore, do not set. |
| `timing1` | Ignore, do not set. |
| `mode` | Working mode: `0` = normal, `1` = listening only. |

#### CANFD Context

| Member | Description |
|--------|-------------|
| `acc_code` | Acceptance code. |
| `acc_mask` | Mask code. |
| `abit_timing` | Ignore, do not set. |
| `dbit_timing` | Ignore, do not set. |
| `brp` | Baud prescaler, set to `0`. |
| `filter` | Filtering method, same as CAN. |
| `mode` | Working mode, same as CAN. |
| `pad` | Data alignment, do not set. |
| `reserved` | Reserved only, not set. |

> **Note:** The baud rate of the device is set by `GetIProperty`. See [Chapter 5.2](#52-sample-code) for details.

### 2.3 can_frame

This structure contains CAN message information.

```c
typedef struct {
    canid_t can_id;  /* 32 bit MAKE_CAN_ID + EFF/RTR/ERR flags */
    BYTE    can_dlc; /* frame payload length in byte (0 .. CAN_MAX_DLEN) */
    BYTE    _pad;    /* padding */
    BYTE    _res0;   /* reserved / padding */
    BYTE    _res1;   /* reserved / padding */
    BYTE    data[CAN_MAX_DLEN] __attribute__((aligned(8)));
} can_frame;
```

| Member | Description |
|--------|-------------|
| `can_id` | Frame ID, 32 bits. Upper 3 bits are flags: bit 31 = extended frame (0=standard, 1=extended), bit 30 = remote frame (0=data, 1=remote), bit 29 = error frame (must be 0). Use `MAKE_CAN_ID` to construct, `GET_ID` to extract. |
| `can_dlc` | Data length. |
| `_pad` | Align, ignore. |
| `_res0` | Reserved only, not set. |
| `_res1` | Reserved only, not set. |
| `data` | Message data, effective length = `can_dlc`. |

### 2.4 canfd_frame

This structure contains CANFD message information.

```c
typedef struct {
    canid_t can_id;  /* 32 bit MAKE_CAN_ID + EFF/RTR/ERR flags */
    BYTE    len;     /* frame payload length in byte */
    BYTE    flags;   /* additional flags for CAN FD, i.e error code */
    BYTE    _res0;   /* reserved / padding */
    BYTE    _res1;   /* reserved / padding */
    BYTE    data[CANFD_MAX_DLEN] __attribute__((aligned(8)));
} canfd_frame;
```

| Member | Description |
|--------|-------------|
| `can_id` | Frame ID, same as [2.3](#23-can_frame). |
| `len` | Data length. |
| `flags` | Additional flags. For CANFD baud rate switch, set to `CANFD_BRS`. |
| `_res0` | Reserved only, not set. |
| `_res1` | Reserved only, not set. |
| `data` | Message data, effective length = `len`. |

### 2.5 ZCAN_Transmit_Data

Contains CAN send message information, used in `ZCAN_Transmit`.

```c
typedef struct tagZCAN_Transmit_Data {
    can_frame frame;
    UINT transmit_type;
} ZCAN_Transmit_Data;
```

| Member | Description |
|--------|-------------|
| `frame` | Message data information, see [2.3](#23-can_frame). |
| `transmit_type` | Sending type: `0` = normal (auto-retry), `1` = single (no retry), `2` = spontaneous self-reception, `3` = single spontaneous self-reception. |

**Sending type descriptions:**

- **Normal sending (`0`):** CAN controller automatically retries until successful, timeout, or bus off.
- **Single sending (`1`):** No automatic retransmission on arbitration loss or error. Used when fixed-time-interval sending is required.
- **Spontaneous self reception (`2`):** Normal transmission with self-reception; sent message can be read from receive buffer.
- **Single spontaneous self reception (`3`):** Single transmission with self-reception, no retry.

### 2.6 ZCAN_TransmitFD_Data

Contains CANFD send message information, used in `ZCAN_TransmitFD`.

```c
typedef struct tagZCAN_TransmitFD_Data {
    canfd_frame frame;
    UINT transmit_type;
} ZCAN_TransmitFD_Data;
```

| Member | Description |
|--------|-------------|
| `frame` | Message data information, see [2.4](#24-canfd_frame). |
| `transmit_type` | Sending type, same as [2.5](#25-zcan_transmit_data). |

### 2.7 ZCAN_Receive_Data

Contains CAN receive message information, used in `ZCAN_Receive`.

```c
typedef struct tagZCAN_Receive_Data {
    can_frame frame;
    UINT64 timestamp; // microseconds
} ZCAN_Receive_Data;
```

| Member | Description |
|--------|-------------|
| `frame` | Message data information, see [2.3](#23-can_frame). |
| `timestamp` | Timestamp in microseconds, based on device startup time. |

### 2.8 ZCAN_ReceiveFD_Data

Contains CANFD receive message information, used in `ZCAN_ReceiveFD`.

```c
typedef struct tagZCAN_ReceiveFD_Data {
    canfd_frame frame;
    UINT64 timestamp; // microseconds
} ZCAN_ReceiveFD_Data;
```

| Member | Description |
|--------|-------------|
| `frame` | Message data information, see [2.4](#24-canfd_frame). |
| `timestamp` | Timestamp in microseconds. |

### 2.9 IProperty

Used to obtain/set device parameter information. For example code, refer to [Chapter 5.2](#52-sample-code).

```c
typedef struct tagIProperty {
    SetValueFunc SetValue;
    GetValueFunc GetValue;
    GetPropertysFunc GetPropertys;
} IProperty;
```

| Member | Description |
|--------|-------------|
| `SetValue` | Set equipment attribute values, see [Chapter 4](#4-attribute-list). |
| `GetValue` | Get attribute values. |
| `GetPropertys` | Return all attributes contained in the device. |

---

## 3 APIs Description

### 3.1 ZCAN_OpenDevice

Opens the device. A device can only be opened once.

```c
DEVICE_HANDLE ZCAN_OpenDevice(UINT device_type, UINT device_index, UINT reserved);
```

| Parameter | Description |
|-----------|-------------|
| `device_type` | Device type, see macro definition in `zlgcan.h`. |
| `device_index` | Device index number. First device = `0`, second = `1`, etc. |
| `reserved` | Reserved only. |

**Returns:** `INVALID_DEVICE_HANDLE` on failure, otherwise device handle.

### 3.2 ZCAN_CloseDevice

Closes the device. Each open has one close.

```c
UINT ZCAN_CloseDevice(DEVICE_HANDLE device_handle);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle returned by `ZCAN_OpenDevice`. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.3 ZCAN_GetDeviceInf

Obtains device information.

```c
UINT ZCAN_GetDeviceInf(DEVICE_HANDLE device_handle, ZCAN_DEVICE_INFO* pInfo);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `pInfo` | Device information structure, see [2.1](#21-zcan_device_info). |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.4 ZCAN_IsDeviceOnLine

Checks whether the device is online.

```c
UINT ZCAN_IsDeviceOnLine(DEVICE_HANDLE device_handle);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |

**Returns:** `STATUS_ONLINE` if online, `STATUS_OFFLINE` if not.

### 3.5 ZCAN_InitCAN

Initializes CAN.

```c
CHANNEL_HANDLE ZCAN_InitCAN(DEVICE_HANDLE device_handle, UINT can_index, ZCAN_CHANNEL_INIT_CONFIG* pInitConfig);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `can_index` | Channel index: channel 0 = `0`, channel 1 = `1`, etc. |
| `pInitConfig` | Initialization structure, see [2.2](#22-zcan_channel_init_config). |

**Returns:** `INVALID_CHANNEL_HANDLE` on failure, otherwise channel handle.

### 3.6 ZCAN_StartCAN

Starts the CAN channel.

```c
UINT ZCAN_StartCAN(CHANNEL_HANDLE channel_handle);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.7 ZCAN_ResetCAN

Resets the CAN channel. Recovery via `ZCAN_StartCAN`.

```c
UINT ZCAN_ResetCAN(CHANNEL_HANDLE channel_handle);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.8 ZCAN_ClearBuffer

Clears the library receive buffer.

```c
UINT ZCAN_ClearBuffer(CHANNEL_HANDLE channel_handle);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.9 ZCAN_Transmit

Sends CAN frames.

```c
UINT ZCAN_Transmit(CHANNEL_HANDLE channel_handle, ZCAN_Transmit_Data* pTransmit, UINT len);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `pTransmit` | Pointer to `ZCAN_Transmit_Data` array. |
| `len` | Number of frames. |

**Returns:** Actual number of successfully sent frames.

### 3.10 ZCAN_TransmitFD

Sends CANFD frames.

```c
UINT ZCAN_TransmitFD(CHANNEL_HANDLE channel_handle, ZCAN_TransmitFD_Data* pTransmit, UINT len);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `pTransmit` | Pointer to `ZCAN_TransmitFD_Data` array. |
| `len` | Number of frames. |

**Returns:** Actual number of successfully sent frames.

### 3.11 ZCAN_GetReceiveNum

Obtains the number of CAN or CANFD messages in the buffer.

```c
UINT ZCAN_GetReceiveNum(CHANNEL_HANDLE channel_handle, BYTE type);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `type` | `0` = CAN, `1` = CANFD. |

**Returns:** Number of frames.

### 3.12 ZCAN_Receive

Receives CAN frames. Recommended to use `ZCAN_GetReceiveNum` first to ensure the buffer has data.

```c
UINT ZCAN_Receive(CHANNEL_HANDLE channel_handle, ZCAN_Receive_Data* pReceive, UINT len, int wait_time);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `pReceive` | Pointer to `ZCAN_Receive_Data` array. |
| `len` | Array length (max frames this call; actual return ≤ this). |
| `wait_time` | Block wait time in ms. `-1` = wait forever (default). |

**Returns:** Actual number of received frames.

### 3.13 ZCAN_ReceiveFD

Receives CANFD frames. Recommended to use `ZCAN_GetReceiveNum` first.

```c
UINT ZCAN_ReceiveFD(CHANNEL_HANDLE channel_handle, ZCAN_ReceiveFD_Data* pReceive, UINT len, int wait_time);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `pReceive` | Pointer to `ZCAN_ReceiveFD_Data` array. |
| `len` | Array length (max frames this call; actual return ≤ this). |
| `wait_time` | Block wait time in ms. `-1` = wait forever (default). |

**Returns:** Actual number of received frames.

### 3.14 GetIProperty

Returns the property configuration interface.

```c
IProperty GetIProperty(DEVICE_HANDLE device_handle);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |

**Returns:** Pointer to `IProperty` (see [2.9](#29-iproperty)). Returns empty on failure.

### 3.15 ReleaseIProperty

Releases the property interface. Pairs with `GetIProperty`.

```c
UINT ReleaseIProperty(IProperty pIProperty);
```

| Parameter | Description |
|-----------|-------------|
| `pIProperty` | Value returned by `GetIProperty`. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.16 ZCAN_SetAbitBaud

Sets the baudrate of the CANFD arbitration domain. Use when attribute `n/canfd_abit_baud_rate` fails.

```c
UINT ZCAN_SetAbitBaud(DEVICE_HANDLE device_handle, UINT can_index, UINT abitbaud);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `can_index` | Channel index: channel 0 = `0`, channel 1 = `1`, etc. |
| `abitbaud` | Arbitration domain baudrate value, see [Attribute List](#4-attribute-list). |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.17 ZCAN_SetDbitBaud

Sets the baudrate of the CANFD data domain. Use when attribute `n/canfd_dbit_baud_rate` fails.

```c
UINT ZCAN_SetDbitBaud(DEVICE_HANDLE device_handle, UINT can_index, UINT dbitbaud);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `can_index` | Channel index: channel 0 = `0`, channel 1 = `1`, etc. |
| `dbitbaud` | Data domain baudrate value, see [Attribute List](#4-attribute-list). |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.18 ZCAN_SetBaudRateCustom

Sets the CANFD custom baudrate. Use when attribute `n/baud_rate_custom` fails.

```c
UINT ZCAN_SetBaudRateCustom(DEVICE_HANDLE device_handle, UINT can_index, char* RateCustom);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `can_index` | Channel index: channel 0 = `0`, channel 1 = `1`, etc. |
| `RateCustom` | Custom baudrate string, see [Attribute List](#4-attribute-list). |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.19 ZCAN_SetCANFDStandard

Sets the CANFD standard type. Use when attribute `n/canfd_standard` fails.

```c
UINT ZCAN_SetCANFDStandard(DEVICE_HANDLE device_handle, UINT can_index, UINT canfd_standard);
```

| Parameter | Description |
|-----------|-------------|
| `device_handle` | Device handle value. |
| `can_index` | Channel index: channel 0 = `0`, channel 1 = `1`, etc. |
| `canfd_standard` | `0` = CANFD ISO, `1` = CANFD BOSCH. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

### 3.20 ZCAN_SetResistanceEnable

> **This function is not used.**

### 3.21 ZCAN_ClearFilter

Clears channel filtering settings. Use when attribute `n/filter_clear` fails. Must be called as part of the filter configuration sequence.

```c
UINT ZCAN_ClearFilter(CHANNEL_HANDLE channel_handle);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

> **Note:** Filter configuration order: clear filter → set mode → set start ID → set end ID → activate filter. Each filter group must be set as a complete sequence.

### 3.22 ZCAN_SetFilterMode

Configures the channel filtering mode. Use when attribute `n/filter_mode` fails.

```c
UINT ZCAN_SetFilterMode(CHANNEL_HANDLE channel_handle, UINT mode);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `mode` | `0` = Standard Frame, `1` = Extended Frame. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

> **Note:** Must be called as part of the filter configuration sequence (see [3.21](#321-zcan_clearfilter)).

### 3.23 ZCAN_SetFilterStartID

Configures the channel filtering start ID. Use when attribute `n/filter_start` fails.

```c
UINT ZCAN_SetFilterStartID(CHANNEL_HANDLE channel_handle, UINT startID);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `startID` | Start ID value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

> **Note:** Must be called as part of the filter configuration sequence (see [3.21](#321-zcan_clearfilter)).

### 3.24 ZCAN_SetFilterEndID

Configures the channel filtering end ID. Use when attribute `n/filter_end` fails.

```c
UINT ZCAN_SetFilterEndID(CHANNEL_HANDLE channel_handle, UINT EndID);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |
| `EndID` | End ID value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

> **Note:** Must be called as part of the filter configuration sequence (see [3.21](#321-zcan_clearfilter)).

### 3.25 ZCAN_AckFilter

Validates channel filtering settings. Use when attribute `n/filter_ack` fails. Must be the final step in the filter configuration sequence.

```c
UINT ZCAN_AckFilter(CHANNEL_HANDLE channel_handle);
```

| Parameter | Description |
|-----------|-------------|
| `channel_handle` | Channel handle value. |

**Returns:** `STATUS_OK` on success, `STATUS_ERR` on failure.

> **Note:** Must be called as the final step in the filter configuration sequence (see [3.21](#321-zcan_clearfilter)).

---

## 4 Attribute List

| Parameter | Path | Value |
|-----------|------|-------|
| **Arbitration domain baudrate** | `n/canfd_abit_baud_rate` | `10000000` (10Mbps), `8000000` (8Mbps), `5000000` (5Mbps), `4000000` (4Mbps), `2000000` (2Mbps), `1000000` (1Mbps), `800000` (800kbps), `500000` (500kbps), `250000` (250kbps), `125000` (125kbps), `100000` (100kbps), `50000` (50kbps) |
| | | *n = channel number (0 = channel 1, 1 = channel 2). Set before `ZCAN_InitCAN`.* |
| **Data domain baudrate** | `n/canfd_dbit_baud_rate` | `5000000` (5Mbps), `4000000` (4Mbps), `2000000` (2Mbps), `1000000` (1Mbps), `800000` (800kbps), `500000` (500kbps), `250000` (250kbps), `125000` (125kbps), `100000` (100kbps) |
| | | *n = channel number. Set before `ZCAN_InitCAN`.* |
| **Custom baudrate** | `n/baud_rate_custom` | Custom baudrate string |
| | | *Set before `ZCAN_InitCAN`.* |
| **Filter mode** | `n/filter_mode` | `"0"` = standard frame, `"1"` = extended frame |
| | | *Set after `ZCAN_InitCAN`.* |
| **Filter start frame** | `n/filter_start` | Hex char, e.g. `"0x00000000"` |
| | | *Set after `ZCAN_InitCAN`.* |
| **Filter end frame** | `n/filter_end` | Hex char, e.g. `"0x00000000"` |
| | | *Set after `ZCAN_InitCAN`.* |
| **Clear filter** | `n/filter_clear` | `"0"` |
| | | *Set after `ZCAN_InitCAN`.* |
| **Filter activate** | `n/filter_ack` | `"0"` |
| | | *Set after `ZCAN_InitCAN`.* |
| **CANFD standard type** | `n/canfd_standard` | `"0"` = CANFD ISO, `"1"` = CANFD BOSCH |
| | | *Set before `ZCAN_InitCAN`.* |

---

## 5 Flow of Using API

### 5.1 Flow

#### General Execution Flow

```
OPEN DEVICE        ZCAN_OpenDevice           [Necessary]
SET BAUDRATE       GetIProperty->SetValue    [Necessary]
INIT CHANNEL       ZCAN_InitCAN              [Necessary]
SET FILTER         (Optional)
START CHANNEL      ZCAN_StartCAN             [Necessary]
DATA OPERATION     ZCAN_Transmit             [Optional]
                   ZCAN_TransmitFD           [Optional]
                   ZCAN_Receive              [Optional]
                   ZCAN_ReceiveFD            [Optional]
                   ZCAN_ResetCAN             [Optional]
CLOSE DEVICE       ZCAN_CloseDevice          [Necessary]
```

#### Filter Setup Flow

Each channel supports up to 64 standard frame filtering groups or 32 extended frame filtering groups.

```
CLEAR FILTER       GetIProperty->SetValue  or  ZCAN_ClearFilter
SET FILTER 1       GetIProperty->SetValue  or  ZCAN_SetFilterMode
                   GetIProperty->SetValue  or  ZCAN_SetFilterStartID
                   GetIProperty->SetValue  or  ZCAN_SetFilterEndID
FILTER ACTIVE      GetIProperty->SetValue  or  ZCAN_AckFilter
```

> **Note:** These functions/attributes for filtering settings must be called in groups; invoking them alone is meaningless.

### 5.2 Sample Code

#### Open & Close Device

```c
m_DevType = ZCAN_USBCANFD_2000;
m_DevIndex = 0;
DWORD Reserved = 0;

m_dev = ZCAN_OpenDevice(m_DevType, m_DevIndex, Reserved);
if (INVALID_DEVICE_HANDLE == m_dev)
{
    MessageBox("open failed");
    return;
}

if (STATUS_OK != ZCAN_CloseDevice(m_dev))
{
    MessageBox("Close failed!");
    return;
}
MessageBox("Close successful!");
```

#### Set Baudrate 1 — Via Interface Property

```c
IProperty* pPro = GetIProperty(m_dev);
if (pPro == NULL)
{
    MessageBox("Property's NULL!");
    return;
}

if (STATUS_OK != pPro->SetValue("0/canfd_abit_baud_rate", "500000"))
{
    MessageBox("Set ch0 rateA failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("0/canfd_dbit_baud_rate", "1000000"))
{
    MessageBox("Set ch0 rateD failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("1/canfd_abit_baud_rate", "500000"))
{
    MessageBox("Set ch1 rateA failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("1/canfd_dbit_baud_rate", "10000000"))
{
    MessageBox("Set ch1 rated failed!");
    ReleaseIProperty(pPro);
    return;
}
```

#### Set Baudrate 2 — Via Direct APIs

```c
if (STATUS_OK != ZCAN_SetAbitBaud(m_dev, 0, 500000))
{
    MessageBox("Set ch0 rateA failed!");
    return;
}

if (STATUS_OK != ZCAN_SetDbitBaud(m_dev, 0, 1000000))
{
    MessageBox("Set ch0 rated failed!");
    return;
}

if (STATUS_OK != ZCAN_SetAbitBaud(m_dev, 1, 500000))
{
    MessageBox("Set ch1 rateA failed!");
    return;
}

if (STATUS_OK != ZCAN_SetDbitBaud(m_dev, 1, 1000000))
{
    MessageBox("Set ch1 rated failed!");
    return;
}
```

#### Set Channel Filter 1 — Via Interface Property

```c
if (STATUS_OK != pPro->SetValue("0/filter_clear", "0"))
{
    MessageBox("clear ch0 filter failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("0/filter_mode", "0"))
{
    MessageBox("set ch0 filter mode failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("0/filter_start", "0x000100"))
{
    MessageBox("set ch0 filter start failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("0/filter_end", "0x000200"))
{
    MessageBox("set ch0 filter end failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != pPro->SetValue("0/filter_ack", "0"))
{
    MessageBox("set ch0 filter ack failed!");
    ReleaseIProperty(pPro);
    return;
}
```

#### Set Channel Filter 2 — Via Direct APIs

```c
if (STATUS_OK != ZCAN_ClearFilter(dev_ch1))
{
    MessageBox("clear ch0 filter failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != ZCAN_SetFilterMode(dev_ch1, 0))
{
    MessageBox("set ch0 filter mode failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != ZCAN_SetFilterStartID(dev_ch1, 0x100))
{
    MessageBox("set ch0 filter start failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != ZCAN_SetFilterEndID(dev_ch1, 0x200))
{
    MessageBox("set ch0 filter end failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_OK != ZCAN_AckFilter(dev_ch1))
{
    MessageBox("set ch0 filter ack failed!");
    ReleaseIProperty(pPro);
    return;
}
```

#### Init and Start Channel

```c
ZCAN_CHANNEL_INIT_CONFIG cfg;
memset(&cfg, 0, sizeof(cfg));
cfg.can_type = TYPE_CANFD;
cfg.canfd.mode = 0;
cfg.canfd.filter = 0;
cfg.canfd.pad = 0;
cfg.canfd.brp = 0;
cfg.canfd.acc_code = 0;
cfg.canfd.acc_mask = 0xffffffff;
cfg.canfd.reserved = 0;

dev_ch1 = ZCAN_InitCAN(m_dev, 0, &cfg);
if (INVALID_CHANNEL_HANDLE == dev_ch1)
{
    MessageBox("Init-CAN0 failed!");
    ReleaseIProperty(pPro);
    return;
}

if (STATUS_ERR == ZCAN_StartCAN(dev_ch1))
{
    MessageBox("Start-CAN0 failed!");
    ReleaseIProperty(pPro);
    return;
}
```

#### Send Frame

```c
ZCAN_Transmit_Data can_data;
can_data.frame.can_id = MAKE_CAN_ID(0x100, 0, 0, 0);
can_data.frame.can_dlc = 8;
for (i = 0; i < can_data.frame.can_dlc; i++)
{
    can_data.frame.data[i] = i;
}
can_data.transmit_type = 0;
if (!ZCAN_Transmit(dev_ch1, &can_data, 1))
{
    MessageBox("send failed");
    return;
}

ZCAN_TransmitFD_Data canfd_data;
canfd_data.frame.can_id = MAKE_CAN_ID(0x200, 0, 0, 0);
canfd_data.frame.len = 64;
for (i = 0; i < canfd_data.frame.len; i++)
{
    canfd_data.frame.data[i] = i;
}
canfd_data.transmit_type = 0;
if (1 != ZCAN_TransmitFD(dev_ch1, &canfd_data, 1))
{
    MessageBox("sendFD failed");
    return;
}
```

#### Receive Frame

```c
ZCAN_Receive_Data pCanObj0[2500];
ZCAN_ReceiveFD_Data pCanObjFD0[2500];

can0_num = ZCAN_GetReceiveNum(dev_ch1, 0);
if (can0_num)
{
    UINT ReadLen = 0;
    ReadLen = ZCAN_Receive(dev_ch1, pCanObj0, can0_num, 50);
    RV_CAN0_NUMS += ReadLen;
    can0_num = 0;
}

can0fd_num = ZCAN_GetReceiveNum(dev_ch1, 1);
if (can0fd_num)
{
    UINT ReadLen = 0;
    ReadLen = ZCAN_ReceiveFD(dev_ch1, pCanObjFD0, can0fd_num, 50);
    RV_CANFDO_NUMS += ReadLen;
    can0fd_num = 0;
}
```

---

## 6 Compatible with ZLG ControlCAN.dll API Library Manual

If this device uses standard CAN, it is compatible with ZLG CANtest and CANPro protocol analysis software. This chapter provides an overview of its data structure and functions.

For detailed instructions on how to use CANtest and CANPro software, please refer to:
- "How to Compatible with the Use of Zhou Ligong CANTest Software"
- "How to Compatible with the Use of Zhou Ligong CANPro Protocol Analysis Platform V1.50.pdf"

To use the `ControlCAN.dll`, `ControlCAN.lib`, and `ControlCAN.h` files mentioned in the document, simply replace the relevant files provided by this driver library with the corresponding file names.

### 6.1 Data Structure Definition

#### 6.1.1 VCI_BOARD_INFO

The structure `VCI_BOARD_INFO` contains the device information of the USB-CAN series interface card. Filled by `VCI_ReadBoardInfo`.

```c
typedef struct _VCI_BOARD_INFO {
    USHORT hw_Version;
    USHORT fw_Version;
    USHORT dr_Version;
    USHORT in_Version;
    USHORT irq_Num;
    BYTE   can_Num;
    CHAR   str_Serial_Num[20];
    CHAR   str_hw_Type[40];
    USHORT Reserved[4];
} VCI_BOARD_INFO, *PVCI_BOARD_INFO;
```

| Member | Description |
|--------|-------------|
| `hw_Version` | Hardware version number (hex). `0x0100` = V1.00. |
| `fw_Version` | Firmware version number (hex). |
| `dr_Version` | Driver version number (hex). |
| `in_Version` | Interface library version number (hex). |
| `irq_Num` | Retention parameter. |
| `can_Num` | Number of CAN channels. |
| `str_Serial_Num` | Serial number of the board. |
| `str_hw_Type` | Hardware type, e.g. `"USBCANFD0002"` (includes `'\0'`). |
| `Reserved` | Reserved. |

#### 6.1.2 VCI_CAN_OBJ

The structure `VCI_CAN_OBJ` is a CAN frame structure. One structure represents one frame. Used in `VCI_Transmit` and `VCI_Receive`.

```c
typedef struct _VCI_CAN_OBJ {
    UINT ID;
    UINT TimeStamp;
    BYTE TimeFlag;
    BYTE SendType;
    BYTE RemoteFlag;
    BYTE ExternFlag;
    BYTE DataLen;
    BYTE Data[8];
    BYTE Reserved[3];
} VCI_CAN_OBJ, *PVCI_CAN_OBJ;
```

| Member | Description |
|--------|-------------|
| `ID` | Frame ID. 32-bit, right-aligned. |
| `TimeStamp` | Time identifier from device power-on. Unit: 0.1ms. |
| `TimeFlag` | `1` = `TimeStamp` is valid. Only meaningful for received frames. |
| `SendType` | `0` = normal (auto-retry for 4 seconds), `1` = single (no retry). |
| `RemoteFlag` | `0` = data frame, `1` = remote frame (data empty). |
| `ExternFlag` | `0` = standard frame (11-bit ID), `1` = extended frame (29-bit ID). |
| `DataLen` | DLC (≤ 8), number of valid bytes in `Data`. |
| `Data[8]` | CAN frame data. Valid bytes constrained by `DataLen`. |
| `Reserved` | Reserved. |

#### 6.1.3 VCI_INIT_CONFIG

Defines the CAN configuration. Filled before `VCI_InitCan`.

```c
typedef struct _INIT_CONFIG {
    DWORD AccCode;
    DWORD AccMask;
    DWORD Reserved;
    UCHAR Filter;
    UCHAR Timing0;
    UCHAR Timing1;
    UCHAR Mode;
} VCI_INIT_CONFIG, *PVCI_INIT_CONFIG;
```

| Member | Description |
|--------|-------------|
| `AccCode` | Acceptance code. Set to `0`. |
| `AccMask` | Mask code. Recommended: `0xFFFFFF` to receive all, or `0`. |
| `Reserved` | Reserved. |
| `Filter` | Not used for this device. |
| `Timing0` | Use `VCI_SetReference` interface to set baudrate. |
| `Timing1` | Use `VCI_SetReference` interface to set baudrate. |
| `Mode` | `0` = normal, `1` = listening, `2` = loopback (spontaneous self-reception). |

### 6.2 API Illustrate

#### 6.2.1 VCI_OpenDevice

Opens the device. Note that one device can only be opened once.

```c
DWORD __stdcall VCI_OpenDevice(DWORD DeviceType, DWORD DeviceInd, DWORD Reserved);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type (see adapter device type macros). |
| `DeviceInd` | Device index. First device = `0`, second = `1`, etc. |
| `Reserved` | Reserved parameter, usually `0`. |

**Returns:** `1` on success, `0` on failure.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

dwRel = VCI_OpenDevice(nDeviceType, nDeviceInd, 0);
if (dwRel != 1)
{
    MessageBox(_T("Opening device failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.2 VCI_CloseDevice

Closes the device.

```c
DWORD __stdcall VCI_CloseDevice(DWORD DeviceType, DWORD DeviceInd);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |

**Returns:** `1` on success, `0` on failure.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
DWORD dwRel;

dwRel = VCI_CloseDevice(nDeviceType, nDeviceInd);
if (dwRel != 1)
{
    MessageBox(_T("Closing device failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.3 VCI_InitCAN

Initializes the specified CAN channel. Multiple calls required for multiple channels.

```c
DWORD __stdcall VCI_InitCAN(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd, PVCI_INIT_CONFIG pInitConfig);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index: CAN1 = `0`, CAN2 = `1`. |
| `pInitConfig` | Initialization parameter structure. |

**Returns:** `1` on success, `0` on failure.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

dwRel = VCI_OpenDevice(nDeviceType, nDeviceInd, 0);
if (dwRel != 1)
{
    MessageBox(_T("Opening device failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}

VCI_INIT_CONFIG vic;
vic.AccCode = 0;
vic.AccMask = 0;
vic.Filter = 0;
vic.Timing0 = 0;
vic.Timing1 = 0;
vic.Mode = 0;

dwRel = VCI_InitCAN(nDeviceType, nDeviceInd, nCANInd, &vic);
if (dwRel != 1)
{
    VCI_CloseDevice(nDeviceType, nDeviceInd);
    MessageBox(_T("Initializing channel failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.4 VCI_ReadBoardInfo

Obtains device information.

```c
DWORD __stdcall VCI_ReadBoardInfo(DWORD DeviceType, DWORD DeviceInd, PVCI_BOARD_INFO pInfo);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `pInfo` | Structure pointer for storing device information. |

**Returns:** `1` on success, `0` on failure.

```c
#include "ControlCan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
DWORD dwRel;

VCI_BOARD_INFO vbi;
dwRel = VCI_ReadBoardInfo(nDeviceType, nDeviceInd, &vbi);
if (dwRel != 1)
{
    MessageBox(_T("Getting board info failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.5 VCI_GetReceiveNum

Retrieves the number of frames received but not yet read in the receive buffer.

```c
ULONG __stdcall VCI_GetReceiveNum(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index: CAN1 = `0`, CAN2 = `1`. |

**Returns:** Number of unread frames. Returns `-1` if device does not exist or USB is disconnected.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

dwRel = VCI_GetReceiveNum(nDeviceType, nDeviceInd, nCANInd);
```

#### 6.2.6 VCI_ClearBuffer

Clears the buffer of the specified CAN channel. Clears both receive and send buffers.

```c
DWORD __stdcall VCI_ClearBuffer(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index. |

**Returns:** `1` on success, `0` on failure.

```c
#include "ControlCan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

dwRel = VCI_ClearBuffer(nDeviceType, nDeviceInd, nCANInd);
```

#### 6.2.7 VCI_StartCAN

Activates a CAN channel. Multiple calls required for multiple channels.

```c
DWORD __stdcall VCI_StartCAN(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index. |

**Returns:** `1` on success, `0` on failure.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

VCI_INIT_CONFIG vic;

if (VCI_OpenDevice(nDeviceType, nDeviceInd, 0) != 1)
{
    MessageBox(_T("Opening device failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}

vic.AccCode = 0;
vic.AccMask = 0;
vic.Filter = 0;
vic.Timing0 = 0;
vic.Timing1 = 0;
vic.Mode = 0;

if (VCI_InitCAN(nDeviceType, nDeviceInd, nCANInd, &vic) != 1)
{
    VCI_CloseDevice(nDeviceType, nDeviceInd);
    MessageBox(_T("Initializing channel failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}

if (VCI_StartCAN(nDeviceType, nDeviceInd, nCANInd) != 1)
{
    VCI_CloseDevice(nDeviceType, nDeviceInd);
    MessageBox(_T("Starting CAN failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.8 VCI_ResetCAN

Resets CAN. Used with `VCI_StartCAN` to restore normal state without re-initialization. For example, when the CAN card enters bus off state.

```c
DWORD __stdcall VCI_ResetCAN(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index. |

**Returns:** `1` on success, `0` on failure.

```c
#include "ControlCan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

dwRel = VCI_ResetCAN(nDeviceType, nDeviceInd, nCANInd);
if (dwRel != 1)
{
    MessageBox(_T("Reset failed!"), _T("Warning"), MB_OK|MB_ICONQUESTION);
    return FALSE;
}
```

#### 6.2.9 VCI_Transmit

Send function. Returns the actual number of frames successfully sent.

```c
ULONG __stdcall VCI_Transmit(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd, PVCI_CAN_OBJ pSend, DWORD Length);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index. |
| `pSend` | Pointer to `VCI_CAN_OBJ` array. |
| `Length` | Number of frames to send. Max 1000; recommended 1 per call for efficiency. |

**Returns:** Actual number of frames sent. `-1` if device does not exist or USB is disconnected.

```c
#include "Controlcan.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

VCI_CAN_OBJ vco[48];
for (int i = 0; i < 48; i++)
{
    vco[i].ID = i;
    vco[i].RemoteFlag = 0;
    vco[i].ExternFlag = 0;
    vco[i].DataLen = 8;
    for (int j = 0; j < 8; j++)
        vco[i].Data[j] = j;
}

dwRel = VCI_Transmit(nDeviceType, nDeviceInd, nCANInd, vco, 48);
```

#### 6.2.10 VCI_Receive

Receive function. Reads data from the receive buffer of the specified CAN channel.

```c
ULONG __stdcall VCI_Receive(DWORD DevType, DWORD DevIndex, DWORD CANIndex, PVCI_CAN_OBJ pReceive, ULONG Len, INT WaitTime);
```

| Parameter | Description |
|-----------|-------------|
| `DevType` | Device type. |
| `DevIndex` | Device index. |
| `CANIndex` | CAN channel index. |
| `pReceive` | Pointer to `VCI_CAN_OBJ` array for reception. **Array must be larger than `Len` to avoid memory errors.** |
| `Len` | Length of receive array. Device has ~2000 frame cache per channel. Recommended: set array and `Len` to ≥ 2000 (e.g. 2500) to prevent overflow. |
| `WaitTime` | Retention parameter. |

**Returns:** Actual number of frames read. `-1` if device does not exist or USB is disconnected.

> **Tip:** Call `VCI_Receive` every 30ms for optimal balance of timeliness and efficiency.

```c
#include "ControlCANFD.h"

int nDeviceType = 41; /* USBCANFD */
int nDeviceInd = 0;   /* First device */
int nCANInd = 0;      /* First channel */
DWORD dwRel;

VCI_CAN_OBJ vco[2500];
dwRel = VCI_Receive(nDeviceType, nDeviceInd, nCANInd, vco, 2500, 0);
if (dwRel > 0)
{
    /* Process incoming frames here */
}
else if (dwRel == -1)
{
    /* USBCANFD device does not exist or USB is offline.
       Call VCI_CloseDevice then VCI_OpenDevice for hot-plugging support. */
}
```

#### 6.2.11 VCI_SetReference

Attribute setting function for baud rate and filter configuration.

```c
DWORD __stdcall VCI_SetReference(DWORD DeviceType, DWORD DeviceInd, DWORD CANInd, DWORD RefType, PVOID pData);
```

| Parameter | Description |
|-----------|-------------|
| `DeviceType` | Device type. |
| `DeviceInd` | Device index. |
| `CANInd` | CAN channel index. |
| `RefType` | Attribute type (see table below). |
| `pData` | Data pointer for the attribute type (see table below). |

**Baudrate setting (`RefType = 0`):**

`pData` points to a `DWORD` value. Correspondence:

| Value | Baudrate |
|-------|----------|
| `0x060003` | 1 Mbps |
| `0x060004` | 800 kbps |
| `0x060007` | 500 kbps |
| `0x1C0008` | 250 kbps |
| `0x1C0011` | 125 kbps |
| `0x160023` | 100 kbps |
| `0x1C002C` | 50 kbps |
| `0x1600B3` | 20 kbps |
| `0x1C00E0` | 10 kbps |
| `0x1C01C1` | 5 kbps |

> **Note:** Only the arbitration domain baudrate is set. Data domain and its baudrate is fixed to 1 Mbps.

**Filter setting (`RefType = 1, 2, 3`):**

| RefType | pData |
|---------|-------|
| `3` | Clear filter. Any value. |
| `1` | Add filter item. Pointer to `VCI_FILTER_RECORD`. |
| `2` | Activate filter. Any value. |

```c
typedef struct _VCI_FILTER_RECORD {
    DWORD ExtFrame; // extended frame or not
    DWORD Start;
    DWORD End;
} VCI_FILTER_RECORD, *PVCI_FILTER_RECORD;
```

### 6.3 Flow of Using API

#### 6.3.1 Flow

```
OPEN DEVICE         VCI_OpenDevice
SET BAUDRATE        VCI_SetReference
INIT CHANNEL        VCI_InitCAN
SET FILTER          (Optional)
START CHANNEL       VCI_StartCAN
SEND CAN FRAME      VCI_Transmit
RECV CAN FRAME      VCI_Receive
RESET CHANNEL       VCI_ResetCAN
CLOSE DEVICE        VCI_CloseDevice
```

#### 6.3.2 Filter Setup Flow

Filter operations must be executed as a unified sequence. Individual commands are invalid on their own. Each channel supports up to 64 standard or 32 extended filtering groups.

```
CLEAR FILTER        VCI_SetReference, RefType = 3
SET FILTER 1        VCI_SetReference, RefType = 1
SET FILTER 2        VCI_SetReference, RefType = 1
SET FILTER N        VCI_SetReference, RefType = 1
FILTER ACTIVE       VCI_SetReference, RefType = 2
```

> **Note:** Filter settings must be called in groups as shown above; invoking them alone is meaningless.
