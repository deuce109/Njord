package readers

import (
	"deuce109/7dtd-map-server/v2/models"
	"encoding/xml"
	"os"
	"path/filepath"
	"strconv"
)

type persistentPlayerData struct {
	XMLName xml.Name     `xml:"persistentplayerdata"`
	Players []playerData `xml:"player"`
}

type playerData struct {
	// XMLName    xml.Name `xml:"player"`
	Username   string `xml:"playername,attr"`
	Id         string `xml:"userid,attr"`
	Platform   string `xml:"nativeplatform,attr"`
	PlatformId string `xml:"nativeuserid,attr"`
}

type adminTools struct {
	XMLName     xml.Name    `xml:"adminTools"`
	AdminsBlock adminsBlock `xml:"users"`
}

type adminsBlock struct {
	// XMLName xml.Name    `xml:"adminTools>admins"`
	Admins []adminInfo `xml:"user"`
}

type adminInfo struct {
	// XMLName  xml.Name `xml:"admins>admin"`
	Platform string `xml:"platform,attr"`
	Id       string `xml:"userid,attr"`
	Name     string `xml:"name,attr"`
	Level    string `xml:"permission_level,attr"`
}

var PlayerData []models.UserInfo

func parseXml[T interface{}](xmlPath string, t *T) error {

	_xml, err := os.ReadFile(xmlPath)

	if err != nil {
		return err
	} else {
		err = xml.Unmarshal(_xml, &t)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPlayerData(playerXmlPath string) ([]playerData, error) {
	var data persistentPlayerData

	absPath, err := filepath.Abs(filepath.Join(playerXmlPath, "players.xml"))

	if err != nil {
		return nil, err
	}

	err = parseXml(absPath, &data)
	if err != nil {
		return nil, err
	}

	return data.Players, nil
}

func getAdminInfo(adminXmlPath string) ([]adminInfo, error) {
	var tools adminTools

	absPath, err := filepath.Abs(filepath.Join(adminXmlPath, "serveradmin.xml"))

	if err != nil {
		return nil, err
	}

	err = parseXml(absPath, &tools)
	if err != nil {
		return nil, err
	}

	return tools.AdminsBlock.Admins, nil
}

func GetUserInfo(playerDataXmlPath, adminDataXmlPath string) error {

	playerData, err := getPlayerData(playerDataXmlPath)

	if err != nil {
		return err
	}

	admins, err := getAdminInfo(adminDataXmlPath)

	if err != nil {
		return err
	}

	for _, player := range playerData {
		for _, admin := range admins {
			if player.PlatformId == admin.Id {

				permLevel, err := strconv.Atoi(admin.Level)

				if err != nil {
					return err
				}

				user := &models.UserInfo{
					Username:   player.Username,
					UserId:     player.Id,
					Platform:   player.Platform,
					PlatformId: player.PlatformId,
					Level:      permLevel,
				}

				PlayerData = append(PlayerData, *user)
				break
			}

		}
	}

	return err
}
