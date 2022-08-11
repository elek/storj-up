// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd/files/templates"
	"storj.io/storj-up/pkg/common"
)

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <selector>",
		Short: "add more services to the docker compose file. " + selectorHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}
			templateProject, err := common.LoadComposeFromBytes(templates.ComposeTemplate)
			if err != nil {
				return err
			}
			updatedComposeProject, err := addToCompose(composeProject, templateProject, args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(addCmd())
}

func addToCompose(compose *types.Project, template *types.Project, services []string) (*types.Project, error) {
	if compose == nil {
		compose = &types.Project{Services: []types.ServiceConfig{}}
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}
	for _, service := range resolvedServices {
		if !common.ContainsService(compose.Services, service) {
			newService, err := template.GetService(service)
			if err != nil {
				return nil, err
			}
			compose.Services = append(compose.Services, newService)
		}
		if service == "prometheus" {
			err := ExtractFile("", "prometheus.yml", templates.PrometheusYaml)
			if err != nil {
				return nil, err
			}

		}
		if service == "geth" {
			files, err := templates.BlockchainFiles.ReadDir("test-blockchain")
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				fileContent, err := templates.BlockchainFiles.ReadFile("test-blockchain/" + file.Name())
				if err != nil {
					return nil, err
				}
				err = ExtractFile("test-blockchain", file.Name(), fileContent)
				if err != nil {
					return nil, err
				}
			}
			files, err = templates.CethFiles.ReadDir("cethacea")
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				fileContent, err := templates.CethFiles.ReadFile("cethacea/" + file.Name())
				if err != nil {
					return nil, err
				}
				err = ExtractFile("", "."+file.Name(), fileContent)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return compose, nil
}
