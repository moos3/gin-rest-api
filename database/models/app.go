package models

import "github.com/moos3/gin-rest-api/lib/common"

type App struct {
	Name string
	Slug string
	Description string
	ClientID string
	ClientSecret string
	CallbackURL string
	AppURL string
}



func (a *App) Read(m common.JSON) {
	o.ID = uint(m["id"].(float64))
	o.Token = m["token"].(string)
	o.AppID = m["app_id"].(uint)
	o.UserID = m["user_id"].(uint)
}

func (a *App) Serialize() common.JSON{
	return common.JSON{
		"id":     o.ID,
		"token":  o.Token,
		"app_id": o.AppID,
	}
}

