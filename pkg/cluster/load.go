package cluster

import (
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"

	"github.com/aws/eks-anywhere/pkg/types"
)

type kubeConfigCluster struct {
	Name string `json:"name"`
}

type kubeConfigYAML struct {
	Clusters []*kubeConfigCluster `json:"clusters"`
}

func LoadManagement(kubeconfig string) (*types.Cluster, error) {
	if kubeconfig == "" {
		return nil, nil
	}
	kubeConfigBytes, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return nil, err
	}
	kc := &kubeConfigYAML{}
	kc.Clusters = []*kubeConfigCluster{}
	err = yaml.Unmarshal(kubeConfigBytes, &kc)
	if err != nil {
		return nil, fmt.Errorf("error parsing kubeconfig file: %v", err)
	}
	return &types.Cluster{
		Name:               kc.Clusters[0].Name,
		KubeconfigFile:     kubeconfig,
		ExistingManagement: true,
	}, nil
}
