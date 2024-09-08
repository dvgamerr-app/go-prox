package envs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var Version string
var BuildTime string

var IsDev bool = false

func Load() error {
	const dotenv string = ".env"
	if _, err := os.Stat(dotenv); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	IsDev = true
	return godotenv.Load(dotenv)
}
