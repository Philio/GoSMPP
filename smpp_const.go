// GoSMPP - An SMPP library for Go
// Copyright 2010 Phil Bayfield
// This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
// Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
package smpp

type smppCommand uint32

const (
	CMD_GENERIC_NACK		= 0x80000000
	CMD_BIND_RECEIVER 		= 0x00000001
	CMD_BIND_RECEIVER_RESP		= 0x80000001
	CMD_BIND_TRANSMITTER		= 0x00000002
	CMD_BIND_TRANSMITTER_RESP	= 0x80000002
	CMD_QUERY_SM			= 0x00000003
	CMD_QUERY_SM_RESP		= 0x80000003
	CMD_SUBMIT_SM			= 0x00000004
	CMD_SUBMIT_SM_RESP		= 0x80000004
	CMD_DELIVER_SM			= 0x00000005
	CMD_DELIVER_SM_RESP		= 0x80000005
	CMD_UNBIND			= 0x00000006
	CMD_UNBIND_RESP			= 0x80000006
	CMD_REPLACE_SM			= 0x00000007
	CMD_REPLACE_SM_RESP		= 0x80000007
	CMD_CANCEL_SM			= 0x00000008
	CMD_CANCEL_SM_RESP		= 0x80000008
	CMD_BIND_TRANSCEIVER		= 0x00000009
	CMD_BIND_TRANSCEIVER_RESP	= 0x80000009
	CMD_OUTBIND			= 0x0000000b
	CMD_ENQUIRE_LINK		= 0x00000015
	CMD_ENQUIRE_LINK_RESP		= 0x80000015
	CMD_SUBMIT_MULTI		= 0x00000021
	CMD_SUBMIT_MULTI_RESP		= 0x80000021
	CMD_DATA_SM			= 0x00000103
	CMD_DATA_SM_RESP		= 0x80000103
)

type smppCommandStatus uint32

const (
	STATUS_ESME_ROK			= 0x00000000	// No Error
	STATUS_ESME_RINVMSGLEN		= 0x00000001	// Message Length is invalid
	STATUS_ESME_RINVCMDLEN		= 0x00000002	// Command Length is invalid
	STATUS_ESME_RINVCMDID		= 0x00000003	// Invalid Command ID
	STATUS_ESME_RINVBNDSTS		= 0x00000004	// Incorrect BIND Status for given command
	STATUS_ESME_RALYBND		= 0x00000005	// ESME Already in Bound State
	STATUS_ESME_RINVPRTFLG		= 0x00000006	// Invalid Priority Flag
	STATUS_ESME_RINVREGDLVFLG	= 0x00000007	// Invalid Registered Delivery Flag
	STATUS_ESME_RSYSERR		= 0x00000008	// System Error
	STATUS_ESME_RINVSRCADR		= 0x0000000a	// Invalid Source Address
	STATUS_ESME_RINVDSTADR		= 0x0000000b	// Invalid Dest Addr
	STATUS_ESME_RINVMSGID		= 0x0000000c	// Message ID is invalid
	STATUS_ESME_RBINDFAIL		= 0x0000000d	// Bind Failed
	STATUS_ESME_RINVPASWD		= 0x0000000e	// Invalid Password
	STATUS_ESME_RINVSYSID		= 0x0000000f	// Invalid System ID
	STATUS_ESME_RCANCELFAIL		= 0x00000011	// Cancel SM Failed
	STATUS_ESME_RREPLACEFAIL	= 0x00000013	// Replace SM Failed
	STATUS_ESME_RMSGQFUL		= 0x00000014	// Message Queue Full
	STATUS_ESME_RINVSERTYP		= 0x00000015	// Invalid Service Type
	STATUS_ESME_RINVNUMDESTS	= 0x00000033	// Invalid number of destinations
	STATUS_ESME_RINVDLNAME		= 0x00000034	// Invalid Distribution List name
	STATUS_ESME_RINVDESTFLAG	= 0x00000040	// Destination flag is invalid
	STATUS_ESME_RINVSUBREP		= 0x00000042	// Invalid ‘submit with replace’ request
	STATUS_ESME_RINVESMCLASS	= 0x00000043	// Invalid esm_class field data
	STATUS_ESME_RCNTSUBDL		= 0x00000044	// Cannot Submit to Distribution List
	STATUS_ESME_RSUBMITFAIL		= 0x00000045	// Submit_sm or submit_multi failed
	STATUS_ESME_RINVSRCTON		= 0x00000048	// Invalid Source address TON
	STATUS_ESME_RINVSRCNPI		= 0x00000049	// Invalid Source address NPI
	STATUS_ESME_RINVDSTTON		= 0x00000050	// Invalid Destination address TON
	STATUS_ESME_RINVDSTNPI		= 0x00000051	// Invalid Destination address NPI
	STATUS_ESME_RINVSYSTYP		= 0x00000053	// Invalid system_type field
	STATUS_ESME_RINVREPFLAG		= 0x00000054	// Invalid replace_if_present flag
	STATUS_ESME_RINVNUMMSGS		= 0x00000055	// Invalid number of messages
	STATUS_ESME_RTHROTTLED		= 0x00000058	// Throttling error (ESME has exceeded allowed message limits)
	STATUS_ESME_RINVSCHED		= 0x00000061	// Invalid Scheduled Delivery Time
	STATUS_ESME_RINVEXPIRY		= 0x00000062	// Invalid message validity period (Expiry time)
	STATUS_ESME_RINVDFTMSGID	= 0x00000063	// Predefined Message Invalid or Not Found
	STATUS_ESME_RX_T_APPN		= 0x00000064	// ESME Receiver Temporary App Error Code
	STATUS_ESME_RX_P_APPN		= 0x00000065	// ESME Receiver Permanent App Error Code
	STATUS_ESME_RX_R_APPN		= 0x00000066	// ESME Receiver Reject Message Error Code
	STATUS_ESME_RQUERYFAIL		= 0x00000067	// Query_sm request failed
	STATUS_ESME_RINVOPTPARSTREAM	= 0x000000c0	// Error in the optional part of the PDU Body.
	STATUS_ESME_ROPTPARNOTALLWD	= 0x000000c1	// Optional Parameter not allowed
	STATUS_ESME_RINVPARLEN		= 0x000000c2	// Invalid Parameter Length.
	STATUS_ESME_RMISSINGOPTPARAM	= 0x000000c3	// Expected Optional Parameter missing
	STATUS_ESME_RINVOPTPARAMVAL	= 0x000000c4	// Invalid Optional Parameter Value
	STATUS_ESME_RDELIVERYFAILURE	= 0x000000fe	// Delivery Failure (used for data_sm_resp)
	STATUS_ESME_RUNKNOWNERR		= 0x000000ff	// Unknown Error
)
