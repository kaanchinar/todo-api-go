# Observability
In this stage, setting up the monitoring infrastructure within the cluster using kube-prometheus-stack, collecting metrics from the application with ServiceMonitor, and creating a dashboard on Grafana that reflects key performance indicators (SLOs) are implemented.

## Environment
* Kubernetes cluster
* Helm CLI v4.1.4
- Helm Chart `kube-prometheus-stack`

## Task steps
1. Deploy the application to the cluster and verify that the `/metrics` endpoint returns the correct metrics.
2. Add the `kube-prometheus-stack` Helm chart to the repository.
3. Set up the monitoring infrastructure in the `monitoring` namespace and set `grafana.adminPassword` via Helm.
4. Log in to the Grafana UI and verify that the default Kubernetes dashboards are working.
5. Write a custom `ServiceMonitor` manifest for the application and apply it to the `monitoring` namespace.
6. Verify in the Prometheus UI that the `http_requests_total` metric is being collected and is visible.
7. Create a new empty dashboard in Grafana.
8. Set up the following panels and PromQL queries within the dashboard:
    * **Request rate** → `rate(http_requests_total[5m])`
    * **Error rate %** → `rate(http_errors_total[5m]) / rate(http_requests_total[5m]) * 100`
    * **P99 latency** → `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))`
    * **CPU usage** → `container_cpu_usage_seconds_total`
    * **Memory usage** → `container_memory_working_set_bytes`
    * **Availability %** → `(1 - error_rate) * 100` (Color rule: green >99.9%, yellow >99%, red below)
9. Export the prepared dashboard in JSON format and add it to the repository under the name `grafana/dashboards/slo.json`.
