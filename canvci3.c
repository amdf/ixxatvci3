#if _WIN32

#include <vcinpl.h>
#include "vciguid.h"
#include "canvci3.h"

#define MAX_MSGID_11BIT 0x7FFU
#define MAX_MSGID_29BIT 0x1FFFFFFFU

typedef struct _CANDEVHANDLES
{
	HANDLE hDevice;       // device handle
	HANDLE hCanCtl;       // controller handle 
	HANDLE hCanChn;       // channel handle
	BYTE uCanOpMode;      // CAN_OPMODE_* at cantype.h
} CANDEVHANDLES, *PCANDEVHANDLES;

#define CAN_DEV_MAX 10
static CANDEVHANDLES can_dev[CAN_DEV_MAX] = { 0 };

static void    DisplayError(HRESULT hResult);

HRESULT CAN_VCI3_SelectDevice(UINT8 bUserSelect, UINT8 uAssignNumber)
{
	HRESULT hResult;
	
	if (uAssignNumber >= CAN_DEV_MAX)
		return VCI_E_INVALIDARG;

	if (bUserSelect == FALSE)
	{
		HANDLE        hEnum;
		VCIDEVICEINFO sInfo;

		hResult = vciEnumDeviceOpen(&hEnum);

		if (hResult == VCI_OK)
		{
			hResult = vciEnumDeviceNext(hEnum, &sInfo);
		}

		vciEnumDeviceClose(hEnum);

		if (hResult == VCI_OK)
		{
			hResult = vciDeviceOpen(&sInfo.VciObjectId, &can_dev[0].hDevice);
		}
	}
	else
	{
		VCIID VciObjectId;
		hResult = vciSelectDeviceDlg(0, &VciObjectId);
		if (hResult == VCI_OK)
		{
			hResult = vciDeviceOpen(&VciObjectId, &can_dev[uAssignNumber].hDevice);
		}
	}

	if (hResult == VCI_OK)
	{
		// default value is 11-bit standard mode
		can_dev[uAssignNumber].uCanOpMode = CAN_OPMODE_STANDARD;
	}

#ifdef _DEBUG
	DisplayError(hResult);
#endif
	return hResult;
}

HRESULT CAN_VCI3_SetOperatingMode(UINT8 uDevNum, BYTE uCanOpMode)
{
	if (uDevNum >= CAN_DEV_MAX)
		return VCI_E_INVALIDARG;

	can_dev[uDevNum].uCanOpMode = uCanOpMode; // CAN_OPMODE_* at cantype.h
	return VCI_OK;
}

HRESULT CAN_VCI3_OpenConnection(UINT8 uDevNum, UINT8 uBtr0, UINT8 uBtr1)
{
	HRESULT hResult;

	if (uDevNum >= CAN_DEV_MAX)
		return VCI_E_INVALIDARG;

	if (can_dev[uDevNum].hDevice != NULL)
	{
		hResult = canChannelOpen(can_dev[uDevNum].hDevice, 0, FALSE, &can_dev[uDevNum].hCanChn);

		if (hResult == VCI_OK)
		{
			UINT16 wRxFifoSize = 1024;
			UINT16 wRxThreshold = 1;
			UINT16 wTxFifoSize = 128;
			UINT16 wTxThreshold = 1;

			hResult = canChannelInitialize(can_dev[uDevNum].hCanChn,
				wRxFifoSize, wRxThreshold,
				wTxFifoSize, wTxThreshold);
		}

		if (hResult == VCI_OK)
		{
			hResult = canChannelActivate(can_dev[uDevNum].hCanChn, TRUE);
		}

		if (hResult == VCI_OK)
		{
			hResult = canControlOpen(can_dev[uDevNum].hDevice, 0, &can_dev[uDevNum].hCanCtl);
		}

		if (hResult == VCI_OK)
		{
			hResult = canControlInitialize(can_dev[uDevNum].hCanCtl, can_dev[uDevNum].uCanOpMode,
				uBtr0, uBtr1);
			if (VCI_E_ACCESSDENIED == hResult) { //initialized by someone else			
				// get current bitrate. is it equal to desired one?
				CANLINESTATUS st;
				canControlGetStatus(can_dev[uDevNum].hCanCtl, &st);
				if (st.bBtReg0 != uBtr0 || st.bBtReg1 != uBtr1)
					hResult = VCI_E_BUSY;
				else {
					hResult = VCI_OK;
				}									
			} else {
				if (hResult == VCI_OK)
				{
					hResult = canControlSetAccFilter(can_dev[uDevNum].hCanCtl, FALSE,
						CAN_ACC_CODE_ALL, CAN_ACC_MASK_ALL);
				}

				if (hResult == VCI_OK)
				{
					hResult = canControlStart(can_dev[uDevNum].hCanCtl, TRUE);
				}
			}
		}
	}
	else
	{
		hResult = VCI_E_INVHANDLE;
	}
#ifdef _DEBUG
	DisplayError(hResult);
#endif
	return hResult;
}

HRESULT CAN_VCI3_OpenConnectionDetectBitrate(UINT8 uDevNum, UINT16 uTimeoutMs, UINT32 uArrayElementCount, BYTE * ArrayBtr0, BYTE * ArrayBtr1, INT32 * pIndexArray) {
	HRESULT hResult;

	if ((uDevNum >= CAN_DEV_MAX) || (NULL == ArrayBtr0) || (NULL == ArrayBtr1) || (NULL == pIndexArray))
	{
		return VCI_E_INVALIDARG;
	}

	if (can_dev[uDevNum].hDevice != NULL)
	{
		hResult = canChannelOpen(can_dev[uDevNum].hDevice, 0, FALSE, &can_dev[uDevNum].hCanChn);

		if (hResult == VCI_OK)
		{
			UINT16 wRxFifoSize = 1024;
			UINT16 wRxThreshold = 1;
			UINT16 wTxFifoSize = 128;
			UINT16 wTxThreshold = 1;

			hResult = canChannelInitialize(can_dev[uDevNum].hCanChn,
				wRxFifoSize, wRxThreshold,
				wTxFifoSize, wTxThreshold);
		}

		if (hResult == VCI_OK)
		{
			hResult = canChannelActivate(can_dev[uDevNum].hCanChn, TRUE);
		}

		if (hResult == VCI_OK)
		{
			hResult = canControlOpen(can_dev[uDevNum].hDevice, 0, &can_dev[uDevNum].hCanCtl);
		}

		if (hResult == VCI_OK)
		{
			hResult = canControlDetectBitrate(can_dev[uDevNum].hCanCtl, uTimeoutMs, uArrayElementCount, ArrayBtr0, ArrayBtr1, pIndexArray);
						
			UINT8 uBtr0, uBtr1;
			if (hResult == VCI_OK) {
				if ((*pIndexArray < (INT32)uArrayElementCount) && (*pIndexArray >= 0))
				{
					uBtr0 = ArrayBtr0[*pIndexArray];
					uBtr1 = ArrayBtr1[*pIndexArray];
				}
			}
			else {
				return hResult;
			}

			hResult = canControlInitialize(can_dev[uDevNum].hCanCtl, can_dev[uDevNum].uCanOpMode,
				uBtr0, uBtr1);

			if (VCI_E_ACCESSDENIED == hResult) { //initialized by someone else			
				// get current bitrate. is it equal to desired one?
				CANLINESTATUS st;
				canControlGetStatus(can_dev[uDevNum].hCanCtl, &st);
				if (st.bBtReg0 != uBtr0 || st.bBtReg1 != uBtr1)
					hResult = VCI_E_BUSY;
				else {
					hResult = VCI_OK;
				}
			}
			else {
				if (hResult == VCI_OK)
				{
					hResult = canControlSetAccFilter(can_dev[uDevNum].hCanCtl, FALSE,
						CAN_ACC_CODE_ALL, CAN_ACC_MASK_ALL);
				}

				if (hResult == VCI_OK)
				{
					hResult = canControlStart(can_dev[uDevNum].hCanCtl, TRUE);
				}
			}
		}
	}
	else
	{
		hResult = VCI_E_INVHANDLE;
	}
#ifdef _DEBUG
	DisplayError(hResult);
#endif
	return hResult;
}

HRESULT CAN_VCI3_TxData(UINT8 uDevNum, UINT32 uMsgId, UINT8 bRtr, BYTE * MsgData, UINT8 uMsgDataSize)
{
	HRESULT hResult;
	CANMSG  sCanMsg = { 0 };
	UINT8   i;
	BOOL bSendExtended = FALSE;

	if (uDevNum >= CAN_DEV_MAX)
	{
		return VCI_E_INVALIDARG;
	}		

	if (uMsgId > MAX_MSGID_11BIT)
	{
		bSendExtended = TRUE;
	}

	sCanMsg.dwTime = 0;
	sCanMsg.dwMsgId = uMsgId & MAX_MSGID_29BIT;
	sCanMsg.uMsgInfo.Bytes.bType = CAN_MSGTYPE_DATA;
	sCanMsg.uMsgInfo.Bytes.bFlags = CAN_MAKE_MSGFLAGS(0, 0, 0, 0, 0);

	if (!bRtr) {
		if (uMsgDataSize > 8)
		{
			return VCI_E_INVALIDARG;
		}
				
		sCanMsg.uMsgInfo.Bits.dlc = uMsgDataSize;

		if (NULL != MsgData) {
			for (i = 0; i < uMsgDataSize; i++)
			{
				sCanMsg.abData[i] = MsgData[i];
			}
		}
	}
	 else
	{		
		sCanMsg.uMsgInfo.Bits.rtr = 1;
	}

	if (bSendExtended)
	{
		sCanMsg.uMsgInfo.Bits.ext = 1;
	}

	hResult = canChannelSendMessage(can_dev[uDevNum].hCanChn, INFINITE, &sCanMsg);

	return hResult;
}

HRESULT CAN_VCI3_RxWaitData(UINT8 uDevNum, UINT32 * uMsgId, UINT8 * bRtr, BYTE * MsgData, UINT8 * uMsgDataSize)
{
	HRESULT hResult;
	CANMSG  sCanMsg;
	if ((NULL == uMsgId) 
		|| (NULL == bRtr)
		|| (NULL == MsgData)
		|| (NULL == uMsgDataSize)
		|| (uDevNum >= CAN_DEV_MAX)
	)
	{
		return VCI_E_INVALIDARG;
	}

	hResult = canChannelReadMessage(can_dev[uDevNum].hCanChn, 100, &sCanMsg);

	if (hResult == VCI_OK)
	{
		if (sCanMsg.uMsgInfo.Bytes.bType == CAN_MSGTYPE_DATA)
		{
			*bRtr = (sCanMsg.uMsgInfo.Bits.rtr == 0) ? 0 : 1;

			UINT8 j;

			*uMsgId = sCanMsg.dwMsgId;
			*uMsgDataSize = sCanMsg.uMsgInfo.Bits.dlc;

			for (j = 0; j < sCanMsg.uMsgInfo.Bits.dlc; j++)
			{
				MsgData[j] = sCanMsg.abData[j];
			}

			return VCI_OK;
		}
		else
		{
			return VCI_E_NO_DATA;
		}
	}
	
	return hResult;
}

HRESULT CAN_VCI3_RxData(UINT8 uDevNum, UINT32 * uMsgId, UINT8 * bRtr, BYTE * MsgData, UINT8 * uMsgDataSize)
{
	HRESULT hResult;
	CANMSG  sRxCanMsg;
	if ((NULL == uMsgId)
		|| (NULL == bRtr)
		|| (NULL == MsgData)
		|| (NULL == uMsgDataSize)
		|| (uDevNum >= CAN_DEV_MAX)
		)
	{
		return VCI_E_INVALIDARG;
	}

	hResult = canChannelPeekMessage(can_dev[uDevNum].hCanChn, &sRxCanMsg);

	if (hResult == VCI_OK)
	{
		if (sRxCanMsg.uMsgInfo.Bytes.bType == CAN_MSGTYPE_DATA)
		{
			*bRtr = (sRxCanMsg.uMsgInfo.Bits.rtr == 0) ? 0 : 1;

			UINT8 j;

			*uMsgId = sRxCanMsg.dwMsgId;
			*uMsgDataSize = sRxCanMsg.uMsgInfo.Bits.dlc;

			for (j = 0; j < sRxCanMsg.uMsgInfo.Bits.dlc; j++)
			{
				MsgData[j] = sRxCanMsg.abData[j];
			}

			return VCI_OK;
		}
		else
		{
			return VCI_E_NO_DATA;
		}
	}

	return hResult;
}

HRESULT CAN_VCI3_GetStatus(UINT8 uDevNum, PCANCHANSTATUS pCanStat)
{
	if ( (NULL == pCanStat)
		|| (uDevNum >= CAN_DEV_MAX)
		)
	{
		return VCI_E_INVALIDARG;
	}

	return canChannelGetStatus(can_dev[uDevNum].hCanChn, pCanStat);
}

HRESULT CAN_VCI3_CloseDevice(UINT8 uDevNum)
{
	if (uDevNum >= CAN_DEV_MAX)		
	{
		return VCI_E_INVALIDARG;
	}

	canControlReset(can_dev[uDevNum].hCanCtl);
	canChannelClose(can_dev[uDevNum].hCanChn);
	canControlClose(can_dev[uDevNum].hCanCtl);

	vciDeviceClose(can_dev[uDevNum].hDevice);

	return VCI_OK;
}

static void DisplayError(HRESULT hResult)
{
	char szError[VCI_MAX_ERRSTRLEN];

	if (hResult != NO_ERROR)
	{
		if (hResult == -1)
			hResult = GetLastError();

		szError[0] = 0;
		vciFormatError(hResult, szError, VCI_MAX_ERRSTRLEN);
		MessageBoxA(NULL, szError, "VCI3", MB_OK | MB_ICONSTOP);
	}
}

void CAN_VCI3_FormatError(HRESULT hrError, PCHAR pszText, UINT32 dwSize) {
	vciFormatError(hrError, pszText, dwSize);
}

#endif