package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type JobManager interface {
	CreateJob(params *JobParams) (string, error)
}

type NomadJobManager struct {
	cli *api.Client
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

func (s *NomadJobManager) ObserveDeployment(jobID string, jobModifyInd uint64) error {
	stream := s.cli.EventStream()
	q := &api.QueryOptions{}
	topics := map[api.Topic][]string{
		api.TopicDeployment: {jobID},
	}

	eChan, err := stream.Stream(context.TODO(), topics, jobModifyInd, q)
	if err != nil {
		return err
	}

	for events := range eChan {
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

	return fmt.Errorf("events channel was closed")
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
	return fmt.Sprintf("%v:%v", network.IP, network.DynamicPorts[0].Value), nil
}

func (s *NomadJobManager) CreateJob(params *JobParams) (string, error) {
	job := NewNginxJob(params)

	jr, _, err := s.cli.Jobs().Register(job, nil)
	if err != nil {
		return "", err
	}

	err = s.ObserveDeployment(*job.ID, jr.JobModifyIndex)
	if err != nil {
		return "", err
	}

	return s.FindServerAddress(*job.ID)
}
