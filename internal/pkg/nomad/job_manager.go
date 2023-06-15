package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type JobManager interface {
	CreateJob(ctx context.Context, params *JobParams) (string, error)
}

type NomadJobManager struct {
	cli *api.Client // ideally I would use interface here but upstream doesn't have any for client
}

func NewDefaultJobManager() (*NomadJobManager, error) {
	conf := api.DefaultConfig()
	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return NewNomadJobManager(cli), err
}

func NewNomadJobManager(cli *api.Client) *NomadJobManager {
	return &NomadJobManager{cli: cli}
}

func (s *NomadJobManager) ObserveDeploymentEvents(ctx context.Context, eChan <-chan *api.Events) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("observe deployment timeout")
		case events := <-eChan:
			for _, event := range events.Events {
				deployment, err := event.Deployment()
				if err != nil {
					return err
				}

				if deployment.Status == "successful" {
					return nil
				} else if deployment.Status == "failed" {
					return fmt.Errorf("deployment failed: %v", deployment.StatusDescription)
				}
			}
		}
	}
}

func (s *NomadJobManager) ObserveDeployment(ctx context.Context, jobID string, jobModifyInd uint64) error {
	stream := s.cli.EventStream()
	topics := map[api.Topic][]string{
		api.TopicDeployment: {jobID},
	}

	eventsChan, err := stream.Stream(ctx, topics, jobModifyInd, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, TASK_PROGRESS_DEADLINE)
	defer cancel()

	return s.ObserveDeploymentEvents(ctx, eventsChan)
}

func (s *NomadJobManager) GetAllocationNetworkDetails(allocation *api.Allocation) (*api.NetworkResource, error) {
	if allocation.AllocatedResources == nil {
		return nil, fmt.Errorf("missing allocated resource details in allocation")
	}

	if allocation.AllocatedResources.Tasks == nil {
		return nil, fmt.Errorf("missing tasks details in allocated resource")
	}

	if _, ok := allocation.AllocatedResources.Tasks["nginx"]; !ok {
		return nil, fmt.Errorf("missing nginx task details in allocation")
	}

	if allocation.AllocatedResources.Tasks["nginx"].Networks == nil {
		return nil, fmt.Errorf("missing network details in allocation")
	}

	if len(allocation.AllocatedResources.Tasks["nginx"].Networks) == 0 {
		return nil, fmt.Errorf("network list is empty in allocation")
	}

	return allocation.AllocatedResources.Tasks["nginx"].Networks[0], nil
}

func (s *NomadJobManager) FindServerAddress(jobID string) (string, error) {
	allocs, _, err := s.cli.Jobs().Allocations(jobID, false, nil)
	if err != nil {
		return "", err
	}

	if len(allocs) == 0 {
		return "", fmt.Errorf("empty allocation list for jobID %v", jobID)
	}

	allocation, _, err := s.cli.Allocations().Info(allocs[0].ID, nil)
	if err != nil {
		return "", err
	}

	network, err := s.GetAllocationNetworkDetails(allocation)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v:%v", network.IP, network.DynamicPorts[0].Value), nil
}

// It's rather anti-pattern that /service/:name has to wait for deployment of service due
// to asynchronous nature of Nomad.
//
// In real case scenario I think it would be better if
// /service/:name would only schedule a job and it would be client side responsibility
// to send another query(ies) to get status with IP and port of the service (similarly as it's done with
// `nomad deployment status -monitor`)
//
// ofc this would be handled by some client side tool which would abstracts away this workflow
func (s *NomadJobManager) CreateJob(ctx context.Context, params *JobParams) (string, error) {
	job := NewNginxJob(params)

	jr, _, err := s.cli.Jobs().Register(job, nil)
	if err != nil {
		return "", err
	}

	err = s.ObserveDeployment(ctx, *job.ID, jr.JobModifyIndex)
	if err != nil {
		return "", err
	}

	return s.FindServerAddress(*job.ID)
}
