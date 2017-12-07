package ack

const ACK_DEFAULT_TIMEOUT = 3 //ack默认超时时间3000ms

type AckTaskQueue struct {
	queue map[int]*AckTask
}

/**
 * 创建任务队列
 */
func NewAckTaskQueue() *AckTaskQueue {

	return &AckTaskQueue{
		queue: make(map[int]*AckTask),
	}
}

/**
 * 添加任务
 */
func (me *AckTaskQueue) Add(task *AckTask, timeout int) {
	task.ackTaskQueue = me
	me.queue[task.AckMessageId] = task
	//任务执行器
	if timeout <= 0 {
		timeout = ACK_DEFAULT_TIMEOUT
	}

	//计时
	go task.Run(timeout)
}

/**
 * 获取并删除任务
 */
func (me *AckTaskQueue) GetAndRemove(sessionId int) *AckTask {

	task := me.queue[sessionId]

	if task != nil {
		delete(me.queue, sessionId)
	}

	return task
}
