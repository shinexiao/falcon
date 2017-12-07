package ack

/**
 * 任务回调接口
 */
type AckCallback interface {
	OnSuccess(context *AckTask);
	OnTimeout(context *AckTask);
}
