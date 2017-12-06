package dispatcher

import (
	"connector/handler"
	"connector/protocol"
	"connector/connection"
	"connector/message"
	"common"
	"github.com/gogap/logrus"
	"fmt"
	"runtime"
)

const (
	POLICY_REJECT = 2
	POLICY_LOG    = 1
	POLICY_IGNORE = 0
)

type MessageDispatcher struct {
	handlers          map[int8]handler.MessageHandler
	unsupportedPolicy int
}

/**
 * 创建消息派发器
 */
func NewMessageDispatcher() *MessageDispatcher {

	return &MessageDispatcher{
		handlers:          make(map[int8]handler.MessageHandler),
		unsupportedPolicy: POLICY_REJECT,
	}
}

/**
 * 注册命令处理
 */
func (me *MessageDispatcher) Register(cmd int8, messageHandler handler.MessageHandler) {
	me.handlers[cmd] = messageHandler
}

/**
 * 接收消息
 */
func (me *MessageDispatcher) OnReceive(packet *protocol.Packet, conn *connection.Conn) {

	defer func() {
		if err := recover(); err != nil {
			errPacket := &protocol.Packet{
				Cmd:       protocol.CMD_ERROR,
				SessionId: packet.SessionId,
			}
			errMsg := message.NewErrorMessage(packet.Cmd, errPacket, conn)
			errMsg.Code = common.ERROR_DISPATCH_ERROR
			errMsg.Reson = "handle message error"
			errMsg.Close()

			//错误追踪
			const depth = 32
			var pcs [depth]uintptr
			n := runtime.Callers(3, pcs[:])
			for _, ptr := range pcs[0:n] {
				fn := runtime.FuncForPC(ptr)
				if fn == nil {
					return
				}
				fmt.Println(fn.Name())
			}

			me.error(fmt.Sprintf("handle message error: %v", err))
		}
	}()

	messageHandler := me.handlers[packet.Cmd]

	if messageHandler != nil {
		messageHandler.Handle(packet, conn)
	} else {
		if me.unsupportedPolicy > POLICY_IGNORE {
			if me.unsupportedPolicy == POLICY_REJECT { //拒绝处理
				errPacket := &protocol.Packet{
					Cmd:       protocol.CMD_ERROR,
					SessionId: packet.SessionId,
				}
				errMsg := message.NewErrorMessage(errPacket.Cmd, errPacket, conn)
				errMsg.Code = common.ERROR_UNSUPPORTED_CMD
				errMsg.Reson = "Param invalid"
				errMsg.Close()
				me.info("reject message.")
			} else {
				me.info("ignore message.")
			}
		}
	}

}

//输出日志
func (me *MessageDispatcher) info(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "MessageDispatcher").Info(str)
}

//输出日志
func (me *MessageDispatcher) error(str string) {
	logrus.
		WithField("system", "falcon").
		WithField("struct", "MessageDispatcher").Error(str)
}
