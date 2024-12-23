package globals

type Config struct {
	TcpPort    int    `json:"tcpPort"`
	TLSTcpPort int    `json:"tlsTcpPort"`
	TLSPemFile string `json:"tlsPemFile"`
	TLSKeyFile string `json:"tlsKeyFile"`
	DbPath     string `json:"dbPath"`
	UseTLS     bool   `json:"useTls"`
}
