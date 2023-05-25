### INTRODUCTION

This document describes how to test out the otel_Kubearmor_receiver. There are two ways to deploy kubearmor, on bare metal and in a kubernetes environment.
I would be explaining the two ways in which you can try out this example in both environments.

### KUBEAMOR ON BARE METAL

#### Requirements:

We would be creating an opentelemetry collector to test out the receiver. The OpenTelemetry Collector offers a vendor-agnostic implementation of how to receive, process and export telemetry data. Read more about it in the [docs](https://opentelemetry.io/docs/collector/). There are different versions:

1. [Collector-core collector](https://github.com/open-telemetry/opentelemetry-collector)
   The components that are a part of this collector are fixed that i.e. components are not contributed to this collector. It is maintained by the opentelemetry community
2. [Collector contrib collector](https://github.com/open-telemetry/opentelemetry-collector-contrib)
    This consists of a growing number of components contributed by the community, observability vendors and any one in general with a need to create custom components for a specific use,
3. Custom collector
   This is created by users for specific use case. Only needed components are included, unneeded ones are not included. Custom collectors can easily be created using the [opentelemetry collector builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder). This is what we would be using for our tutorial.

#### Steps:

> ** Easy Fix **
> - Pull custom collector container
> ```
>  docker run -d --net=host chinwendu20/receiver
>  ```
>  No need to follow the steps below if you follow this route

- ##### Create a custom opentelemetry collector distribution.

1. Go to [opentelemetry collector's release page](https://github.com/open-telemetry/opentelemetry-collector/releases), download the "ocb" binary compatible with your system's architecture

2. Use the collector builder to create the custom collector

Note: Take a look at the [collector-builder.yml](../collector-builder.yml). Note that we have included the kubearmor receiver under the receivers map.

Run the command below:

```bash
/path/to/ocb/binary --config=collector-builder.yml

```

Note: 
- Please replace /path/to/ocb/binary with the actual path to the ocb binary you downloaded
- The collector-builder.ymlfile is located in this repo at /example/collector-builder.yml. Use the actual path as the value to --config flag

If everything went correctly, you should have an otel-custom folder containing an otel-custom binary. That is our collector distribution. We may proceed to testing the collector.

Note: This step is important as the binary compiled by ocb would not work in the container
3. Create collector binary

```bash
cd otel-custom
GO111MODULE=on CGO_ENABLED=0 go build .

```
- ##### Testing the kubearmor receiver in the collector distribution

1. Run the collector

Run the command below:

```bash
/path/to/otel-custom --config=config.yml

```

Note: 
- Please replace /path/to/otel-custom with the actual path to the otel-custom binary you downloaded
- The config.yml file is located in this repo at /example/config.yml. Use the actual path as the value to --config flag

Examine the logs to see that it is properly running.

### KUBEAMOR ON KUBERNETES ENVIRONMENT

For this tutorial we would be making use of the minikube kubernetes environment

#### Steps:
- ##### Follow [previous step](https://github.com/Chinwendu20/OTel-receiver/blob/otel/example/tutorial.md#create-a-custom-opentelemetry-collector-distribution) on creating custom opentelemetry collector.
- ##### Follow the steps in this [markdown](https://github.com/kubearmor/KubeArmor/tree/main/deployments/k3s) to deploy kubearmor in k3s environemnt
- ##### Install opentelemetry operator. Follow these steps:

1. Ensure [cert manager is installed](https://cert-manager.io/docs/installation/) in the cluster.
2. Deploy the operator:

```bash
    kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
```
(Reference: [Opentelemetry operator readme](https://github.com/open-telemetry/opentelemetry-operator))

- Create custom collector container

**Note: I have created a container already and have included it in the kubernetes manifest file in this tutorial. You can skip this step if you want and use that instead. Go to step [4](#Deploy-the-collector-in-your-cluster)**

1. Build custom collector docker image. We would be using the [Dockerfile](../../Dockerfile) to build the image. Ensure you are in the `example` directory. Run this command:
     ```
     docker build . -t=<docker username>/<image name>
     ```
     Note: Replace `docker username ` with your docker username and `image name` with your preferred image name.
    
2. Push to docker hub:

```
docker push <docker username>/<image name>

```
3. Replace the container name in this [line](../collector-k8-manifest.yml) with your container name.

4. #### Deploy the collector in your cluster

```bash
kubectl apply -f collector-k8-manifest.yml
```
5. View the logs of the container to note that it runs fine.

### OpenTelemetry KubeArmor Logs pattern

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
