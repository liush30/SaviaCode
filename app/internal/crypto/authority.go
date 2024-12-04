package crypto

import (
	"eldercare_health/app/internal/db"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"log"
	"strings"
)

var mapAuth = make(map[string]*abe.MAABEAuth)

func init() {
	//查询授权以及属性信息
	dbClient, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	attributes, err := db.GetAuthoritiesNameAndAttributes(dbClient)
	log.Println(attributes)
	if err != nil {
		log.Fatal(err)
	}

	initAuthority(attributes)
}

func initAuthority(authorities []db.AuthorityAttributes) {
	maabe := abe.NewMAABE()
	for _, auth := range authorities {
		log.Printf("authority: %s, attributes: %s\n", auth.AuthorityName, auth.Attributes)
		authAttributes := strings.Split(auth.Attributes, ",")
		authority, err := createAuthority(maabe, auth.AuthorityName, authAttributes)
		if err != nil {
			log.Fatal(err)
		}
		mapAuth[auth.AuthorityName] = authority
	}
}

// 创建授权机构
func createAuthority(maabe *abe.MAABE, name string, attribs []string) (*abe.MAABEAuth, error) {
	auth, err := maabe.NewMAABEAuth(name, attribs)
	if err != nil {
		return nil, fmt.Errorf("failed generation authority %s: %v", name, err)
	}
	return auth, nil
}
