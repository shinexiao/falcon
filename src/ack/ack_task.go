package ack

import "time"

type AckTask struct {
	AckMessageId int
	ackTaskQueue *AckTaskQueue
	callback     AckCallback
	timer        *time.Timer
}

/**
 * 创建任务
 */
func NewAckTask(ackMessageId int) *AckTask {

	return &AckTask{
		AckMessageId: ackMessageId,
	}
}

/**
 * 设置消息队列
 */
func (me *AckTask) SetAckTaskQueue(ackTaskQueue *AckTaskQueue) *AckTask {

	me.ackTaskQueue = ackTaskQueue
	return me
}

/**
 * 设置回调
 */
func (me *AckTask) SetCallback(callback AckCallback) *AckTask {
	me.callback = callback
	return me
}

/**
 * 尝试完成任务
 */
func (me *AckTask) tryDone() bool {

	return me.timer.Stop()
}

/**
 * 任务完成
 */
func (me *AckTask) OnResponse() {
	if me.tryDone() {
		me.callback.OnSuccess(me)
		me.callback = nil
	}
}

/**
 * 任务超时
 */
func (me *AckTask) OnTimeout() {
	task := me.ackTaskQueue.GetAndRemove(me.AckMessageId)
	if nil != task && me.tryDone() {
		me.callback.OnTimeout(me)
		me.callback = nil
	}
}

/**
 * 运行
 */
func (me *AckTask) Run(timeout int) {

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	me.timer = timer

	select {
	case <-timer.C:
		//超时
		me.OnTimeout()
	}

}
