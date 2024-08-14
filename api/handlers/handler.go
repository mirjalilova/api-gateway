package handlers
import(
	"github.com/mirjalilova/api-gateway/client"
)
type Handlers struct {
	Clients client.Clients
}

func NewHandlerStruct() *Handlers {
	return &Handlers{
		Clients: *client.NewClients(),
	}
}