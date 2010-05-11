// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

// Imports
import (
	"os"
	"net"
	"bufio"
	"strconv"
)

// Used for all outbound connections
type smpp struct {
	conn		net.Conn
	reader		*bufio.Reader
	writer		*bufio.Writer
	connected	bool
	bound		bool
	async		bool
	sequence	uint32
}

// Connect to server
func (smpp *smpp) connect(host string, port int) (err os.Error) {
	// Create TCP connection
	smpp.conn, err = net.Dial("tcp", "", host + ":" + strconv.Itoa(port))
	if err != nil {
		return
	}
	smpp.connected = true
	// Setup buffered reader/writer
	smpp.reader = bufio.NewReader(smpp.conn)
	smpp.writer = bufio.NewWriter(smpp.conn)
	return
}

// Close connection
func (smpp *smpp) close() (err os.Error) {
	err = smpp.conn.Close()
	smpp.connected = false	
	return
}

// Send bind request (called via NewTransmitter/NewReceiver/NewTransceiver) always synchronous
func (smpp *smpp) bind(cmd, rcmd SMPPCommand, params Params) (err os.Error) {
	// Sequence number starts at 1
	smpp.sequence ++
	// PDU header
	hdr := new(PDUHeader)
	hdr.CmdLength = 23 // Min length
	hdr.CmdId     = cmd
	hdr.CmdStatus = STATUS_ESME_ROK
	hdr.Sequence  = smpp.sequence
	// Create bind PDU
	pdu := new(PDUBind)
	// Mising params cause panic, this provides a clean error/exit
	paramOK := false
	defer func() {
		if !paramOK && recover() != nil {
			err = os.NewError("Bind: Panic, invalid params")
			return
		}
	}()
	// Populate params
	pdu.SystemId     = params["systemId"].(string)
	pdu.Password     = params["password"].(string)
	pdu.SystemType   = params["systemType"].(string)
	pdu.IfVersion    = SMPP_INTERFACE_VER
	pdu.AddrTon      = params["addrTon"].(SMPPTypeOfNumber)
	pdu.AddrNpi      = params["addrNpi"].(SMPPNumericPlanIndicator)
	pdu.AddressRange = params["addressRange"].(string)
	// Add length of strings to pdu length
	hdr.CmdLength += uint32(len(pdu.SystemId))
	hdr.CmdLength += uint32(len(pdu.Password))
	hdr.CmdLength += uint32(len(pdu.SystemType))
	hdr.CmdLength += uint32(len(pdu.AddressRange))
	// Params were fine 'disable' the recover
	paramOK = true
	// Send PDU
	pdu.setHeader(hdr)
	err = pdu.write(smpp.writer)
	if err != nil {
		return
	}
	// Get response
	_, err = smpp.GetResp(rcmd, smpp.sequence)
	return
}

// Set async commands on/offer (trancsceiver is always async)
func (smpp *smpp) Async(async bool) {
	smpp.async = async
}

// Send unbind request
func (smpp *smpp) Unbind() (sequence uint32, err os.Error) {
	// Check connected and bound
	if !smpp.connected || !smpp.bound {
		err = os.NewError("Unbind: A bound connection is required to unbind")
		return
	}
	// Increment sequence number
	smpp.sequence ++
	// PDU header
	hdr := new(PDUHeader)
	hdr.CmdLength = 16
	hdr.CmdId     = CMD_UNBIND
	hdr.CmdStatus = STATUS_ESME_ROK
	hdr.Sequence  = smpp.sequence
	// Create bind PDU
	pdu := new(PDUUnbind)
	pdu.setHeader(hdr)
	// Send PDU
	err = pdu.write(smpp.writer)
	if err != nil {
		return
	}
	// If not async get the response
	if smpp.async {
		sequence = smpp.sequence
	} else {
		_, err = smpp.GetResp(CMD_UNBIND_RESP, smpp.sequence)
	}
	return
}

// Get response PDU 
func (smpp *smpp) GetResp(cmd SMPPCommand, sequence uint32) (rpdu PDU, err os.Error) {
	// Read the header
	hdr := new(PDUHeader)
	err = hdr.read(smpp.reader)
	if err != nil {
		return nil, err
	}
	// Has packet been read
	pduRead := false	
	// Defer reading rest of packet from buffer on error
	defer func() {
		if err != nil && !pduRead && hdr.CmdLength > 16 {
			p := make([]byte, hdr.CmdLength - 16)
			smpp.reader.Read(p)
		}
	}()
	// Check cmd and/or sequence if not 0
	if cmd != CMD_NONE && hdr.CmdId != cmd {
		err = os.NewError("Get Response: Invalid command")
		return nil, err
	}
	// Check sequence number if not 0
	if sequence > 0 && hdr.Sequence != sequence {
		err = os.NewError("Get Response: Invalid sequence number")
		return nil, err
	}
	// Check for error response
	if hdr.CmdStatus != STATUS_ESME_ROK {
		err = os.NewError("Get Response: PDU contains an error")
		return nil, err
	}
	// Set PDU as read (to disable the deferred read)
	pduRead = true
	// Get response PDU
	switch hdr.CmdId {
		// Default unhandled PDU
		default:
			err = os.NewError("Get Response: Unknown or unhandled PDU received")
			return nil, err
		// Bind responses
		case CMD_BIND_RECEIVER_RESP, CMD_BIND_TRANSMITTER_RESP, CMD_BIND_TRANSCEIVER_RESP:
			rpdu = new(PDUBindResp)
			rpdu.setHeader(hdr)
			err = rpdu.read(smpp.reader)
			if err != nil {
				return nil, err
			}
			// Set connection as bound
			smpp.bound = true
		// Unbind response
		case CMD_UNBIND_RESP:
			rpdu = new(PDUUnbindResp)
			rpdu.setHeader(hdr)
			err = rpdu.read(smpp.reader)
			if err != nil {
				return nil, err
			}
			// Set connection as unbound and disconnect
			smpp.bound = false
			smpp.close()
		// SubmitSM response
		case CMD_SUBMIT_SM_RESP:
			rpdu = new(PDUSubmitSMResp)
			rpdu.setHeader(hdr)
			err = rpdu.read(smpp.reader)
			if err != nil {
				return nil, err
			}
	}
	return
}

// Create a new Transmitter
func NewTransmitter(host string, port int, params Params) (tx *Transmitter, err os.Error) {
	// Merge params with defaults
	allParams := mergeParams(params, defaultsBind)
	// Create new transmitter
	tx = new(Transmitter)
	// Connect to server
	err = tx.connect(host, port)
	if err != nil {
		return nil, err
	}
	// Close connection on error
	defer func() {
		if err != nil {
			tx.close()
		}
	}()
	// Bind with server
	err = tx.bind(CMD_BIND_TRANSMITTER, CMD_BIND_TRANSMITTER_RESP, allParams)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Receiver
func NewReceiver(host string, port int, params Params) (rx *Receiver, err os.Error) {
	// Merge params with defaults
	allParams := mergeParams(params, defaultsBind)
	// Create new receiver
	rx = new(Receiver)
	// Connect to server
	err = rx.connect(host, port)
	if err != nil {
		return nil, err
	}
	// Close connection on error
	defer func() {
		if err != nil {
			rx.close()
		}
	}()
	// Bind with server
	err = rx.bind(CMD_BIND_RECEIVER, CMD_BIND_RECEIVER_RESP, allParams)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Transceiver
func NewTransceiver(host string, port int, params Params) (trx *Transceiver, err os.Error) {
	// Merge params with defaults
	allParams := mergeParams(params, defaultsBind)
	// Create new receiver
	trx = new(Transceiver)
	// Connect to server
	err = trx.connect(host, port)
	if err != nil {
		return nil, err
	}
	// Close connection on error
	defer func() {
		if err != nil {
			trx.close()
		}
	}()
	// Bind with server
	err = trx.bind(CMD_BIND_TRANSCEIVER, CMD_BIND_TRANSCEIVER_RESP, allParams)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Server
func NewServer() {
	return
}
