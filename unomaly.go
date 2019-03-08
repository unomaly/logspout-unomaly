package unomaly

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gliderlabs/logspout/router"
	ingest "github.com/unomaly/ingest-go"
)

func init() {
	router.AdapterFactories.Register(NewUnomalyAdapter, "unomaly")
}

// UnomalyAdapter is an adapter that streams JSON to Logstash.
type UnomalyAdapter struct {
	conn   net.Conn
	route  *router.Route
	ingest *ingest.Ingest
}

// NewUnomalyAdapter creates a UnomalyAdapter
func NewUnomalyAdapter(route *router.Route) (router.LogAdapter, error) {
	host := route.Options["ingestionHost"]
	if host == "" {
		host = os.Getenv("UNOMALY_INGESTION")
	}

	ingest := ingest.Init(host)

	return &UnomalyAdapter{ingest: ingest}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *UnomalyAdapter) Stream(logstream chan *router.Message) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("error getting hostname", err)
	}

	for m := range logstream {
		var data map[string]interface{}
		data["stream"] = m.Source
		data["logspout_container"] = m.Container.Name
		data["logspout_container_id"] = m.Container.ID
		data["logspout_hostname"] = m.Container.Config.Hostname
		data["logspout_docker_image"] = m.Container.Config.Image
		data["router_hostname"] = hostname

		ev := &ingest.Event{
			Message:   m.Data,
			Source:    hostname,
			Timestamp: time.Now(),
			Metadata:  data,
		}

		a.ingest.Send(ev)
	}
}
