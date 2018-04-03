package gtdata

type EntityKey struct {
	KeyUID            string
	KeyAppData        string
	KeyGroup          string
	KeyFriend         string
	KeyFriendRequest  string
	KeyBlack          string
	KeyMessageOffline string
}

type App struct {
	Appid   uint64 `redis:"appid" json:"appid"`
	Name    string `redis:"name" json:"name"`
	Owner   uint64 `redis:"owner" json:"owner"`
	Desc    string `redis:"desc" json:"desc"`
	Regdate int64  `redis:"regdate" json:"regdate"`
	Type    string `redis:"type" json:"type"`
	Share   uint64 `redis:"share" json:"share"`
}
