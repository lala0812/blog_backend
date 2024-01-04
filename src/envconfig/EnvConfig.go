package envconfig

import (
	"os"
)

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		if key == "DB_WITH_SSL" {
			if value == "true" || value == "require" {
				return "require"
			} else {
				return "disable"
			}
		}
		return value
	} else {
		switch {
		case key == "USERS":
			return "Li,Tzu-Hao"
		case key == "PORT":
			return "8080"
		case key == "NAME":
			return "Default User"
		case key == "DB_HOST":
			return "localhost"
		case key == "DB_USER":
			return "postgres"
		case key == "DB_PASSWORD":
			return "1234"
		case key == "DB_NAME":
			return "blog_project"
		case key == "DB_PORT":
			return "5432"
		case key == "DB_WITH_SSL":
			return "disable"
		default:
			return ""
		}
	}

}
