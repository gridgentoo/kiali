package appender

import (
	"strings"

	"github.com/kiali/kiali/graph/tree"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/services/business/checkers"
)

type SidecarsCheckAppender struct{}

func (a SidecarsCheckAppender) AppendGraph(trees *[]tree.ServiceNode, namespaceName string) {
	k8s, err := kubernetes.NewClient()
	checkError(err)

	for _, tree := range *trees {
		a.applySidecarsChecks(&tree, namespaceName, k8s)
	}
}

func (a SidecarsCheckAppender) applySidecarsChecks(n *tree.ServiceNode, namespaceName string, k8s *kubernetes.IstioClient) {
	serviceName := strings.Split(n.Name, ".")[0]
	serviceVersion := n.Version

	pods, err := k8s.GetServicePods(namespaceName, serviceName, serviceVersion)
	checkError(err)

	checker := checkers.PodChecker{Pods: pods.Items}
	typeValidations := checker.Check()
	nameValidations := (*typeValidations)["pod"]

	sidecarsOk := true
	if nameValidations != nil {
		for _, check := range *nameValidations {
			sidecarsOk = sidecarsOk && check.Valid
		}
	}
	n.Metadata["hasMissingSidecars"] = !sidecarsOk

	for _, child := range n.Children {
		a.applySidecarsChecks(child, namespaceName, k8s)
	}
}
