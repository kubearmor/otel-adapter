<div align="center">
<img src="https://github.com/kubearmor/KubeArmor/blob/main/.gitbook/assets/logo.png" width="400"/>
  
  
<br />
    <h1>OpenTelemetry KubeArmor Adapter<img src="https://opentelemetry.io/img/logos/opentelemetry-logo-nav.png" alt="OpenTelemetry Icon" width="45" height=""></h1>
  <p>
KubeArmor receiver component for openTelemetry collector.

</p>
<a href="https://join.slack.com/t/kubearmor/shared_invite/zt-1ltmqdbc6-rSHw~LM6MesZZasmP2hAcA"><img src="https://img.shields.io/badge/Join%20Our%20Community-Slack-blue" height="auto" width="auto" /></a>
</div>



## About

The openTelemetry KubeArmor receiver converts KubeArmor telemetry data (logs, visibilty events, policy violations) to the openTelemetry format. This [adds opentelemetry support to Kubearmor](https://github.com/kubearmor/KubeArmor/issues/894)
providing a vendor agnostic means of exporting kubearmor's telemetry data to various observability backend such as [elastic search](https://www.elastic.co/guide/en/apm/guide/current/open-telemetry-direct.html#connect-open-telemetry-collector), [grafana](https://grafana.com/docs/opentelemetry/collector/), [signoz](https://signoz.io/blog/opentelemetry-apm/) and a bunch of other [opentelemetry adopters](https://github.com/open-telemetry/community/blob/main/ADOPTERS.md)!

## Documentation :notebook:

* :point_right: [Design Doc](DESIGNDOC.md)
* :dart: [Tutorial](example/tutorials/tutorial.md)
* :heavy_check_mark: [Learn about KubeArmor](https://github.com/kubearmor/KubeArmor#readme)

<!-- ## License

[MIT](https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/blob/main/LICENSE) -->


