// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

import (
	"os"
)

// Transmitter type
type Transmitter struct {
	smppConn
	sequence	uint32
}

// Bind transmitter
func (tx *Transmitter) bind(params []interface{}) (err os.Error) {
	// Sequence number starts at 1
	tx.sequence = 1
	// Create bind PDU
	pdu := new(pduBindTransmitter)
	// PDU header
	pdu.header = new(pduHeader)
	pdu.header.cmdLength = 23 // Min length
	pdu.header.cmdId     = CMD_BIND_TRANSMITTER
	pdu.header.cmdStatus = STATUS_ESME_ROK
	pdu.header.sequence  = tx.sequence
	// System id (username)
	if len(params) > 0 {
		pdu.systemId = params[0].(string)
		pdu.header.cmdLength += uint32(len(pdu.systemId))
	}
	// Password
	if len(params) > 1 {
		pdu.password = params[1].(string)
		pdu.header.cmdLength += uint32(len(pdu.password))
	}
	// System type
	if len(params) > 2 {
		pdu.systemType = params[2].(string)
		pdu.header.cmdLength += uint32(len(pdu.systemType))
	}
	// Interface version
	pdu.ifVersion = SMPP_INTERFACE_VER
	// TON
	if len(params) > 3 {
		pdu.addrTon = params[3].(SMPPTypeOfNumber)
	} else {
		pdu.addrTon = TON_UNKNOWN
	}
	// NPI
	if len(params) > 4 {
		pdu.addrNpi = params[4].(SMPPNumericPlanIndicator)
	} else {
		pdu.addrNpi = NPI_UNKNOWN
	}
	// Address range
	if len(params) > 5 {
		pdu.addressRange = params[5].(string)
		pdu.header.cmdLength += uint32(len(pdu.addressRange))
	}
	// Send PDU
	err = pdu.write(tx.writer)
	if err != nil {
		return
	}
	// Get response
	err = tx.bindResp()
	return
}

// Bind transmitter response
func (tx *Transmitter) bindResp() (err os.Error) {
	// Create bind response PDU
	pdu := new(pduBindTransmitterResp)
	err = pdu.read(tx.reader)
	return
}
