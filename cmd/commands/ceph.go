/*
Copyright 2023 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"fmt"
	"strings"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// CephCmd represents the ceph command
var CephCmd = &cobra.Command{
	Use:                "ceph",
	Short:              "call a 'ceph' CLI command with arbitrary args",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clientsets := GetClientsets(cmd.Context())
		VerifyOperatorPodIsRunning(cmd.Context(), clientsets, OperatorNamespace, CephClusterNamespace)
		logging.Info("running 'ceph' command with args: %v", args)
		if args[0] == "daemon" {
			c := fmt.Sprintf("CEPH_ARGS='' %s", cmd.Use)
			if len(args) > 1 {
				if strings.HasPrefix(args[1], "osd.") {
					exec.RunCommandInOsdPod(cmd.Context(), clientsets, args[1], c, args, CephClusterNamespace, false, true)
				}
			}
		}
		exec.RunCommandInOperatorPod(cmd.Context(), clientsets, cmd.Use, args, OperatorNamespace, CephClusterNamespace, false, true)
	},
}
