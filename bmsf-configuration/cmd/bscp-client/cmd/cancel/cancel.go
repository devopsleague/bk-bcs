/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cancel

import (
	"github.com/spf13/cobra"
)

var cancelCmd *cobra.Command

// init all resource create sub command.
func init() {
	cancelCmd = &cobra.Command{
		Use:     "cancel",
		Aliases: []string{"cc"},
		Short:   "Cancel specified release",
		Long:    "Cancel wrongly submitted or no longer used specified release",
	}
}

// InitCommands init all cancel commands.
func InitCommands() []*cobra.Command {
	// init all sub resource command.
	cancelCmd.AddCommand(cancelMultiReleaseCmd())
	return []*cobra.Command{cancelCmd}
}
