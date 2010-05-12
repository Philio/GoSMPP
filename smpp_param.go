// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

// Default params
var (
	// Bind defaults
	defaultsBind = Params{"systemId": "", "password": "", "systemType": "", "addrTon": SMPPTypeOfNumber(TON_UNKNOWN), "addrNpi": SMPPNumericPlanIndicator(NPI_UNKNOWN), "addressRange": ""}
	
	// SubmitSM defaults
	defaultsSubmitSM = Params{"serviceType": "", "sourceAddrTon": SMPPTypeOfNumber(TON_UNKNOWN), "sourceAddrNpi": SMPPNumericPlanIndicator(NPI_UNKNOWN), "sourceAddr": "", "destAddrTon": SMPPTypeOfNumber(TON_UNKNOWN), "destAddrNpi": SMPPNumericPlanIndicator(NPI_UNKNOWN), "esmClass": SMPPEsmClassESME(ESME_MSG_MODE_DEFAULT), "protocolId":	uint8(0x00), "priorityFlag": SMPPPriority(PRIORITY_NORMAL), "schedDelTime": "", "validityPeriod": "", "regDelivery":	SMPPDelivery(DELIVERY_NONE), "replaceFlag": uint8(0x00), "dataCoding": SMPPDataCoding(CODING_LATIN1), "smDefaultMsgId": uint8(0x00)}
	
	// SubmitMulti defaults
	defaultsSubmitMulti = Params{"serviceType": "", "sourceAddrTon": SMPPTypeOfNumber(TON_UNKNOWN), "sourceAddrNpi": SMPPNumericPlanIndicator(NPI_UNKNOWN), "sourceAddr": "", "destAddrTon": SMPPTypeOfNumber(TON_UNKNOWN), "destAddrNpi": SMPPNumericPlanIndicator(NPI_UNKNOWN), "esmClass": SMPPEsmClassESME(ESME_MSG_MODE_DEFAULT), "protocolId":	uint8(0x00), "priorityFlag": SMPPPriority(PRIORITY_NORMAL), "schedDelTime": "", "validityPeriod": "", "regDelivery": SMPPDelivery(DELIVERY_NONE), "replaceFlag": uint8(0x00), "dataCoding": SMPPDataCoding(CODING_LATIN1), "smDefaultMsgId": uint8(0x00)}
)

// Params definitions
type Params map[string]interface{}
type OptParams map[SMPPOptionalParamTag]interface{}

// Merge params
func mergeParams(params, defaults Params) (res Params) {
	res = defaults
	for key, val := range params {
		res[key] = val
	}
	return
}
