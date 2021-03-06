/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package edge_controller

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
)

type createServiceOptions struct {
	commonOptions
	terminatorStrategy string
	tags               map[string]string
	roleAttributes     []string
	configs            []string
	requireEncryption  bool
	optionalEncryption bool
}

// newCreateServiceCmd creates the 'edge controller create service local' command for the given entity type
func newCreateServiceCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &createServiceOptions{
		commonOptions: commonOptions{
			CommonOptions: common.CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
		tags: make(map[string]string),
	}

	cmd := &cobra.Command{
		Use:   "service <name>",
		Short: "creates a service managed by the Ziti Edge Controller",
		Long:  "creates a service managed by the Ziti Edge Controller",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runCreateService(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringToStringVarP(&options.tags, "tags", "t", nil, "Add tags to service definition")
	cmd.Flags().StringSliceVarP(&options.roleAttributes, "role-attributes", "a", nil, "Role attributes of the new service")
	cmd.Flags().StringSliceVarP(&options.configs, "configs", "c", nil, "Configuration id or names to be associated with the new service")
	cmd.Flags().StringVar(&options.terminatorStrategy, "terminator-strategy", "", "Specifies the terminator strategy for the service")
	cmd.Flags().BoolVarP(&options.optionalEncryption, "encryption-optional", "o", false, "Sets end-to-end encryption for the service to be optional, defaults to required")
	options.AddCommonFlags(cmd)

	return cmd
}

// runCreateService implements the command to create a service
func runCreateService(o *createServiceOptions) (err error) {

	configs, err := mapNamesToIDs("configs", o.configs...)
	if err != nil {
		return err
	}

	entityData := gabs.New()
	setJSONValue(entityData, o.Args[0], "name")
	if o.terminatorStrategy != "" {
		setJSONValue(entityData, o.terminatorStrategy, "terminatorStrategy")
	}

	setJSONValue(entityData, !o.optionalEncryption, "encryptionRequired")

	setJSONValue(entityData, o.roleAttributes, "roleAttributes")
	setJSONValue(entityData, configs, "configs")
	setJSONValue(entityData, o.tags, "tags")

	result, err := createEntityOfType("services", entityData.String(), &o.commonOptions)

	if err != nil {
		panic(err)
	}

	serviceId := result.S("data", "id").Data()

	if _, err = fmt.Fprintf(o.Out, "%v\n", serviceId); err != nil {
		panic(err)
	}

	return err
}
