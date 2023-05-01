### Overview

This receiver is created to fulfill the purpose of [adding opentelemetry support to Kubearmor](https://github.com/kubearmor/KubeArmor/issues/894). This receiver would convert the existing logs in opentelemetry to the [plog.logs format](https://github.com/open-telemetry/opentelemetry-collector/tree/main/pdata), this is the pipeline format in which logs are transported in memory in the collector.

[Stanza](https://github.com/observIQ/stanza) is a fast and lightweight log transport and processing agent. It has been donated to opentelemetry. Log based components in the collector contrib ( a repository for OpenTelemetry Collector components) use stanza as an intermediary to transform logs to plog.logs. Examples include:

- [filelogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver)
- [syslogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/syslogreceiver)
- [tcplogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/tcplogreceiver)
- [udplogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/udplogreceiver)

I would be leveraging the same approach in creating the receiver.

To support Kubearmor, I would be creating the kubearmorlog (we could come up with a better name later) stanza input operator. An operator in Stanza is a task that helps us read from a file, parse the log, filter it and then push it to another log stream pipeline (similarly to the forwarding plugin of FluentD or Fluent Bit) and directly to the observability backend of your choice.

Similarly to the other agents on the market, there are several types of operators:

- Input

- Parser

- Transform

- Output (Link to source: [isitobservable.com](https://isitobservable.io/open-telemetry/what-is-stanza-and-what-does-it-do))

The stanza adapter in the [opentelemetry contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza/adapter) takes care of intergrating stanza which does the conversion of logs to plog.logs format. To create the input operator that would work with this, work has to be done to implement the stanza.LogReceiverType.


#### Proposed kubearmor receiver config

All opentelemetry components have configuration that modifies how they function. 

Below is an initial config.yaml for the kubearmor receiver:

```yaml
# Proposwed ID for kubearmor receiver
kubearmorreceiver:
  #Specifies the kubearmor relay server endpoint, this would be optional by default it would be the value of the KUBEARMOR_SERVICE or in a k8 environment, the value
  # of the kubearmor relay service endpoint
  endpoint: https://127.0.0.1:32767 
  # By default, all of kubearmor telemetry data would be forwarded but users can modify which logs to forward by making of the options: policy,system or all
  logfilter: all
    
```
### Design
![image](https://user-images.githubusercontent.com/59079323/234896206-c223391f-f4aa-44d9-97d5-9b9b063464ab.png)

 

### TODO

- Decide if this component should be part of the [opentelemetry collector contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib). The collector 
