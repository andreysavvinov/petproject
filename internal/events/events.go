package events

const (
	ExchangeName = "music.events"
	RoutingKey   = "user.deleted"
)

type UserDeleted struct {
	UserID    string `json:"user_id"`
	MyMusicID int64  `json:"my_music_id"`
}
