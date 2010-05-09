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

func (tx *Transmitter) SubmitSM(dest, msg string, params Params, optional ...OptParams) (sequence uint32, msgId string, err os.Error) {
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
	pdu := new(PDUSubmitSM)
	// PDU header
	pdu.Header = new(PDUHeader)
	pdu.Header.CmdLength = 34
	pdu.Header.CmdId     = CMD_SUBMIT_SM
	pdu.Header.CmdStatus = STATUS_ESME_ROK
	pdu.Header.Sequence  = tx.sequence
	// Mising params cause panic, this provides a clean error/exit
	paramOK := false
	defer func() {
		if !paramOK && recover() != nil {
			err = os.NewError("SubmitSM: Panic, invalid params")
			return
		}
	}()
	// Populate params
	pdu.ServiceType     = allParams["serviceType"].(string)
	pdu.SourceAddrTon   = allParams["sourceAddrTon"].(SMPPTypeOfNumber)
	pdu.SourceAddrNpi   = allParams["sourceAddrNpi"].(SMPPNumericPlanIndicator)
	pdu.SourceAddr      = allParams["sourceAddr"].(string)
	pdu.DestAddrTon     = allParams["destAddrTon"].(SMPPTypeOfNumber)
	pdu.DestAddrNpi     = allParams["destAddrNpi"].(SMPPNumericPlanIndicator)
	pdu.DestAddr        = dest
	pdu.EsmClass        = allParams["esmClass"].(SMPPEsmClassESME)
	pdu.ProtocolId      = allParams["protocolId"].(uint8)
	pdu.PriorityFlag    = allParams["priorityFlag"].(SMPPPriority)
	pdu.SchedDelTime    = allParams["schedDelTime"].(string)
	pdu.ValidityPeriod  = allParams["validityPeriod"].(string)
	pdu.RegDelivery     = allParams["regDelivery"].(SMPPDelivery)
	pdu.ReplaceFlag     = allParams["replaceFlag"].(uint8)
	pdu.DataCoding      = allParams["dataCoding"].(SMPPDataCoding)
	pdu.SmDefaultMsgId  = allParams["smDefaultMsgId"].(uint8)
	pdu.SmLength        = uint8(len(msg))
	pdu.ShortMessage    = msg
	// Add length of strings to pdu length
	pdu.Header.CmdLength += uint32(len(pdu.ServiceType))
	pdu.Header.CmdLength += uint32(len(pdu.SourceAddr))
	pdu.Header.CmdLength += uint32(len(pdu.DestAddr))
	pdu.Header.CmdLength += uint32(len(pdu.SchedDelTime))
	pdu.Header.CmdLength += uint32(len(pdu.ValidityPeriod))
	pdu.Header.CmdLength += uint32(len(pdu.ShortMessage))
	// Calculate size of optional params
	if len(optional) > 0 && len(optional[0]) > 0 {
		pdu.Optional = optional[0]
		for _, val := range optional[0] {
			v := reflect.NewValue(val)
			switch t := v.(type) {
				default:
					err = os.NewError("SubmitSM: Invalid optional param format")
					return
				case *reflect.StringValue:
					pdu.Header.CmdLength += uint32(len(val.(string)))
					pdu.OptionalLen += uint32(len(val.(string)))
				case *reflect.Uint8Value:
					pdu.Header.CmdLength ++
					pdu.OptionalLen ++
				case *reflect.Uint16Value:
					pdu.Header.CmdLength += 2
					pdu.OptionalLen += 2
				case *reflect.Uint32Value:
					pdu.Header.CmdLength += 4
					pdu.OptionalLen += 4
			}
			// Add 4 bytes for optional param header
			pdu.Header.CmdLength += 4
			pdu.OptionalLen += 4
		}
	}
	// Params were fine 'disable' the recover
	paramOK = true
	// Send PDU
	err = pdu.write(tx.writer)
	if err != nil {
		return
	}
	// If not async get the response
	if tx.async {
		sequence = tx.sequence
	} else {
		rpdu, err := tx.GetResp(CMD_SUBMIT_SM_RESP, tx.sequence)
		if err != nil {
			return
		}
		//msgId = rpdu.MessageId
	}
	return
}
