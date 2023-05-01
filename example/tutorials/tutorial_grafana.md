### About
This tutorial shows how to export kubearmor telemetry data to grafana using the kubearmor receiver and [loki exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/lokiexporter)

### Steps
1. Follow [this tutorial](https://github.com/kubearmor/KubeArmor/blob/ce18fee4f87be7786dc1275aeb94ab7096c8b590/getting-started/deployment_guide.md#L20-L19) to set up Kubearmor
2. Follow this [tutorial](https://github.com/Chinwendu20/OTel-receiver/blob/otel/example/tutorial.md#create-a-custom-opentelemetry-collector-distribution) to create an opentelemetry custom collector.
3. This [tutorial](https://grafana.com/docs/opentelemetry/collector/send-logs-to-loki/) shows how to set up grafana and grafana loki. It also shows how to set up configuration for the collector.
    Note:
    - You would be using the kubearmor receiver for this tutorial. Ensure [this line](https://github.com/Chinwendu20/OTel-receiver/blob/e1b84530dd186b065275156fd80d5c8819295108/example/config.yml#L2-L4) and this [other line](https://github.com/Chinwendu20/OTel-receiver/blob/e1b84530dd186b065275156fd80d5c8819295108/example/config.yml#L13-L14) is in your configuration file.
    - For kubearmor kubernetes deployment, you can make use of [this file](https://github.com/Chinwendu20/OTel-receiver/blob/otel/example/collector-k8-manifest.yml) for the collector configuration. Replace in this [line](https://github.com/Chinwendu20/OTel-receiver/blob/e1b84530dd186b065275156fd80d5c8819295108/example/collector-k8-manifest.yml#L16) [your](#loki-endpoint) loki endpoint. 
    - For Kubearmor bare metal deployment, you can make use of [this file](https://github.com/Chinwendu20/OTel-receiver/blob/otel/example/config.yml) for deployment. Add the loki exporter by adding [this line](https://github.com/Chinwendu20/OTel-receiver/blob/e1b84530dd186b065275156fd80d5c8819295108/example/collector-k8-manifest.yml#L15-L16) and [this line](https://github.com/Chinwendu20/OTel-receiver/blob/e1b84530dd186b065275156fd80d5c8819295108/example/collector-k8-manifest.yml#L25) to the file. Replace loki endpoint with [yours](#loki-endpoint)


> ###### Loki endpoint
> If you followed the grafana loki tutorial referenced above, you might have to, run a `docker inspect <network name used in docker compose file>` to get the endpoint.



![image](https://user-images.githubusercontent.com/59079323/235289951-6842da6f-a020-4723-81f6-02bae0987d1c.png)
