package messages

import "go.temporal.io/sdk/workflow"

type ResendClaimCodesSignal struct {
	Email string
}

func GetSignalChannelForResendClaimCodes(ctx workflow.Context) workflow.ReceiveChannel {
	return workflow.GetSignalChannel(ctx, "ResendClaimCodesSignal")
}

func GetSignalChannelForCancelSubscription(ctx workflow.Context) workflow.ReceiveChannel {
	return workflow.GetSignalChannel(ctx, "CancelSubscriptionSignal")
}
