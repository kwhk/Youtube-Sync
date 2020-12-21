/* Start up script to set up web server */

package bin

import (
	"fmt"
	"net/http"
)

const (
	connPort = "8000"
	connHost = "localhost"
	connType = "tcp"
	connUrl = connHost + ":" + connPort
)