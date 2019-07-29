### Introduction
In K8s, when a namespace is deleted, pods and other resources that belong to this namespace will be deleted soon. 

However, the namespace deleting operation may be a mistake. 

This project is used to prevent important resources from being deleted.

### How it works
1. hook the namespace deleting operation by k8s webhook

1. get resources in this namespace

1. if something important exists,the delete operation will be rejected

### How to use
1. enable webhook by adding `ValidatingAdmissionWebhook` to kube-apiserver flag `--admission-control`
1. get all resources, `kubectl api-resources -namespaced=true`
1. modify the config.json to ignore the unimportant resources.
1. deploy this project. using k8s Deployment and Service is recommended. ServiceAccount is also needed.
1. to create ValidatingWebhookConfiguration, see ValidatingWebhookConfiguration.yaml.example
1. caBundle can be gotten

    <code> kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n' </code>
1. have a try