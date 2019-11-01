package tmsh

import (
	"fmt"
	"strings"
)

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

// GetAllClientSSLProfiles returns a list of all Client SSL Profiles
func (bigip *BigIP) GetAllClientSSLProfiles() ([]ClientSSLProfile, error) {
	ret, err := bigip.ExecuteCommand("list ltm profile client-ssl")
	if err != nil {
		return nil, err
	}

	var profs []ClientSSLProfile
	for _, p := range splitLtmOutput(ret) {
		var prof ClientSSLProfile
		if err := Unmarshal(p, &prof); err != nil {
			return nil, err
		}
		profs = append(profs, prof)
	}

	return profs, nil
}

// GetClientSSLProfile gets a Client SSL Profile by name. Return nil if the profile does not found.
func (bigip *BigIP) GetClientSSLProfile(name string) (*ClientSSLProfile, error) {
	ret, _ := bigip.ExecuteCommand("list ltm profile client-ssl " + name)
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	var prof ClientSSLProfile
	if err := Unmarshal(ret, &prof); err != nil {
		return nil, err
	}

	return &prof, nil
}
