package repositories

import (
	"deuce109/7dtd-map-server/v2/models"
	"deuce109/7dtd-map-server/v2/readers"
)

func GetPlayerInfo(username string) *models.UserInfo {
	var userData *models.UserInfo
	for _, p := range readers.PlayerData {
		if p.Username == username {
			userData = &p
			break
		}
	}
	return userData
}
