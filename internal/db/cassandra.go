package db

import (
	"chat-system/internal/utils"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitCassandra() {
	host := os.Getenv("CASSANDRA_HOST")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	portStr := os.Getenv("CASSANDRA_PORT")
	consistencyStr := os.Getenv("CASSANDRA_CONSISTENCY")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		utils.GetLogger().Error("Invalid port  : ", err)
		fmt.Println("Invalid port: %v", err)
	}

	consistency := gocql.ParseConsistency(consistencyStr)

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Consistency = consistency
	cluster.Port = port

	Session, err = cluster.CreateSession()
	if err != nil {
		utils.GetLogger().Error("Failed to connect to Cassandra   : ", err)
		fmt.Println("Failed to connect to Cassandra: %v", err)
	}

	executeCQLFile("scripts/setup_cassandra.cql")
}

func executeCQLFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.GetLogger().Error("Failed to read CQL file : ", err)
		fmt.Println("Failed to read CQL file: %v", err)
	}

	commands := strings.Split(string(content), ";")
	for _, command := range commands {
		if strings.TrimSpace(command) != "" {
			err := Session.Query(command).Exec()
			if err != nil {
				utils.GetLogger().Error("Failed to execute command : ", err)
				fmt.Println("Failed to execute command: %v", err)
			}
		}
	}
}
