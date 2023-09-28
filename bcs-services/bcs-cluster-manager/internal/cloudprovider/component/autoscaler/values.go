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
 */

package autoscaler

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	cmproto "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/api/clustermanager"
	cmoptions "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/options"
)

const (
	templateName = "AutoScalerValues"

	encryptYes = "yes"
	encryptNo  = "no"
)

var (
	defaultTotalCpu = fmt.Sprintf("%d:%d", 0, 30000)   // nolint
	defaultTotalMem = fmt.Sprintf("%d:%d", 0, 6400000) // nolint
)

var valuesTemplate = `# Generated by bcs-cluster-manager

replicaCount: {{.ReplicaCount}}
 
namespace: {{.Namespace}}

{{- if .Registry }}
image:
  registry: {{.Registry}}
{{- end }}

command:
- ./bcs-cluster-autoscaler
- --v=4
- --stderrthreshold=info
- --namespace={{.Namespace}}
- --cloud-provider=bcs
- --estimator=clusterresource
{{- range .Nodes}}
- --nodes={{.}}
{{- end}}
- --expander={{.Expander}}
- --skip-nodes-with-local-storage={{.SkipNodesWithLocalStorage}}
- --skip-nodes-with-system-pods={{.SkipNodesWithSystemPods}}
- --scale-down-enabled={{.IsScaleDownEnable}}
- --max-empty-bulk-delete={{.MaxEmptyBulkDelete}}
- --scale-down-unneeded-time={{.ScaleDownUnneededTime}}
- --scale-down-utilization-threshold={{.ScaleDownUtilizationThreahold}}
- --scale-down-gpu-utilization-threshold={{.ScaleDownGpuUtilizationThreshold}}
- --ok-total-unready-count={{.OkTotalUnreadyCount}}
- --max-total-unready-percentage={{.MaxTotalUnreadyPercentage}}
- --scale-down-unready-time={{.ScaleDownUnreadyTime}}
- --buffer-resource-ratio={{.BufferResourceRatio}}
- --buffer-cpu-ratio={{.BufferResourceCpuRatio}}
- --buffer-mem-ratio={{.BufferResourceMemRatio}}
- --max-graceful-termination-sec={{.MaxGracefulTerminationSec}}
- --scan-interval={{.ScanInterval}}
- --max-node-provision-time={{.MaxNodeProvisionTime}}
- --max-node-startup-time={{.MaxNodeProvisionTime}}
- --max-node-start-schedule-time={{.MaxNodeProvisionTime}}
- --scale-up-from-zero={{.ScaleUpFromZero}}
- --scale-down-delay-after-add={{.ScaleDownDelayAfterAdd}}
- --scale-down-delay-after-delete={{.ScaleDownDelayAfterDelete}}
- --scale-down-delay-after-failure={{.ScaleDownDelayAfterFailure}}
- --ignore-daemonsets-utilization=true
- --new-pod-scale-up-delay={{.NewPodScaleUpDelay}}
- --expendable-pods-priority-cutoff={{.PodsPriorityCutoff}}

webhook:
  mode: "{{.WebhookMode}}"
  config: "{{.WebhookServer}}"
  token: "{{.WebhookToken}}"

env:
  apiAddress: "{{.APIAddress}}"
  token: "{{.Token}}"
  operator: "bcs"
  encryption: "{{.Encryption}}"

nodeSelector: null

resources: 
  requests:
    cpu: 1 
    memory: 2Gi
  limits:
    cpu: 4
    memory: 8Gi

tolerations:
- effect: "NoSchedule"
  key: "node-role.kubernetes.io/master"
  operator: "Exists"

affinity:
  nodeAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      preference:
        matchExpressions:
        - key: node-role.kubernetes.io/master
          operator: Exists
`

// AutoScalerValues is the values for the autoscaler application
type AutoScalerValues struct { // nolint
	Namespace                        string
	APIAddress                       string
	Token                            string
	ReplicaCount                     int
	Nodes                            []string
	Expander                         string
	SkipNodesWithLocalStorage        bool
	SkipNodesWithSystemPods          bool
	IsScaleDownEnable                bool
	MaxEmptyBulkDelete               uint32
	ScaleDownUnneededTime            time.Duration
	ScaleDownUtilizationThreahold    float64
	ScaleDownGpuUtilizationThreshold float64
	OkTotalUnreadyCount              uint32
	MaxTotalUnreadyPercentage        uint32
	ScaleDownUnreadyTime             time.Duration
	BufferResourceRatio              float64
	BufferResourceCpuRatio           float64
	BufferResourceMemRatio           float64
	MaxGracefulTerminationSec        uint32
	ScanInterval                     time.Duration
	MaxNodeProvisionTime             time.Duration
	ScaleUpFromZero                  bool
	ScaleDownDelayAfterAdd           time.Duration
	ScaleDownDelayAfterDelete        time.Duration
	ScaleDownDelayAfterFailure       time.Duration
	Encryption                       string
	WebhookMode                      string
	WebhookServer                    string
	WebhookToken                     string
	Registry                         string
	NewPodScaleUpDelay               uint32
	PodsPriorityCutoff               int32
}

// AutoScaler component paras
type AutoScaler struct {
	// NodeGroups cluster nodeGroup list
	NodeGroups []cmproto.NodeGroup
	// AutoScalingOption autoScaling deploy paras
	AutoScalingOption *cmproto.ClusterAutoScalingOption
	// Replicas
	Replicas int
}

// GetValues get autoScaler values
func (as *AutoScaler) GetValues() (string, error) {
	if len(as.NodeGroups) == 0 {
		as.Replicas = 0
	}
	tmpl, err := template.New(templateName).Parse(valuesTemplate)
	if err != nil {
		return "", err
	}
	// set vars
	op := cmoptions.GetGlobalCMOptions()
	values := AutoScalerValues{
		Namespace:    op.ComponentDeploy.AutoScaler.ReleaseNamespace,
		APIAddress:   op.ComponentDeploy.BCSAPIGateway,
		Token:        op.ComponentDeploy.Token,
		ReplicaCount: as.Replicas,
		Encryption:   encryptNo,
		Registry:     op.ComponentDeploy.Registry,
	}

	if cmoptions.GetEditionInfo().IsInnerEdition() {
		values.Encryption = encryptYes
	}

	values.Expander = as.AutoScalingOption.Expander
	values.SkipNodesWithLocalStorage = as.AutoScalingOption.SkipNodesWithLocalStorage
	values.SkipNodesWithSystemPods = as.AutoScalingOption.SkipNodesWithSystemPods
	values.IsScaleDownEnable = as.AutoScalingOption.IsScaleDownEnable
	values.MaxEmptyBulkDelete = as.AutoScalingOption.MaxEmptyBulkDelete
	values.ScaleDownUnneededTime = time.Duration(as.AutoScalingOption.ScaleDownUnneededTime) * time.Second
	values.ScaleDownUtilizationThreahold = float64(as.AutoScalingOption.ScaleDownUtilizationThreahold) / 100
	values.ScaleDownGpuUtilizationThreshold = float64(as.AutoScalingOption.ScaleDownGpuUtilizationThreshold) / 100
	values.OkTotalUnreadyCount = as.AutoScalingOption.OkTotalUnreadyCount
	values.MaxTotalUnreadyPercentage = as.AutoScalingOption.MaxTotalUnreadyPercentage
	values.ScaleDownUnreadyTime = time.Duration(as.AutoScalingOption.ScaleDownUnreadyTime) * time.Second
	values.BufferResourceRatio = float64(100-as.AutoScalingOption.BufferResourceRatio) / 100
	values.BufferResourceCpuRatio = float64(100-as.AutoScalingOption.BufferResourceCpuRatio) / 100
	values.BufferResourceMemRatio = float64(100-as.AutoScalingOption.BufferResourceMemRatio) / 100
	values.MaxGracefulTerminationSec = as.AutoScalingOption.MaxGracefulTerminationSec
	values.ScanInterval = time.Duration(as.AutoScalingOption.ScanInterval) * time.Second
	values.MaxNodeProvisionTime = time.Duration(as.AutoScalingOption.MaxNodeProvisionTime) * time.Second
	values.ScaleUpFromZero = as.AutoScalingOption.ScaleUpFromZero
	values.ScaleDownDelayAfterAdd = time.Duration(as.AutoScalingOption.ScaleDownDelayAfterAdd) * time.Second
	values.ScaleDownDelayAfterDelete = time.Duration(as.AutoScalingOption.ScaleDownDelayAfterDelete) * time.Second
	values.ScaleDownDelayAfterFailure = time.Duration(as.AutoScalingOption.ScaleDownDelayAfterFailure) * time.Second
	values.NewPodScaleUpDelay = as.AutoScalingOption.NewPodScaleUpDelay
	values.PodsPriorityCutoff = func() int32 {
		if as.AutoScalingOption.ExpendablePodsPriorityCutoff == 0 {
			return -10
		}

		return as.AutoScalingOption.ExpendablePodsPriorityCutoff
	}()

	if as.AutoScalingOption != nil && as.AutoScalingOption.Webhook != nil && as.AutoScalingOption.Webhook.Mode != "" {
		values.WebhookMode = as.AutoScalingOption.Webhook.Mode
	}
	if as.AutoScalingOption != nil && as.AutoScalingOption.Webhook != nil && as.AutoScalingOption.Webhook.Server != "" {
		values.WebhookServer = as.AutoScalingOption.Webhook.Server
	}
	if as.AutoScalingOption != nil && as.AutoScalingOption.Webhook != nil && as.AutoScalingOption.Webhook.Token != "" {
		values.WebhookToken = as.AutoScalingOption.Webhook.Token
	}

	for _, v := range as.NodeGroups {
		if len(v.NodeGroupID) != 0 && v.AutoScaling != nil {
			values.Nodes = append(values.Nodes,
				fmt.Sprintf("%d:%d:%s", v.AutoScaling.MinSize, v.AutoScaling.MaxSize, v.NodeGroupID))
		}
	}
	// parse
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		return "", err
	}
	return buf.String(), nil

}
