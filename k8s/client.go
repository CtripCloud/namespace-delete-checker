package k8s

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// KubeConfigFile is a string flag which indicates kubeconfig filepath
	KubeConfigFile string

	// Clientset init on start, used by others to create k8s resources
	Clientset *kubernetes.Clientset
	// KubeConfig init on start, used by others to create k8s rest client
	KubeConfig *rest.Config

	HttpClient *http.Client
)

// MustInit init sharedinformers, clients, etc.
func MustInit() {
	var err error
	KubeConfig, err = clientcmd.BuildConfigFromFlags("", KubeConfigFile)
	if err != nil {
		logrus.Panic(err)
	}
	Clientset = kubernetes.NewForConfigOrDie(KubeConfig)

	transport, err := rest.TransportFor(KubeConfig)
	if err != nil {
		logrus.Panic(err)
	}
	if transport != http.DefaultTransport {
		HttpClient = &http.Client{Transport: transport}
		if KubeConfig.Timeout > 0 {
			HttpClient.Timeout = KubeConfig.Timeout
		}
	}
}
