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
	AccountIDs   []string     `json:"accountIDs,omitempty"`
	Archived     bool         `json:"archived,omitempty"`
	Capabilities Capabilities `json:"capabilities,omitempty"`
	Cluster      string       `json:"cluster,omitempty"`
	Collections  []string     `json:"collections,omitempty"`
	Created      string       `json:"created,omitempty"`
	Entrypoint   string       `json:"entrypoint,omitempty"`
	// Events                       []HistoryEvent        `json:"events,omitempty"`
	Filesystem                   ProfileFilesystem     `json:"filesystem,omitempty"`
	Hash                         int64                 `json:"hash,omitempty"`
	HostNetwork                  bool                  `json:"hostNetwork,omitempty"`
	HostPid                      bool                  `json:"hostPid,omitempty"`
	Image                        string                `json:"image,omitempty"`
	ImageID                      string                `json:"imageID,omitempty"`
	Infra                        bool                  `json:"infra,omitempty"`
	Istio                        bool                  `json:"istio,omitempty"`
	K8s                          ProfileKubernetesData `json:"k8s,omitempty"`
	Label                        string                `json:"label,omitempty"`
	LastUpdate                   string                `json:"lastUpdate,omitempty"`
	LearnedStartup               bool                  `json:"learnedStartup,omitempty"`
	Namespace                    string                `json:"namespace,omitempty"`
	Network                      ProfileNetwork        `json:"network,omitempty"`
	OS                           string                `json:"os,omitempty"`
	Processes                    ProfileProcesses      `json:"processes,omitempty"`
	RelearningCause              string                `json:"relearningCause,omitempty"`
	RemainingLearningDurationSec float64               `json:"remainingLearningDurationSec,omitempty"`
	State                        string                `json:"state,omitempty"`
}

type Capabilities struct {
	CI                     bool `json:"ci,omitempty"`
	CloudMetadata          bool `json:"cloudMetadata,omitempty"`
	DNSCache               bool `json:"dnsCache,omitempty"`
	DynamicDNSQuery        bool `json:"dynamicDNSQuery,omitempty"`
	DynamicFileCreation    bool `json:"dynamicFileCreation,omitempty"`
	DynamicProcessCreation bool `json:"dynamicProcessCreation,omitempty"`
	K8s                    bool `json:"k8s,omitempty"`
	Proxy                  bool `json:"proxy,omitempty"`
	PullImage              bool `json:"pullImage,omitempty"`
	Sshd                   bool `json:"sshd,omitempty"`
	Unpacker               bool `json:"unpacker,omitempty"`
}

type HistoryEvent struct {
	ID       string `json:"_id,omitempty"`
	Command  string `json:"command,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Time     string `json:"time,omitempty"`
}

type ProfileFilesystem struct {
	Behavioral []FilesystemEntry `json:"behavioral,omitempty"`
	Static     []FilesystemEntry `json:"static,omitempty"`
}

type FilesystemEntry struct {
	Mount   bool   `json:"mount,omitempty"`
	Path    string `json:"path,omitempty"`
	Process string `json:"process,omitempty"`
	Time    string `json:"time,omitempty"`
}

type ProfileKubernetesData struct {
	ClusterRoles   []K8sRole `json:"clusterRoles,omitempty"`
	Roles          []K8sRole `json:"roles,omitempty"`
	ServiceAccount string    `json:"serviceAccount,omitempty"`
}

type K8sRole struct {
	Labels      []Label   `json:"labels,omitempty"`
	Name        string    `json:"name,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	RoleBinding string    `json:"roleBinding,omitempty"`
	Rules       []K8sRule `json:"rules,omitempty"`
	Version     string    `json:"version,omitempty"`
}

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type K8sRule struct {
	APIGroups       []string `json:"apiGroups,omitempty"`
	NonResourceURLs []string `json:"nonResourceURLs,omitempty"`
	ResourceNames   []string `json:"resourceNames,omitempty"`
	Resources       []string `json:"resources,omitempty"`
	Verbs           []string `json:"verbs,omitempty"`
}

type ProfileNetwork struct {
	Behavioral ProfileNetworkBehavioral `json:"behavioral,omitempty"`
	GeoIP      ProfileNetworkGeoIP      `json:"geoip,omitempty"`
	Static     ProfileNetworkStatic     `json:"static,omitempty"`
}

type ProfileNetworkBehavioral struct {
	DNSQueries     []DNSQuery      `json:"dnsQueries,omitempty"`
	ListeningPorts []ListeningPort `json:"listeningPorts,omitempty"`
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
	All   bool   `json:"all,omitempty"`
	Ports []Port `json:"ports,omitempty"`
}

type Port struct {
	Port int    `json:"port,omitempty"`
	Time string `json:"time,omitempty"`
}

type ProfileNetworkGeoIP struct {
	Countries []GeoIPCountry `json:"countries,omitempty"`
	Modified  string         `json:"modified,omitempty"`
}

type GeoIPCountry struct {
	Code     string `json:"code,omitempty"`
	IP       string `json:"ip,omitempty"`
	Modified string `json:"modified,omitempty"`
}

type ProfileNetworkStatic struct {
	ListeningPorts []ListeningPort `json:"listeningPorts,omitempty"`
}

type ProfileProcesses struct {
	Behavioral []ProcessEntry `json:"behavioral,omitempty"`
	Static     []ProcessEntry `json:"static,omitempty"`
}

type ProcessEntry struct {
	Command     string `json:"command,omitempty"`
	Interactive bool   `json:"interactive,omitempty"`
	MD5         string `json:"md5,omitempty"`
	Modified    bool   `json:"modified,omitempty"`
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
