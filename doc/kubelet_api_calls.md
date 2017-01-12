#### A summary of the code locations where `kubelet` uses the `kubeClient` to make essential API calls to the API-server:

* During the creation of the api-server pod source: [here](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/kubelet.go#L258)
    * [here](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/config/apiserver.go#L31-L34) a `ListerWatcher` is instantiated
    * [here](https://github.com/kubernetes/kubernetes/blob/master/pkg/client/cache/listwatch.go#L60-L80) is the definition of what a `ListerWatcher` does.
    * [and here](https://github.com/kubernetes/kubernetes/blob/master/pkg/client/restclient/client.go) is where the basic `RESTClient` for the k8s api-server is defined.
* TODO: add more ...
   