// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

import (
	"os"
	"bufio"
	"fmt"
)

// PDU header
type pduHeader struct {
	cmdLength	uint32
	cmdId		SMPPCommand
	cmdStatus	SMPPCommandStatus
	sequence	uint32
}

// Read PDU header
func (hdr *pduHeader) read(r *bufio.Reader) (err os.Error) {
	// Read all 16 header bytes
	p := make([]byte, 16)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	// Convert bytes into header vars
	hdr.cmdLength = uint32(unpackUint(p[0:4]))
	hdr.cmdId     = SMPPCommand(unpackUint(p[4:8]))
	hdr.cmdStatus = SMPPCommandStatus(unpackUint(p[8:12]))
	hdr.sequence  = uint32(unpackUint(p[12:16]))
	return
}

// Write PDU header
func (hdr *pduHeader) write(w *bufio.Writer) (err os.Error) {
	// Convert header into byte array
	p := make([]byte, 16)
	copy(p[0:4],   packUint(uint64(hdr.cmdLength), 4))
	copy(p[4:8],   packUint(uint64(hdr.cmdId), 4))
	copy(p[8:12],  packUint(uint64(hdr.cmdStatus), 4))
	copy(p[12:16], packUint(uint64(hdr.sequence), 4))
	// Write header
	_, err = w.Write(p)
	if err != nil {
		return
	}
	// Flush write buffer
	err = w.Flush()
	return
}

// Bind Transmitter PDU
type pduBindTransmitter struct {
	header		*pduHeader
	systemId	string
	password	string
	systemType	string
	ifVersion	uint8
	addrTon		SMPPTypeOfNumber
	addrNpi		SMPPNumericPlanIndicator
	addressRange	string
}

// Read Bind Transmitter PDU
func (pdu *pduBindTransmitter) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write Bind Transmitter PDU
func (pdu *pduBindTransmitter) write(w *bufio.Writer) (err os.Error) {
	// Write header
	err = pdu.header.write(w)
	if err != nil {
		return
	}
	// Create byte array the size of the PDU
	p := make([]byte, pdu.header.cmdLength - 16)
	pos := 0
	// Copy system id
	if len(pdu.systemId) > 0 {
		copy(p[pos:len(pdu.systemId)], []byte(pdu.systemId))
		pos += len(pdu.systemId)
	}
	pos ++ // Null terminator
	// Copy password
	if len(pdu.password) > 0 {
		copy(p[pos:pos + len(pdu.password)], []byte(pdu.password))
		pos += len(pdu.password)
	}
	pos ++ // Null terminator
	// Copy system type
	if len(pdu.systemType) > 0 {
		copy(p[pos:pos + len(pdu.systemType)], []byte(pdu.systemType))
		pos += len(pdu.systemType)
	}
	pos ++ // Null terminator
	// Add interface version
	p[pos] = byte(pdu.ifVersion)
	pos ++
	// Add TON
	p[pos] = byte(pdu.addrTon)
	pos ++
	// Add NPI
	p[pos] = byte(pdu.addrNpi)
	pos ++
	// Copy address range
	if len(pdu.addressRange) > 0 {
		copy(p[pos:pos + len(pdu.addressRange)], []byte(pdu.addressRange))
		pos += len(pdu.addressRange)
	}
	// Write to buffer
	_, err = w.Write(p)
	if err != nil {
		return
	}
	// Flush write buffer
	err = w.Flush()
	return
}

// Bind Transmitter Response PDU
type pduBindTransmitterResp struct {
	header		*pduHeader
	systemId	string
	ifVersion	uint8		// Optional
}

// Read Bind Transmitter Response PDU
func (pdu *pduBindTransmitterResp) read(r *bufio.Reader) (err os.Error) {
	// Read header
	pdu.header = new(pduHeader)
	err = pdu.header.read(r)
	if err != nil {
		return
	}
	fmt.Printf("Response header: %#v\n", pdu.header)
	p := make([]byte, pdu.header.cmdLength - 16)
	r.Read(p)
	fmt.Printf("Response body: %#v\n", p)
	return
}

// Write Bind Transmitter Response PDU
func (pdu *pduBindTransmitterResp) write(w *bufio.Writer) (err os.Error) {
	return
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
	systemId	string
	ifVersion	uint8		// Optional
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
	systemId	string
	ifVersion	uint8		// Optional
}

// Unbind PDU
type pduUnbind struct {
	header		*pduHeader
}

// Unbind Response PDU
type pduUnbindResp struct {
	header		*pduHeader
}

// Generic Nack PDU
type pduGenericNack struct {
	header		*pduHeader
}

// Submit SM PDU
type pduSubmitSM struct {
	header		*pduHeader
	serviceType	string
	sourceAddrTon	uint8
	sourceAddrNpi	uint8
	sourceAddr	string
	destAddrTon	uint8
	destAddrNpi	uint8
	destAddr	string
	esmClass	uint8
	protocolId	uint8
	priorityFlag	uint8
	schedDelTime	string
	validityPeriod	string
	regDelivery	uint8
	replaceFlag	uint8
	dataCoding	uint8
	smDefaultMsgId	uint8
	smLength	uint8
	shortMessage	string
	userMsgRef	*pduOptParam	// Optional
	sourcePort	*pduOptParam	// Optional
	sourceAddrSub	*pduOptParam	// Optional
	destPort	*pduOptParam	// Optional
	destAddrSub	*pduOptParam	// Optional
	sarMsgRef	*pduOptParam	// Optional
	sarTotalSegs	*pduOptParam	// Optional
	sarSegSeqnum	*pduOptParam	// Optional
	moreMsgsToSend	*pduOptParam	// Optional
	payloadType	*pduOptParam	// Optional
	msgPayload	*pduOptParam	// Optional
	privacyInd	*pduOptParam	// Optional
	callbkNum	*pduOptParam	// Optional
	callbkNumPreInd	*pduOptParam	// Optional
	callbkNumAtag	*pduOptParam	// Optional
	sourceSubaddr	*pduOptParam	// Optional
	destSubaddr	*pduOptParam	// Optional
	userResCode	*pduOptParam	// Optional
	displayTime	*pduOptParam	// Optional
	smsSignal	*pduOptParam	// Optional
	msValidity	*pduOptParam	// Optional
	msMsgWaitFac	*pduOptParam	// Optional
	numOfMsgs	*pduOptParam	// Optional
	alertOnMsgDel	*pduOptParam	// Optional
	langInd		*pduOptParam	// Optional
	itsReplyType	*pduOptParam	// Optional
	itsSessInfo	*pduOptParam	// Optional
	ussdServiceOp	*pduOptParam	// Optional
}

// Submit SM Response PDU
type pduSubmitSMResp struct {
	header		*pduHeader
	messageId	string
}

// Submit Mutli PDU
type pduSubmitMulti struct {
	header		*pduHeader
	serviceType	string
	sourceAddrTon	uint8
	sourceAddrNpi	uint8
	sourceAddr	string
	numberOfDests	uint8
	destAddrs	[]*pduDestAddr
	esmClass	uint8
	protocolId	uint8
	priorityFlag	uint8
	schedDelTime	string
	validityPeriod	string
	regDelivery	uint8
	replaceFlag	uint8
	dataCoding	uint8
	smDefaultMsgId	uint8
	smLength	uint8
	shortMessage	string
	userMsgRef	*pduOptParam	// Optional
	sourcePort	*pduOptParam	// Optional
	sourceAddrSub	*pduOptParam	// Optional
	destPort	*pduOptParam	// Optional
	destAddrSub	*pduOptParam	// Optional
	sarMsgRef	*pduOptParam	// Optional
	sarTotalSegs	*pduOptParam	// Optional
	sarSegSeqnum	*pduOptParam	// Optional
	payloadType	*pduOptParam	// Optional
	msgPayload	*pduOptParam	// Optional
	privacyInd	*pduOptParam	// Optional
	callbkNum	*pduOptParam	// Optional
	callbkNumPreInd	*pduOptParam	// Optional
	callbkNumAtag	*pduOptParam	// Optional
	sourceSubaddr	*pduOptParam	// Optional
	destSubaddr	*pduOptParam	// Optional
	displayTime	*pduOptParam	// Optional
	smsSignal	*pduOptParam	// Optional
	msValidity	*pduOptParam	// Optional
	msMsgWaitFac	*pduOptParam	// Optional
	alertOnMsgDel	*pduOptParam	// Optional
	langInd		*pduOptParam	// Optional
	destFlag	uint8
}

// Submit Multi Response PDU
type pduSubmitMultiResp struct {
	header		*pduHeader
	messageId	string
	noUnsuccess	uint8
	unsuccessSmes	[]*pduUnsuccessSme
}

// Deliver SM PDU
type pduDeliverSM struct {
	header		*pduHeader
	serviceType	string
	sourceAddrTon	uint8
	sourceAddrNpi	uint8
	sourceAddr	string
	destAddrTon	uint8
	destAddrNpi	uint8
	destAddr	string
	esmClass	uint8
	protocolId	uint8
	priorityFlag	uint8
	schedDelTime	string
	validityPeriod	string
	regDelivery	uint8
	replaceFlag	uint8
	dataCoding	uint8
	smDefaultMsgId	uint8
	smLength	uint8
	shortMessage	string
	userMsgRef	*pduOptParam	// Optional
	sourcePort	*pduOptParam	// Optional
	destPort	*pduOptParam	// Optional
	sarMsgRef	*pduOptParam	// Optional
	sarTotalSegs	*pduOptParam	// Optional
	sarSegSeqnum	*pduOptParam	// Optional
	userResCode	*pduOptParam	// Optional
	privacyInd	*pduOptParam	// Optional
	payloadType	*pduOptParam	// Optional
	msgPayload	*pduOptParam	// Optional
	callbkNum	*pduOptParam	// Optional
	sourceSubaddr	*pduOptParam	// Optional
	destSubaddr	*pduOptParam	// Optional
	langInd		*pduOptParam	// Optional
	itsSessInfo	*pduOptParam	// Optional
	netErrorCode	*pduOptParam	// Optional
	messageState	*pduOptParam	// Optional
	recMessageId	*pduOptParam	// Optional
}

// Deliver SM Response PDU
type pduDeliverSMResp struct {
	header		*pduHeader
	messageId	string
}

// Optional paramater
type pduOptParam struct {
	tag		uint16
	length		uint16
	value		interface{}
}

// Destination address
type pduDestAddr struct {
	destFlag	uint8
	destAddr	*pduSMEDestAddr		// Either
	distList	*pduDistributionList	// Or
}	

// SME Destination address
type pduSMEDestAddr struct {
	destAddrTon	uint8
	destAddrNpi	uint8
	destAddr	string
}

// Distribution list
type pduDistributionList struct {
	dlName		string
}

// SME Destination address
type pduUnsuccessSme struct {
	destAddrTon	uint8
	destAddrNpi	uint8
	destAddr	string
	errorCode	uint32
}

// Unpack uint from l bytes (big endian)
func unpackUint(p []byte) (n uint64) {
	l := uint8(len(p))
	for i := uint8(0); i < l; i ++ {
		n |= uint64(p[i]) << ((l - i - 1) * 8)
	}
	return
}

// Pack uint into l bytes (big endian) 
func packUint(n uint64, l uint8) (p []byte) {
	p = make([]byte, l)
	for i := uint8(0); i < l; i ++ {
		p[i] = byte(n >> ((l - i - 1) * 8))
	}
	return
}

// Read an unsigned int of l bytes length
func readUint(r *bufio.Reader, l uint8) (n uint64, err os.Error) {
	// Read bytes
	p := make([]byte, l)
	_, err = r.Read(p)
	if err != nil {
		return 0, err
	}
	// Unpack bytes
	n = unpackUint(p)
	return
}

// Write an unsigned int of l bytes length
func writeUint(w *bufio.Writer, n uint64, l uint8) (err os.Error) {
	// Pack n into bytes
	p := packUint(n, l)
	// Write bytes
	_, err = w.Write(p)
	return
}
