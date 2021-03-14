// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)

// Methods for configuring dnsmasq on the FMS server

package network

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/Team254/cheesy-arena/model"
)

type DnsMasq struct {
	mutex sync.Mutex
}

func NewDnsMasq() *DnsMasq {
	return &DnsMasq{}
}

func (dm *DnsMasq) ConfigureTeamEthernet(teams [6]*model.Team) error {
	// Make sure multiple configurations aren't being set at the same time.
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Determine what new team VLANs are needed and build the commands to set them up.
	oldTeamVlans, err := dm.getTeamVlans()
	if err != nil {
		return err
	}
	replaceTeamVlan := func(team *model.Team, vlan int) {
		if team == nil {
			return
		}
		if oldTeamVlans[team.Id] == vlan {
			delete(oldTeamVlans, team.Id)
		} else {
			contents := []byte(fmt.Sprintf(
				"# Options for VLAN%d\n" +
				"# Team %d\n" +
				"\n" +
				"dhcp-range=set:vlan%d,10.%d.%d.101,10.%d.%d.199,255.255.255.0,12h\n" +
				"dhcp-option=tag:vlan%d,3,10.%d.%d.61\n",
				vlan, team.Id, vlan, team.Id/100, team.Id%100, team.Id/100, team.Id%100,
				vlan, team.Id/100, team.Id%100))
			err := ioutil.WriteFile(fmt.Sprintf("/etc/dnsmasq.d/vlan%d.conf", vlan), contents, 0664)
			if err != nil {
				log.Printf("Failed to configure VLAN%d for team %d: %s", vlan, team.Id, err.Error())
				return
			}
		}
	}
	replaceTeamVlan(teams[0], red1Vlan)
	replaceTeamVlan(teams[1], red2Vlan)
	replaceTeamVlan(teams[2], red3Vlan)
	replaceTeamVlan(teams[3], blue1Vlan)
	replaceTeamVlan(teams[4], blue2Vlan)
	replaceTeamVlan(teams[5], blue3Vlan)

	// Remove configuration files for VLANs no longer needed
	for _, vlan := range oldTeamVlans {
		os.Remove(fmt.Sprintf("/etc/dnsmasq.d/vlan%d.conf", vlan))
	}

	// Restart the dnsmasq service
	cmd := exec.Command("/usr/bin/sudo", "/usr/bin/systemctl", "restart", "dnsmasq")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (dm *DnsMasq) getTeamVlans() (map[int]int, error) {
	files, err := ioutil.ReadDir("/etc/dnsmasq.d")
	if err != nil {
		return nil, err
	}

	teamVlans := make(map[int]int)

	for _, file := range files {
		fn := file.Name()
		if fn == "vlan100.conf" {
			// Skip vlan 100
			continue
		}
		if !strings.HasPrefix(fn, "vlan") && !strings.HasSuffix(fn, ".conf") {
			// skip files that don't match vlan*.conf
			continue
		}

		vlan, err := strconv.Atoi(fn[4:6])
		if err != nil {
			return nil, err
		}

		fh, err := os.Open("/etc/dnsmasq.d/" + fn)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "# Team") {
				team, err := strconv.Atoi(line[7:])
				if err != nil {
					fh.Close()
					return nil, err
				}
				teamVlans[team] = vlan
			}
		}

		fh.Close()
	}

	return teamVlans, nil
}
