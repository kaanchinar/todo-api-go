## Task steps
1. Create the `.github/workflows/ci.yaml` file.
2. Set up **Job 1 (lint-and-test)** within the workflow: it should execute the `helm lint` check.
3. Set up **Job 2 (build-and-push)** (only on the main branch and if job 1 passes successfully): build the image and push it to GitHub Container Registry (GHCR). Use `latest` + `sha-${{ github.sha }}` as the tag.
4. Set up **Job 3 (update-chart)** (after job 2 finishes): change the `image.tag` value in the `values-prod.yaml` file to the new SHA, commit it with the `[skip ci]` message to avoid infinite commit loops, and push it back to the repository.
5. Add the GHCR token to GitHub Secrets and use it via `secrets.GITHUB_TOKEN` within the workflow.
6. Add the pipeline status badge to the `README.md` file.
7. Install ArgoCD into the Kubernetes cluster.
8. Log in to the ArgoCD UI and change the default admin password.
9. Write the `argocd/apps/dev-app.yaml` manifest — `source: helm/todo-app`, `valueFiles: [values.yaml, values-dev.yaml]`, `namespace: todo-dev`, `syncPolicy: automated + selfHeal + prune`.
10. Write the `argocd/apps/prod-app.yaml` manifest — with the same structure, but using `values-prod.yaml` and `namespace: todo-prod`.
11. Write the `argocd/root-app.yaml` manifest (App-of-Apps pattern) and target the `argocd/apps/` directory.
12. Execute only the `kubectl apply -f argocd/root-app.yaml` command in the cluster; leave the management of all other resources to ArgoCD.