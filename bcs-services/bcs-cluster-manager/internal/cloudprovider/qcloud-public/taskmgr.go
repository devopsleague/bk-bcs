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

package qcloud

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	proto "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/api/clustermanager"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/common"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/qcloud-public/tasks"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/template"
	icommon "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/common"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/utils"
)

var taskMgr sync.Once

func init() {
	taskMgr.Do(func() {
		cloudprovider.InitTaskManager(cloudName, newtask())
	})
}

func newtask() *Task {
	task := &Task{
		works: make(map[string]interface{}),
	}

	// init qcloud cluster-manager task, may be call bkops interface to call extra operation

	// import task
	task.works[importClusterNodesStep.StepMethod] = tasks.ImportClusterNodesTask
	task.works[registerClusterKubeConfigStep.StepMethod] = tasks.RegisterClusterKubeConfigTask

	// create cluster task
	task.works[createTKEClusterStep.StepMethod] = tasks.CreateTkeClusterTask
	task.works[checkTKEClusterStatusStep.StepMethod] = tasks.CheckTkeClusterStatusTask
	task.works[checkCreateClusterNodeStatusStep.StepMethod] = tasks.CheckCreateClusterNodeStatusTask
	task.works[registerTkeClusterKubeConfigStep.StepMethod] = tasks.RegisterTkeClusterKubeConfigTask
	task.works[updateCreateClusterDBInfoStep.StepMethod] = tasks.UpdateCreateClusterDBInfoTask

	// delete cluster task
	task.works[deleteTKEClusterStep.StepMethod] = tasks.DeleteTKEClusterTask
	task.works[cleanClusterDBInfoStep.StepMethod] = tasks.CleanClusterDBInfoTask

	// add node to cluster
	task.works[addNodesToClusterStep.StepMethod] = tasks.AddNodesToClusterTask
	task.works[checkAddNodesStatusStep.StepMethod] = tasks.CheckAddNodesStatusTask
	task.works[updateAddNodeDBInfoStep.StepMethod] = tasks.UpdateNodeDBInfoTask

	// remove node from cluster
	task.works[removeNodesFromClusterStep.StepMethod] = tasks.RemoveNodesFromClusterTask
	task.works[updateRemoveNodeDBInfoStep.StepMethod] = tasks.UpdateRemoveNodeDBInfoTask

	// init qcloud node-group task

	// autoScaler task
	// task.works[ensureAutoScalerStep.StepMethod] = tasks.EnsureAutoScalerTask

	// create nodeGroup task
	task.works[createCloudNodeGroupStep.StepMethod] = tasks.CreateCloudNodeGroupTask
	task.works[checkCloudNodeGroupStatusStep.StepMethod] = tasks.CheckCloudNodeGroupStatusTask
	// task.works[updateCreateNodeGroupDBInfoTask] = tasks.UpdateCreateNodeGroupDBInfoTask

	// delete nodeGroup task
	task.works[deleteNodeGroupStep.StepMethod] = tasks.DeleteCloudNodeGroupTask
	// task.works[updateDeleteNodeGroupDBInfoTask] = tasks.UpdateDeleteNodeGroupDBInfoTask

	// clean node in nodeGroup task
	task.works[cleanNodeGroupNodesStep.StepMethod] = tasks.CleanNodeGroupNodesTask
	task.works[checkClusterCleanNodsStep.StepMethod] = tasks.CheckClusterCleanNodsTask
	// task.works[checkCleanNodeGroupNodesStatusTask] = tasks.CheckCleanNodeGroupNodesStatusTask
	// task.works[updateCleanNodeGroupNodesDBInfoTask] = tasks.UpdateCleanNodeGroupNodesDBInfoTask

	// update desired nodes task
	task.works[applyInstanceMachinesStep.StepMethod] = tasks.ApplyInstanceMachinesTask
	task.works[checkClusterNodesStatusStep.StepMethod] = tasks.CheckClusterNodesStatusTask

	// move nodes to nodeGroup task

	return task
}

// Task background task manager
type Task struct {
	works map[string]interface{}
}

// Name get cloudName
func (t *Task) Name() string {
	return cloudName
}

// GetAllTask register all backgroup task for worker running
func (t *Task) GetAllTask() map[string]interface{} {
	return t.works
}

// BuildCreateClusterTask build create cluster task
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildCreateClusterTask(cls *proto.Cluster, opt *cloudprovider.CreateClusterOption) ( // nolint
	*proto.Task, error) {
	// create cluster currently only has three steps:
	// 0. check if need to generate master instance. you need to call cvm api to produce master instance if necessary.
	//    but we only support add existed instance to cluster as master currently.
	// 1. call qcloud CreateTKECluster to create tke cluster
	// 2. call GetTKECluster to check cluster run status(cluster status: Running Creating Abnormal))
	// 3. update cluster DB info when create cluster successful
	// may be need to call external previous or behind operation by bkops

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildCreateClusterTask cluster info empty")
	}
	if opt == nil || opt.Cloud == nil {
		return nil, fmt.Errorf("BuildCreateClusterTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.CreateCluster),
		TaskName:       cloudprovider.CreateClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	// generate taskName
	taskName := fmt.Sprintf(createClusterTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	createClusterTask := &CreateClusterTaskOption{
		Cluster: cls, MasterNodes: opt.MasterNodes, WorkerNodes: opt.WorkerNodes, NodeTemplate: opt.NodeTemplate}

	// step1: createTKECluster and return clusterID inject common paras
	createClusterTask.BuildCreateClusterStep(task)
	// step2: check cluster status by clusterID
	createClusterTask.BuildCheckClusterStatusStep(task)
	// step3: check cluster nodes status
	createClusterTask.BuildCheckClusterNodesStatusStep(task)
	// step4: qcloud-public register cluster kubeConfig
	createClusterTask.BuildRegisterClsKubeConfigStep(task)
	// step5: qcloud-public import cluster nodes
	createClusterTask.BuildImportClusterNodesStep(task)
	// step5: install cluster watch component
	common.BuildWatchComponentTaskStep(task, cls, "")
	// step4: 若需要则设置节点注解
	common.BuildNodeAnnotationsTaskStep(task, cls.ClusterID, nil, func() map[string]string {
		if opt.NodeTemplate != nil && len(opt.NodeTemplate.GetAnnotations()) > 0 {
			return opt.NodeTemplate.GetAnnotations()
		}
		return nil
	}())

	// step6: install gse agent
	common.BuildInstallGseAgentTaskStep(task, &common.GseInstallInfo{
		ClusterId:  cls.ClusterID,
		BusinessId: cls.BusinessID,
		CloudArea:  cls.GetClusterBasicSettings().GetArea(),
		User:       cls.GetNodeSettings().GetWorkerLogin().GetInitLoginUsername(),
		Passwd:     cls.GetNodeSettings().GetWorkerLogin().GetInitLoginPassword(),
		KeyInfo:    cls.GetNodeSettings().GetWorkerLogin().GetKeyPair(),
		Port: func() string {
			exist := checkClusterOsNameInWhiteImages(cls, &opt.CommonOption)
			if exist {
				return fmt.Sprintf("%v", utils.ConnectPort)
			}

			return ""
		}(),
		AllowReviseCloudId: icommon.True,
	}, cloudprovider.WithStepAllowSkip(true))

	// step7: transfer host module
	common.BuildTransferHostModuleStep(task, cls.BusinessID, cls.GetClusterBasicSettings().GetModule().
		GetWorkerModuleID(), cls.GetClusterBasicSettings().GetModule().GetMasterModuleID())

	// step8: 业务后置自定义流程: 支持标准运维任务 或者 后置脚本
	if opt.NodeTemplate != nil && len(opt.NodeTemplate.UserScript) > 0 {
		common.BuildJobExecuteScriptStep(task, common.JobExecParas{
			ClusterID: cls.ClusterID,
			Content:   opt.NodeTemplate.UserScript,
			// dynamic node ips
			NodeIps:   "",
			Operator:  opt.Operator,
			StepName:  common.PostInitStepJob,
			Translate: common.PostInitJob,
		})
	}
	// business post define sops task or script
	if opt.NodeTemplate != nil && opt.NodeTemplate.ScaleOutExtraAddons != nil {
		err := template.BuildSopsFactory{
			StepName: template.UserAfterInit,
			Cluster:  cls,
			Extra: template.ExtraInfo{
				// dynamic node ips
				NodeIPList:      "",
				NodeOperator:    opt.Operator,
				ShowSopsUrl:     true,
				TranslateMethod: template.UserPostInit,
			}}.BuildSopsStep(task, opt.NodeTemplate.ScaleOutExtraAddons, false)
		if err != nil {
			return nil, fmt.Errorf("BuildScalingNodesTask business BuildBkSopsStepAction failed: %v", err)
		}
	}

	// step9: update DB info by cluster data
	createClusterTask.BuildUpdateTaskStatusStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildCreateClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.CreateClusterJob.String()

	if len(opt.WorkerNodes) > 0 {
		task.CommonParams[cloudprovider.WorkerNodeIPsKey.String()] = strings.Join(opt.WorkerNodes, ",")
	}
	if len(opt.MasterNodes) > 0 {
		task.CommonParams[cloudprovider.MasterNodeIPsKey.String()] = strings.Join(opt.MasterNodes, ",")
	}

	return task, nil
}

// BuildCreateVirtualClusterTask build create virtual cluster task
func (t *Task) BuildCreateVirtualClusterTask(cls *proto.Cluster,
	opt *cloudprovider.CreateVirtualClusterOption) (*proto.Task, error) {
	// create virtual cluster by host cluster namespace
	// 1. hostCluster create namespace or exist in cluster
	// 2. hostCluster deploy vcluster/agent component
	// 3. wait subCluster kube-agent deployed
	// 4. subCluster deploy k8s-watch component

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildCreateVirtualClusterTask cluster info empty")
	}
	if opt == nil || opt.Cloud == nil || opt.HostCluster == nil || opt.Namespace == nil {
		return nil, fmt.Errorf("BuildCreateVirtualClusterTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.CreateVirtualCluster),
		TaskName:       cloudprovider.CreateVirtualClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	// generate taskName
	taskName := fmt.Sprintf(createVirtualClusterTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	createTask := CreateVirtualClusterTask{
		Cluster:     cls,
		HostCluster: opt.HostCluster,
		Namespace:   opt.Namespace,
	}
	createTask.BuildCreateNamespaceStep(task)
	createTask.BuildCreateResourceQuotaStep(task)
	createTask.BuildInstallVclusterStep(task)
	createTask.BuildCheckAgentStatusStep(task)
	createTask.BuildInstallWatchStep(task)

	// step6: 系统初始化 postAction bkops, platform run default steps
	if opt.Cloud != nil && opt.Cloud.ClusterManagement != nil && opt.Cloud.ClusterManagement.CreateCluster != nil {
		err := template.BuildSopsFactory{
			StepName: template.SystemInit,
			Cluster:  cls,
			Extra: template.ExtraInfo{
				NodeOperator: opt.Operator,
			}}.BuildSopsStep(task, opt.Cloud.ClusterManagement.CreateCluster, false)
		if err != nil {
			return nil, fmt.Errorf("BuildCreateVirtualClusterTask BuildBkSopsStepAction failed: %v", err)
		}
	}

	createTask.BuildUpdateTaskStatusStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildCreateVirtualClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.CreateVirtualClusterJob.String()

	return task, nil
}

// BuildImportClusterTask build import cluster task
func (t *Task) BuildImportClusterTask(cls *proto.Cluster, opt *cloudprovider.ImportClusterOption) (*proto.Task, error) {
	// import cluster currently only has two steps:
	// 0. import cluster: call TKEInterface import cluster master and node instances from cloud(clusterID or kubeConfig)
	// 1. internal install bcs-k8s-watch & agent service; external import qcloud kubeConfig
	// may be need to call external previous or behind operation by bkops

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildImportClusterTask cluster info empty")
	}
	if opt == nil || opt.Cloud == nil {
		return nil, fmt.Errorf("BuildImportClusterTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.ImportCluster),
		TaskName:       cloudprovider.ImportClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	// generate taskName
	taskName := fmt.Sprintf(importClusterTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	importTask := ImportClusterTaskOption{Cluster: cls}
	// step0: import cluster nodes step
	importTask.BuildImportClusterNodesStep(task)
	// step1: import cluster registerKubeConfigStep
	importTask.BuildRegisterKubeConfigStep(task)
	// step2: install cluster watch component
	common.BuildWatchComponentTaskStep(task, cls, "")
	// step3: install image pull secret addon if config
	common.BuildInstallImageSecretAddonTaskStep(task, cls)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildImportClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.ImportClusterJob.String()

	return task, nil
}

// BuildDeleteVirtualClusterTask build delete virtual cluster task
func (t *Task) BuildDeleteVirtualClusterTask(cls *proto.Cluster,
	opt *cloudprovider.DeleteVirtualClusterOption) (*proto.Task, error) {
	// delete cluster has three steps:
	// 1. delete virtual cluster
	// 2. delete vcluster namespace in hostCluster
	// 3. delete cluster relative data && cluster credential

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildDeleteVirtualClusterTask cluster info empty")
	}
	if opt == nil || opt.Operator == "" || opt.Cloud == nil || opt.HostCluster == nil || opt.Namespace == nil {
		return nil, fmt.Errorf("BuildDeleteVirtualClusterTask TaskOptions is lost")
	}

	// init task information
	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.DeleteVirtualCluster),
		TaskName:       cloudprovider.DeleteVirtualClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	taskName := fmt.Sprintf(deleteVirtualClusterTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName
	task.CommonParams[cloudprovider.UserKey.String()] = opt.Operator

	// setting all steps details
	deleteVirtualClusterTask := &DeleteVirtualClusterTaskOption{
		Cluster:     cls,
		Cloud:       opt.Cloud,
		HostCluster: opt.HostCluster,
		Namespace:   opt.Namespace,
	}
	deleteVirtualClusterTask.BuildUninstallVClusterStep(task)
	deleteVirtualClusterTask.BuildDeleteNamespaceStep(task)
	deleteVirtualClusterTask.BuildCleanClusterDBInfoStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildDeleteVirtualClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.DeleteVirtualClusterJob.String()
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator

	return task, nil
}

// BuildDeleteClusterTask build deleteCluster task
func (t *Task) BuildDeleteClusterTask(cls *proto.Cluster, opt *cloudprovider.DeleteClusterOption) (*proto.Task, error) {
	// delete cluster has three steps:
	// 1. clean nodeGroup nodes and delete nodeGroup Info
	// 2. call qcloud DeleteTKECluster to delete tke cluster
	// 3. clean DB cluster info and associated data info when delete successful

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildDeleteClusterTask cluster info empty")
	}
	if opt == nil || opt.Operator == "" || opt.Cloud == nil || opt.Cluster == nil {
		return nil, fmt.Errorf("BuildDeleteClusterTask TaskOptions is lost")
	}

	// init task information
	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.DeleteCluster),
		TaskName:       cloudprovider.DeleteClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	taskName := fmt.Sprintf(deleteClusterTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName
	task.CommonParams[cloudprovider.UserKey.String()] = opt.Operator

	// setting all steps details
	deleteClusterTask := &DeleteClusterTaskOption{
		Cluster:           cls,
		DeleteMode:        opt.DeleteMode.String(),
		LastClusterStatus: opt.LatsClusterStatus,
	}
	// step1: DeleteTKECluster delete tke cluster
	deleteClusterTask.BuildDeleteTKEClusterStep(task)
	// step2: update cluster DB info and associated data
	deleteClusterTask.BuildCleanClusterDBInfoStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildDeleteClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.DeleteClusterJob.String()
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator

	return task, nil
}

// BuildAddNodesToClusterTask build addNodes task
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildAddNodesToClusterTask(cls *proto.Cluster, nodes []*proto.Node, // nolint
	opt *cloudprovider.AddNodesOption) (*proto.Task, error) {
	// addNodesToCluster has three steps:
	// 1. call qcloud AddExistedInstancesToCluster to add node
	// 2. call qcloud QueryTkeClusterInstances to check instance status(running initializing failed))
	// 3. update node DB info when add node successful
	// may be need to call external previous or behind operation by bkops

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildAddNodesToClusterTask cluster info empty")
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("BuildAddNodesToClusterTask lost nodes info")
	}

	if opt == nil || opt.Cloud == nil || opt.Operator == "" {
		return nil, fmt.Errorf("BuildAddNodesToClusterTask TaskOptions is lost")
	}

	if opt.Login == nil || (opt.Login.GetInitLoginPassword() == "" && opt.Login.GetKeyPair().GetKeyID() == "") {
		return nil, fmt.Errorf("BuildAddNodesToClusterTask login info empty")
	}

	// format node IPs
	nodeIDs := make([]string, 0)
	nodeIPs := make([]string, 0)
	for i := range nodes {
		nodeIPs = append(nodeIPs, nodes[i].InnerIP)
		nodeIDs = append(nodeIDs, nodes[i].NodeID)
	}

	// init task information
	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.AddNodesToCluster),
		TaskName:       cloudprovider.AddNodesToClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeIPList:     nodeIPs,
	}
	taskName := fmt.Sprintf(tkeAddNodeTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	addNodesTask := &AddNodesToClusterTaskOption{
		Cluster:      cls,
		Cloud:        opt.Cloud,
		NodeTemplate: opt.NodeTemplate,
		NodeIPs:      nodeIPs,
		NodeIDs:      nodeIDs,
		Operator:     opt.Operator,
		NodeSchedule: opt.NodeSchedule,
		Login:        opt.Login,
	}
	// step0: addNodesToTKECluster add node to cluster
	addNodesTask.BuildAddNodesToClusterStep(task)
	// step1: check cluster add node status
	addNodesTask.BuildCheckAddNodesStatusStep(task)
	// step2: update DB node info by instanceIP
	addNodesTask.BuildUpdateAddNodeDBInfoStep(task)
	// step3: install gse agent
	common.BuildInstallGseAgentTaskStep(task, &common.GseInstallInfo{
		ClusterId:  cls.ClusterID,
		BusinessId: cls.BusinessID,
		User:       opt.Login.GetInitLoginUsername(),
		Passwd:     opt.Login.GetInitLoginPassword(),
		KeyInfo:    opt.Login.GetKeyPair(),
		Port: func() string {
			exist := checkIfWhiteImageOsNames(&cloudprovider.ClusterGroupOption{
				CommonOption: opt.CommonOption,
				Cluster:      cls,
			})
			if exist {
				return fmt.Sprintf("%v", utils.ConnectPort)
			}

			return ""
		}(),
	})
	if cls.GetBusinessID() != "" && cls.GetClusterBasicSettings().GetModule().GetWorkerModuleID() != "" {
		common.BuildTransferHostModuleStep(task, cls.GetBusinessID(),
			cls.GetClusterBasicSettings().GetModule().GetWorkerModuleID(), "")
	}

	// step3: 业务后置自定义流程: 支持标准运维任务 或者 后置脚本
	if opt.NodeTemplate != nil && len(opt.NodeTemplate.UserScript) > 0 {
		common.BuildJobExecuteScriptStep(task, common.JobExecParas{
			ClusterID:        cls.ClusterID,
			Content:          opt.NodeTemplate.UserScript,
			NodeIps:          strings.Join(nodeIPs, ","),
			Operator:         opt.Operator,
			StepName:         common.PostInitStepJob,
			AllowSkipJobTask: opt.NodeTemplate.AllowSkipScaleOutWhenFailed,
			Translate:        common.PostInitJob,
		})
	}

	// business post define sops task or script
	if opt.NodeTemplate != nil && opt.NodeTemplate.ScaleOutExtraAddons != nil {
		err := template.BuildSopsFactory{
			StepName: template.UserAfterInit,
			Cluster:  cls,
			Extra: template.ExtraInfo{
				NodeIPList:      strings.Join(nodeIPs, ","),
				NodeOperator:    opt.Operator,
				ShowSopsUrl:     true,
				TranslateMethod: template.UserPostInit,
			}}.BuildSopsStep(task, opt.NodeTemplate.ScaleOutExtraAddons, false)
		if err != nil {
			return nil, fmt.Errorf("BuildScalingNodesTask business BuildBkSopsStepAction failed: %v", err)
		}
	}

	// step4: 若需要则设置节点注解
	addNodesTask.BuildNodeAnnotationsStep(task)
	// step5: 设置平台公共标签
	// addNodesTask.BuildNodeLabelsStep(task)
	// step6: 设置节点可调度状态
	addNodesTask.BuildUnCordonNodesStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildAddNodesToClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator
	task.CommonParams[cloudprovider.UserKey.String()] = opt.Operator

	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.AddNodeJob.String()
	task.CommonParams[cloudprovider.NodeIPsKey.String()] = strings.Join(nodeIPs, ",")
	task.CommonParams[cloudprovider.NodeIDsKey.String()] = strings.Join(nodeIDs, ",")

	return task, nil
}

// BuildRemoveNodesFromClusterTask build removeNodes task
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildRemoveNodesFromClusterTask(cls *proto.Cluster, nodes []*proto.Node, // nolint
	opt *cloudprovider.DeleteNodesOption) (*proto.Task, error) {
	// removeNodesFromCluster has two steps:
	// 1. call qcloud DeleteTkeClusterInstance to delete node
	// 2. update node DB info when delete node successful
	// may be need to call external previous or behind operation by bkops

	// validate request params
	if cls == nil {
		return nil, fmt.Errorf("BuildRemoveNodesFromClusterTask cluster info empty")
	}
	if opt == nil || opt.Cloud == nil {
		return nil, fmt.Errorf("BuildRemoveNodesFromClusterTask TaskOptions is lost")
	}

	// format all nodes InnerIP
	var (
		nodeIPs        []string
		nodeIDs        []string
		terminateNodes []string
		retainNodes    []string
	)
	for _, node := range nodes {
		nodeIPs = append(nodeIPs, node.InnerIP)
		nodeIDs = append(nodeIDs, node.NodeID)
		// sort different node by charge type
		switch opt.DeleteMode {
		case cloudprovider.Terminate.String():
			if node.ChargeType == icommon.POSTPAIDBYHOUR {
				terminateNodes = append(terminateNodes, node.InnerIP)
			} else {
				retainNodes = append(retainNodes, node.InnerIP)
			}
		case cloudprovider.Retain.String():
			retainNodes = append(retainNodes, node.InnerIP)
		default:
		}
	}

	// init task information
	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.RemoveNodesFromCluster),
		TaskName:       cloudprovider.RemoveNodesFromClusterTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      cls.ClusterID,
		ProjectID:      cls.ProjectID,
		Creator:        opt.Operator,
		Updater:        opt.Operator,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeIPList:     nodeIPs,
	}
	// generate taskName
	taskName := fmt.Sprintf(tkeCleanNodeTaskTemplate, cls.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	removeNodesTask := &RemoveNodesFromClusterTaskOption{
		Cluster:        cls,
		Cloud:          opt.Cloud,
		DeleteMode:     opt.DeleteMode,
		NodeIPs:        nodeIPs,
		NodeIDs:        nodeIDs,
		terminateNodes: terminateNodes,
		retainNodes:    retainNodes,
	}

	// step0: cordon nodes
	removeNodesTask.BuildCordonNodesStep(task)

	// step1: 业务自定义缩容流程: 支持 缩容节点前置脚本和前置标准运维流程
	if opt.NodeTemplate != nil && len(opt.NodeTemplate.ScaleInPreScript) > 0 {
		common.BuildJobExecuteScriptStep(task, common.JobExecParas{
			ClusterID:        cls.ClusterID,
			Content:          opt.NodeTemplate.ScaleInPreScript,
			NodeIps:          strings.Join(nodeIPs, ","),
			Operator:         opt.Operator,
			StepName:         common.PreInitStepJob,
			AllowSkipJobTask: opt.NodeTemplate.AllowSkipScaleInWhenFailed,
			Translate:        common.PreInitJob,
		})
	}
	// business define sops task
	if opt.NodeTemplate != nil && opt.NodeTemplate.ScaleInExtraAddons != nil {
		err := template.BuildSopsFactory{
			StepName: template.UserPreInit,
			Cluster:  cls,
			Extra: template.ExtraInfo{
				NodeIPList:      strings.Join(nodeIPs, ","),
				NodeOperator:    opt.Operator,
				ShowSopsUrl:     true,
				TranslateMethod: template.UserBeforeInit,
			}}.BuildSopsStep(task, opt.NodeTemplate.ScaleInExtraAddons, true)
		if err != nil {
			return nil, fmt.Errorf("BuildRemoveNodesFromClusterTask business "+
				"BuildBkSopsStepAction failed: %v", err)
		}
	}

	// step2: removeNodesFromTKECluster remove nodes
	removeNodesTask.BuildRemoveNodesFromClusterStep(task)
	// step3: check deleting node status
	removeNodesTask.BuildCheckClusterCleanNodesStep(task)
	// step4: update node DB info
	removeNodesTask.BuildUpdateRemoveNodeDBInfoStep(task)
	// step5: remove nodes from cmdb
	// common.BuildRemoveHostStep(task, cls.BusinessID, nodeIPs)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildRemoveNodesFromClusterTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.DeleteNodeJob.String()
	task.CommonParams[cloudprovider.NodeIPsKey.String()] = strings.Join(nodeIPs, ",")
	task.CommonParams[cloudprovider.NodeIDsKey.String()] = strings.Join(nodeIDs, ",")
	task.CommonParams[cloudprovider.OperatorKey.String()] = opt.Operator

	return task, nil
}

// BuildCreateNodeGroupTask build create node group task
func (t *Task) BuildCreateNodeGroupTask(group *proto.NodeGroup, opt *cloudprovider.CreateNodeGroupOption) (
	*proto.Task, error) {
	// validate request params
	if group == nil {
		return nil, fmt.Errorf("BuildCreateNodeGroupTask group info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildCreateNodeGroupTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.CreateNodeGroup),
		TaskName:       cloudprovider.CreateNodeGroupTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
	}
	// generate taskName
	taskName := fmt.Sprintf(createNodeGroupTaskTemplate, group.ClusterID, group.Name)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	createNodeGroupTask := &CreateNodeGroupTaskOption{Group: group}
	// step1. call qcloud create node group
	createNodeGroupTask.BuildCreateCloudNodeGroupStep(task)
	// step2. wait qcloud create node group complete
	createNodeGroupTask.BuildCheckCloudNodeGroupStatusStep(task)
	// step3. ensure autoscaler(安装/更新CA组件) in cluster
	common.BuildEnsureAutoScalerTaskStep(task, group.ClusterID, group.Provider)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildCreateNodeGroupTask task StepSequence empty")
	}

	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.CreateNodeGroupJob.String()

	return task, nil
}

// BuildCleanNodesInGroupTask clean specified nodes in NodeGroup
// including remove nodes from NodeGroup, clean data in nodes
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildCleanNodesInGroupTask(nodes []*proto.Node, group *proto.NodeGroup, // nolint
	opt *cloudprovider.CleanNodesOption) (*proto.Task, error) {

	// clean nodeGroup nodes in cloud only has two steps:
	// 1. call asg RemoveInstances to clean cluster nodes
	// because cvms return to cloud asg resource pool, all clean works are handle by asg
	// we do little task here

	// validate request params
	if nodes == nil {
		return nil, fmt.Errorf("BuildCleanNodesInGroupTask nodes info empty")
	}
	if group == nil {
		return nil, fmt.Errorf("BuildCleanNodesInGroupTask group info empty")
	}
	if opt == nil || len(opt.Operator) == 0 || opt.Cluster == nil {
		return nil, fmt.Errorf("BuildCleanNodesInGroupTask TaskOptions is lost")
	}

	var (
		nodeIPs, nodeIDs, deviceIDs = make([]string, 0), make([]string, 0), make([]string, 0)
	)
	for _, node := range nodes {
		nodeIPs = append(nodeIPs, node.InnerIP)
		nodeIDs = append(nodeIDs, node.NodeID)
		deviceIDs = append(deviceIDs, node.DeviceID)
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.CleanNodeGroupNodes),
		TaskName:       cloudprovider.CleanNodesInGroupTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
		NodeIPList:     nodeIPs,
	}
	// generate taskName
	taskName := fmt.Sprintf(cleanNodeGroupNodesTaskTemplate, group.ClusterID, group.Name)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// instance passwd
	passwd := group.LaunchTemplate.InitLoginPassword
	task.CommonParams[cloudprovider.PasswordKey.String()] = passwd

	// setting all steps details
	cleanNodeGroupNodes := &CleanNodeInGroupTaskOption{
		Group:     group,
		NodeIPs:   nodeIPs,
		NodeIds:   nodeIDs,
		DeviceIDs: deviceIDs,
		Operator:  opt.Operator,
	}

	// step0: cordon nodes
	cleanNodeGroupNodes.BuildCordonNodesStep(task)

	// step1. business user define flow
	if group.NodeTemplate != nil && len(group.NodeTemplate.ScaleInPreScript) > 0 {
		common.BuildJobExecuteScriptStep(task, common.JobExecParas{
			ClusterID:        opt.Cluster.ClusterID,
			Content:          group.NodeTemplate.ScaleInPreScript,
			NodeIps:          strings.Join(nodeIPs, ","),
			Operator:         opt.Operator,
			StepName:         common.PreInitStepJob,
			AllowSkipJobTask: group.NodeTemplate.AllowSkipScaleInWhenFailed,
			Translate:        common.PreInitJob,
		})
	}

	if group.NodeTemplate != nil && group.NodeTemplate.ScaleInExtraAddons != nil &&
		len(group.NodeTemplate.ScaleInExtraAddons.PreActions) > 0 {
		err := template.BuildSopsFactory{
			StepName: template.UserPreInit,
			Cluster:  opt.Cluster,
			Extra: template.ExtraInfo{
				InstancePasswd:  passwd,
				NodeIPList:      strings.Join(nodeIPs, ","),
				NodeOperator:    opt.Operator,
				ShowSopsUrl:     true,
				TranslateMethod: template.UserBeforeInit,
			}}.BuildSopsStep(task, group.NodeTemplate.ScaleInExtraAddons, true)
		if err != nil {
			return nil, fmt.Errorf("BuildCleanNodesInGroupTask ScaleInExtraAddons.PreActions "+
				"BuildBkSopsStepAction failed: %v", err)
		}
	}

	// step2: delete cluster group nodes
	cleanNodeGroupNodes.BuildCleanNodeGroupNodesStep(task)
	// step3: check deleting node status
	cleanNodeGroupNodes.BuildCheckClusterCleanNodesStep(task)
	// step4: remove nodes from cmdb
	common.BuildRemoveHostStep(task, opt.Cluster.BusinessID, nodeIPs)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildCleanNodesInGroupTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]

	// set global task paras
	task.CommonParams[cloudprovider.NodeIDsKey.String()] = strings.Join(nodeIDs, ",")
	task.CommonParams[cloudprovider.NodeIPsKey.String()] = strings.Join(nodeIPs, ",")

	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.CleanNodeGroupNodesJob.String()
	return task, nil
}

// BuildDeleteNodeGroupTask when delete nodegroup, we need to create background
// task to clean all nodes in nodegroup, release all resource in cloudprovider,
// finnally delete nodes information in local storage.
// @param group: need to delete
func (t *Task) BuildDeleteNodeGroupTask(group *proto.NodeGroup, nodes []*proto.Node,
	opt *cloudprovider.DeleteNodeGroupOption) (*proto.Task, error) {
	// validate request params
	if group == nil {
		return nil, fmt.Errorf("BuildDeleteNodeGroupTask group info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildDeleteNodeGroupTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.DeleteNodeGroup),
		TaskName:       cloudprovider.DeleteNodeGroupTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
	}
	// generate taskName
	taskName := fmt.Sprintf(deleteNodeGroupTaskTemplate, group.ClusterID, group.Name)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	deleteNodeGroupTask := &DeleteNodeGroupTaskOption{
		Group:                  group,
		CleanInstanceInCluster: opt.CleanInstanceInCluster,
	}
	// step1. call qcloud delete node group
	deleteNodeGroupTask.BuildDeleteNodeGroupStep(task)

	// step2. ensure autoscaler to remove this nodegroup
	if group.EnableAutoscale {
		common.BuildEnsureAutoScalerTaskStep(task, group.ClusterID, group.Provider)
	}

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildDeleteNodeGroupTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]

	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.DeleteNodeGroupJob.String()
	return task, nil
}

// BuildMoveNodesToGroupTask build move nodes to group task
func (t *Task) BuildMoveNodesToGroupTask(nodes []*proto.Node, group *proto.NodeGroup,
	opt *cloudprovider.MoveNodesOption) (*proto.Task, error) {
	return nil, cloudprovider.ErrCloudNotImplemented
}

// BuildUpdateDesiredNodesTask build update desired nodes task
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildUpdateDesiredNodesTask(desired uint32, group *proto.NodeGroup, // nolint
	opt *cloudprovider.UpdateDesiredNodeOption) (*proto.Task, error) {
	// validate request params
	if desired == 0 {
		return nil, fmt.Errorf("BuildUpdateDesiredNodesTask desired is zero")
	}
	if group == nil {
		return nil, fmt.Errorf("BuildUpdateDesiredNodesTask group info empty")
	}
	if opt == nil || len(opt.Operator) == 0 || opt.Cluster == nil {
		return nil, fmt.Errorf("BuildUpdateDesiredNodesTask TaskOptions is lost")
	}

	// normal or external nodePool
	// isExternal := cloudprovider.IsExternalNodePool(group)

	// generate main task
	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.UpdateNodeGroupDesiredNode),
		TaskName:       cloudprovider.UpdateDesiredNodesTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
	}
	// generate taskName
	taskName := fmt.Sprintf(updateNodeGroupDesiredNodeTemplate, group.ClusterID, group.Name)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	passwd := group.LaunchTemplate.InitLoginPassword
	task.CommonParams[cloudprovider.PasswordKey.String()] = passwd

	// setting all steps details
	updateDesiredNodesTask := &UpdateDesiredNodesTaskOption{
		Group:        group,
		Desired:      desired,
		Operator:     opt.Operator,
		NodeSchedule: opt.NodeSchedule,
	}

	// step1. call qcloud interface to set desired nodes
	updateDesiredNodesTask.BuildApplyInstanceMachinesStep(task)
	// step2. check cluster nodes and all nodes status is running
	updateDesiredNodesTask.BuildCheckClusterNodeStatusStep(task)
	// step3: install gse agent
	common.BuildInstallGseAgentTaskStep(task, &common.GseInstallInfo{
		ClusterId:   opt.Cluster.ClusterID,
		NodeGroupId: group.NodeGroupID,
		BusinessId:  opt.Cluster.BusinessID,
		User:        group.GetLaunchTemplate().GetInitLoginUsername(),
		Passwd:      passwd,
		KeyInfo:     group.GetLaunchTemplate().GetKeyPair(),
		Port: func() string {
			exist := checkIfWhiteImageOsNames(&cloudprovider.ClusterGroupOption{
				CommonOption: opt.CommonOption,
				Cluster:      opt.Cluster,
				Group:        opt.NodeGroup,
			})
			if exist {
				return fmt.Sprintf("%v", utils.ConnectPort)
			}

			return ""
		}(),
	})
	// step4: transfer host module
	moduleID := cloudprovider.GetTransModuleInfo(opt.Cluster, opt.AsOption, opt.NodeGroup)
	if moduleID != "" {
		common.BuildTransferHostModuleStep(task, opt.Cluster.BusinessID, moduleID, "")
	}

	// step5: business define sops task 支持脚本和标准运维流程
	if group.NodeTemplate != nil && len(group.NodeTemplate.UserScript) > 0 {
		common.BuildJobExecuteScriptStep(task, common.JobExecParas{
			ClusterID:        group.ClusterID,
			Content:          group.NodeTemplate.UserScript,
			NodeIps:          "",
			Operator:         opt.Operator,
			StepName:         common.PostInitStepJob,
			AllowSkipJobTask: group.NodeTemplate.GetAllowSkipScaleOutWhenFailed(),
			Translate:        common.PostInitJob,
		})
	}
	if group.NodeTemplate != nil && group.NodeTemplate.ScaleOutExtraAddons != nil {
		err := template.BuildSopsFactory{
			StepName: template.UserAfterInit,
			Cluster:  opt.Cluster,
			Extra: template.ExtraInfo{
				InstancePasswd:     passwd,
				NodeIPList:         "",
				NodeOperator:       opt.Operator,
				ShowSopsUrl:        true,
				ExternalNodeScript: "",
				NodeGroupID:        group.NodeGroupID,
				TranslateMethod:    template.UserPostInit,
			}}.BuildSopsStep(task, group.NodeTemplate.ScaleOutExtraAddons, false)
		if err != nil {
			return nil, fmt.Errorf("BuildScalingNodesTask business BuildBkSopsStepAction failed: %v", err)
		}
	}

	// step6: set node annotations
	common.BuildNodeAnnotationsTaskStep(task, opt.Cluster.ClusterID, nil,
		cloudprovider.GetAnnotationsByNg(opt.NodeGroup))
	// step7. set node scheduler by nodeIPs
	updateDesiredNodesTask.BuildUnCordonNodesStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildUpdateDesiredNodesTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]

	// must set job-type
	task.CommonParams[cloudprovider.ScalingNodesNumKey.String()] = strconv.Itoa(int(desired))
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.UpdateNodeGroupDesiredNodeJob.String()
	task.CommonParams[cloudprovider.ManualKey.String()] = strconv.FormatBool(opt.Manual)

	return task, nil
}

// BuildSwitchNodeGroupAutoScalingTask ensure auto scaler status and update nodegroup status to normal
func (t *Task) BuildSwitchNodeGroupAutoScalingTask(group *proto.NodeGroup, enable bool,
	opt *cloudprovider.SwitchNodeGroupAutoScalingOption) (*proto.Task, error) {
	// validate request params
	if group == nil {
		return nil, fmt.Errorf("BuildSwitchNodeGroupAutoScalingTask nodegroup info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildSwitchNodeGroupAutoScalingTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.SwitchNodeGroupAutoScaling),
		TaskName:       cloudprovider.SwitchNodeGroupAutoScalingTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
	}
	// generate taskName
	taskName := fmt.Sprintf(switchNodeGroupAutoScalingTaskTemplate, group.ClusterID, group.Name)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// step1. ensure auto scaler
	common.BuildEnsureAutoScalerTaskStep(task, group.ClusterID, group.Provider)
	// step2. update node group info in DB
	// switchNodeGroupTask.BuildUpdateNodeGroupAutoScalingDBStep(task)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildSwitchNodeGroupAutoScalingTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.SwitchNodeGroupAutoScalingJob.String()
	return task, nil
}

// BuildUpdateAutoScalingOptionTask build update auto scaler option task
func (t *Task) BuildUpdateAutoScalingOptionTask(scalingOption *proto.ClusterAutoScalingOption,
	opt *cloudprovider.UpdateScalingOption) (*proto.Task, error) {
	// validate request params
	if scalingOption == nil {
		return nil, fmt.Errorf("BuildUpdateAutoScalingOptionTask scaling option info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildUpdateAutoScalingOptionTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.UpdateAutoScalingOption),
		TaskName:       cloudprovider.UpdateAutoScalingOptionTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      scalingOption.ClusterID,
		ProjectID:      scalingOption.ProjectID,
		Creator:        scalingOption.Creator,
		Updater:        scalingOption.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	// generate taskName
	taskName := fmt.Sprintf(updateAutoScalingOptionTemplate, scalingOption.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	// step1. ensure auto scaler
	common.BuildEnsureAutoScalerTaskStep(task, scalingOption.ClusterID, scalingOption.Provider)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildUpdateAutoScalingOptionTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.UpdateAutoScalingOptionJob.String()
	return task, nil
}

// BuildSwitchAsOptionStatusTask build switch auto scaler option status task - 开启/关闭集群自动扩缩容
func (t *Task) BuildSwitchAsOptionStatusTask(scalingOption *proto.ClusterAutoScalingOption, enable bool,
	opt *cloudprovider.CommonOption) (*proto.Task, error) {
	// validate request params
	if scalingOption == nil {
		return nil, fmt.Errorf("BuildSwitchAutoScalingOptionStatusTask scalingOption info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildSwitchAutoScalingOptionStatusTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.SwitchAutoScalingOptionStatus),
		TaskName:       cloudprovider.SwitchAutoScalingOptionStatusTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      scalingOption.ClusterID,
		ProjectID:      scalingOption.ProjectID,
		Creator:        scalingOption.Creator,
		Updater:        scalingOption.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
	}
	// generate taskName
	taskName := fmt.Sprintf(switchAutoScalingOptionStatusTemplate, scalingOption.ClusterID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	// step1. ensure auto scaler
	common.BuildEnsureAutoScalerTaskStep(task, scalingOption.ClusterID, scalingOption.Provider)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildSwitchAutoScalingOptionStatusTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.SwitchAutoScalingOptionStatusJob.String()
	return task, nil
}

// BuildAddExternalNodeToCluster add external to cluster
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func (t *Task) BuildAddExternalNodeToCluster(group *proto.NodeGroup, nodes []*proto.Node, // nolint
	opt *cloudprovider.AddExternalNodesOption) (*proto.Task, error) {
	// AddExternalNodeToCluster has three steps:
	// 1. call qcloud getExternalNodeScript get addNodes script
	// 2. call bksops add nodes to cluster
	// may be need to call external previous or behind operation by bkops

	return nil, cloudprovider.ErrCloudNotImplemented
}

// BuildDeleteExternalNodeFromCluster remove external node from cluster
func (t *Task) BuildDeleteExternalNodeFromCluster(group *proto.NodeGroup, nodes []*proto.Node,
	opt *cloudprovider.DeleteExternalNodesOption) (*proto.Task, error) {
	// DeleteExternalNodeFromCluster has two steps:
	// 1. call qcloud DeleteExternalNodes
	// 2. call bksops clean node
	// may be need to call external previous or behind operation by bkops

	return nil, cloudprovider.ErrCloudNotImplemented
}

// BuildUpdateNodeGroupTask when update nodegroup, we need to create background task,
func (t *Task) BuildUpdateNodeGroupTask(group *proto.NodeGroup, opt *cloudprovider.CommonOption) (*proto.Task, error) {
	// validate request params
	if group == nil {
		return nil, fmt.Errorf("BuildUpdateNodeGroupTask group info empty")
	}
	if opt == nil {
		return nil, fmt.Errorf("BuildUpdateNodeGroupTask TaskOptions is lost")
	}

	nowStr := time.Now().Format(time.RFC3339)
	task := &proto.Task{
		TaskID:         uuid.New().String(),
		TaskType:       cloudprovider.GetTaskType(cloudName, cloudprovider.UpdateNodeGroup),
		TaskName:       cloudprovider.UpdateNodeGroupTask.String(),
		Status:         cloudprovider.TaskStatusInit,
		Message:        "task initializing",
		Start:          nowStr,
		Steps:          make(map[string]*proto.Step),
		StepSequence:   make([]string, 0),
		ClusterID:      group.ClusterID,
		ProjectID:      group.ProjectID,
		Creator:        group.Creator,
		Updater:        group.Updater,
		LastUpdate:     nowStr,
		CommonParams:   make(map[string]string),
		ForceTerminate: false,
		NodeGroupID:    group.NodeGroupID,
	}
	// generate taskName
	taskName := fmt.Sprintf(updateNodeGroupTaskTemplate, group.ClusterID, group.NodeGroupID)
	task.CommonParams[cloudprovider.TaskNameKey.String()] = taskName

	// setting all steps details
	// step1. ensure auto scaler
	common.BuildEnsureAutoScalerTaskStep(task, group.ClusterID, group.Provider)

	// set current step
	if len(task.StepSequence) == 0 {
		return nil, fmt.Errorf("BuildUpdateNodeGroupTask task StepSequence empty")
	}
	task.CurrentStep = task.StepSequence[0]
	task.CommonParams[cloudprovider.JobTypeKey.String()] = cloudprovider.UpdateNodeGroupJob.String()
	return task, nil
}

// BuildSwitchClusterNetworkTask switch cluster network mode
func (t *Task) BuildSwitchClusterNetworkTask(cls *proto.Cluster,
	subnet *proto.SubnetSource, opt *cloudprovider.SwitchClusterNetworkOption) (*proto.Task, error) {
	return nil, cloudprovider.ErrCloudNotImplemented
}
