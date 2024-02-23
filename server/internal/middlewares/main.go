package middlewares


var log logging.Log

func init() {
	log = logging.GetLogger("middleware")
}
