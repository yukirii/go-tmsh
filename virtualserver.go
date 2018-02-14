package tmsh

import (
	"fmt"
	"strconv"
	"strings"
)

// VirtualServer contains information about each virtual server
type VirtualServer struct {
	Name        string             `ltm:"name"`
	Destination string             `ltm:"destination"`
	IpProtocol  string             `ltm:"ip-protocol"`
	Mask        string             `ltm:"mask"`
	Partition   string             `ltm:"partition"`
	Pool        string             `ltm:"pool"`
	Profiles    map[string]Profile `ltm:"profiles"`
}

type Profile struct {
	Context string `ltm:"context"`
}

// GetAllVirtualServers returns a list of all virtual servers
func (bigip *BigIP) GetAllVirtualServers() ([]VirtualServer, error) {
	ret, err := bigip.ExecuteCommand("list ltm virtual all-properties")
	if err != nil {
		return nil, err
	}

	var vss []VirtualServer
	for _, s := range splitLtmOutput(ret) {
		var vs VirtualServer
		if err := Unmarshal(s, &vs); err != nil {
			return nil, err
		}
		vss = append(vss, vs)
	}

	return vss, nil
}

// GetVirtualServer gets a virtual server by name. Rerutn nil if the virtual server does not found.
func (bigip *BigIP) GetVirtualServer(name string) (*VirtualServer, error) {
	ret, _ := bigip.ExecuteCommand("list ltm virtual " + name)
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	var vs VirtualServer
	if err := Unmarshal(ret, &vs); err != nil {
		return nil, err
	}

	return &vs, nil
}

// CreateVirtualServer creates a virtual server.
func (bigip *BigIP) CreateVirtualServer(vsName, poolName, targetVIP, defaultProfileName string, targetPort int) error {
	destination := targetVIP + ":" + strconv.Itoa(targetPort)
	cmd := "create ltm virtual " + vsName + " { destination " + destination + " ip-protocol tcp mask 255.255.255.255 pool " + poolName + " profiles add { " + defaultProfileName + " } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// DeleteVirtualServer removes a virtual server.
func (bigip *BigIP) DeleteVirtualServer(vsName string) error {
	ret, _ := bigip.ExecuteCommand("delete ltm virtual " + vsName)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// AddVirtualServerProfile adds a profile to a virtual server.
func (bigip *BigIP) AddVirtualServerProfile(vsName, profileName, context string) error {
	cmd := "modify ltm virtual " + vsName + " profiles add { " + profileName + " { context " + context + " } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// DeleteVirtualServerProfile removes a profile from a virtual server.
func (bigip *BigIP) DeleteVirtualServerProfile(vsName, profileName, context string) error {
	cmd := "modify ltm virtual " + vsName + " profiles delete { " + profileName + " }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// ApplyPolicyToVirtualServer applies a policy to a virtual server.
func (bigip *BigIP) ApplyPolicyToVirtualServer(vsName, policyName string) error {
	cmd := "modify ltm virtual " + vsName + " fw-enforced-policy " + policyName
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
