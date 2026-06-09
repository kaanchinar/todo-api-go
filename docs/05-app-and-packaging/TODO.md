## Task steps
1. Write a minimal application (app) in any language (Node.js / Python / Go).
2. Set up 3 endpoints in the application: `/` (version information in JSON format), `/health`, `/metrics` (using the Prometheus client library).
3. Implement these metrics in the `/metrics` endpoint: `http_requests_total` (counter), `http_request_duration_seconds` (histogram), `http_errors_total` (counter).
4. Write a multi-stage Dockerfile (builder + runtime stage).
5. Add a `.dockerignore` file.
6. Ensure the final image size is less than 100MB.
7. Create a chart using the `helm create gopher` command and delete unnecessary default files.
8. Fill in the `Chart.yaml` file with the appropriate information.
9. Define these values in the `values.yaml` file: `image.repository`, `image.tag`, `replicaCount`, `service`, `ingress`, `hpa`, `resources`, `serviceMonitor`.
10. Write the `values-dev.yaml` file: set 1 replica and an ingress host suitable for the dev environment.
11. Write the `values-prod.yaml` file: set 3 replicas, HPA active, prod ingress host, and serviceMonitor active.
12. Make the image tag dynamic in the `templates/deployment.yaml` files in the format `{{ .Values.image.tag }}`.
13. Make the `templates/hpa.yaml` and `templates/servicemonitor.yaml` files conditional (dependent on conditions) using the `{{ if .Values.hpa.enabled }}` condition.
14. Check the chart using the `helm lint helm/gopher` command (there should be no errors).
15. Run the `helm template gopher helm/gopher -f helm/gopher/values-dev.yaml` command to verify the manifest output.