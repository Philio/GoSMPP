include $(GOROOT)/src/Make.$(GOARCH)
 
TARG=smpp
GOFILES=smpp.go smpp_transmitter.go smpp_receiver.go smpp_transceiver.go smpp_server.go smpp_pdu.go
 
include $(GOROOT)/src/Make.pkg 
