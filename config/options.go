package config

import "flag"

var (
	DatabaseName string
	DatabasePath string
)

func init() {
	flag.StringVar(&DatabaseName, "db-name", "database", "Name of the database file")
	flag.StringVar(&DatabasePath, "db-path", "./", "Path to the database file")
}
