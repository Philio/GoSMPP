// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

// Optional paramater
type pduOptParam struct {
	tag		uint16
	length		uint16
	value		interface{}
}

// PDU header
type pduHeader struct {
	cmdLength	uint32
	cmdId		uint32
	cmdStatus	uint32
	sequence	uint32
}

// Bind Transmitter PDU
type pduBindTransmitter struct {
	header		*pduHeader
	systemId	string
	password	string
	systemType	string
	ifVersion	uint8
	addrTon		uint8
	addrNpi		uint8
	addressRange	string
}

// Bind Transmitter Response PDU
type pduBindTransmitterResp struct {
	header		*pduHeader
	systemId	string		// Optional
	ifVersion	*pduOptParam	// Optional
}

// Bind Receiver PDU
type pduBindReceiver struct {
	header		*pduHeader
	systemId	string
	password	string
	systemType	string
	ifVersion	uint8
	addrTon		uint8
	addrNpi		uint8
	addressRange	string
}

// Bind Receiver Response PDU
type pduBindReceiverResp struct {
	header		*pduHeader
	systemId	string		// Optional
	ifVersion	*pduOptParam	// Optional
}

// Bind Transceiver PDU
type pduBindTransceiver struct{
	header		*pduHeader
	systemId	string
	password	string
	systemType	string
	ifVersion	uint8
	addrTon		uint8
	addrNpi		uint8
	addressRange	string
}

// Bind Transceiver Response PDU
type pduBindTransceiverResp struct {
	header		*pduHeader
	systemId	string		// Optional
	ifVersion	*pduOptParam	// Optional
}


