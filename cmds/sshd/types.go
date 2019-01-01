package main

// This kind of thing should be exported from the ssh package, I think?
type sshRemoteForward struct {
	Op                byte   //SSH_MSG_CHANNEL_OPEN
	Type              string //"forwarded-tcpip"
	SenderChannel     uint32
	InitialWindowSize uint32
	MaximumPacketSize uint32
	Addr              string
	Port              uint32
	OriginatorAddr    string
	OriginatorPort    uint32
}
