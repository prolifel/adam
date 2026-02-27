package main

type Config struct {
	AccessKeyId     string `env:"ACCESS_KEY_ID,required"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,required"`
	SMTPHost        string `env:"SMTP_HOST"`
	SMTPPort        int    `env:"SMTP_PORT"`
	SMTPUsername    string `env:"SMTP_USERNAME"`
	SMTPPassword    string `env:"SMTP_PASSWORD"`
	EmailFrom       string `env:"EMAIL_FROM"`
	EmailTo         string `env:"EMAIL_TO"` // Comma-separated email addresses
	Token           string `env:"TOKEN"`
}

type AuthenticateRequest struct {
	AccessKeyId     string `json:"username"`
	SecretAccessKey string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

type ContainerProfile struct {
	ID           string       `json:"_id"`
	AccountIDs   []string     `json:"accountIDs"`
	Archived     bool         `json:"archived"`
	Capabilities Capabilities `json:"capabilities,omitempty"`
	Cluster      string       `json:"cluster,omitempty"`
	Collections  []string     `json:"collections"`
	Created      string       `json:"created,omitempty"`
	Entrypoint   string       `json:"entrypoint,omitempty"`
	// Events                       []HistoryEvent        `json:"events,omitempty"`
	Filesystem                   ProfileFilesystem     `json:"filesystem,omitempty"`
	Hash                         int64                 `json:"hash,omitempty"`
	HostNetwork                  bool                  `json:"hostNetwork"`
	HostPid                      bool                  `json:"hostPid"`
	Image                        string                `json:"image,omitempty"`
	ImageID                      string                `json:"imageID,omitempty"`
	Infra                        bool                  `json:"infra"`
	Istio                        bool                  `json:"istio"`
	K8s                          ProfileKubernetesData `json:"k8s,omitempty"`
	Label                        string                `json:"label,omitempty"`
	LastUpdate                   string                `json:"lastUpdate,omitempty"`
	LearnedStartup               bool                  `json:"learnedStartup"`
	Namespace                    string                `json:"namespace,omitempty"`
	Network                      ProfileNetwork        `json:"network,omitempty"`
	OS                           string                `json:"os,omitempty"`
	Processes                    ProfileProcesses      `json:"processes,omitempty"`
	RelearningCause              string                `json:"relearningCause,omitempty"`
	RemainingLearningDurationSec float64               `json:"remainingLearningDurationSec,omitempty"`
	State                        string                `json:"state,omitempty"`
}

type Capabilities struct {
	CI                     bool `json:"ci"`
	CloudMetadata          bool `json:"cloudMetadata"`
	DNSCache               bool `json:"dnsCache"`
	DynamicDNSQuery        bool `json:"dynamicDNSQuery"`
	DynamicFileCreation    bool `json:"dynamicFileCreation"`
	DynamicProcessCreation bool `json:"dynamicProcessCreation"`
	K8s                    bool `json:"k8s"`
	Proxy                  bool `json:"proxy"`
	PullImage              bool `json:"pullImage"`
	Sshd                   bool `json:"sshd"`
	Unpacker               bool `json:"unpacker"`
}

type HistoryEvent struct {
	ID       string `json:"_id,omitempty"`
	Command  string `json:"command,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Time     string `json:"time,omitempty"`
}

type ProfileFilesystem struct {
	Behavioral []FilesystemEntry `json:"behavioral"`
	Static     []FilesystemEntry `json:"static"`
}

type FilesystemEntry struct {
	Mount   bool   `json:"mount"`
	Path    string `json:"path,omitempty"`
	Process string `json:"process,omitempty"`
	Time    string `json:"time,omitempty"`
}

type ProfileKubernetesData struct {
	ClusterRoles   []K8sRole `json:"clusterRoles"`
	Roles          []K8sRole `json:"roles"`
	ServiceAccount string    `json:"serviceAccount,omitempty"`
}

type K8sRole struct {
	Labels      []Label   `json:"labels"`
	Name        string    `json:"name,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	RoleBinding string    `json:"roleBinding,omitempty"`
	Rules       []K8sRule `json:"rules"`
	Version     string    `json:"version,omitempty"`
}

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type K8sRule struct {
	APIGroups       []string `json:"apiGroups"`
	NonResourceURLs []string `json:"nonResourceURLs"`
	ResourceNames   []string `json:"resourceNames"`
	Resources       []string `json:"resources"`
	Verbs           []string `json:"verbs"`
}

type ProfileNetwork struct {
	Behavioral ProfileNetworkBehavioral `json:"behavioral,omitempty"`
	GeoIP      ProfileNetworkGeoIP      `json:"geoip,omitempty"`
	Static     ProfileNetworkStatic     `json:"static,omitempty"`
}

type ProfileNetworkBehavioral struct {
	DNSQueries     []DNSQuery      `json:"dnsQueries"`
	ListeningPorts []ListeningPort `json:"listeningPorts"`
	OutboundPorts  ProfilePortData `json:"outboundPorts,omitempty"`
}

type DNSQuery struct {
	DomainName string `json:"domainName,omitempty"`
	DomainType string `json:"domainType,omitempty"`
}

type ListeningPort struct {
	App       string          `json:"app,omitempty"`
	PortsData ProfilePortData `json:"portsData,omitempty"`
}

type ProfilePortData struct {
	All   bool   `json:"all"`
	Ports []Port `json:"ports"`
}

type Port struct {
	Port int    `json:"port,omitempty"`
	Time string `json:"time,omitempty"`
}

type ProfileNetworkGeoIP struct {
	Countries []GeoIPCountry `json:"countries"`
	Modified  string         `json:"modified,omitempty"`
}

type GeoIPCountry struct {
	Code     string `json:"code,omitempty"`
	IP       string `json:"ip,omitempty"`
	Modified string `json:"modified,omitempty"`
}

type ProfileNetworkStatic struct {
	ListeningPorts []ListeningPort `json:"listeningPorts"`
}

type ProfileProcesses struct {
	Behavioral []ProcessEntry `json:"behavioral"`
	Static     []ProcessEntry `json:"static"`
}

type ProcessEntry struct {
	Command     string `json:"command,omitempty"`
	Interactive bool   `json:"interactive"`
	MD5         string `json:"md5,omitempty"`
	Modified    bool   `json:"modified"`
	Path        string `json:"path,omitempty"`
	PPath       string `json:"ppath,omitempty"`
	Time        string `json:"time,omitempty"`
	User        string `json:"user,omitempty"`
}

type VerdictRecord struct {
	ID             int
	CollectionName string
	Key            string
	Value          string
	Verdict        string
	Remarks        string
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// header := []string{"id", "collection_name", "key", "value", "verdict", "remarks"}
type CapabilitiesCSVHeader struct {
	ID             string `prep:"trim" validate:"required"`
	CollectionName string `prep:"trim" validate:"required"`
	Key            string `prep:"trim" validate:"required"`
	Value          string `prep:"trim" validate:"required"`
	Verdict        string `prep:"trim" validate:"required"`
	Remarks        string `prep:"trim" validate:"required"`
}

// ContainerPolicy represents the Prisma Cloud Runtime Container Policy
type ContainerPolicy struct {
	ID               string          `json:"_id,omitempty"`
	LearningDisabled bool            `json:"learningDisabled"`
	Rules            []ContainerRule `json:"rules"`
}

type ContainerRule struct {
	AdvancedProtectionEffect       string         `json:"advancedProtectionEffect,omitempty"`
	CloudMetadataEnforcementEffect string         `json:"cloudMetadataEnforcementEffect,omitempty"`
	Collections                    []Collection   `json:"collections"`
	CustomRules                    []CustomRule   `json:"customRules"`
	Disabled                       bool           `json:"disabled"`
	DNS                            DNSRule        `json:"dns,omitempty"`
	Filesystem                     FileSystemRule `json:"filesystem,omitempty"`
	KubernetesEnforcementEffect    string         `json:"kubernetesEnforcementEffect,omitempty"`
	Modified                       string         `json:"modified,omitempty"`
	Name                           string         `json:"name,omitempty"`
	Network                        NetworkRule    `json:"network,omitempty"`
	Owner                          string         `json:"owner,omitempty"`
	PreviousName                   string         `json:"previousName,omitempty"`
	Processes                      ProcessRule    `json:"processes,omitempty"`
	SkipExecSessions               bool           `json:"skipExecSessions"`
	WildFireAnalysisEffect         string         `json:"wildFireAnalysis,omitempty"`
}

// Collection represents a collection within a rule - full Prisma Cloud schema
type Collection struct {
	AccountIDs  []string `json:"accountIDs"`
	AppIDs      []string `json:"appIDs"`
	Clusters    []string `json:"clusters"`
	Color       string   `json:"color,omitempty"`
	Containers  []string `json:"containers"`
	Description string   `json:"description,omitempty"`
	Functions   []string `json:"functions"`
	Hosts       []string `json:"hosts"`
	Images      []string `json:"images"`
	Labels      []string `json:"labels"`
	Modified    string   `json:"modified,omitempty"`
	Name        string   `json:"name,omitempty"`
	Namespaces  []string `json:"namespaces"`
	Owner       string   `json:"owner,omitempty"`
	Prisma      bool     `json:"prisma"`
	System      bool     `json:"system"`
}

// DNSRule contains DNS-related runtime rules (matches fetch manual.json)
type DNSRule struct {
	DefaultEffect string     `json:"defaultEffect,omitempty"`
	Disabled      bool       `json:"disabled"`
	DomainList    DomainList `json:"domainList,omitempty"`
}

// DomainList contains allowed/blocked domains
type DomainList struct {
	Allowed []string `json:"allowed"`
	Denied  []string `json:"denied"`
	Effect  string   `json:"effect,omitempty"`
}

// FileSystemRule contains filesystem-related runtime rules (matches fetch manual.json)
type FileSystemRule struct {
	AllowedList                []string             `json:"allowedList"`
	BackdoorFilesEffect        string               `json:"backdoorFilesEffect,omitempty"`
	DefaultEffect              string               `json:"defaultEffect,omitempty"`
	DeniedList                 FileSystemDeniedList `json:"deniedList,omitempty"`
	Disabled                   bool                 `json:"disabled"`
	EncryptedBinariesEffect    string               `json:"encryptedBinariesEffect,omitempty"`
	NewFilesEffect             string               `json:"newFilesEffect,omitempty"`
	SuspiciousELFHeadersEffect string               `json:"suspiciousELFHeadersEffect,omitempty"`
}

type FileSystemDeniedList struct {
	Effect string   `json:"effect,omitempty"`
	Paths  []string `json:"paths"`
}

// NetworkRule contains network-related runtime rules (matches fetch manual.json)
type NetworkRule struct {
	AllowedIPs         []string      `json:"allowedIPs"`
	DefaultEffect      string        `json:"defaultEffect,omitempty"`
	DeniedIPs          []string      `json:"deniedIPs"`
	DeniedIPsEffect    string        `json:"deniedIPsEffect,omitempty"`
	Disabled           bool          `json:"disabled"`
	ListeningPorts     ContainerPort `json:"listeningPorts,omitempty"`
	ModifiedProcEffect string        `json:"modifiedProcEffect,omitempty"`
	OutboundPorts      ContainerPort `json:"outboundPorts,omitempty"`
	PortScanEffect     string        `json:"portScanEffect,omitempty"`
	RawSocketsEffect   string        `json:"rawSocketsEffect,omitempty"`
}

type ContainerPort struct {
	Allowed []ContainerPortObject `json:"allowed"`
	Denied  []ContainerPortObject `json:"denied"`
	Effect  string                `json:"effect,omitempty"`
}

type ContainerPortObject struct {
	Deny  bool `json:"deny"`
	End   int  `json:"end,omitempty"`
	Start int  `json:"start,omitempty"`
}

// ProcessRule contains process-related runtime rules (matches fetch manual.json)
type ProcessRule struct {
	AllowedList           []string   `json:"allowedList"`
	CheckParentChild      bool       `json:"checkParentChild"`
	CryptoMinersEffect    string     `json:"cryptoMinersEffect,omitempty"`
	DefaultEffect         string     `json:"defaultEffect,omitempty"`
	DeniedList            DeniedList `json:"deniedList,omitempty"`
	Disabled              bool       `json:"disabled"`
	LateralMovementEffect string     `json:"lateralMovementEffect,omitempty"`
	ModifiedProcessEffect string     `json:"modifiedProcessEffect,omitempty"`
	ReverseShellEffect    string     `json:"reverseShellEffect,omitempty"`
	SuidBinariesEffect    string     `json:"suidBinariesEffect,omitempty"`
}

type DeniedList struct {
	Effect string   `json:"effect,omitempty"`
	Paths  []string `json:"paths"`
}

// CustomRule contains custom runtime rule definitions
type CustomRule struct {
	ID     int    `json:"_id,omitempty"`
	Effect string `json:"effect,omitempty"`
	Action string `json:"action,omitempty"`
}
