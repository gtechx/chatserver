package config

var ServerAddr string = "127.0.0.1:9090"
var ServerNet string = "ws"

var RedisAddr string = "127.0.0.1:6379"
var RedisPassword string = ""
var RedisDefaultDB uint64 = 2

var MysqlAddr string = "127.0.0.1:3306"
var MysqlUserPassword string = "root:ztgame@123"
var MysqlDefaultDB string = "gtchat"
var MysqlTablePrefix string = "gtchat_"

var StartUID uint64 = 1000
var StartAPPID uint64 = 0

var DefaultGroupName = "MyFriends"

var DBConfig = `{
	"RedisAddr": "127.0.0.1:6379",
	"RedisPassword": "",
	"RedisDefaultDB": 2,
	"RedisMaxConn": 10,

	"MysqlAddr": "127.0.0.1:3306",
	"MysqlUserPassword": "root:ztgame@123",
	"MysqlDefaultDB": "gtchat",
	"MysqlTablePrefix": "gtchat_",
	"MysqlLogMode": true,
	"MysqlMaxConn": 10,

	"DefaultGroupName": "MyFriends"
}`
