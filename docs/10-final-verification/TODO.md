# Final verification
Verification of the project's structural integrity, security rules, and documentation based on the final verification criteria.

## Environment
* Git Repository
* Documentation (README.md)

## Task steps
1. Verify that the repository structure is exactly as follows:
   ├── app/
   ├── helm/gopher/
   ├── argocd/
   └── .github/workflows/
2. Ensure that the steps for setting up and running the project from scratch are fully and clearly written in the `README.md` file.
3. Verify that no sensitive information (secrets, tokens, passwords) remains in plain-text format in the repository.
4. Demo/test the entire CI/CD and GitOps flow live once, starting from the `git push` of the code to the repository up to the complete synchronization of the environment by ArgoCD.
5. Confirm one last time that Prometheus monitoring is working, all metrics show data in Grafana panels, and the dashboard JSON file exists in the repo.
