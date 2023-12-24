package orderhandler

import (
	"net/http"
)

type OrderRestAPIProvider interface {
	Save(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
}
