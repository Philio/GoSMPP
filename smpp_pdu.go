// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

import (
	"os"
	"bufio"
	"reflect"
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

// Optional paramater
type pduOptParam struct {
	tag		uint16
	length		uint16
	value		interface{}
}

// Read optional param
func (op *pduOptParam) read(r *bufio.Reader) (err os.Error) {
	// Read first 4 descripter bytes
	p := make([]byte, 4)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	op.tag    = uint16(unpackUint(p[0:2]))
	op.length = uint16(unpackUint(p[2:4]))
	// Read value data
	vp := make([]byte, op.length)
	_, err = r.Read(vp)
	if err != nil {
		return
	}
	// Determine data type of value
	switch op.tag {
		case TAG_SC_INTERFACE_VERSION:
			op.value = vp[0]
	}
	return
}

// Write optional param
func (op *pduOptParam) write(w *bufio.Writer) (err os.Error) {
	// Create byte array
	p := make([]byte, 4 + op.length)
	copy(p[0:2], packUint(uint64(op.tag), 2))
	copy(p[2:4], packUint(uint64(op.length), 2))
	// Determine data type of value
	v := reflect.NewValue(op.value)
	switch t := v.(type) {
		case *reflect.StringValue:
			copy(p[4:op.length], []byte(op.value.(string)))
		case *reflect.BoolValue:
			if op.value.(bool) {
				p[4] = byte(1)
			} else {
				p[4] = byte(0)
			}
		case *reflect.Uint8Value:
			p[4] = byte(op.value.(uint8))
		case *reflect.Uint16Value:
			copy(p[4:6], packUint(uint64(op.value.(uint16)), 2))
		case *reflect.Uint32Value:
			copy(p[4:8], packUint(uint64(op.value.(uint32)), 4))
		case *reflect.Uint64Value:
			copy(p[4:12], packUint(uint64(op.value.(uint64)), 8))
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

// Bind PDU
type pduBind struct {
	header		*pduHeader
	systemId	string
	password	string
	systemType	string
	ifVersion	uint8
	addrTon		SMPPTypeOfNumber
	addrNpi		SMPPNumericPlanIndicator
	addressRange	string
}

// Read Bind PDU
// @todo used for server
func (pdu *pduBind) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write Bind PDU
func (pdu *pduBind) write(w *bufio.Writer) (err os.Error) {
	// Write header
	err = pdu.header.write(w)
	if err != nil {
		err = os.NewError("Bind: Error writing header")
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
		err = os.NewError("Bind: Error writing to buffer")
		return
	}
	// Flush write buffer
	err = w.Flush()
	if err != nil {
		err = os.NewError("Bind: Error flushing write buffer")
	}
	return
}

// Bind Response PDU
type pduBindResp struct {
	header		*pduHeader
	systemId	string
	optional	OptParams
}

// Read Bind Response PDU
func (pdu *pduBindResp) read(r *bufio.Reader) (err os.Error) {
	// Read header
	pdu.header = new(pduHeader)
	err = pdu.header.read(r)
	if err != nil {
		err = os.NewError("Bind Response: Error reading header")
		return
	}
	// Read system id (null terminated string or null)
	line, err := r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind Response: Error reading SMSC system id")
		return
	}
	if len(line) > 1 {
		pdu.systemId = string(line[0:len(line) - 1])
	}
	// Read optional param
	if pdu.header.cmdLength > 16 + uint32(len(pdu.systemId)) + 1 {
		op := new(pduOptParam)
		err = op.read(r)
		if err != nil {
			err = os.NewError("Bind Response: Error reading optional param")
			return
		}
		pdu.optional = OptParams{SMPPOptionalParamTag(op.tag): op.value}
	}
	return
}

// Write Bind Response PDU
// @todo used for server
func (pdu *pduBindResp) write(w *bufio.Writer) (err os.Error) {
	return
}

// Unbind PDU
type pduUnbind struct {
	header		*pduHeader
}

// Read Unbind PDU
func (pdu *pduUnbind) read(r *bufio.Reader) (err os.Error) {
	// Read header
	pdu.header = new(pduHeader)
	err = pdu.header.read(r)
	if err != nil {
		err = os.NewError("Unbind: Error reading header")
	}
	return
}

// Write Unbind PDU
func (pdu *pduUnbind) write(w *bufio.Writer) (err os.Error) {
	// Write header
	err = pdu.header.write(w)
	if err != nil {
		err = os.NewError("Unbind: Error writing header")
	}
	return
}

// Unbind Response PDU
type pduUnbindResp struct {
	header		*pduHeader
}

// Read Unbind Response PDU
func (pdu *pduUnbindResp) read(r *bufio.Reader) (err os.Error) {
	// Read header
	pdu.header = new(pduHeader)
	err = pdu.header.read(r)
	if err != nil {
		err = os.NewError("Unbind Response: Error reading header")
	}
	return
}

// Write Unbind Response PDU
func (pdu *pduUnbindResp) write(w *bufio.Writer) (err os.Error) {
	// Write header
	err = pdu.header.write(w)
	if err != nil {
		err = os.NewError("Unbind Response: Error writing header")
	}
	return
}

// Generic Nack PDU
type pduGenericNack struct {
	header		*pduHeader
}

// Submit SM PDU
type pduSubmitSM struct {
	header		*pduHeader
	serviceType	string
	sourceAddrTon	SMPPTypeOfNumber
	sourceAddrNpi	SMPPNumericPlanIndicator
	sourceAddr	string
	destAddrTon	SMPPTypeOfNumber
	destAddrNpi	SMPPNumericPlanIndicator
	destAddr	string
	esmClass	SMPPEsmClassESME
	protocolId	uint8
	priorityFlag	SMPPPriority
	schedDelTime	string
	validityPeriod	string
	regDelivery	SMPPDelivery
	replaceFlag	uint8
	dataCoding	SMPPDataCoding
	smDefaultMsgId	uint8
	smLength	uint8
	shortMessage	string
	optional	OptParams
	optionalLen	uint32
}

// Read SubmitSM PDU
func (pdu *pduSubmitSM) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write SubmitSM PDU
func (pdu *pduSubmitSM) write(w *bufio.Writer) (err os.Error) {
	// Write header
	err = pdu.header.write(w)
	if err != nil {
		err = os.NewError("SubmitSM: Error writing header")
		return
	}
	// Create byte array the size of the PDU
	p := make([]byte, pdu.header.cmdLength - 16 - pdu.optionalLen)
	pos := 0
	// Copy service type
	if len(pdu.serviceType) > 0 {
		copy(p[pos:len(pdu.serviceType)], []byte(pdu.serviceType))
		pos += len(pdu.serviceType)
	}
	pos ++ // Null terminator
	// Source TON
	p[pos] = byte(pdu.sourceAddrTon)
	pos ++
	// Source NPI
	p[pos] = byte(pdu.sourceAddrNpi)
	pos ++
	// Source Address
	if len(pdu.sourceAddr) > 0 {
		copy(p[pos:pos + len(pdu.sourceAddr)], []byte(pdu.sourceAddr))
		pos += len(pdu.sourceAddr)
	}
	pos ++ // Null terminator
	// Destination TON
	p[pos] = byte(pdu.destAddrTon)
	pos ++
	// Destination NPI
	p[pos] = byte(pdu.destAddrNpi)
	pos ++
	// Destination Address
	if len(pdu.destAddr) > 0 {
		copy(p[pos:pos + len(pdu.destAddr)], []byte(pdu.destAddr))
		pos += len(pdu.destAddr)
	}
	pos ++ // Null terminator
	// ESM Class
	p[pos] = byte(pdu.esmClass)
	pos ++
	// Protocol Id
	p[pos] = byte(pdu.protocolId)
	pos ++
	// Priority Flag
	p[pos] = byte(pdu.priorityFlag)
	pos ++
	// Sheduled Delivery Time
	if len(pdu.schedDelTime) > 0 {
		copy(p[pos:pos + len(pdu.schedDelTime)], []byte(pdu.schedDelTime))
		pos += len(pdu.schedDelTime)
	}
	pos ++ // Null terminator
	// Validity Period
	if len(pdu.validityPeriod) > 0 {
		copy(p[pos:pos + len(pdu.validityPeriod)], []byte(pdu.validityPeriod))
		pos += len(pdu.validityPeriod)
	}
	pos ++ // Null terminator
	// Registered Delivery
	p[pos] = byte(pdu.regDelivery)
	pos ++
	// Replace Flag
	p[pos] = byte(pdu.replaceFlag)
	pos ++
	// Data Coding
	p[pos] = byte(pdu.dataCoding)
	pos ++
	// Default Msg Id
	p[pos] = byte(pdu.smDefaultMsgId)
	pos ++
	// Msg Length
	p[pos] = byte(pdu.smLength)
	pos ++
	// Message
	if len(pdu.shortMessage) > 0 {
		copy(p[pos:pos + len(pdu.shortMessage)], []byte(pdu.shortMessage))
		pos += len(pdu.shortMessage)
	}
	// Write to buffer
	_, err = w.Write(p)
	if err != nil {
		err = os.NewError("SubmitSM: Error writing to buffer")
		return
	}
	// Flush write buffer
	err = w.Flush()
	if err != nil {
		err = os.NewError("SubmitSM: Error flushing write buffer")
		return
	}
	// Optional params
	if len(pdu.optional) > 0 {
		for key, val := range pdu.optional {
			op := new(pduOptParam)
			op.tag = uint16(key)
			op.value = val
			v := reflect.NewValue(val)
			switch t := v.(type) {
				case *reflect.StringValue:
					op.length = uint16(len(val.(string)))
				case *reflect.BoolValue:
					op.length = 1
				case *reflect.Uint8Value:
					op.length = 1
				case *reflect.Uint16Value:
					op.length = 2
				case *reflect.Uint32Value:
					op.length = 4
				case *reflect.Uint64Value:
					op.length = 8
			}
			err = op.write(w)
			if err != nil {
				return
			}
		}
	}
	return
}

// SubmitSM Response PDU
type pduSubmitSMResp struct {
	header		*pduHeader
	messageId	string
}

// Read SubmitSM Response PDU
func (pdu *pduSubmitSMResp) read(r *bufio.Reader) (err os.Error) {
	// Read header
	pdu.header = new(pduHeader)
	err = pdu.header.read(r)
	if err != nil {
		err = os.NewError("SubmitSM Response: Error reading header")
	}
	// Read message id (null terminated string or null)
	line, err := r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("SubmitSM Response: Error reading message id")
		return
	}
	if len(line) > 1 {
		pdu.messageId = string(line[0:len(line) - 1])
	}
	// Check entire packet read
	if pdu.header.cmdLength > uint32(len(line)) + 16 {
		err = os.NewError("SubmitSM Response: Unknown data at end of PDU")
	}
	return
}

// Submit Mutli PDU
type pduSubmitMulti struct {
	header		*pduHeader
	serviceType	string
	sourceAddrTon	uint8
	sourceAddrNpi	uint8
	sourceAddr	string
	numberOfDests	uint8
	destAddrs	[]string
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
