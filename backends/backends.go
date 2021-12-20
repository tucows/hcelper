package backends

import (
	"fmt"
	"regexp"
	"strings"

	api "github.com/hashicorp/vault/api"
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
			// NO. BAD. LIST CREDS AND PROMPT ON NEXT REVISION
			backendToken, err := vc.Logical().Read(fmt.Sprintf("%v/creds/operators", backendpath.Path))
			if err != nil {
				return types.ExportResponse{}, err
			}

			var nomadToken types.NomadToken
			mapstructure.Decode(backendToken.Data, &nomadToken)
			tokenValue = nomadToken.Secret_ID

		case "consul":
			// NO. BAD. LIST CREDS AND PROMPT ON NEXT REVISION
			backendToken, err := vc.Logical().Read(fmt.Sprintf("%v/creds/operators", backendpath.Path))
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
