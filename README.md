# Kubernetes OnDemand Sidecar Injector
This repo contains the Kubernetes OnDemand Sidecar Injector APIs and the helm packages for a straightforward deployment on a Kubernetes cluster.
The cluster can be of any flavor and hosted on-premises or in a cloud provider of your choice.

The APIs permit the retrieval and the update of Kubernetes resource manifests of type Deployment by means of web calls using the standard swagger UI exposed by API service itself (at /swagger/index.html url relative to service endpoint) or by means of other custom external scripts or applications.

The objective is essentially to make easy for application owners to add sidecar containers to existing Pods. This can be useful for multiple reasons like, for instance, troubleshooting or whenever you need to perform other tasks locally on a running container that usually runs under low privileged user.

The sidecar image can be any public docker image or a customized one you need to specifically create for your requirements and you can specify the command to invoke at container startup or the default image's ENTRYPOINT is invoked.

## How to use this code
The application exposing the APIs can be deployed to target Kubernetes cluster by using the following helm packages:
- **kube-ondemand-sidecar-injector**
    - This is the main helm chart for deploying the application serving the APIs and the ClusterRole along with the RoleBinding required to give application kubernetes service account permission to get and update Deployment manifests on the namespace where you deploy the chart
- **kube-ondemand-sidecar-injector-rolebinding**
    - This is a supporting helm chart to grant access to service account of the main chart application to a different namespace. This need to be deployed on all namespaces you want the sidecar injector be able to act on

As usual with helm chart you may further configure the chart deployment itself by enabling, for instance, ingress or adjusting replica count among other settings. It's suggested to look at [values.yaml](src/k8s/kube-ondemand-sidecar-injector/values.yaml) to see other configurable settings.

### Configure the security
All the APIs require the API Key authentication as http request header with name **X-API-KEY** and value is retrieved by the environment variable named SECRET_API_KEY configured with parameter **secretApiKey** on deployment of main helm chart.

As mentioned before the service account receive permission with the ClusterRole-RoleBinding mapping and this is automatically configured on the namespace you deploy the main helm chart on. For the other namespaces, you need to use the supporting chart kube-ondemand-sidecar-injector-rolebinding where you have to configure the name of the service account, if customized, and the containing namespace respectively by means of the parameters **kubeOndemandSidecarInjectorServiceName** and **kubeOndemandSidecarInjectorReleaseNamespace**

### Running in a Kubernetes Cluster

To deploy (or update) the last version of main chart to a Kubernetes use the following helm command (to be installed before if not already available)

```
helm upgrade --install kube-ondemand-sidecar-injector oci://ghcr.io/alesspanms/kube-ondemand-sidecar-injector/charts/kube-ondemand-sidecar-injector -n [your-hosting-namespace] --create-namespace --set secretApiKey=[your-secret-api-key] --set sidecarNamePrefix=[your-preferred-prefix]
```

If you want a specific version you can use the additional parameter with the desired version

```
--version 1.0.0 
```

#### Note
The parameter **sidecarNamePrefix** is for adding a known prefix to container sidecar name not to having chances to collide with existing container name already running on the same Pod.   

#### Additional namespaces
As mentioned, for granting access to other namespaces you have to deploy the Role Binding with the supporting chart with the following command:

```
helm upgrade --install kube-ondemand-sidecar-injector-rolebinding oci://ghcr.io/alesspanms/kube-ondemand-sidecar-injector/charts/kube-ondemand-sidecar-injector-rolebinding -n [your-additional-namespace] --set kubeOndemandSidecarInjectorReleaseNamespace=[your-hosting-namespace] --create-namespace
```

Altough we'll probably not change this supporting chart, if you want to deploy a specific version you can use the additional parameter with the desired version

```
--version 1.0.0 
```

### Running standalone with local user .kube/config file

If you want to use the Kubernetes OnDemand Sidecar Injector remotely using docker and leveraging launching user .kube/config (and, in case, any additional client certificate folder) you can use the following command using the correct parameters and desired version

```
docker run --name [your-preferred-container-name] -d -v [your_user_home_folder]/.kube/:/home/ondemandsidecarinjector/.kube/ -p [desired-local-port]:8080 -e SECRET_API_KEY=[your-api-key] -e SIDECAR_NAME_PREFIX=[your-preferred-container-name-prefix] ghcr.io/alesspanms/kube-ondemand-sidecar-injector:1.0.0
```

In case of startup failures, you mai omit the -d parameter in order to see the complete output error.

#### Note
Current version of the Kubernetes OnDemand Sidecar Injector doesn't support kube context switching so in case you have many context configured it could be better to run each time in an interactive way so you can restart it whenever you change context to use that current-context when initializing.

```
docker run -it -v [your_user_home_folder]/.kube/:/home/ondemandsidecarinjector/.kube/ -p [desired-local-port]:8080 -e SECRET_API_KEY=[your-api-key] -e SIDECAR_NAME_PREFIX=[your-preferred-container-name-prefix] ghcr.io/alesspanms/kube-ondemand-sidecar-injector:1.0.0
```


## How to contribute
Thank you for considering contributing to our project! We welcome contributions from everyone, whether you're fixing a typo, adding a feature, or proposing a major overhaul.

For guidelines and information on how to contribute, please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License
This project adopts the MIT license.  
For more information please refers to [LICENSE](LICENSE) file.