package tmsh

import (
	"fmt"
	"strconv"
	"strings"
)

func splitLtmOutput(str string) []string {
	sep := "ltm"
	strs := []string{}
	for _, s := range strings.Split(str, sep) {
		if s == "" {
			continue
		}
		s = sep + s
		strs = append(strs, s)
	}
	return strs
}

func (bigip *BigIP) RouteAdvertisementEnabled(targetVIP string) error {
	cmd := "modify ltm virtual-address " + targetVIP + " route-advertisement enabled"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) CreateSecurityFWPolicy(firewallPolicyName string) error {
	cmd := "create security firewall policy " + firewallPolicyName
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) AddSecurityFWPolicyRule(firewallPolicyName, ruleName, action string, targetPort int) error {
	if action != "accept" && action != "accept-decisively" &&
		action != "drop" && action != "reject" {
		return fmt.Errorf("Invalid action name: " + action)
	}

	cmd := "modify security firewall policy " + firewallPolicyName + " rules add { " + ruleName +
		" { action " + action + " ip-protocol tcp destination { ports add { " + strconv.Itoa(targetPort) + " } } place-after last } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
