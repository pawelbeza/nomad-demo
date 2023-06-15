package nomad

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func NewEventsChan(status string) chan *api.Events {
	c := make(chan *api.Events, 1)
	c <- &api.Events{
		Events: []api.Event{
			{
				Payload: map[string]interface{}{
					"Deployment": &api.Deployment{
						Status: status,
					},
				},
			},
		},
	}

	return c
}

func TestObserveDeploymentEvents_FailedEvent(t *testing.T) {
	jobManager := NewNomadJobManager(nil)
	c := NewEventsChan("failed")

	err := jobManager.ObserveDeploymentEvents(context.Background(), c)
	assert.Error(t, err)
}

func TestObserveDeploymentEvents_SuccessfulEvent(t *testing.T) {
	jobManager := NewNomadJobManager(nil)
	c := NewEventsChan("successful")

	err := jobManager.ObserveDeploymentEvents(context.Background(), c)
	assert.NoError(t, err)
}

func TestObserveDeploymentEvents_Timeout(t *testing.T) {
	jobManager := NewNomadJobManager(nil)
	c := make(chan *api.Events)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	err := jobManager.ObserveDeploymentEvents(ctx, c)
	assert.Error(t, err)
}
