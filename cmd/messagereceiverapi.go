package cmd

import (
	"assesment/api/messagereceiverapi"
	"github.com/spf13/cobra"
)

var messageReceiverApiCmd = &cobra.Command{
	Use:  "messagereceiver-api",
	RunE: messagereceiverapi.Init,
}

func init() {
	RootCmd.AddCommand(messageReceiverApiCmd)
}
