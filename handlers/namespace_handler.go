package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ctripcloud/namespace-delete-check/cfg"
	"github.com/ctripcloud/namespace-delete-check/k8s"
	"github.com/ctripcloud/namespace-delete-check/util"
)

type Resource struct {
}
type ResourceList struct {
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resource `json:"items"`
}

func toResponseAdmissionResponse(allowed bool, msg string) *v1beta1.AdmissionResponse {
	if allowed {
		return &v1beta1.AdmissionResponse{
			Allowed: allowed,
		}
	}
	return &v1beta1.AdmissionResponse{
		Allowed: allowed,
		Result: &metav1.Status{
			Message: msg,
		},
	}
}

func NamespaceDeleteCheck(c *gin.Context) {
	reponseStatus := http.StatusOK
	requestedAdmissionReview := v1beta1.AdmissionReview{}
	responseAdmissionReview := v1beta1.AdmissionReview{}

	defer func(c *gin.Context) {
		c.JSON(reponseStatus, &responseAdmissionReview)
	}(c)
	err := c.BindJSON(&requestedAdmissionReview)
	if err != nil {
		logrus.WithError(err).Error("request namespace is empty")
		responseAdmissionReview.Response = toResponseAdmissionResponse(false, err.Error())
		return
	}

	namespace := requestedAdmissionReview.Request.Name
	if namespace == "" {
		logrus.Error("request namespace is empty")
		responseAdmissionReview.Response = toResponseAdmissionResponse(false, "request namespace is empty")
		return
	}

	namespacedResources, err := k8s.Clientset.DiscoveryClient.ServerPreferredNamespacedResources()
	if err != nil {
		logrus.WithError(err).Error("get namespaced resource failed")
		responseAdmissionReview.Response = toResponseAdmissionResponse(false, fmt.Sprintf("get namespaced resource failed %s", err.Error()))
		return
	}

	blackListResources := cfg.Config().NsResourceCheckBL

	for _, namespacedResource := range namespacedResources {
		for _, resource := range namespacedResource.APIResources {
			if !util.Contains("get", resource.Verbs) {
				continue
			}
			if util.Contains(cfg.ResourceNameGroup{Name: resource.Name, GroupVersion: namespacedResource.GroupVersion}, blackListResources) {
				continue
			}
			success, err := resourceExistInNamespace(namespace, resource, *namespacedResource)
			if err != nil {
				logrus.WithError(err).Errorf("get resource: %s in this namespace: %s error", resource.Name, namespace)
				responseAdmissionReview.Response = toResponseAdmissionResponse(false, fmt.Sprintf("get resource %s in this namespace %s error %s", resource.Name, namespace, err.Error()))
				return
			}
			logrus.Infof(" namespace:%s resource %s checkout result %t", namespace, resource.Name, success)
			if success {
				responseAdmissionReview.Response = toResponseAdmissionResponse(false, fmt.Sprintf("get resource: %s in this namespace: %s", resource.Name, namespace))
				return
			}
		}
	}

	logrus.Infof("namespace:%s delete checkout success", namespace)
	responseAdmissionReview.Response = toResponseAdmissionResponse(true, "")
	return
}

func buildResourceUrl(resource metav1.APIResource, namespace string, namespacedResource metav1.APIResourceList) string {
	queryUrl := k8s.KubeConfig.Host
	gv := strings.Split(namespacedResource.GroupVersion, "/")
	if len(gv) == 1 {
		queryUrl += "/api"
	} else {
		queryUrl += "/apis"
	}
	queryUrl += "/" + strings.ToLower(namespacedResource.GroupVersion)
	if resource.Namespaced {
		queryUrl += "/namespaces/" + namespace
	}
	queryUrl += "/" + strings.ToLower(resource.Name)
	queryUrl += "?limit=1"
	return queryUrl
}
func resourceExistInNamespace(namespace string, resourceName metav1.APIResource, namespacedResource metav1.APIResourceList) (bool, error) {
	queryUrl := buildResourceUrl(resourceName, namespace, namespacedResource)
	rsp, err := k8s.HttpClient.Get(queryUrl)
	if err != nil {
		return false, err
	}
	if rsp.StatusCode != http.StatusOK {
		logrus.Info(err)
		return false, errors.New(fmt.Sprintf("http status is not ok,query url %s", queryUrl))
	}
	defer rsp.Body.Close()
	resourceList := &ResourceList{}
	err = json.NewDecoder(rsp.Body).Decode(resourceList)
	if err != nil {
		return false, nil
	}
	if len(resourceList.Items) == 0 {
		return false, nil
	}
	return true, nil
}
