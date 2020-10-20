package rm

import (
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/transaction-wg/seata-golang/base/meta"
	"github.com/transaction-wg/seata-golang/base/model"
	"github.com/transaction-wg/seata-golang/base/protocal"
	"github.com/transaction-wg/seata-golang/pkg/client"
	"github.com/transaction-wg/seata-golang/pkg/config"
	"github.com/transaction-wg/seata-golang/pkg/context"
)

var (
	DBKEYS_SPLIT_CHAR = ","
)
//seata的rm
type AbstractResourceManager struct {
	RpcClient     *client.RpcRemoteClient
	ResourceCache map[string]model.IResource
}

func NewAbstractResourceManager(client *client.RpcRemoteClient) AbstractResourceManager {
	resourceManager := AbstractResourceManager{
		RpcClient:     client,
		ResourceCache: make(map[string]model.IResource),
	}
	go resourceManager.handleRegisterRM()
	return resourceManager
}

func (resourceManager AbstractResourceManager) RegisterResource(resource model.IResource) {
	resourceManager.ResourceCache[resource.GetResourceId()] = resource
}

func (resourceManager AbstractResourceManager) UnregisterResource(resource model.IResource) {

}

func (resourceManager AbstractResourceManager) GetManagedResources() map[string]model.IResource {
	return resourceManager.ResourceCache
}

//向tc服注册分支事务
func (resourceManager AbstractResourceManager) BranchRegister(branchType meta.BranchType, resourceId string,
	clientId string, xid string, applicationData []byte, lockKeys string) (int64, error) {
	request := protocal.BranchRegisterRequest{
		Xid:             xid,
		BranchType:      branchType,
		ResourceId:      resourceId,
		LockKey:         lockKeys,
		ApplicationData: applicationData,
	}
	resp, err := resourceManager.RpcClient.SendMsgWithResponse(request)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	response := resp.(protocal.BranchRegisterResponse)
	if response.ResultCode == protocal.ResultCodeSuccess {
		return response.BranchId, nil
	} else {
		return 0, response.GetError()
	}
}
//向tc报告分支状态
func (resourceManager AbstractResourceManager) BranchReport(branchType meta.BranchType, xid string, branchId int64,
	status meta.BranchStatus, applicationData []byte) error {
	request := protocal.BranchReportRequest{
		Xid:             xid,
		BranchId:        branchId,
		Status:          status,
		ApplicationData: applicationData,
	}
	resp, err := resourceManager.RpcClient.SendMsgWithResponse(request)
	if err != nil {
		return errors.WithStack(err)
	}
	response := resp.(protocal.BranchReportResponse)
	if response.ResultCode == protocal.ResultCodeFailed {
		return response.GetError()
	}
	return nil
}
//todo ?
func (resourceManager AbstractResourceManager) LockQuery(ctx *context.RootContext, branchType meta.BranchType, resourceId string, xid string,
	lockKeys string) (bool, error) {
	return false, nil
}

func (resourceManager AbstractResourceManager) handleRegisterRM() {
	for {
		serverAddress := <-resourceManager.RpcClient.GettySessionOnOpenChannel
		resourceManager.doRegisterResource(serverAddress)
	}
}
//向tc注册所有资源
func (resourceManager AbstractResourceManager) doRegisterResource(serverAddress string) {
	if resourceManager.ResourceCache == nil || len(resourceManager.ResourceCache) == 0 {
		return
	}
	message := protocal.RegisterRMRequest{
		AbstractIdentifyRequest: protocal.AbstractIdentifyRequest{
			Version:                 config.GetClientConfig().SeataVersion,
			ApplicationId:           config.GetClientConfig().ApplicationId,
			TransactionServiceGroup: config.GetClientConfig().TransactionServiceGroup,
		},
		ResourceIds: resourceManager.getMergedResourceKeys(),
	}

	resourceManager.RpcClient.RegisterResource(serverAddress, message)
}
//那所有资源id，逗号分隔 （ActionName,ActionName...）
func (resourceManager AbstractResourceManager) getMergedResourceKeys() string {
	var builder strings.Builder
	if resourceManager.ResourceCache != nil && len(resourceManager.ResourceCache) > 0 {
		for key, _ := range resourceManager.ResourceCache {
			builder.WriteString(key)
			builder.WriteString(DBKEYS_SPLIT_CHAR)
		}
		resourceKeys := builder.String()
		resourceKeys = resourceKeys[:len(resourceKeys)-1] //todo 为什么不直接返回
		return resourceKeys
	}
	return ""
}
