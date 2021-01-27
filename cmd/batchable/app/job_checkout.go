package app

import (
	"batchable/cmd/batchable/config"
	"encoding/json"
	"io/ioutil"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/stan.go"
)

// JobCheckout is executed when a workload finishes a job, and registers it as completed
func JobCheckout(
	w http.ResponseWriter,
	req *http.Request,
	conf *config.Config,
	jobManifest *JobManifest,
	broker *BrokerShim) {

	event := cloudevents.NewEvent()

	body, readErr := ioutil.ReadAll(req.Body)

	if readErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read request Body"))
		return
	}

	eventMarshalErr := json.Unmarshal(body, &event)

	if eventMarshalErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can not unmarshal cloudevent, make sure you send a cloudevent in structured content mode"))
		return
	}

	eventID := event.Context.GetID()

	if !jobManifest.HasJob(eventID) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Could not publish your event to the broker. Job may have timed out."))
		println("Job ID:", eventID, "does not exists anymore. Publishing blocked.")
		return
	}

	// Fetch the nopublish event context extension. This will prevent publishing the recieved event to our broker.
	// This is normally used, if you want to define the end of a chain of workloads, where the last link of the chain
	// Should not create any new events in the broker anymore
	data, _ := event.Context.GetExtension("nopublish")

	if conf.Debug {
		println("Deleting Job ID:", eventID)
	}

	if data != true {
		if conf.Debug {
			println("Publishing recieved event to broker")
		}
		publishErr := (*broker).PublishResult(*conf, event)
		if publishErr != nil {
			println("Could not publish event to broker")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not publish your event to the broker"))
			return
		}
	}

	jobManifest.DeleteJob(eventID)

	if jobManifest.Size() < conf.MaxConcurrency {
		// Initialize a new subscription should the old one have been closed
		(*broker).Start(func(msg *stan.Msg) {
			MessageHandler(msg, conf, jobManifest, broker)
		})
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("OK"))
}
