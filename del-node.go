// +build linux

package main

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// del-node        host:port node_id

var delNodeCommand = cli.Command{
	Name:      "del-node",
	Usage:     "del a redis node from existed cluster.",
	ArgsUsage: `host:port node_id`,
	Action: func(context *cli.Context) error {
		var addr string
		var nodeid string

		if len(context.Args()) < 2 {
			logrus.Fatalf("Must provide \"host:port node_id\" for del-node command!")
		}

		if addr = context.Args().Get(0); addr == "" {
			logrus.Fatalf("Please check host:port for del-node command!")
		}
		if nodeid = context.Args().Get(1); nodeid == "" {
			logrus.Fatalf("Please check node_id for del-node command!")
		}

		rt := NewRedisTrib()
		if err := rt.DelNodeClusterCmd(addr, nodeid); err != nil {
			//logrus.Errorf("%p", err)
			return err
		}
		return nil
	},
}

func (self *RedisTrib) DelNodeClusterCmd(addr, nodeid string) error {
	if addr == "" {
		return errors.New("Please check host:port for del-node command.")
	}
	logrus.Printf(">>> Removing node %s from cluster %s", nodeid, addr)

	// Load cluster information
	if err := self.LoadClusterInfoFromNode(addr); err != nil {
		return err
	}

	// Check if the node exists and is not empty
	node := self.GetNodeByName(nodeid)

	if node == nil {
		logrus.Fatalf("[ERR] No such node ID %s", nodeid)
	}

	if len(node.Slots()) > 0 {
		logrus.Fatalf("Node %s is not empty! Reshard data away and try again.", node.String())
	}
	// Send CLUSTER FORGET to all the nodes but the node to remove
	logrus.Printf(">>> Sending CLUSTER FORGET messages to the cluster...")
	for _, n := range self.nodes {
		if n == nil {
			continue
		}
		// send cluster forget cmd
	}
	// Finally shutdown the node
	logrus.Printf(">>> SHUTDOWN the node.")
	return nil
}