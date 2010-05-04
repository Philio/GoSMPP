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

// Send bind request
func (smpp *smpp) bind(cmd, rcmd SMPPCommand, params []interface{}) (err os.Error) {
	// Sequence number starts at 1
	smpp.sequence ++
	// Create bind PDU
	pdu := new(pduBind)
	// PDU header
	pdu.header = new(pduHeader)
	pdu.header.cmdLength = 23 // Min length
	pdu.header.cmdId     = cmd
	pdu.header.cmdStatus = STATUS_ESME_ROK
	pdu.header.sequence  = smpp.sequence
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
	err = pdu.write(smpp.writer)
	if err != nil {
		return
	}
	// Create bind response PDU
	rpdu := new(pduBindResp)
	// Read PDU data
	err = rpdu.read(smpp.reader)
	if err != nil {
		return
	}
	// Validate PDU data
	if rpdu.header.cmdId != rcmd {
		err = os.NewError("Bind Response: Invalid command")
		return
	}
	if rpdu.header.cmdStatus != STATUS_ESME_ROK {
		err = os.NewError("Bind Response: Error received from SMSC")
		return
	}
	if rpdu.header.sequence != smpp.sequence {
		err = os.NewError("Bind Response: Invalid sequence number")
		return
	}
	smpp.bound = true
	return
}

// Send unbind request
func (smpp *smpp) Unbind() (err os.Error) {
	// Increment sequence number
	smpp.sequence ++
	// Create bind PDU
	pdu := new(pduUnbind)
	// PDU header
	pdu.header = new(pduHeader)
	pdu.header.cmdLength = 16
	pdu.header.cmdId     = CMD_UNBIND
	pdu.header.cmdStatus = STATUS_ESME_ROK
	pdu.header.sequence  = smpp.sequence
	// Send PDU
	err = pdu.write(smpp.writer)
	if err != nil {
		return
	}
	// Create unbind response PDU
	rpdu := new(pduUnbindResp)
	// Read PDU data
	err = rpdu.read(smpp.reader)
	if err != nil {
		return
	}
	// Validate PDU data
	if rpdu.header.cmdId != CMD_UNBIND_RESP {
		err = os.NewError("Unbind Response: Invalid command")
		return
	}
	if rpdu.header.cmdStatus != STATUS_ESME_ROK {
		err = os.NewError("Unbind Response: Error received from SMSC")
		return
	}
	if rpdu.header.sequence != smpp.sequence {
		err = os.NewError("Unbind Response: Invalid sequence number")
		return
	}
	// Disconnect
	smpp.bound = false
	smpp.close()
	return
}

// Create a new Transmitter
func NewTransmitter(host string, port int, params ...interface{}) (tx *Transmitter, err os.Error) {
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
	err = tx.bind(CMD_BIND_TRANSMITTER, CMD_BIND_TRANSMITTER_RESP, params)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Receiver
func NewReceiver(host string, port int, params ...interface{}) (rx *Receiver, err os.Error) {
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
	err = rx.bind(CMD_BIND_RECEIVER, CMD_BIND_RECEIVER_RESP, params)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Transceiver
func NewTransceiver(host string, port int, params ...interface{}) (trx *Transceiver, err os.Error) {
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
	err = trx.bind(CMD_BIND_TRANSCEIVER, CMD_BIND_TRANSCEIVER_RESP, params)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Server
func NewServer() {
	return
}
