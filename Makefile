include $(GOROOT)/src/Make.$(GOARCH)
 
TARG=smpp
GOFILES=smpp.go smpp_packet.go
 
include $(GOROOT)/src/Make.pkg 
