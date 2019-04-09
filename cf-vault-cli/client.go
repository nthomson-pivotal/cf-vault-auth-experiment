package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	cli "gopkg.in/urfave/cli.v1"
)

const AppsKey = "apps"
const SpacesKey = "spaces"
const OrgsKey = "orgs"

type SpaceInfo struct {
	SpaceName string `json:"name"`
	SpaceGUID string `json:"guid"`
}

type VaultAuthPluginClient struct {
	cliConnection plugin.CliConnection
	VaultClient   *api.Client
}

func (c *VaultAuthPluginClient) GetAppMapping(context *cli.Context) error {
	appName := context.String("app")
	vaultMount := context.GlobalString("mount")

	if appName == "" {
		return cli.NewExitError("You must specify a CF application name get mapping", 3)
	}

	printLoadingMessage()

	app, err := c.cliConnection.GetApp(appName)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	spaceInfo, err := getSpaceInfo(app.SpaceGuid, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	space, err := c.cliConnection.GetSpace(spaceInfo.SpaceName)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	appPolicies, spacePolicies, orgPolicies, err := getAppSpaceOrgPolicies(c.VaultClient, vaultMount, app.Guid, app.SpaceGuid, space.Organization.Guid)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	b := color.New(color.Bold)

	b.Print("App: \t\t")
	fmt.Printf("%s (%s)\n\n", appName, app.Guid)

	b.Print("Mapped policies: \t\t")
	fmt.Printf("%s\n", appPolicies)
	b.Print("Inherited space policies: \t")
	fmt.Printf("%s\n", spacePolicies)
	b.Print("Inherited org policies: \t")
	fmt.Printf("%s\n", orgPolicies)

	return nil
}

func (c *VaultAuthPluginClient) SetAppMapping(context *cli.Context) error {
	appName := context.String("app")
	policies := context.String("policies")
	vaultMount := context.GlobalString("mount")

	if appName == "" {
		return cli.NewExitError("You must specify a CF application name to map", 3)
	}

	if policies == "" {
		return cli.NewExitError("You must specify one or more Vault policy names to map to", 3)
	}

	app, err := c.cliConnection.GetApp(appName)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = writeMapping(c.VaultClient, vaultMount, AppsKey, app.Guid, policies)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func (c *VaultAuthPluginClient) RemoveAppMapping(context *cli.Context) error {
	appName := context.String("app")
	vaultMount := context.GlobalString("mount")

	if appName == "" {
		return cli.NewExitError("You must specify a CF application name to unmap", 3)
	}

	app, err := c.cliConnection.GetApp(appName)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = removeMapping(c.VaultClient, vaultMount, AppsKey, app.Guid)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func (c *VaultAuthPluginClient) GetSpaceMapping(context *cli.Context) error {
	spaceName := context.String("space")
	vaultMount := context.GlobalString("mount")

	printLoadingMessage()

	var err error
	var spaceGUID string
	var orgGUID string

	if spaceName == "" {
		space, err := c.cliConnection.GetCurrentSpace()
		if err != nil {
			return cli.NewExitError(err, 3)
		}

		spaceName = space.Name
		spaceGUID = space.Guid

		org, err := c.cliConnection.GetCurrentOrg()
		if err != nil {
			return cli.NewExitError(err, 3)
		}

		orgGUID = org.Guid
	} else {
		space, err := c.cliConnection.GetSpace(spaceName)
		if err != nil {
			return cli.NewExitError(err, 3)
		}

		spaceGUID = space.Guid
		orgGUID = space.Organization.Guid
	}

	spacePolicies, orgPolicies, err := getSpaceOrgPolicies(c.VaultClient, vaultMount, spaceGUID, orgGUID)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	b := color.New(color.Bold)

	b.Print("Space: \t\t")
	fmt.Printf("%s (%s)\n\n", spaceName, spaceGUID)

	b.Print("Mapped policies: \t\t")
	fmt.Printf("%s\n", spacePolicies)
	b.Print("Inherited org policies: \t")
	fmt.Printf("%s\n", orgPolicies)

	return nil
}

func (c *VaultAuthPluginClient) SetSpaceMapping(context *cli.Context) error {
	spaceName := context.String("space")
	policies := context.String("policies")
	vaultMount := context.GlobalString("mount")

	if policies == "" {
		return cli.NewExitError("You must specify one or more Vault policy names to map to", 3)
	}

	_, spaceGUID, err := getEffectiveSpaceNameGuid(spaceName, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = writeMapping(c.VaultClient, vaultMount, SpacesKey, spaceGUID, policies)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func (c *VaultAuthPluginClient) RemoveSpaceMapping(context *cli.Context) error {
	spaceName := context.String("space")
	vaultMount := context.GlobalString("mount")

	_, spaceGUID, err := getEffectiveSpaceNameGuid(spaceName, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = removeMapping(c.VaultClient, vaultMount, SpacesKey, spaceGUID)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func getEffectiveSpaceNameGuid(spaceName string, cliConnection plugin.CliConnection) (string, string, error) {
	var spaceGUID string

	if spaceName == "" {
		space, err := cliConnection.GetCurrentSpace()
		if err != nil {
			return "", "", err
		}

		spaceGUID = space.Guid
	} else {
		space, err := cliConnection.GetSpace(spaceName)
		if err != nil {
			return "", "", err
		}

		spaceGUID = space.Guid
	}

	return spaceName, spaceGUID, nil
}

func (c *VaultAuthPluginClient) GetOrgMapping(context *cli.Context) error {
	orgName := context.String("org")
	vaultMount := context.GlobalString("mount")

	printLoadingMessage()

	orgName, orgGUID, err := getEffectiveOrgNameGuid(orgName, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	orgPolicies, err := getOrgPolicies(c.VaultClient, vaultMount, orgGUID)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	b := color.New(color.Bold)

	b.Print("Org: \t\t")
	fmt.Printf("%s (%s)\n\n", orgName, orgGUID)

	b.Print("Mapped policies: \t\t")
	fmt.Printf("%s\n", orgPolicies)

	return nil
}

func (c *VaultAuthPluginClient) SetOrgMapping(context *cli.Context) error {
	orgName := context.String("org")
	policies := context.String("policies")
	vaultMount := context.GlobalString("mount")

	if policies == "" {
		return cli.NewExitError("You must specify one or more Vault policy names to map to", 3)
	}

	_, orgGUID, err := getEffectiveOrgNameGuid(orgName, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = writeMapping(c.VaultClient, vaultMount, OrgsKey, orgGUID, policies)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func (c *VaultAuthPluginClient) RemoveOrgMapping(context *cli.Context) error {
	orgName := context.String("org")
	vaultMount := context.GlobalString("mount")

	_, orgGUID, err := getEffectiveOrgNameGuid(orgName, c.cliConnection)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	err = removeMapping(c.VaultClient, vaultMount, OrgsKey, orgGUID)
	if err != nil {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func getEffectiveOrgNameGuid(orgName string, cliConnection plugin.CliConnection) (string, string, error) {
	var orgGUID string

	if orgName == "" {
		org, err := cliConnection.GetCurrentOrg()
		if err != nil {
			return "", "", err
		}

		orgName = org.Name
		orgGUID = org.Guid
	} else {
		org, err := cliConnection.GetOrg(orgName)
		if err != nil {
			return "", "", err
		}

		orgGUID = org.Guid
	}

	return orgName, orgGUID, nil
}

func printLoadingMessage() {
	fmt.Println("Retrieving mapping information from Vault...")
	fmt.Println("")
}

func writeMapping(c *api.Client, mount string, mapType string, guid string, policies string) error {
	mapData := map[string]interface{}{
		"value": policies,
	}

	_, err := c.Logical().Write("auth/"+mount+"/map/"+mapType+"/"+guid, mapData)
	if err != nil {
		return err
	}

	return nil
}

func removeMapping(c *api.Client, mount string, mapType string, guid string) error {
	_, err := c.Logical().Delete("auth/" + mount + "/map/" + mapType + "/" + guid)
	if err != nil {
		return err
	}

	return nil
}

func getMapping(c *api.Client, mount string, mapType string, guid string) (string, error) {
	secretValues, err := c.Logical().Read("auth/" + mount + "/map/" + mapType + "/" + guid)
	if err != nil {
		return "", err
	}

	if secretValues.Data == nil {
		return "", nil
	}

	if secretValues.Data["value"] == nil {
		return "", nil
	}

	return secretValues.Data["value"].(string), nil
}

func getSpaceInfo(space string, cliConnection plugin.CliConnection) (*SpaceInfo, error) {
	data, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v2/spaces/"+space+"/summary")
	if err != nil {
		return nil, err
	}

	var spaceInfo SpaceInfo

	err = json.Unmarshal([]byte(strings.Join(data, "")), &spaceInfo)
	if err != nil {
		return nil, err
	}

	return &spaceInfo, nil
}

func getAppSpaceOrgPolicies(c *api.Client, mount string, app string, space string, org string) (string, string, string, error) {
	appPolicies, err := getMapping(c, mount, AppsKey, app)
	if err != nil {
		return "", "", "", cli.NewExitError(err, 3)
	}

	spacePolicies, orgPolicies, err := getSpaceOrgPolicies(c, mount, space, org)
	if err != nil {
		return "", "", "", cli.NewExitError(err, 3)
	}

	return appPolicies, spacePolicies, orgPolicies, nil
}

func getSpaceOrgPolicies(c *api.Client, mount string, space string, org string) (string, string, error) {
	spacePolicies, err := getMapping(c, mount, SpacesKey, space)
	if err != nil {
		return "", "", cli.NewExitError(err, 3)
	}

	orgPolicies, err := getOrgPolicies(c, mount, org)
	if err != nil {
		return "", "", cli.NewExitError(err, 3)
	}

	return spacePolicies, orgPolicies, nil
}

func getOrgPolicies(c *api.Client, mount string, org string) (string, error) {
	orgPolicies, err := getMapping(c, mount, OrgsKey, org)
	if err != nil {
		return "", cli.NewExitError(err, 3)
	}

	return orgPolicies, nil
}
