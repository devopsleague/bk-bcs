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

package bkmonitor

import (
	"context"
	"time"

	bcsmonitor "github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/component/bcs_monitor"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/storegw/bcs_system/source/base"
	"github.com/prometheus/prometheus/prompb"
)

// 和现在规范的名称保持一致
const provider = "BK-Monitor"

// handleNodeMetric xxx
func (m *BKMonitor) handleNodeMetric(ctx context.Context, projectId, clusterId, nodeName string, promql string, start,
	end time.Time, step time.Duration) ([]*prompb.TimeSeries, error) {
	nodeMatch, _, err := base.GetNodeMatchByName(ctx, clusterId, nodeName)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"clusterId":  clusterId,
		"ip":         nodeMatch,
		"fstype":     DISK_FSTYPE,
		"mountpoint": DISK_MOUNTPOINT,
		"provider":   PROVIDER,
	}

	matrix, _, err := bcsmonitor.QueryRangeMatrix(ctx, projectId, promql, params, start, end, step)
	if err != nil {
		return nil, err
	}

	return base.MatrixToSeries(matrix), nil
}

// GetNodeInfo 节点信息
func (m *BKMonitor) GetNodeInfo(ctx context.Context, projectId, clusterId, nodeName string, t time.Time) (*base.NodeInfo,
	error) {
	nodeMatch, ips, err := base.GetNodeMatchByName(ctx, clusterId, nodeName)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"clusterId":  clusterId,
		"ip":         nodeMatch,
		"fstype":     DISK_FSTYPE,
		"mountpoint": DISK_MOUNTPOINT,
		"provider":   PROVIDER,
	}

	info := &base.NodeInfo{}

	// 节点信息
	infoPromql := `cadvisor_version_info{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}`
	infoLabelSet, err := bcsmonitor.QueryLabelSet(ctx, projectId, infoPromql, params, t)
	if err != nil {
		return nil, err
	}
	info.DockerVersion = infoLabelSet["dockerVersion"]
	info.Release = infoLabelSet["kernelVersion"]
	info.Sysname = infoLabelSet["osVersion"]
	info.Provider = provider
	info.IP = ips

	promqlMap := map[string]string{
		"coreCount": `sum by (bk_instance) (count without(cpu, mode) (node_cpu_seconds_total{cluster_id="%<clusterId>s", mode="idle", bk_instance=~"%<ip>s", %<provider>s}))`,
		"memory":    `sum by (bk_instance) (node_memory_MemTotal_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s})`,
		"disk":      `sum by (bk_instance) (node_filesystem_size_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", fstype=~"%<fstype>s", mountpoint=~"%<mountpoint>s", %<provider>s})`,
	}

	result, err := bcsmonitor.QueryMultiValues(ctx, projectId, promqlMap, params, time.Now())
	if err != nil {
		return nil, err
	}

	info.CPUCount = result["coreCount"]
	info.Memory = result["memory"]
	info.Disk = result["disk"]

	return info, nil
}

// GetNodeCPUUsage 节点CPU使用率
func (m *BKMonitor) GetNodeCPUUsage(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `
		sum(irate(node_cpu_seconds_total{cluster_id="%<clusterId>s", mode!="idle", bk_instance=~"%<ip>s", %<provider>s}[2m])) /
		sum(count without(cpu, mode) (node_cpu_seconds_total{cluster_id="%<clusterId>s", mode="idle", bk_instance=~"%<ip>s", %<provider>s})) *
		100`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeMemoryUsage 内存使用率
func (m *BKMonitor) GetNodeMemoryUsage(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `
		(sum(node_memory_MemTotal_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}) -
        sum(node_memory_MemFree_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}) -
        sum(node_memory_Buffers_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}) -
        sum(node_memory_Cached_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}) +
        sum(node_memory_Shmem_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s})) /
        sum(node_memory_MemTotal_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}) *
        100`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeDiskUsage 节点磁盘使用率
func (m *BKMonitor) GetNodeDiskUsage(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `
		(sum(node_filesystem_size_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", fstype=~"%<fstype>s", mountpoint=~"%<mountpoint>s", %<provider>s}) -
        sum(node_filesystem_free_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", fstype=~"%<fstype>s", mountpoint=~"%<mountpoint>s", %<provider>s})) /
        sum(node_filesystem_size_bytes{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", fstype=~"%<fstype>s", mountpoint=~"%<mountpoint>s", %<provider>s}) *
        100`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeDiskioUsage 节点磁盘IO使用率
func (m *BKMonitor) GetNodeDiskioUsage(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `max(rate(node_disk_io_time_seconds_total{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}[2m]) * 100)`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodePodCount PodCount
func (m *BKMonitor) GetNodePodCount(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	// 注意 k8s 1.19 版本以前的 metrics 是 kubelet_running_pod_count
	promql :=
		`max by (bk_instance) (kubelet_running_pods{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s})`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeContainerCount 容器Count
func (m *BKMonitor) GetNodeContainerCount(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	// 注意 k8s 1.19 版本以前的 metrics 是 kubelet_running_container_count
	// https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.16.md 添加 running/exited/created/unknown label
	// https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md

	// container_state 常量定义 https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/container/runtime.go#L258
	// 使用不等于, 兼容高低版本
	promql :=
		`max by (bk_instance) (kubelet_running_containers{cluster_id="%<clusterId>s", container_state!="exited|created|unknown", bk_instance=~"%<ip>s", %<provider>s})`

	// 按版本兼容逻辑
	// if k8sclient.K8SLessThan(ctx, clusterId, "v1.19") {
	// 	promql = `max by (bk_instance) (kubelet_running_container_count{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s:.*", %<provider>s})`
	// }

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeNetworkTransmit 节点网络发送量
func (m *BKMonitor) GetNodeNetworkTransmit(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `max(rate(node_network_transmit_bytes_total{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}[2m]))`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}

// GetNodeNetworkReceive 节点网络接收
func (m *BKMonitor) GetNodeNetworkReceive(ctx context.Context, projectId, clusterId, nodeName string, start, end time.Time,
	step time.Duration) ([]*prompb.TimeSeries, error) {
	promql := `max(rate(node_network_receive_bytes_total{cluster_id="%<clusterId>s", bk_instance=~"%<ip>s", %<provider>s}[2m]))`

	return m.handleNodeMetric(ctx, projectId, clusterId, nodeName, promql, start, end, step)
}
