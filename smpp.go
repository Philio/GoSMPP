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
type smppConn struct {
	conn		net.Conn
	reader		*bufio.Reader
	writer		*bufio.Writer
	connected	bool
}

// Connect to server
func (smpp *smppConn) connect(host string, port int) (err os.Error) {
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
func (smpp *smppConn) close() (err os.Error) {
	err = smpp.conn.Close()
	smpp.connected = false	
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
	// Bind with server
	err = tx.bind(params)
	if err != nil {
		return nil, err
	}
	// Get bind response
	err = tx.bindResp()
	return
}

// Create a new Receiver
func NewReceiver(host string, port int, username, password string) (rx *Receiver, err os.Error) {
	// Create new receiver
	rx = new(Receiver)
	// Connect to server
	err = rx.connect(host, port)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Transceiver
func NewTransceiver(host string, port int, username, password string) (trx *Transceiver, err os.Error) {
	// Create new transceiver
	trx = new(Transceiver)
	// Connect to server
	err = trx.connect(host, port)
	if err != nil {
		return nil, err
	}
	return
}

// Create a new Server
func NewServer() {

}
