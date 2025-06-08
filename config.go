package config

type ServerConfig struct {
	DNS         string
	PublicKey   string
	PublicIP    string
	Port        int
	AllowedIPs  string
	WGInterface string
}

type ClientConfig struct {
	ClientIP            string
	PrivateKey          string
	PresistentKeepAlive int
}

type GenerationConfig struct {
	Dir string
}
