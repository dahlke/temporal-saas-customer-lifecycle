package messages

import "go.temporal.io/sdk/workflow"

func GetSignalChannelForResendClaimCodes(ctx workflow.Context) workflow.ReceiveChannel {
	return workflow.GetSignalChannel(ctx, "ResendClaimCodesSignal")
}
