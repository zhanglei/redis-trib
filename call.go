// +build linux

package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// call            host:port command arg arg .. arg
var callCommand = cli.Command{
	Name:      "call",
	Usage:     "run command in redis cluster.",
	ArgsUsage: `host:port command arg arg .. arg`,
	Action: func(context *cli.Context) error {
		rt := NewRedisTrib()
		if err := rt.CallClusterCmd(context); err != nil {
			//logrus.Errorf("%p", err)
			return err
		}
		return nil
	},
}

func (self *RedisTrib) CallClusterCmd(context *cli.Context) error {
	var addr string

	if len(context.Args()) < 2 {
		logrus.Fatalf("Must provide \"host:port command\" for call command!")
	} else if addr = context.Args().Get(0); addr == "" {
		logrus.Fatalf("Please check host:port for call command!")
	}

	if err := self.LoadClusterInfoFromNode(addr); err != nil {
		return err
	}

	cmd := context.Args().Get(1)
	cmdArgs := ToInterfaceArray(context.Args()[2:])

	logrus.Printf(">>> Calling %s %s", cmd, cmdArgs)
	_, err := self.EachPrint(cmd, cmdArgs...)
	if err != nil {
		logrus.Println("Command failed:", err)
		return err
	}

	return nil
}