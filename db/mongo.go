package db

import (
	"crypto/tls"
	"database-ms/config"
	"fmt"
	"net"
	"strings"

	"gopkg.in/mgo.v2"
)

var instance *mgo.Session

var err error

// GetInstance return copy of db session
func GetInstance(c *config.Configuration) *mgo.Session {

	dialInfo := mgo.DialInfo{
		Addrs:    strings.Split(c.MongoCluster, ","),
		Username: c.MongoUsername,
		Password: c.MongoPassword,
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig) // add TLS config
		return conn, err
	}

	instance, err = mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return instance
}
