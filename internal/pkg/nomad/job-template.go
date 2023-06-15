package nomad

import (
	"fmt"
	"nomad-demo/internal/pkg/util"
	"time"

	"github.com/hashicorp/nomad/api"
)

func CreateFetcherCmd(url string, script bool) string {
	outPath := "${NOMAD_ALLOC_DIR}/index.html"
	fetchUrlCmd := fmt.Sprintf("wget -T 5 -O - %v", url)

	if !script {
		return fmt.Sprintf("%v > %v", fetchUrlCmd, outPath)
	}

	executeCmd := fmt.Sprintf("sh <(%v) 2>&1 > %v", fetchUrlCmd, outPath)
	return fmt.Sprintf("%v && %v", fetchUrlCmd, executeCmd)
}

type JobParams struct {
	ServiceName string
	Url         string
	Script      bool
}

func NewNginxJob(params *JobParams) *api.Job {
	return &api.Job{
		Name:        util.PointerOf(params.ServiceName),
		ID:          util.PointerOf(params.ServiceName),
		Datacenters: []string{"dc1"},
		Type:        util.PointerOf("service"),
		TaskGroups: []*api.TaskGroup{
			{
				Name:  util.PointerOf(params.ServiceName),
				Count: util.PointerOf(1),
				Update: &api.UpdateStrategy{
					ProgressDeadline: util.PointerOf(60 * time.Second),
					HealthyDeadline:  util.PointerOf(30 * time.Second),
				},
				Tasks: []*api.Task{
					{
						Name:   "page-fetcher",
						Driver: "docker",
						Config: map[string]interface{}{
							"image":   "busybox",
							"command": "/bin/sh",
							"args":    []string{"-c", CreateFetcherCmd(params.Url, params.Script)},
						},
						Resources: &api.Resources{
							CPU:      util.PointerOf(100),
							MemoryMB: util.PointerOf(128),
						},
						Lifecycle: &api.TaskLifecycle{
							Hook:    "prestart",
							Sidecar: false,
						},
					},
					{
						Name:   "nginx",
						Driver: "docker",
						Config: map[string]interface{}{
							"image": "nginx",
							"port_map": []map[string]int{
								{"http": 8080},
							},
							"volumes": []string{
								"custom/default.conf:/etc/nginx/conf.d/default.conf",
							},
						},
						Templates: []*api.Template{
							{
								EmbeddedTmpl: util.PointerOf(`
									server {
									listen 8080;
									server_name nginx.service.consul;
									location / {
										root {{env "NOMAD_ALLOC_DIR"}};
									}
									}
								`),
								DestPath: util.PointerOf("custom/default.conf"),
							},
						},
						Resources: &api.Resources{
							CPU:      util.PointerOf(100),
							MemoryMB: util.PointerOf(128),
							Networks: []*api.NetworkResource{
								{
									DynamicPorts: []api.Port{
										{Label: "http", To: 8080},
									},
								},
							},
						},
						Services: []*api.Service{
							{
								Name:      "nginx",
								Tags:      []string{"nginx", "web"},
								PortLabel: "http",
								Checks: []api.ServiceCheck{
									{
										Type:     "tcp",
										Interval: 10 * time.Second,
										Timeout:  2 * time.Second,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
