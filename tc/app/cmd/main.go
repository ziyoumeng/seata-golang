package main

import (
	"github.com/transaction-wg/seata-golang/base/common"
	"github.com/transaction-wg/seata-golang/tc/model"
	_ "net/http/pprof"
	"os"
	"strconv"
)

import (
	gxnet "github.com/dubbogo/gost/net"
	"github.com/urfave/cli/v2"
)

import (
	"github.com/transaction-wg/seata-golang/pkg/util/log"
	"github.com/transaction-wg/seata-golang/pkg/util/uuid"
	"github.com/transaction-wg/seata-golang/tc/config"
	"github.com/transaction-wg/seata-golang/tc/holder"
	"github.com/transaction-wg/seata-golang/tc/lock"
	_ "github.com/transaction-wg/seata-golang/tc/metrics"
	"github.com/transaction-wg/seata-golang/tc/server"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "start seata golang tc server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config, c",
						Usage: "Load configuration from `FILE`",
					},
					&cli.StringFlag{
						Name:  "serverNode, n",
						Value: "1",
						Usage: "server node id, such as 1, 2, 3. default is 1",
					},
				},
				Action: func(c *cli.Context) error {
					configPath := c.String("config")
					serverNode := c.Int("serverNode")
					ip, _ := gxnet.GetLocalIP()

					config.InitConf(configPath)
					model.Migrate()
					panic("migrate ok")
					conf := config.GetServerConfig()
					port, _ := strconv.Atoi(conf.Port)
					common.XID.Init(ip, port)

					uuid.Init(serverNode)
					lock.Init()
					holder.Init()
					srv := server.NewServer()
					srv.Start(conf.Host + ":" + conf.Port)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}

