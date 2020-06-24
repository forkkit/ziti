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

package subcmd

import (
	"fmt"
	"github.com/openziti/fabric/pb/mgmt_pb"
	"github.com/openziti/foundation/channel2"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	snapshotDbClient = NewMgmtClient(snapshotDbCmd)
	Root.AddCommand(snapshotDbCmd)
}

var snapshotDbCmd = &cobra.Command{
	Use:   "snapshot-db",
	Short: "request the controller to create a snapshot of its database",
	Run:   snapshotDb,
}
var snapshotDbClient *mgmtClient

func snapshotDb(cmd *cobra.Command, args []string) {
	ch, err := snapshotDbClient.Connect()
	if err != nil {
		panic(err)
	}

	requestMsg := channel2.NewMessage(int32(mgmt_pb.ContentType_SnapshotDbRequestType), nil)
	responseMsg, err := ch.SendAndWaitWithTimeout(requestMsg, 5*time.Second)
	if err != nil {
		panic(err)
	}
	if responseMsg.ContentType == channel2.ContentTypeResultType {
		result := channel2.UnmarshalResult(responseMsg)
		if result.Success {
			fmt.Printf("\nsuccess\n\n")
		} else {
			fmt.Printf("\nfailure [%s]\n\n", result.Message)
		}
	} else {
		panic(fmt.Errorf("unexpected response type %v", responseMsg.ContentType))
	}
}
