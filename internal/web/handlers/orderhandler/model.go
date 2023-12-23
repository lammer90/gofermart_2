package orderhandler

import (
	"net/http"
)

type OrderRestApiProvider interface {
	Save(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
}
