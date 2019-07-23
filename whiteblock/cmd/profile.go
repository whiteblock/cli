package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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
	return util.GetP("profile", &profile)
}

func GetBiome(org Organization) (map[string]interface{}, error) {
	if len(org.Biomes) == 0 {
		return nil, fmt.Errorf("No available biomes")
	}
	if len(org.Biomes) == 1 { // There is only one biome so just choose that one
		return org.Biomes[0], nil
	}
	//Dont bother searching for biome if biome is not defined
	if !util.Exists("biome") {
		biomeChoices := []string{}
		for _, biome := range org.Biomes {
			biomeChoices = append(biomeChoices, biome["alias"].(string))
		}
		index := util.OptionListPrompt("Please select a biome", biomeChoices)
		util.Set("biome", fmt.Sprint(org.Biomes[index]["id"]))

		return org.Biomes[index], nil

	}
	var biomeName string
	err := util.GetP("biome", &biomeName)
	if err != nil {
		return nil, err
	}
	i := 0
	//Allow for automatic detect of organization id
	biomeId, err := strconv.Atoi(biomeName)
	isBiomeId := (err == nil)

	for i = 0; i < len(org.Biomes); i++ {
		if (isBiomeId && int(org.Biomes[i]["id"].(float64)) == biomeId) ||
			(!isBiomeId && org.Biomes[i]["alias"].(string) == biomeName) {
			return org.Biomes[i], nil
		}
	}
	return nil, fmt.Errorf("could not find biome")
}

func LoadBiomeAddress() error {
	var profile Profile
	var org Organization
	//Grab organization
	if util.Exists("profile") {
		err := util.GetP("profile", &profile)
		if err != nil {
			return err
		}
		org = profile.Organization
	} else {
		return fmt.Errorf("no profile data")
	}
	log.WithFields(log.Fields{"org": org}).Debug("got the org data")
	biome, err := GetBiome(org)
	if err != nil {
		return err
	}
	conf.ServerAddr = biome["host"].(string) + ":5001"
	return nil
}
