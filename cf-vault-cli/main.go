package main

import (
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/hashicorp/vault/api"
	cli "gopkg.in/urfave/cli.v1"
)

type BasicPlugin struct {
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		log.Fatalln("Error: VAULT_ADDR must be set")
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		log.Fatalln("Error: VAULT_TOKEN must be set")
	}

	client, _ := api.NewClient(&api.Config{
		Address: vaultAddr,
	})

	client.SetToken(vaultToken)

	authClient := new(VaultAuthPluginClient)
	authClient.VaultClient = client
	authClient.cliConnection = cliConnection

	app := cli.NewApp()
	app.Name = "vault"
	app.Usage = "Manage how CloudFoundry orgs, spaces and applications map to Vault policies"
	app.Description = "Manage how CloudFoundry orgs, spaces and applications map to Vault policies"
	app.Version = "0.0.1"
	app.HideVersion = false

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mount, m",
			Value: "cf",
			Usage: "name of the Vault mount where the CF auth plugin is installed",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "apps",
			Usage: "operations related to apps",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get mappings for a given app",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "app, a",
							Usage: "Name of the CF app for which to retrieve mappings",
						},
					},
					Action: authClient.GetAppMapping,
				},
				{
					Name:  "set",
					Usage: "set an app mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "app, a",
							Usage: "Name of the CF app to map",
						},
						cli.StringFlag{
							Name:  "policies, p",
							Usage: "Vault policies to map to",
						},
					},
					Action: authClient.SetAppMapping,
				},
				{
					Name:  "remove",
					Usage: "remove an app mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "app, a",
							Usage: "Name of the CF app to unmap",
						},
					},
					Action: authClient.RemoveAppMapping,
				},
			},
		},
		{
			Name:  "spaces",
			Usage: "operations related to spaces",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get a space mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "space, s",
							Usage: "Name of the CF space for which to retrieve mappings",
						},
					},
					Action: authClient.GetSpaceMapping,
				},
				{
					Name:  "set",
					Usage: "set a space mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "space, s",
							Usage: "Name of the CF space to map",
						},
						cli.StringFlag{
							Name:  "policies, p",
							Usage: "Vault policies to map to",
						},
					},
					Action: authClient.SetSpaceMapping,
				},
				{
					Name:  "remove",
					Usage: "remove a space mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "space, s",
							Usage: "Name of the CF space to unmap",
						},
					},
					Action: authClient.RemoveSpaceMapping,
				},
			},
		},
		{
			Name:  "orgs",
			Usage: "operations related to orgs",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get a org mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "org, o",
							Usage: "Name of the CF org for which to retrieve mappings",
						},
					},
					Action: authClient.GetOrgMapping,
				},
				{
					Name:  "set",
					Usage: "set a org mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "org, o",
							Usage: "Name of the CF org to map",
						},
						cli.StringFlag{
							Name:  "policies, p",
							Usage: "Vault policies to map to",
						},
					},
					Action: authClient.SetOrgMapping,
				},
				{
					Name:  "remove",
					Usage: "remove an org mapping",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "org, o",
							Usage: "Name of the CF org to unmap",
						},
					},
					Action: authClient.RemoveOrgMapping,
				},
			},
		},
	}
	app.Action = noArgs

	app.Run(args)
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "vault",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "vault",
				HelpText: "Mapping CF apps to Vault policies",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "vault\n   cf vault",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BasicPlugin))
}

func noArgs(c *cli.Context) error {
	cli.ShowAppHelp(c)

	return cli.NewExitError("no commands provided", 2)
}

func scan(c *cli.Context) error {

	if c.Args().Present() {

		t := c.Args().First()
		fmt.Println("scanning", t)
		return nil
	}

	cli.ShowSubcommandHelp(c)
	return cli.NewExitError("no arguments for subcommand", 3)
}
