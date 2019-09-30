package unomaly

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gliderlabs/logspout/router"
	ingest "github.com/unomaly/ingest-go"
)

var debug = false

func init() {
	router.AdapterFactories.Register(NewUnomalyAdapter, "unomaly")
	if os.Getenv("UNOMALY_DEBUG") != "" {
		debug = true
	}
}

// UnomalyAdapter is an adapter that forwards logs to Unomaly
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

	in := ingest.Init(
		host,
		ingest.SkipTLSVerify(),
		ingest.Gzip(),
		ingest.APIPath("/v1/batch"),
	)
	a := &UnomalyAdapter{ingest: in}

	log.Printf("Adapter created: %+v", a.ingest)

	return a, nil
}

// Stream implements the router.LogAdapter interface.
func (a *UnomalyAdapter) Stream(logstream chan *router.Message) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("error getting hostname", err)
	}

	for m := range logstream {
		data := make(map[string]interface{})
		data["logspout_container"] = m.Container.Name
		data["logspout_container_id"] = m.Container.ID
		data["logspout_hostname"] = m.Container.Config.Hostname
		data["logspout_docker_image"] = m.Container.Config.Image
		data["router_hostname"] = hostname

		ev := &ingest.Event{
			Message:   m.Data,
			Source:    m.Container.Config.Image,
			Timestamp: time.Now(),
			Metadata:  data,
		}

		if debug {
			log.Printf("DEBUG Event: %+v", ev)
			log.Printf("DEBUG stream: %+v", m)
		}

		a.ingest.Send(ev)
	}
}
