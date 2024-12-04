package crypto

import (
	"eldercare_health/app/internal/db"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"strings"
)

// 生成密钥
func generateKeys(auth *abe.MAABEAuth, gid string, attribs []string) ([]*abe.MAABEKey, error) {
	fmt.Println("attribs: ", attribs)
	keys, err := auth.GenerateAttribKeys(gid, attribs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate attribute keys: %v", err)
	}
	return keys, nil
}

// 根据用户属性信息生成密钥
func generateUserKey(userId string) ([]*abe.MAABEKey, error) {
	//根据user id 查询用户属性信息
	dbClient, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	authAttributes, err := db.GetUserAttributesByUserID(dbClient, userId)
	if err != nil {
		return nil, err
	}
	var userKeys []*abe.MAABEKey
	for _, authAttribute := range authAttributes {
		//根据用户属性信息生成密钥
		keys, err := generateKeys(mapAuth[authAttribute.AuthorityName], userId, strings.Split(authAttribute.Attributes, ","))
		if err != nil {
			return nil, err
		}
		fmt.Println(keys)
		userKeys = append(userKeys, keys...)
	}
	return userKeys, nil
}
