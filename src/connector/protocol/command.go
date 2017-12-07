package protocol

const (
	CMD_HEARTBEAT    int8 = 1  //心跳
	CMD_HANDSHAKE    int8 = 2  //握手
	CMD_BIND         int8 = 5  //绑定用户
	CMD_UNBIND       int8 = 6  //解绑用户
	CMD_FAST_CONNECT int8 = 7  //快速重连
	CMD_ERROR        int8 = 10 //错误信息
	CMD_OK           int8 = 11 //成功信息
	CMD_PUSH         int8 = 15 //上行消息
	CMD_ACK          int8 = 23 //ack消息
	CMD_NACK         int8 = 24 //nack消息
)

//LOGIN(3),
//LOGOUT(4),
//PAUSE(8),
//RESUME(9),
//HTTP_PROXY(12),
//KICK(13),
//GATEWAY_KICK(14),
//GATEWAY_PUSH(16),
//NOTIFICATION(17),
//GATEWAY_NOTIFICATION(18),
//CHAT(19),
//GATEWAY_CHAT(20),
//GROUP(21),
//GATEWAY_GROUP(22),
//UNKNOWN(-1);
//Command(int cmd) {
//this.cmd = (byte) cmd;
//}
//
//public final byte cmd;
//
//public static Command toCMD(byte b) {
//if (b > 0 && b < values().length) return values()[b - 1];
//return UNKNOWN;
//}
//}
