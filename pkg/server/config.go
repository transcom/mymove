package server

// HostsConfig gives names for the hosts we respond to and the interface(s) to listen on
type HostsConfig struct {
	ListenInterface string // Possibly the ListenInterface should be pulled out of this into ? NetworkConfig  ?
	MyName          string
	OfficeName      string
	TspName         string
	OrdersName      string
}

// LocalEnvConfig contains environment type and  directories for the running server
type LocalEnvConfig struct {
	Environment string // Local environment
	SiteDir     string // Where the built client site and assets are stored
	ConfigDir   string // Local configuration directory
}
