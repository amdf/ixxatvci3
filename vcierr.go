package ixxatvci3

// Code generated vcierr.go DO NOT EDIT.
//////////////////////////////////////////////////////////////////////////
// IXXAT Automation GmbH
//////////////////////////////////////////////////////////////////////////
/**
  VCI error codes.

  @file "vcierr.h"
*/
//////////////////////////////////////////////////////////////////////////
// (C) 2002-2011 IXXAT Automation GmbH, all rights reserved
//////////////////////////////////////////////////////////////////////////

/*****************************************************************************
 facility codes
*****************************************************************************/

const SEV_INFO = 0x40000000  // informational
const SEV_WARN = 0x80000000  // warnings
const SEV_ERROR = 0xC0000000 // errors
const SEV_MASK = 0xC0000000
const SEV_SUCCESS = 0x00000000

const RESERVED_FLAG = 0x10000000
const CUSTOMER_FLAG = 0x20000000

const STATUS_MASK = 0x0000FFFF
const FACILITY_MASK = 0x0FFF0000

// const SEV_STD_INFO = uint32(SEV_INFO | CUSTOMER_FLAG | FACILITY_STD)
// const SEV_STD_WARN = (SEV_WARN | CUSTOMER_FLAG | FACILITY_STD)
// const SEV_STD_ERROR = (SEV_ERROR | CUSTOMER_FLAG | FACILITY_STD)

const FACILITY_VCI = 0x00010000

const SEV_VCI_INFO = (SEV_INFO | CUSTOMER_FLAG | FACILITY_VCI)
const SEV_VCI_WARN = (SEV_WARN | CUSTOMER_FLAG | FACILITY_VCI)
const SEV_VCI_ERROR = (SEV_ERROR | CUSTOMER_FLAG | FACILITY_VCI)

const FACILITY_DAL = 0x00020000
const SEV_DAL_INFO = (SEV_INFO | CUSTOMER_FLAG | FACILITY_DAL)
const SEV_DAL_WARN = (SEV_WARN | CUSTOMER_FLAG | FACILITY_DAL)
const SEV_DAL_ERROR = (SEV_ERROR | CUSTOMER_FLAG | FACILITY_DAL)

const FACILITY_CCL = 0x00030000
const SEV_CCL_INFO = (SEV_INFO | CUSTOMER_FLAG | FACILITY_CCL)
const SEV_CCL_WARN = (SEV_WARN | CUSTOMER_FLAG | FACILITY_CCL)
const SEV_CCL_ERROR = (SEV_ERROR | CUSTOMER_FLAG | FACILITY_CCL)

const FACILITY_BAL = 0x00040000
const SEV_BAL_INFO = (SEV_INFO | CUSTOMER_FLAG | FACILITY_BAL)
const SEV_BAL_WARN = (SEV_WARN | CUSTOMER_FLAG | FACILITY_BAL)
const SEV_BAL_ERROR = (SEV_ERROR | CUSTOMER_FLAG | FACILITY_BAL)

/*##########################################################################*/
/*##                                                                      ##*/
/*##     VCI error codes                                                  ##*/
/*##                                                                      ##*/
/*##########################################################################*/

//
// MessageId: VCI_SUCCESS
//
// MessageText:
//
//  The operation completed successfully.
//
const VCI_SUCCESS = 0x00000000
const VCI_OK = VCI_SUCCESS

//
// MessageId: VCI_E_UNEXPECTED
//
// MessageText:
//
//  Unexpected failure
//
const VCI_E_UNEXPECTED = (SEV_VCI_ERROR | 0x0001)

//
// MessageId: VCI_E_NOT_IMPLEMENTED
//
// MessageText:
//
//  Not implemented
//
const VCI_E_NOT_IMPLEMENTED = (SEV_VCI_ERROR | 0x0002)

//
// MessageId: VCI_E_OUTOFMEMORY
//
// MessageText:
//
//  Not enough storage is available to complete this operation.
//
const VCI_E_OUTOFMEMORY = (SEV_VCI_ERROR | 0x0003)

//
// MessageId: VCI_E_INVALIDARG
//
// MessageText:
//
//  One or more parameters are invalid.
//
const VCI_E_INVALIDARG = (SEV_VCI_ERROR | 0x0004)

//
// MessageId: VCI_E_NOINTERFACE
//
// MessageText:
//
//  The object does not support the requested interface
//
const VCI_E_NOINTERFACE = (SEV_VCI_ERROR | 0x0005)

//
// MessageId: VCI_E_INVPOINTER
//
// MessageText:
//
//  Invalid pointer
//
const VCI_E_INVPOINTER = (SEV_VCI_ERROR | 0x0006)

//
// MessageId: VCI_E_INVHANDLE
//
// MessageText:
//
//  Invalid handle
//
const VCI_E_INVHANDLE = (SEV_VCI_ERROR | 0x0007)

//
// MessageId: VCI_E_ABORT
//
// MessageText:
//
//  Operation aborted
//
const VCI_E_ABORT = (SEV_VCI_ERROR | 0x0008)

//
// MessageId: VCI_E_FAIL
//
// MessageText:
//
//  Unspecified error
//
const VCI_E_FAIL = (SEV_VCI_ERROR | 0x0009)

//
// MessageId: VCI_E_ACCESSDENIED
//
// MessageText:
//
//  Access is denied.
//
const VCI_E_ACCESSDENIED = (SEV_VCI_ERROR | 0x000A)

//
// MessageId: VCI_E_TIMEOUT
//
// MessageText:
//
//  This operation returned because the timeout period expired.
//
const VCI_E_TIMEOUT = (SEV_VCI_ERROR | 0x000B)

//
// MessageId: VCI_E_BUSY
//
// MessageText:
//
//  The requested resource is in use.
//
const VCI_E_BUSY = (SEV_VCI_ERROR | 0x000C)

//
// MessageId: VCI_E_PENDING
//
// MessageText:
//
//  The data necessary to complete this operation is not yet available.
//
const VCI_E_PENDING = (SEV_VCI_ERROR | 0x000D)

//
// MessageId: VCI_E_NO_DATA
//
// MessageText:
//
//  No more data available.
//
const VCI_E_NO_DATA = (SEV_VCI_ERROR | 0x000E)

//
// MessageId: VCI_E_NO_MORE_ITEMS
//
// MessageText:
//
//  No more entries are available from an enumeration operation.
//
const VCI_E_NO_MORE_ITEMS = (SEV_VCI_ERROR | 0x000F)

//
// MessageId: VCI_E_NOTINITIALIZED
//
// MessageText:
//
//  The component is not initialized.
//
const VCI_E_NOT_INITIALIZED = (SEV_VCI_ERROR | 0x0010)

//
// MessageId: VCI_E_ALREADY_INITIALIZED
//
// MessageText:
//
//  An attempt was made to reinitialize an already initialized component.
//
const VCI_E_ALREADY_INITIALIZED = (SEV_VCI_ERROR | 0x00011)

//
// MessageId: VCI_E_RXQUEUE_EMPTY
//
// MessageText:
//
//  Receive queue empty.
//
const VCI_E_RXQUEUE_EMPTY = (SEV_VCI_ERROR | 0x00012)

//
// MessageId: VCI_E_TXQUEUE_FULL
//
// MessageText:
//
//  Transmit queue full.
//
const VCI_E_TXQUEUE_FULL = (SEV_VCI_ERROR | 0x0013)

//
// MessageId: VCI_E_BUFFER_OVERFLOW
//
// MessageText:
//
//  The data was too large to fit into the specified buffer.
//
const VCI_E_BUFFER_OVERFLOW = (SEV_VCI_ERROR | 0x0014)

//
// MessageId: VCI_E_INVALID_STATE
//
// MessageText:
//
//  The component is not in a valid state to perform this request.
//
const VCI_E_INVALID_STATE = (SEV_VCI_ERROR | 0x0015)

//
// MessageId: VCI_E_OBJECT_ALREADY_EXISTS
//
// MessageText:
//
//  The object already exists.
//
const VCI_E_OBJECT_ALREADY_EXISTS = (SEV_VCI_ERROR | 0x0016)

//
// MessageId: VCI_E_INVALID_INDEX
//
// MessageText:
//
//  Invalid index.
//
const VCI_E_INVALID_INDEX = (SEV_VCI_ERROR | 0x0017)

//
// MessageId: VCI_E_END_OF_FILE
//
// MessageText:
//
//  The end-of-file marker has been reached.
//  There is no valid data in the file beyond this marker.
//
const VCI_E_END_OF_FILE = (SEV_VCI_ERROR | 0x0018)

//
// MessageId: VCI_E_DISCONNECTED
//
// MessageText:
//
// Attempt to send a message to a disconnected communication port.
//
const VCI_E_DISCONNECTED = (SEV_VCI_ERROR | 0x0019)

//
// MessageId: VCI_E_WRONG_FLASHFWVERSION
//
// MessageText:
//
// Invalid flash firmware version or version not supported.
// Check driver version and/or update firmware.
//
const VCI_E_WRONG_FLASHFWVERSION = (SEV_VCI_ERROR | 0x001A)
