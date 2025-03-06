package postgresdb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Spec struct {
		Db_url string `yaml:"db_url"`
	} `yaml:"spec"`
}

func (c *Conf) getConf() *Conf {
	yamlFile, err := os.ReadFile("postgresdb/config.yaml")
	if err != nil {
		fmt.Printf("yamlFile get err: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}
	return c

}

func Connect(c *Conf) (*pgx.Conn, error) {

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PWD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("DB")
	dburl := fmt.Sprintf(c.getConf().Spec.Db_url, user, password, host, port, db)

	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
