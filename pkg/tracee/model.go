package tracee

const (
	SeverityKey = "Severity"
	CategoryKey = "MITRE ATT&CK"
)

type Context struct {
	Timestamp      int    `json:"timestamp"`
	UserID         int    `json:"userId"`
	ProcessName    string `json:"processName"`
	HostName       string `json:"hostName"`
	ContainerID    string `json:"containerId"`
	ContainerImage string `json:"containerImage"`
	ContainerName  string `json:"containerName"`
	PodName        string `json:"podName"`
	PodNamespace   string `json:"podNamespace"`
	PodUID         string `json:"podUID"`
	EventID        int    `json:"eventId,string"`
	EventName      string `json:"eventName"`
	ReturnValue    int    `json:"returnValue"`
}

type SignatureMetadata struct {
	ID          string
	Version     string
	Name        string
	Description string
	Severity    int
	Tags        []string
	Properties  map[string]interface{}
}

type Event struct {
	Data        map[string]interface{}
	Context     Context
	SigMetadata SignatureMetadata
}
