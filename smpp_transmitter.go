// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

import (
	"os"
	"reflect"
)

// Transmitter type
type Transmitter struct {
	smpp
}

func (tx *Transmitter) SubmitSM(dest, msg string, params Params, optional ...OptParams) (msgId string, err os.Error) {
	// Check connected and bound
	if !tx.connected || !tx.bound {
		err = os.NewError("SubmitSM: A bound connection is required to submit a message")
		return
	}
	// Check destination number and message
	if dest == "" {
		err = os.NewError("SubmitSM: A destination number is required and should not be null")
		return
	}
	// Merge params with defaults
	allParams := mergeParams(params, defaultsSubmitSM)
	// Increment sequence number
	tx.sequence ++
	// Create new PDU
	pdu := new(pduSubmitSM)
	// PDU header
	pdu.header = new(pduHeader)
	pdu.header.cmdLength = 34
	pdu.header.cmdId     = CMD_SUBMIT_SM
	pdu.header.cmdStatus = STATUS_ESME_ROK
	pdu.header.sequence  = tx.sequence
	// Mising params cause panic, this provides a clean error/exit
	paramOK := false
	defer func() {
		if !paramOK && recover() != nil {
			err = os.NewError("SubmitSM: Panic, invalid params")
			return
		}
	}()
	// Populate params
	pdu.serviceType     = allParams["serviceType"].(string)
	pdu.sourceAddrTon   = allParams["sourceAddrTon"].(SMPPTypeOfNumber)
	pdu.sourceAddrNpi   = allParams["sourceAddrNpi"].(SMPPNumericPlanIndicator)
	pdu.sourceAddr      = allParams["sourceAddr"].(string)
	pdu.destAddrTon     = allParams["destAddrTon"].(SMPPTypeOfNumber)
	pdu.destAddrNpi     = allParams["destAddrNpi"].(SMPPNumericPlanIndicator)
	pdu.destAddr        = dest
	pdu.esmClass        = allParams["esmClass"].(SMPPEsmClassESME)
	pdu.protocolId      = allParams["protocolId"].(uint8)
	pdu.priorityFlag    = allParams["priorityFlag"].(SMPPPriority)
	pdu.schedDelTime    = allParams["schedDelTime"].(string)
	pdu.validityPeriod  = allParams["validityPeriod"].(string)
	pdu.regDelivery     = allParams["regDelivery"].(SMPPDelivery)
	pdu.replaceFlag     = allParams["replaceFlag"].(uint8)
	pdu.dataCoding      = allParams["dataCoding"].(SMPPDataCoding)
	pdu.smDefaultMsgId  = allParams["smDefaultMsgId"].(uint8)
	pdu.smLength        = uint8(len(msg))
	pdu.shortMessage    = msg
	// Add length of strings to pdu length
	pdu.header.cmdLength += uint32(len(pdu.serviceType))
	pdu.header.cmdLength += uint32(len(pdu.sourceAddr))
	pdu.header.cmdLength += uint32(len(pdu.destAddr))
	pdu.header.cmdLength += uint32(len(pdu.schedDelTime))
	pdu.header.cmdLength += uint32(len(pdu.validityPeriod))
	pdu.header.cmdLength += uint32(len(pdu.shortMessage))
	// Calculate size of optional params
	if len(optional) > 0 && len(optional[0]) > 0 {
		pdu.optional = optional[0]
		for _, val := range optional[0] {
			v := reflect.NewValue(val)
			switch t := v.(type) {
				default:
					err = os.NewError("SubmitSM: Invalid optional param format")
					return
				case *reflect.StringValue:
					pdu.header.cmdLength += uint32(len(val.(string)))
					pdu.optionalLen += uint32(len(val.(string)))
				case *reflect.BoolValue:
					pdu.header.cmdLength ++
					pdu.optionalLen ++
				case *reflect.Uint8Value:
					pdu.header.cmdLength ++
					pdu.optionalLen ++
				case *reflect.Uint16Value:
					pdu.header.cmdLength += 2
					pdu.optionalLen += 2
				case *reflect.Uint32Value:
					pdu.header.cmdLength += 4
					pdu.optionalLen += 4
				case *reflect.Uint64Value:
					pdu.header.cmdLength += 8
					pdu.optionalLen += 8
			}
			// Add 4 bytes for optional param header
			pdu.header.cmdLength += 4
			pdu.optionalLen += 4
		}
	}
	// Params were fine 'disable' the recover
	paramOK = true
	// Send PDU
	err = pdu.write(tx.writer)
	if err != nil {
		return
	}
	// Create SubmitSM Response PDU
	rpdu := new(pduSubmitSMResp)
	// Read PDU data
	err = rpdu.read(tx.reader)
	if err != nil {
		return
	}
	// Validate PDU data
	if rpdu.header.cmdId != CMD_SUBMIT_SM_RESP {
		err = os.NewError("SubmitSM Response: Invalid command")
		return
	}
	if rpdu.header.cmdStatus != STATUS_ESME_ROK {
		err = os.NewError("SubmitSM Response: Error received from SMSC")
		return
	}
	if rpdu.header.sequence != tx.sequence {
		err = os.NewError("SubmitSM Response: Invalid sequence number")
		return
	}
	return rpdu.messageId, nil
}
