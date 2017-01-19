Just deploy with:

```
kubectl create namespace iot
kubectl create -f iot-device-type.yaml
kubectl create -f iot-daemon-set-type.yaml
kubectl create -f iot-pod-type.yaml
kubectl create -f sample-iot-devices.yaml
kubectl create -f sample-iot-daemon-sets.yaml
kubectl create -f sample-iot-pods.yaml
```

API at: http://localhost:8080/apis/fujitsu.com/v1
