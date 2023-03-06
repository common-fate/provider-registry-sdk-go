package msg

type RequestType string

const (
	RequestTypeGrant         RequestType = "grant"
	RequestTypeRevoke        RequestType = "revoke"
	RequestTypeDescribe      RequestType = "describe"
	RequestTypeLoadResources RequestType = "load"
)

// Request is an RPC request made to the Handler.
type Request interface {
	Type() RequestType
}

type Target struct {
	// Kind is defines which behaviour of the provider to use, e.g SSO or Group
	// The kind are defined by the provider schema, and each deployment is registered with its kind configuration in the database
	Kind      string            `json:"kind"`
	Arguments map[string]string `json:"arguments"`
}

type Grant struct {
	Subject string `json:"subject"`
	Target  Target `json:"target"`
}

func (Grant) Type() RequestType { return RequestTypeGrant }

type Revoke struct {
	Subject string `json:"subject"`
	Target  Target `json:"target"`
}

func (Revoke) Type() RequestType { return RequestTypeRevoke }

type LoadResources struct {
	Task string         `json:"task"`
	Ctx  map[string]any `json:"ctx"`
}

func (LoadResources) Type() RequestType { return RequestTypeLoadResources }

type Describe struct{}

func (Describe) Type() RequestType { return RequestTypeDescribe }
