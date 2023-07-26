### About
This tutorial shows how to export KubeArmor telemetry data to grafana using the KubeArmor receiver and [Loki exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/lokiexporter)

### Prerequisites
1. Follow the [KubeArmor deployment guide](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/deployment_guide.md#L20-L19) to set up Kubearmor.
2. Follow the [collector installation guide](./tutorial.md#install-the-collector) to install or build an OpenTelemetry custom collector for KubeArmor.

### Steps
1. Follow Grafana [OpenTelemetry setup guide](https://grafana.com/docs/opentelemetry/collector/send-logs-to-loki/) to setup Grafana and Grafana Loki. It also shows how to set up configuration for the collector.

   - You would be using the kubearmor receiver for this tutorial. Ensure `receivers` are configured in your collector configuration properly.

   - For KubeArmor Kubernetes deployment, you can make use of the [K8s manifest](../collector-k8-manifest.yml) which already has the collector configuration. You just need to replace `exporters.loki.endpoint` under `config` with the Loki endpoint in your cluster.

   - For KubeArmor bare metal deployments, you can make use of the example [config file](../config.yml) for deployment. You'll need to add Loki exporter and endpoint by changing the file
     ```yaml
     exporters:
       <...>
       loki:
         endpoint: "<loki endpoint>"

     service:
       pipelines:
         logs:
         <...>
           exporters:
             - loki
             - logging
     ```
   To find the Loki endpoint in your Docker environment, run `docker inspect <network name used in docker compose file>`
   If successful, your Loki dashboard would start getting logs and would look similar to:
   ![image](https://user-images.githubusercontent.com/59079323/235289951-6842da6f-a020-4723-81f6-02bae0987d1c.png)

2. To create Grafana dashboard use the [Grafana dashboard JSON](../grafana_dashboard.json) and follow [this tutorial](https://grafana.com/docs/grafana/latest/dashboards/manage-dashboards/#import-a-dashboard) to import the Grafana json on your own server.
   Once setup successfully, your Grafana dashboard would look similar to:
   [Video of grafana dashboard](https://1drv.ms/v/s!AqdT9dah_scBkD5QWHz--sK7acwZ?e=cmty14)

    Features:
    - View the amount of each unique value of each log attribute using pie chanrt, guage and table.
    - Dynamically choose the log attribute that you would like to view using the `log attribute` variable
    - Filter through logs using `filter` variable.
