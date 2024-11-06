package messages

import "go.temporal.io/sdk/workflow"

// "AcceptClaimCode" signal channel
func GetSignalChannelForAcceptClaimCode(ctx workflow.Context) workflow.ReceiveChannel {
	return workflow.GetSignalChannel(ctx, "AcceptClaimCodeSignal")
}
