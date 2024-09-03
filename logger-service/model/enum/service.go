package enum

type Service string

const (
	ServiceAuth      Service = "auth"
	ServiceGoogleMap Service = "google_map"
	ServiceSpider    Service = "spider"
	ServiceUser      Service = "user"
	ServiceMail      Service = "mail"
	ServiceLogger    Service = "logger"
)

func (s Service) ToString() string {
	return string(s)
}
