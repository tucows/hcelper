package backends

import (
	"fmt"
	"regexp"
	"strings"

	api "github.com/hashicorp/vault/api"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/mapstructure"
	types "github.com/tucows/hcelper/types"
)

func GetValidMounts(mountdata map[string]interface{}) (types.MountResponse, error) {
	var validResponse types.MountResponse
	for key, value := range mountdata {
		keyname := strings.TrimRight(key, "/")
		var sMount types.Mount

		// Get mount into a useful Struct
		mapstructure.Decode(value, &sMount)

		// Check if mount is nomad or consul
		match, err := regexp.MatchString(`nomad|consul`, sMount.Type)
		if err != nil {
			return types.MountResponse{}, err
		}
		if match == true {
			response := types.ValidMount{Name: keyname, Path: key, Type: sMount.Type}
			validResponse.Mounts = append(validResponse.Mounts, response)
		}
	}
	return validResponse, nil
}

func GetExportValues(vc *api.Client, mounts []types.ValidMount) (types.ExportResponse, error) {
	var compiledResponse types.ExportResponse

	for _, backendpath := range mounts {
		backendConf, err := vc.Logical().Read(fmt.Sprintf("%v/config/access", backendpath.Path))
		if err != nil {
			return types.ExportResponse{}, err
		}

		var mountAddress types.MountConfig
		mapstructure.Decode(backendConf.Data, &mountAddress)
		var tokenValue string
		// Get the appropriate backend path
		switch backendpath.Type {
		case "nomad":
			selectRole, err := RoleMenu(vc, backendpath)
			if err != nil {
				fmt.Printf("Error in Nomad role selection %v\n", err)
			}
			backendToken, err := vc.Logical().Read(fmt.Sprintf("%v/creds/%v", backendpath.Path, selectRole))
			if err != nil {
				return types.ExportResponse{}, err
			}
			var nomadToken types.NomadToken
			mapstructure.Decode(backendToken.Data, &nomadToken)
			tokenValue = nomadToken.Secret_ID

		case "consul":
			selectRole, err := RoleMenu(vc, backendpath)
			if err != nil {
				fmt.Printf("Error in Consul role selection %v\n", err)
			}
			backendToken, err := vc.Logical().Read(fmt.Sprintf("%v/creds/%v", backendpath.Path, selectRole))
			if err != nil {
				return types.ExportResponse{}, err
			}

			var mountToken types.ConsulToken
			mapstructure.Decode(backendToken.Data, &mountToken)
			tokenValue = mountToken.Token
		}
		compiledExport := types.CompiledExport{URL: mountAddress.Address, Token: tokenValue, Type: backendpath.Type}
		compiledResponse.Exports = append(compiledResponse.Exports, compiledExport)
	}

	return compiledResponse, nil
}

func GetEngineRoles(vc *api.Client, mount types.ValidMount) ([]string, error) {
	var rolecall *api.Secret
	var err error
	switch mount.Type {
	case "nomad":
		rolecall, err = vc.Logical().List(fmt.Sprintf("%v/role", mount.Path))
	case "consul":
		rolecall, err = vc.Logical().List(fmt.Sprintf("%v/roles", mount.Path))
	}

	if err != nil {
		return nil, err
	}
	var roles string
	roles = fmt.Sprintf("%v", rolecall.Data["keys"])
	roles = strings.TrimLeft(roles, "[")
	roles = strings.TrimRight(roles, "]")
	roleSlice := strings.Split(roles, " ")

	return roleSlice, nil
}

func RoleMenu(vc *api.Client, backend types.ValidMount) (string, error) {
	roles, _ := GetEngineRoles(vc, backend)
	label := fmt.Sprintf("Select your %v role:", backend.Type)
	nomadRolePrompt := promptui.Select{
		Label: label,
		Items: roles,
	}
	_, selectRole, err := nomadRolePrompt.Run()
	if err != nil {
		fmt.Printf("Env input failed %v\n", err)
	}

	return selectRole, err
}
