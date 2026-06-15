# Resilience & Validation
In the final stage, conducting post-deployment tests via Helm Hook, testing rollback mechanisms by intentionally misconfiguring, and validating the entire flow from the git push step to ArgoCD synchronization are carried out.

## Environment
* Kubernetes cluster
* ArgoCD UI / CLI
* Git Repository

## Task steps
1. Create a Helm hook manifest named `templates/post-install-test.yaml`.
2. Define the annotations section within the manifest: `helm.sh/hook: post-install,post-upgrade` and `helm.sh/hook-delete-policy: hook-succeeded`.
3. Run the script `curl -sf http://todo.local/health || exit 1` inside the hook container to check the health of the application.
4. Verify that the hook pod runs successfully and completes after the deployment is executed.
5. Intentionally write a non-existent, incorrect image tag in the `values-prod.yaml` file, commit, and push to the repo.
6. Observe in the cluster that ArgoCD synchronizes (syncs) this change and the pods fall into `ImagePullBackOff` or `CrashLoopBackOff` status.
7. Immediately rollback the application to the previous stable revision via the ArgoCD UI.
8. Run the command `helm history todo-app -n todo-prod` in the cluster to check the release history and changes.
9. Send a request to the `/health` endpoint of the application to confirm that the system has returned to a fully operational state.
