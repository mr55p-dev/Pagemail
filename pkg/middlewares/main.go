package middlewares

import "github.com/mr55p-dev/pagemail/pkg/logging"

var log logging.Log

func init() {
	log = logging.GetLogger("middleware")
}
