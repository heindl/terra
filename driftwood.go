package terra

import "github.com/op/go-logging"

var log = logging.MustGetLogger("geostore")

func init() {
	var format = "%{color}%{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}"
	logging.SetFormatter(logging.MustStringFormatter(format))
}
