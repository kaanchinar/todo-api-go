# Alerting & Load testing
In this stage, setting up the Alertmanager alerting mechanism (PrometheusRule) against infrastructure and application errors, Slack/Discord integration, and verifying the system's response to these alerts via load testing (`hey` / `ab`) are carried out.

## Environment
* Prometheus Alertmanager
* Slack Webhook
* hey v0.1.5

## Task steps
1. Create a `PrometheusRule` manifest and define these alerts within it:
    * **HighErrorRate**: if the error rate is over 5% and continues for 2 minutes.
    * **HighLatency**: if the P99 latency is over 500ms and continues for 5 minutes.
    * **PodCrashLooping**: if the pod restart rate is greater than zero and continues for 5 minutes.
2. Update the AlertManager configuration to add the Slack/Discord webhook.
3. Send a test alert to the system and verify that it appears in the corresponding alert channel (Slack/Discord).
4. Install the `hey` or `ab` tool in your local environment or test server.
5. Run the command `hey -z 30s -c 50 http://<app-url>/` for normal load.
6. During the load test, observe the real-time increase of the request rate in the Grafana panel at the same time.
7. To increase the error rate, apply load to a non-existent endpoint: `hey -z 30s -c 100 http://<app-url>/nonexistent`.
8. Visually show that the alert status changes from `Pending` to `Firing` in the Prometheus and Alertmanager interfaces.
9. Stop the load test and verify that the metrics return to a normal/stable state.
