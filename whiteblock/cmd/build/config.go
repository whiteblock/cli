package build

type Config struct {
	Servers      []int                  `json:"servers"`
	Blockchain   string                 `json:"blockchain"`
	Nodes        int                    `json:"nodes"`
	Images       []string               `json:"images"`
	Resources    []Resources            `json:"resources"`
	Params       map[string]interface{} `json:"params"`
	Environments []map[string]string    `json:"environments"`
	Files        []map[string]string    `json:"files"`
	Logs         []map[string]string    `json:"logs"`
	Extras       map[string]interface{} `json:"extras"`
	Meta         map[string]interface{} `json:"__meta"`
}

type Resources struct {
	Cpus      string   `json:"cpus"`
	Memory    string   `json:"memory"`
	Ports     []string `json:"ports"`
	Volumes   []string `json:"volumes"`
	BoundCPUs []int    `json:"boundCPUs,omitonempty"`
}
