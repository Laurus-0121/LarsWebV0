package dao

import "LarsWebV0/model"

func GetApiKey() (string, error) {
	var key model.Apikey
	err := db.First(&key).Error
	if err != nil {
		return "", err
	}
	return key.Key, nil

}
