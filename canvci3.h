#ifndef _CANVCI3_H_
#define _CANVCI3_H_

HRESULT CAN_VCI3_SelectDevice(UINT8 bUserSelect, UINT8 uAssignNumber);
HRESULT CAN_VCI3_SetOperatingMode(UINT8 uDevNum, BYTE uCanOpMode);
HRESULT CAN_VCI3_OpenConnection(UINT8 uDevNum, UINT8 uBtr0, UINT8 uBtr1);
HRESULT CAN_VCI3_OpenConnectionDetectBitrate(UINT8 uDevNum, UINT16 uTimeoutMs, UINT32 uArrayElementCount, BYTE * ArrayBtr0, BYTE * ArrayBtr1, INT32 * pIndexArray);
HRESULT CAN_VCI3_TxData(UINT8 uDevNum, UINT32 uMsgId, UINT8 bRtr, BYTE * MsgData, UINT8 uMsgDataSize);
HRESULT CAN_VCI3_RxWaitData(UINT8 uDevNum, UINT32 * uMsgId, UINT8 * bRtr, BYTE * MsgData, UINT8 * uMsgDataSize);
HRESULT CAN_VCI3_RxData(UINT8 uDevNum, UINT32 * uMsgId, UINT8 * bRtr, BYTE * MsgData, UINT8 * uMsgDataSize);
HRESULT CAN_VCI3_GetStatus(UINT8 uDevNum, PCANCHANSTATUS pCanStat);
HRESULT CAN_VCI3_CloseDevice(UINT8 uDevNum);
void CAN_VCI3_FormatError(HRESULT hrError, PCHAR pszText, UINT32 dwSize);

#endif //_CANVCI3_H_
