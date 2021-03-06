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
	"io"

	"github.com/pkg/errors"

	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
)

type updateServicePolicyOptions struct {
	commonOptions
	name          string
	serviceRoles  []string
	identityRoles []string
}

func newUpdateServicePolicyCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &updateServicePolicyOptions{
		commonOptions: commonOptions{
			CommonOptions: common.CommonOptions{Factory: f, Out: out, Err: errOut},
		},
	}

	cmd := &cobra.Command{
		Use:   "service-policy <idOrName>",
		Short: "updates a service-policy managed by the Ziti Edge Controller",
		Long:  "updates a service-policy managed by the Ziti Edge Controller",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runUpdateServicePolicy(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Set the name of the service policy")
	cmd.Flags().StringSliceVarP(&options.serviceRoles, "service-roles", "s", nil, "Service roles of the service policy")
	cmd.Flags().StringSliceVarP(&options.identityRoles, "identity-roles", "i", nil, "Identity roles of the service policy")
	options.AddCommonFlags(cmd)

	return cmd
}

func runUpdateServicePolicy(o *updateServicePolicyOptions) error {
	id, err := mapNameToID("service-policies", o.Args[0])
	if err != nil {
		return err
	}

	serviceRoles, err := convertNamesToIds(o.serviceRoles, "services")
	if err != nil {
		return err
	}

	identityRoles, err := convertNamesToIds(o.identityRoles, "identities")
	if err != nil {
		return err
	}

	entityData := gabs.New()
	change := false

	if o.Cmd.Flags().Changed("name") {
		setJSONValue(entityData, o.name, "name")
		change = true
	}

	if o.Cmd.Flags().Changed("service-roles") {
		setJSONValue(entityData, serviceRoles, "serviceRoles")
		change = true
	}

	if o.Cmd.Flags().Changed("identity-roles") {
		setJSONValue(entityData, identityRoles, "identityRoles")
		change = true
	}

	if !change {
		return errors.New("no change specified. must specify at least one attribute to change")
	}

	_, err = patchEntityOfType(fmt.Sprintf("service-policies/%v", id), entityData.String(), &o.commonOptions)
	return err
}
