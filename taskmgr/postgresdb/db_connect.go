package postgresdb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
)

type DbConnect struct {
	Spec struct {
		Db_url string `yaml:"db_url"`
	} `yaml:"spec"`
}

func (c *DbConnect) init() {
	yamlFile, err := os.ReadFile("postgresdb/config.yaml")
	if err != nil {
		fmt.Printf("yamlFile get err: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}
}

func (c *DbConnect) Connect() (*pgx.Conn, error) {
	password := os.Getenv("POSTGRES_PWD")
	c.init()
	dburl := fmt.Sprintf(c.Spec.Db_url, password)

	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
