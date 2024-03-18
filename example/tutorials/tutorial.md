## INTRODUCTION

We would be creating an OpenTelemetry collector to test out the receiver. The OpenTelemetry Collector offers a vendor-agnostic implementation of how to receive, process and export telemetry data. Read more about it in the [docs](https://opentelemetry.io/docs/collector/). There are different versions:

1. [Collector-core collector](https://github.com/open-telemetry/opentelemetry-collector)
   The components that are a part of this collector are fixed i.e. components are not contributed to this collector. It is maintained by the OpenTelemetry community
2. [Collector contrib collector](https://github.com/open-telemetry/opentelemetry-collector-contrib)
   This consists of a growing number of components contributed by the community, observability vendors and any one in general with a need to create custom components for a specific use,
3. Custom collector
   This is created by users for specific use case. Only needed components are included, unneeded ones are not included. Custom collectors can easily be created using the [OpenTelemetry collector builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder). This is what we would be using for our tutorial.

This document describes how to test out the `kubearmor_receiver`. There are two ways to deploy kubearmor, on bare metal and in a kubernetes environment. Therefore, we will be explaining how to deploy the collector in both environments.

## PREREQUISITES
Install Kubearmor:
- [KubeArmor in K8s tutorial](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/deployment_guide.md)
- [Kubearmor bare metal deployment tutorial](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/kubearmor_vm.md)

## INSTALL THE COLLECTOR
Depending on your deployment mode, you can follow:
* [Collector in Kubernetes](#collector-in-kubernetes-environment)
* [Collector on bare metal](#collector-on-bare-metal)

### COLLECTOR IN KUBERNETES ENVIRONMENT

#### Prerequisites
1. Ensure [cert-manager is installed](https://cert-manager.io/docs/installation/) in your cluster.
2. Wait for a few minutes till cert-manager gets ready and then deploy the [OpenTelemetry operator](https://github.com/open-telemetry/opentelemetry-operator):
   ```bash
   kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
   ```

#### Run pre-built OpenTelemetry collector in K8s
If you want to skip building the collector yourself, deploy the [example manifest](../collector-k8-manifest.yml) which pulls pre-built `kubearmor/otel-adapter` image:
```bash
kubectl apply -f example/collector-k8-manifest.yml
```

<details>
  <summary><h4> Build and install OpenTelemetry collector in K8s</h4></summary>

1. Build custom collector docker image. We would be using the [Dockerfile](../../../Dockerfile) to build the image.
   ```bash
   docker build -t=<docker username>/<image name> .
   ```
   Note: Replace `docker username ` with your docker username and `image name` with your preferred image name.
2. Push to docker hub:
   ```bash
   docker push <docker username>/<image name>
   ```
3. Replace the `image` in the example [K8s manifest](../collector-k8-manifest.yml) accordingly.
4. Deploy the collector **as a daemonset** in your cluster
   ```bash
   kubectl apply -f example/collector-k8-manifest.yml
   ```
   Checkout other [deployment patterns](https://opentelemetry.io/docs/collector/deployment/) for OpenTelemetry collectors.
5. View the logs of the daemonset to check if it runs fine.
    ```bash
    kubectl logs -n kubearmor deployment/kubearmor-collector-collector -f
    ```

***Learn about the receiver's configuration [here](tutorial.md#kubearmor-receiver-config).***

</details>


#### Cleanup
```bash
# delete the collector
kubectl delete -f example/collector-k8-manifest.yml

# delete the OpenTelemetry operator
kubectl delete -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml

# delete cert-manager
kubectl delete -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
```

### COLLECTOR ON BARE METAL

#### Run pre-built OpeneTelemetry collector
If you want to skip building the example collector yourselves, you can use the pre-built one with:
```bash
docker run -d --net=host --name=kubearmor-otel-adapter kubearmor/otel-adapter
```

#### Build a custom OpenTelemetry collector distribution.
1. Go to [OpenTelemetry collector's release page](https://github.com/open-telemetry/opentelemetry-collector/releases), download the "ocb" binary compatible with your system's architecture.
   Alternatively, if you have Go installed on your system, you can get the latest ocb binary with:
   ```bash
   GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@latest
   ```

2. Use the collector builder to create the custom collector.
   Note: Take a look at the [collector-builder.yml](../collector-builder.yml). Note that we have included the kubearmor receiver under the receivers map.
   Run the command below:
   ```bash
   GO111MODULE=on CGO_ENABLED=0 /path/to/ocb/binary --config=example/collector-builder.yml
   ```
   Note:
   - `/path/to/ocb/binary` is path to the ocb binary you downloaded.
   - `collector-builder.yml` file is located in this repo at `example/collector-builder.yml`.

   If everything went correctly, you should have an `otel-custom` folder containing an otel-custom binary. That is our collector distribution. We may proceed to testing the collector.

##### Run the built collector
Run the collector with:
```bash
/path/to/otel-custom --config=example/config.yml
```
Note:
- `/path/to/otel-custom` is the path to the otel-custom binary built in previous step
- `config.yml` file is located in this repo at `example/config.yml`.
Examine the logs to see that it is properly running.


***Learn about the receiver's configuration options [here](tutorial.md#kubearmor-receiver-config).***

#### Cleanup
```bash
# stop and remove the collector container
docker stop kubearmor-otel-adapter; docker rm kubearmor-otel-adapter
```

### Kubearmor receiver config.

There are two configuration options for the receiver:

- **endpoint:**
  
   This specifies kubearmor's server API URL.
  
- **logfilter**
  
   This is used to specify which logs one is interested in. There are three filters:
  
   - kubearmorLogs:
     
     Use this if you want to see Kubearmor's internal logs only.
     
   - policy
     
     Use this if you want to see alerts only.
  
   - system
     
     Use this if you want to see logs about insights gotten by kubearmor about the host system only.
     
   - all
     
      Use this if you want to see internal logs, insights and alerts.
     
Refer to [kubearmor_receiver/testdata/config.yml](kubearmor_receiver/testdata/config.yml) for a visual example on how to
place the options in your configuration file.  

## OpenTelemetry KubeArmor Logs pattern
```log
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915426000","observedTimeUnixNano":"1679915426487671942","body":{"kvlistValue":{"values":[{"key":"HostPID","value":{"doubleValue":261}},{"key":"PPID","value":{"doubleValue":1}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/var/log/journal/b09389c7d40f420982b5facb1f6e1686"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY|O_NONBLOCK|O_DIRECTORY|O_CLOEXEC"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:26.485913Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PID","value":{"doubleValue":261}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/usr/lib/systemd/systemd-journald"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915426000","observedTimeUnixNano":"1679915426771337238","body":{"kvlistValue":{"values":[{"key":"Resource","value":{"stringValue":"/home/chinwendu/.local/share/JetBrains/consentOptions/accepted"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"HostPID","value":{"doubleValue":9527}},{"key":"PID","value":{"doubleValue":9527}},{"key":"UID","value":{"doubleValue":1000}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/snap/goland/224/jbr/bin/java"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:26.770893Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427853489107","body":{"kvlistValue":{"values":[{"key":"HostPID","value":{"doubleValue":9426}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.h7pJqH"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDWR|O_CREAT|O_EXCL"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.852810Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9426}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427853511342","body":{"kvlistValue":{"values":[{"key":"PPID","value":{"doubleValue":2396}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853020Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"PID","value":{"doubleValue":9426}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.h7pJqH"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"Result","value":{"stringValue":"Passed"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427853539241","body":{"kvlistValue":{"values":[{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853045Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9426}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.h7pJqH"}},{"key":"Data","value":{"stringValue":"syscall=SYS_UNLINK"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427853953047","body":{"kvlistValue":{"values":[{"key":"Result","value":{"stringValue":"Passed"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.5iCIgX"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9426}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853670Z"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDWR|O_CREAT|O_EXCL"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427854182196","body":{"kvlistValue":{"values":[{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853724Z"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.5iCIgX"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PID","value":{"doubleValue":9426}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Data","value":{"stringValue":"syscall=SYS_UNLINK"}},{"key":"Result","value":{"stringValue":"Passed"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427854214680","body":{"kvlistValue":{"values":[{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853760Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9426}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.decOuE"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDWR|O_CREAT|O_EXCL"}},{"key":"Result","value":{"stringValue":"Passed"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427854236825","body":{"kvlistValue":{"values":[{"key":"Data","value":{"stringValue":"syscall=SYS_UNLINK"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9426}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9426}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.decOuE"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.853778Z"}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427863323715","body":{"kvlistValue":{"values":[{"key":"PPID","value":{"doubleValue":9663}},{"key":"PID","value":{"doubleValue":9811}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.862946Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"Resource","value":{"stringValue":"sa_family=AF_NETLINK"}},{"key":"Data","value":{"stringValue":"syscall=SYS_BIND fd=41"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"Network"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427863343314","body":{"kvlistValue":{"values":[{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.862922Z"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PID","value":{"doubleValue":9811}},{"key":"Operation","value":{"stringValue":"Network"}},{"key":"Resource","value":{"stringValue":"domain=AF_NETLINK type=SOCK_RAW|SOCK_CLOEXEC protocol=0"}},{"key":"Data","value":{"stringValue":"syscall=SYS_SOCKET"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427863360182","body":{"kvlistValue":{"values":[{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY|O_CLOEXEC"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.863017Z"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PID","value":{"doubleValue":9811}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/etc/hosts"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427863371906","body":{"kvlistValue":{"values":[{"key":"Operation","value":{"stringValue":"Network"}},{"key":"Data","value":{"stringValue":"syscall=SYS_SOCKET"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PID","value":{"doubleValue":9811}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Resource","value":{"stringValue":"domain=AF_INET type=SOCK_DGRAM|SOCK_NONBLOCK|SOCK_CLOEXEC protocol=0"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.863068Z"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"UID","value":{"doubleValue":1000}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915427000","observedTimeUnixNano":"1679915427863383487","body":{"kvlistValue":{"values":[{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"Network"}},{"key":"Resource","value":{"stringValue":"sa_family=AF_INET sin_port=53 sin_addr=127.0.0.53"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"PID","value":{"doubleValue":9811}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Data","value":{"stringValue":"syscall=SYS_CONNECT fd=41"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:27.863129Z"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915428000","observedTimeUnixNano":"1679915428772864421","body":{"kvlistValue":{"values":[{"key":"HostPID","value":{"doubleValue":9527}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/snap/goland/224/jbr/bin/java"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/home/chinwendu/.local/share/JetBrains/consentOptions/accepted"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:28.772087Z"}},{"key":"PID","value":{"doubleValue":9527}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PPID","value":{"doubleValue":2396}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915429000","observedTimeUnixNano":"1679915429172577282","body":{"kvlistValue":{"values":[{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PID","value":{"doubleValue":9811}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Resource","value":{"stringValue":"/home/chinwendu/.cache/google-chrome/Default/Cache/Cache_Data/todelete_3c4c59a5fecb98f5_0_1"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:29.171674Z"}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Data","value":{"stringValue":"syscall=SYS_UNLINK"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915429000","observedTimeUnixNano":"1679915429950356162","body":{"kvlistValue":{"values":[{"key":"Source","value":{"stringValue":"/snap/goland/224/jbr/bin/java"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/home/chinwendu/.local/share/JetBrains/consentOptions/accepted"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"HostPID","value":{"doubleValue":9527}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"PID","value":{"doubleValue":9527}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:29.949563Z"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915430000","observedTimeUnixNano":"1679915430017628140","body":{"kvlistValue":{"values":[{"key":"UID","value":{"doubleValue":1000}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Data","value":{"stringValue":"syscall=SYS_CONNECT fd=43"}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}},{"key":"Operation","value":{"stringValue":"Network"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:30.017141Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"PID","value":{"doubleValue":9811}},{"key":"Resource","value":{"stringValue":"sa_family=AF_INET sin_port=443 sin_addr=142.250.187.202"}},{"key":"Result","value":{"stringValue":"Passed"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915430000","observedTimeUnixNano":"1679915430773698408","body":{"kvlistValue":{"values":[{"key":"Resource","value":{"stringValue":"/home/chinwendu/.local/share/JetBrains/consentOptions/accepted"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9527}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/snap/goland/224/jbr/bin/java"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:30.772846Z"}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9527}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"Result","value":{"stringValue":"Passed"}}]}},"traceId":"","spanId":""}]}]}]}
{"resourceLogs":[{"resource":{},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1679915431000","observedTimeUnixNano":"1679915431273738373","body":{"kvlistValue":{"values":[{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Source","value":{"stringValue":"/snap/goland/224/jbr/bin/java"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/home/chinwendu/.local/share/JetBrains/consentOptions/accepted"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDONLY"}},{"key":"HostPID","value":{"doubleValue":9527}},{"key":"PPID","value":{"doubleValue":2396}},{"key":"PID","value":{"doubleValue":9527}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:31.272970Z"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}}]}},"traceId":"","spanId":""},{"timeUnixNano":"1679915431000","observedTimeUnixNano":"1679915431359445575","body":{"kvlistValue":{"values":[{"key":"UpdatedTime","value":{"stringValue":"2023-03-27T11:10:31.358770Z"}},{"key":"PPID","value":{"doubleValue":9663}},{"key":"Type","value":{"stringValue":"HostLog"}},{"key":"Operation","value":{"stringValue":"File"}},{"key":"Resource","value":{"stringValue":"/shm/.com.google.Chrome.60n9lZ"}},{"key":"Data","value":{"stringValue":"syscall=SYS_OPENAT fd=-100 flags=O_RDWR|O_CREAT|O_EXCL"}},{"key":"Result","value":{"stringValue":"Passed"}},{"key":"HostName","value":{"stringValue":"babe-chinwendum"}},{"key":"HostPID","value":{"doubleValue":9811}},{"key":"PID","value":{"doubleValue":9811}},{"key":"UID","value":{"doubleValue":1000}},{"key":"Source","value":{"stringValue":"/opt/google/chrome/chrome"}}]}},"traceId":"","spanId":""}]}]}]}
```
