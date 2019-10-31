package tmsh

// ClientSSLProffile contains information about Client SSL profile
type ClientSSLProfile struct {
	Name      string `ltm:"name"`
	Component string `ltm:"component"`

	Cert                string                       `ltm:"cert"`
	CertKeyChain        map[string]map[string]string `ltm:"cert-key-chain"`
	Chain               string                       `ltm:"chain"`
	DefaultsFrom        string                       `ltm:"defaults-from"`
	InheritCertkeychain bool                         `ltm:"inherit-certkeychain"`
	Key                 string                       `ltm:"key"`
}
