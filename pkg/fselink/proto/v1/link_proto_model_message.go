package v1

type ExchangeMessage struct {
	Type    ExchangeMessageType
	AuthKey string
	Payload []byte
}

type ExchangeMessageType int

const (
	message_type_start ExchangeMessageType = iota

	MessageTypeAuthReq
	MessageTypeAuthAns
	MessageTypeEvent
	MessageTypeFullSyncRequest
	MessageTypeFullSyncReply
	MessageTypeError
	MessageTypeGetFile
	MessageTypePushFile
	MessageTypePeerReady
	MessageTypeFileParts
	MessageTypeFileEnd
	MessageTypeClose
	MessageTypeDeleteObject

	message_type_end
)

func (t ExchangeMessageType) String() string {
	ts := [...]string{
		"MessageTypeAuthReq",
		"MessageTypeAuthAns",
		"MessageTypeEvent",
		"MessageTypeFullSyncRequest",
		"MessageTypeFullSyncReply",
		"MessageTypeError",
		"MessageTypeGetFile",
		"MessageTypePeerReady",
		"MessageTypePushFile",
		"MessageTypeFileParts",
		"MessageTypeFileEnd",
		"MessageTypeClose",
		"MessageTypeDeleteObject",
	}

	if t <= message_type_start || t >= message_type_end {
		return "unknown message type"
	}

	return ts[t-1]
}
