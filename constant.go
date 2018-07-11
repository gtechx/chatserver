package main

const (
	ERR_NONE uint16 = iota
	//common error
	ERR_NAME_NOT_VALID
	ERR_NAME_MAX_LEN
	ERR_TIME_OUT
	ERR_MSG_INVALID
	ERR_UNKNOWN
	ERR_INVALID_JSON
	ERR_JSON_SERIALIZE
	ERR_PRESENCE_EXISTS
	ERR_PRESENCE_NOT_EXISTS

	//redis error
	ERR_REDIS
	ERR_DB

	//user error
	ERR_ACCOUNT
	ERR_ACCOUNT_EXISTS
	ERR_ACCOUNT_NOT_EXISTS
	ERR_PASSWORD_INVALID
	ERR_APPDATAID_NOT_EXISTS
	ERR_NICKNAME_EXISTS

	//privilege
	ERR_NO_PRIVILEGE

	//app
	ERR_APP_EXIST
	ERR_APP_NOT_EXISTS

	//user error
	ERR_FRIEND
	ERR_FRIEND_SELF
	ERR_MESSAGE_SELF
	ERR_FRIEND_GROUP_NOT_EXISTS
	ERR_FRIEND_GROUP_EXISTS
	ERR_FRIEND_GROUP_MAX_COUNT
	ERR_FRIEND_GROUP_USER_NOT_EMPTY
	ERR_FRIEND_ADD_NEED_REQ
	ERR_FRIEND_ADD_REFUSE_ALL
	ERR_IN_BLACKLIST
	ERR_FRIEND_MAX
	ERR_FRIEND_EXISTS
	ERR_FRIEND_NOT_EXISTS

	//group
	ERR_GROUP_EXISTS
	ERR_GROUP_NOT_EXISTS
	ERR_OLD_GROUP_NOT_EXISTS
	ERR_NEW_GROUP_EXISTS
	ERR_GROUP_NOT_EMPTY

	ERR_COUNT
)

const (
	VERIFY_TYPE_ALLOW_ALL = iota
	VERIFY_TYPE_NEED_AGREE
	VERIFY_TYPE_REFUSE_ALL
	VERIFY_TYPE_ERR
)

const NAME_MAX_LEN = 64
const GROUP_MAX_COUNT = 128

const (
	PRIVILEGE_ADD_ADMIN = iota
	PRIVILEGE_DEL_ADMIN
	PRIVILEGE_GET_ADMIN

	PRIVILEGE_ADD_USER
	PRIVILEGE_DEL_USER
	PRIVILEGE_GET_USER
	PRIVILEGE_GET_ONLINE_USER
)
