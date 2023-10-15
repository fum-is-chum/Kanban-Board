package config

import (
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	SERVERPORT int
	DBPORT     int
	DBHOST     string
	DBUSER     string
	DBPASS     string
	DBNAME     string
}

func InitConfig() *AppConfig {
	var res = new(AppConfig)
	res = loadConfig(".env")
	if res == nil {
		logrus.Fatal("Config: Cannot start program, failed to load configuration")
		return nil
	}

	return res
}

func InitConfigTest() *AppConfig {
	var res = new(AppConfig)

	const projectDirName = "RESTful-API-testing"
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	res = loadConfig(string(rootPath) + "/.env.test")
	if res == nil {
		logrus.Fatal("Config: Cannot start program, failed to load configuration")
		return nil
	}

	return res
}

func loadConfig(envPath string) *AppConfig {
	var res = new(AppConfig)

	err := godotenv.Load(envPath)

	if err != nil {
		logrus.Error("Config: Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("SERVERPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config: Invalid port value,", err.Error())
			return nil
		}

		res.SERVERPORT = port
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config: Invalid DB port value,", err.Error())
			return nil
		}

		res.DBPORT = port
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHOST = val
	}

	if val, found := os.LookupEnv("DBUSER"); found {
		res.DBUSER = val
	}

	if val, found := os.LookupEnv("DBPASS"); found {
		res.DBPASS = val
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBNAME = val
	}

	return res
}
