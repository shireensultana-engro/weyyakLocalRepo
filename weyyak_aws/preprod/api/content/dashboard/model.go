package dashboard

type WatchedContent struct {
	ContentTitle string `json:"contentTitle"`
	WatchedCount int    `json:"watchedCount"`
}

type UserDevicesresponse struct {
	Lable string `json:"lable"`
	Value int    `json:"value"`
}

type UserByRegion struct {
	RegionName string `json:"regionName"`
	UserCount  int    `json:"userCount"`
}

type Dashboard struct {
	WatchedContent      []WatchedContent
	UserDevicesresponse []UserDevicesresponse
	UserByRegion        []UserByRegion
	ActiveUsers         ActiveUsers
}
type ActiveUsers struct {
	ActiveUserCount string `json:"activeUserCount"`
}
type ApplicationSetting struct {
	Name  string `json:"name"`
	Value string `json:"Value"`
}
type Content struct {
	Id          int    `json:"id"`
	Parent      int    `json:"parent"`
	Droppable   bool   `json:"droppable"`
	Text        string `json:"text"`
	ContentType int    `json:"content_type"`
	Data        string `json:"data"`
	//temporary purpose for xms
	IsSeason bool `json:"isSeason"`
}
