package envs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	AppName   string
	Version   string
	BuildTime string
	CommitID  string
)

var IsDev bool = false

func Load() error {
	const dotenv string = ".env"
	if _, err := os.Stat(dotenv); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	IsDev = true
	return godotenv.Load(dotenv)
}
