package common

const (
	ERROR_OFFLINE             = 1   //user offline
	ERROR_PUSH_CLIENT_FAILURE = 2   //push to client failure
	ERROR_ROUTER_CHANGE       = 3   //router change
	ERROR_ACK_TIMEOUT         = 4   //ack timeout
	ERROR_DISPATCH_ERROR      = 100 //handle message error
	ERROR_UNSUPPORTED_CMD     = 101 //unsupported command
	ERROR_UNKNOWN             = -1  //unknown
)
