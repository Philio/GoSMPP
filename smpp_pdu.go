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

// PDU interface which all PDU types should implement
type PDU interface {
	// Read the PDU from the buffer
	read(r *bufio.Reader) (err os.Error)
	
	// Write the PDU to the buffer
	write(w *bufio.Writer) (err os.Error)
	
	// Set the packet header
	setHeader(hdr *PDUHeader)
	
	// Get the packet header
	GetHeader() *PDUHeader
	
	// Get the struct
	GetStruct() interface{}
}

// Common PDU functions & fields
type PDUCommon struct {
	Header		*PDUHeader
	Optional	OptParams
	OptionalLen	uint32
}

// Set header
func (pdu *PDUCommon) setHeader(hdr *PDUHeader) {
	pdu.Header = hdr
}

// Get header
func (pdu *PDUCommon) GetHeader() *PDUHeader {
	return pdu.Header
}

// Get Struct
func (pdu *PDUCommon) GetStruct() interface{} {
	return *pdu
}

// Write Optional Params
func (pdu *PDUCommon) writeOptional(w *bufio.Writer) (err os.Error) {
	if len(pdu.Optional) > 0 {
		for key, val := range pdu.Optional {
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

// Bind PDU
type PDUBind struct {
	PDUCommon
	SystemId	string
	Password	string
	SystemType	string
	IfVersion	uint8
	AddrTon		SMPPTypeOfNumber
	AddrNpi		SMPPNumericPlanIndicator
	AddressRange	string
}

// Read Bind PDU
func (pdu *PDUBind) read(r *bufio.Reader) (err os.Error) {
	// Read system id (null terminated string or null)
	line, err := r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind: Error reading system id")
		return
	}
	if len(line) > 1 {
		pdu.SystemId = string(line[0:len(line) - 1])
	}
	// Read Password (null terminated string or null)
	line, err = r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind: Error reading Password")
		return
	}
	if len(line) > 1 {
		pdu.Password = string(line[0:len(line) - 1])
	}
	// Read system type
	line, err = r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind: Error reading system type")
		return
	}
	if len(line) > 1 {
		pdu.SystemType = string(line[0:len(line) - 1])
	}
	// Read interface version
	c, err := r.ReadByte()
	if err != nil {
		err = os.NewError("Bind: Error reading interface version")
		return
	}
	pdu.IfVersion = uint8(c)
	// Read TON
	c, err = r.ReadByte()
	if err != nil {
		err = os.NewError("Bind: Error reading default type of number")
		return
	}
	pdu.AddrTon = SMPPTypeOfNumber(c)
	// Read NPI
	c, err = r.ReadByte()
	if err != nil {
		err = os.NewError("Bind: Error reading default number plan indicator")
		return
	}
	pdu.AddrNpi = SMPPNumericPlanIndicator(c)
	// Read Address range
	line, err = r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind Response: Error reading system type")
		return
	}
	if len(line) > 1 {
		pdu.AddressRange = string(line[0:len(line) - 1])
	}
	return
}

// Write Bind PDU
func (pdu *PDUBind) write(w *bufio.Writer) (err os.Error) {
	// Write Header
	err = pdu.Header.write(w)
	if err != nil {
		err = os.NewError("Bind: Error writing Header")
		return
	}
	// Create byte array the size of the PDU
	p := make([]byte, pdu.Header.CmdLength - pdu.OptionalLen - 16)
	pos := 0
	// Copy system id
	if len(pdu.SystemId) > 0 {
		copy(p[pos:len(pdu.SystemId)], []byte(pdu.SystemId))
		pos += len(pdu.SystemId)
	}
	pos ++ // Null terminator
	// Copy Password
	if len(pdu.Password) > 0 {
		copy(p[pos:pos + len(pdu.Password)], []byte(pdu.Password))
		pos += len(pdu.Password)
	}
	pos ++ // Null terminator
	// Copy system type
	if len(pdu.SystemType) > 0 {
		copy(p[pos:pos + len(pdu.SystemType)], []byte(pdu.SystemType))
		pos += len(pdu.SystemType)
	}
	pos ++ // Null terminator
	// Add interface version
	p[pos] = byte(pdu.IfVersion)
	pos ++
	// Add TON
	p[pos] = byte(pdu.AddrTon)
	pos ++
	// Add NPI
	p[pos] = byte(pdu.AddrNpi)
	pos ++
	// Copy Address range
	if len(pdu.AddressRange) > 0 {
		copy(p[pos:pos + len(pdu.AddressRange)], []byte(pdu.AddressRange))
		pos += len(pdu.AddressRange)
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

// Get Struct
func (pdu *PDUBind) GetStruct() interface{} {
	return *pdu
}

// Bind Response PDU
type PDUBindResp struct {
	PDUCommon
	SystemId	string
}

// Read Bind Response PDU
func (pdu *PDUBindResp) read(r *bufio.Reader) (err os.Error) {
	// Read system id (null terminated string or null)
	line, err := r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("Bind Response: Error reading system id")
		return
	}
	if len(line) > 1 {
		pdu.SystemId = string(line[0:len(line) - 1])
	}
	// Read Optional param
	if pdu.Header.CmdLength > uint32(len(pdu.SystemId)) + 17 {
		op := new(pduOptParam)
		err = op.read(r)
		if err != nil {
			err = os.NewError("Bind Response: Error reading Optional param")
			return
		}
		pdu.Optional = OptParams{SMPPOptionalParamTag(op.tag): op.value}
	}
	return
}

// Write Bind Response PDU
func (pdu *PDUBindResp) write(w *bufio.Writer) (err os.Error) {
	// Write Header
	err = pdu.Header.write(w)
	if err != nil {
		err = os.NewError("Bind Response: Error writing Header")
		return
	}
	// Create byte array the size of the PDU
	p := make([]byte, pdu.Header.CmdLength - pdu.OptionalLen - 16)
	pos := 0
	// Copy system id
	if len(pdu.SystemId) > 0 {
		copy(p[pos:len(pdu.SystemId)], []byte(pdu.SystemId))
		pos += len(pdu.SystemId)
	}
	pos ++ // Null terminator
	// Write to buffer
	_, err = w.Write(p)
	if err != nil {
		err = os.NewError("Bind Response: Error writing to buffer")
		return
	}
	// Flush write buffer
	err = w.Flush()
	if err != nil {
		err = os.NewError("Bind Response: Error flushing write buffer")
	}
	// Optional params
	err = pdu.writeOptional(w)
	if err != nil {
		err = os.NewError("Bind Response: Error writing optional params")
	}
	return
}

// Get Struct
func (pdu *PDUBindResp) GetStruct() interface{} {
	return *pdu
}

// Unbind PDU
type PDUUnbind struct {
	PDUCommon
}

// Read Unbind PDU
func (pdu *PDUUnbind) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write Unbind PDU
func (pdu *PDUUnbind) write(w *bufio.Writer) (err os.Error) {
	// Write Header
	err = pdu.Header.write(w)
	if err != nil {
		err = os.NewError("Unbind: Error writing Header")
	}
	return
}

// Get Struct
func (pdu *PDUUnbind) GetStruct() interface{} {
	return *pdu
}

// Unbind Response PDU
type PDUUnbindResp struct {
	PDUCommon
}

// Read Unbind Response PDU
func (pdu *PDUUnbindResp) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write Unbind Response PDU
func (pdu *PDUUnbindResp) write(w *bufio.Writer) (err os.Error) {
	// Write Header
	err = pdu.Header.write(w)
	if err != nil {
		err = os.NewError("Unbind Response: Error writing Header")
	}
	return
}

// Get Struct
func (pdu *PDUUnbindResp) GetStruct() interface{} {
	return *pdu
}

// Submit SM PDU
type PDUSubmitSM struct {
	PDUCommon
	ServiceType	string
	SourceAddrTon	SMPPTypeOfNumber
	SourceAddrNpi	SMPPNumericPlanIndicator
	SourceAddr	string
	DestAddrTon	SMPPTypeOfNumber
	DestAddrNpi	SMPPNumericPlanIndicator
	DestAddr	string
	EsmClass	SMPPEsmClassESME
	ProtocolId	uint8
	PriorityFlag	SMPPPriority
	SchedDelTime	string
	ValidityPeriod	string
	RegDelivery	SMPPDelivery
	ReplaceFlag	uint8
	DataCoding	SMPPDataCoding
	SmDefaultMsgId	uint8
	SmLength	uint8
	ShortMessage	string
}

// Read SubmitSM PDU
func (pdu *PDUSubmitSM) read(r *bufio.Reader) (err os.Error) {
	return
}

// Write SubmitSM PDU
func (pdu *PDUSubmitSM) write(w *bufio.Writer) (err os.Error) {
	// Write Header
	err = pdu.Header.write(w)
	if err != nil {
		err = os.NewError("SubmitSM: Error writing Header")
		return
	}
	// Create byte array the size of the PDU
	p := make([]byte, pdu.Header.CmdLength - 16 - pdu.OptionalLen)
	pos := 0
	// Copy service type
	if len(pdu.ServiceType) > 0 {
		copy(p[pos:len(pdu.ServiceType)], []byte(pdu.ServiceType))
		pos += len(pdu.ServiceType)
	}
	pos ++ // Null terminator
	// Source TON
	p[pos] = byte(pdu.SourceAddrTon)
	pos ++
	// Source NPI
	p[pos] = byte(pdu.SourceAddrNpi)
	pos ++
	// Source Address
	if len(pdu.SourceAddr) > 0 {
		copy(p[pos:pos + len(pdu.SourceAddr)], []byte(pdu.SourceAddr))
		pos += len(pdu.SourceAddr)
	}
	pos ++ // Null terminator
	// Destination TON
	p[pos] = byte(pdu.DestAddrTon)
	pos ++
	// Destination NPI
	p[pos] = byte(pdu.DestAddrNpi)
	pos ++
	// Destination Address
	if len(pdu.DestAddr) > 0 {
		copy(p[pos:pos + len(pdu.DestAddr)], []byte(pdu.DestAddr))
		pos += len(pdu.DestAddr)
	}
	pos ++ // Null terminator
	// ESM Class
	p[pos] = byte(pdu.EsmClass)
	pos ++
	// Protocol Id
	p[pos] = byte(pdu.ProtocolId)
	pos ++
	// Priority Flag
	p[pos] = byte(pdu.PriorityFlag)
	pos ++
	// Sheduled Delivery Time
	if len(pdu.SchedDelTime) > 0 {
		copy(p[pos:pos + len(pdu.SchedDelTime)], []byte(pdu.SchedDelTime))
		pos += len(pdu.SchedDelTime)
	}
	pos ++ // Null terminator
	// Validity Period
	if len(pdu.ValidityPeriod) > 0 {
		copy(p[pos:pos + len(pdu.ValidityPeriod)], []byte(pdu.ValidityPeriod))
		pos += len(pdu.ValidityPeriod)
	}
	pos ++ // Null terminator
	// Registered Delivery
	p[pos] = byte(pdu.RegDelivery)
	pos ++
	// Replace Flag
	p[pos] = byte(pdu.ReplaceFlag)
	pos ++
	// Data Coding
	p[pos] = byte(pdu.DataCoding)
	pos ++
	// Default Msg Id
	p[pos] = byte(pdu.SmDefaultMsgId)
	pos ++
	// Msg Length
	p[pos] = byte(pdu.SmLength)
	pos ++
	// Message
	if len(pdu.ShortMessage) > 0 {
		copy(p[pos:pos + len(pdu.ShortMessage)], []byte(pdu.ShortMessage))
		pos += len(pdu.ShortMessage)
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
	err = pdu.writeOptional(w)
	if err != nil {
		err = os.NewError("SubmitSM: Error writing optional params")
	}
	return
}

// Get Struct
func (pdu *PDUSubmitSM) GetStruct() interface{} {
	return *pdu
}

// SubmitSM Response PDU
type PDUSubmitSMResp struct {
	PDUCommon
	MessageId	string
}

// Read SubmitSM Response PDU
func (pdu *PDUSubmitSMResp) read(r *bufio.Reader) (err os.Error) {
	// Read message id (null terminated string or null)
	line, err := r.ReadBytes(0x00)
	if err != nil {
		err = os.NewError("SubmitSM Response: Error reading message id")
		return
	}
	if len(line) > 1 {
		pdu.MessageId = string(line[0:len(line) - 1])
	}
	// Check entire packet read
	if pdu.Header.CmdLength > uint32(len(line)) + 16 {
		err = os.NewError("SubmitSM Response: Unknown data at end of PDU")
	}
	return
}

// Write SubmitSM Response PDU
func (pdu *PDUSubmitSMResp) write(r *bufio.Writer) (err os.Error) {
	return
}

// Get Struct
func (pdu *PDUSubmitSMResp) GetStruct() interface{} {
	return *pdu
}

// PDU Header
type PDUHeader struct {
	CmdLength	uint32
	CmdId		SMPPCommand
	CmdStatus	SMPPCommandStatus
	Sequence	uint32
}

// Read PDU Header
func (hdr *PDUHeader) read(r *bufio.Reader) (err os.Error) {
	// Read all 16 Header bytes
	p := make([]byte, 16)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	// Convert bytes into Header vars
	hdr.CmdLength = uint32(unpackUint(p[0:4]))
	hdr.CmdId     = SMPPCommand(unpackUint(p[4:8]))
	hdr.CmdStatus = SMPPCommandStatus(unpackUint(p[8:12]))
	hdr.Sequence  = uint32(unpackUint(p[12:16]))
	return
}

// Write PDU Header
func (hdr *PDUHeader) write(w *bufio.Writer) (err os.Error) {
	// Convert Header into byte array
	p := make([]byte, 16)
	copy(p[0:4],   packUint(uint64(hdr.CmdLength), 4))
	copy(p[4:8],   packUint(uint64(hdr.CmdId), 4))
	copy(p[8:12],  packUint(uint64(hdr.CmdStatus), 4))
	copy(p[12:16], packUint(uint64(hdr.Sequence), 4))
	// Write Header
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

// Read Optional param
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
	if op.length > 0 {
		vp := make([]byte, op.length)
		_, err = r.Read(vp)
		if err != nil {
			return
		}
		// Determine data type of value
		switch op.tag {
			case TAG_ADDITIONAL_STATUS_INFO_TEXT, TAG_RECEIPTED_MESSAGE_ID, TAG_SOURCE_SUBADDRESS, TAG_DEST_SUBADDRESS, TAG_NETWORK_ERROR_CODE, TAG_MESSAGE_PAYLOAD, TAG_CALLBACK_NUM, TAG_CALLBACK_NUM_ATAG, TAG_ITS_SESSION_INFO:
				op.value = string(vp)
			case TAG_DEST_ADDR_SUBUNIT, TAG_SOURCE_ADDR_SUBUNIT, TAG_DEST_NETWORK_TYPE, TAG_SOURCE_NETWORK_TYPE, TAG_DEST_BEARER_TYPE, TAG_SOURCE_BEARER_TYPE, TAG_SOURCE_TELEMATICS_ID, TAG_PAYLOAD_TYPE, TAG_MS_MSG_WAIT_FACILITIES, TAG_PRIVACY_INDICATOR, TAG_USER_RESPONSE_CODE, TAG_LANGUAGE_INDICATOR, TAG_SAR_TOTAL_SEGMENTS, TAG_SAR_SEGMENT_SEQNUM, TAG_SC_INTERFACE_VERSION, TAG_DISPLAY_TIME, TAG_MS_VALIDITY, TAG_DPF_RESULT, TAG_SET_DPF, TAG_MS_AVAILABILITY_STATUS, TAG_DELIVERY_FAILURE_REASON, TAG_MORE_MESSAGES_TO_SEND, TAG_MESSAGE_STATE, TAG_CALLBACK_NUM_PRES_IND, TAG_NUMBER_OF_MESSAGES, TAG_SMS_SIGNAL, TAG_ITS_REPLY_TYPE, TAG_USSD_SERVICE_OP:
				op.value = uint8(vp[0])
			case TAG_DEST_TELEMATICS_ID, TAG_USER_MESSAGE_REFERENCE, TAG_SOURCE_PORT, TAG_DESTINATION_PORT, TAG_SAR_MSG_REF_NUM:
				op.value = uint16(unpackUint(vp))
			case TAG_QOS_TIME_TO_LIVE:
				op.value = uint32(unpackUint(vp))
		}
	} else {
		op.value = nil
	}
	return
}

// Write Optional param
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
		case *reflect.Uint8Value:
			p[4] = byte(op.value.(uint8))
		case *reflect.Uint16Value:
			copy(p[4:6], packUint(uint64(op.value.(uint16)), 2))
		case *reflect.Uint32Value:
			copy(p[4:8], packUint(uint64(op.value.(uint32)), 4))
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
