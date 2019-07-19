package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/whiteblock/cli/whiteblock/util"
	"strconv"
)

var (
	profile Profile
)

type Organization struct {
	Id        int                      `json:"id"`
	Name      string                   `json:"name"`
	CreatedAt string                   `json:"created_at"`
	UpdatedAt string                   `json:"updated_at"`
	Biomes    []map[string]interface{} `json:"biomes"`
}

type Profile struct {
	Id           int          `json:"id"`
	Organization Organization `json:"organization"`
}

/*
type Profile struct {
	Id            int                      `json:"id"`
	Email         string                   `json:"email"`
	EmailVerified interface{}              `json:"email_verified"`
	Name          string                   `json:"name"`
	Picture       string                   `json:"picture"`
	CreatedAt     string                   `json:"created_at"`
	UpdatedAt     string                   `json:"updated_at"`
	SshKeys       []map[string]interface{} `json:"ssh_keys"`
	Organizations []Organization           `json:"organizations"`
}
*/
/*func LoadOrganizationApiKey() error {
	rawKey, err := util.ReadStore("org_key")
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawKey, &org_key)
	if err != nil {
		return err
	}

	return nil
}*/

func LoadProfile() error {
	rawProfile, err := util.ReadStore("profile")
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawProfile, &profile)
	if err != nil {
		return err
	}

	return nil
}

func LoadBiomeAddress() error {
	var profile Profile
	var org Organization
	//Grab organization
	if util.StoreExists("profile") {
		rawOrgKey, err := util.ReadStore("profile")
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawOrgKey, &profile)
		if err != nil {
			return err
		}
		org = profile.Organization
	} else {
		return fmt.Errorf("no profile data")
	}

	var biome map[string]interface{}
	if len(org.Biomes) == 0 {
		return fmt.Errorf("No available biomes")
	}
	//Dont bother searching for biome if biome is not defined
	if !util.StoreExists("biome") {
		//TODO, improve this
		biome = org.Biomes[0]
	} else {
		rawBiomeName, err := util.ReadStore("biome")
		if err != nil {
			return err
		}
		biomeName := string(rawBiomeName)
		i := 0
		//Allow for automatic detect of organization id
		biomeId, err := strconv.Atoi(biomeName)
		isBiomeId := (err == nil)

		for i = 0; i < len(org.Biomes); i++ {
			if (isBiomeId && int(org.Biomes[i]["id"].(float64)) == biomeId) ||
				(!isBiomeId && org.Biomes[i]["alias"].(string) == biomeName) {
				biome = org.Biomes[i]
				break
			}
		}
		if i == len(org.Biomes) {
			return fmt.Errorf("Could not find biome")
		}
	}
	conf.ServerAddr = biome["host"].(string) + ":5001"
	return nil
}
