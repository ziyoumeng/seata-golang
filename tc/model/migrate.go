package model

import (
	"fmt"
	"github.com/transaction-wg/seata-golang/tc/config"
)

func Migrate(){
	err := config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(GlobalTransactionDO))
	check(err)
	err = config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(BranchTransactionDO))
	check(err)
	err = config.GetServerConfig().StoreConfig.DBStoreConfig.Engine.Sync(new(LockDO))
	check(err)
	fmt.Println("migrate ok")
}

func check(err error){
	if err != nil {
		panic(err)
	}
}