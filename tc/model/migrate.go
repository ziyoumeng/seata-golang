package model

import "github.com/transaction-wg/seata-golang/tc/config"

func Migrate(){
	_ = config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(GlobalTransactionDO))
	_ = config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(BranchTransactionDO))
	_ = config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(LockDO))
}