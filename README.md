# Cerberus
 
> Cerberus is a Kubernetes-native platform that autonomously monitors, heals, and audits production ML model fleets in real time — no human intervention required.
 
*If it's not automated, it's not done.*
 
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-operator-blue)](https://kubernetes.io)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8)](https://golang.org)
[![Python](https://img.shields.io/badge/Python-3.11-3776AB)](https://python.org)
[![Status](https://img.shields.io/badge/status-active--development-orange)]()
 
---
 
## The Problem
 
Most companies run ML models that degrade silently — nobody knows until the business feels it.
 
Imagine a fraud detection model trained 6 months ago. Since then, fraudsters changed their behavior. Your model is still confidently making predictions — but it's quietly getting worse. Nobody notices until the business loses money.
 
Now multiply that by 10 models running in production simultaneously, some of which feed into each other.
 
Fixing it today requires a human to:
 
1. Notice something is wrong (usually from a business metric, not a technical one)
2. Manually kick off a retraining job
3. Test the new model
4. Deploy it carefully
5. Watch it for regressions
6. Roll back if it gets worse
That process takes days. Sometimes weeks. Cerberus collapses it to minutes — autonomously.
 
---
 
## How It Works
 
Cerberus implements a continuous **observe → decide → act** loop across your entire model fleet:
 
```
┌─────────────────────────────────────┐
│           Model Fleet               │
│  [Model A] → [Model B] → [Model C] │
└──────────────┬──────────────────────┘
               │ metrics
               ▼
┌─────────────────────────────────────┐
│         Observability Layer         │
│   drift detection · latency · acc  │
└──────────────┬──────────────────────┘
               │ signals
               ▼
┌─────────────────────────────────────┐
│          Decision Engine            │
│  policy eval · cascade analysis ·  │
│  LLM-generated plain-English audit  │
└──────────────┬──────────────────────┘
               │ actions
               ▼
┌─────────────────────────────────────┐
│        Kubernetes Operator          │
│  retrain · canary · promote · roll  │
└──────────────┬──────────────────────┘
               │ lineage
               ▼
┌─────────────────────────────────────┐
│           Audit Graph               │
│  Neo4j — complete causal history    │
│  of every decision ever made        │
└─────────────────────────────────────┘
```
 
Every action Cerberus takes is explained in plain English and stored as a causal graph node:
 
> *"Model v2 was quarantined at 14:32 because drift score exceeded 0.85 on the income feature. Retraining triggered on last 30 days of clean data. Challenger v3 promoted after 15 minutes of canary with 2.1% higher F1."*
 
---
 
## Key Capabilities
 
**Autonomous Remediation**
Detects degradation and acts — retrain, quarantine, promote, or roll back — without waiting for a human.
 
**Cascade-Aware**
Understands model dependency graphs. If Model A degrades, Cerberus knows which downstream models are at risk and responds intelligently.
 
**Kubernetes-Native**
Built as a proper Kubernetes operator using custom resource definitions (`MLModel`, `MLPolicy`). Declarative, GitOps-friendly, installable on any cluster.
 
**Complete Audit Trail**
Every decision stored as a causal graph in Neo4j. Required for regulated industries (finance, healthcare). Most ML systems have zero of this.
 
**Adversarial Simulation**
Built-in chaos engine to inject drift, poisoned batches, and traffic spikes — test your remediation policies before they matter in production.
 
**War Room Dashboard**
Live mission-control UI showing model health, event feed, decision timeline, and dependency map in real time.
 
---
 
## Architecture
 
```
cerberus/
├── operator/                  # Kubernetes operator (Go, controller-runtime)
│   ├── api/v1alpha1/          # CRD type definitions (MLModel, MLPolicy)
│   └── internal/controller/   # Reconciliation logic
├── engine/                    # Decision + adversarial engine (Python)
│   ├── model-server/          # FastAPI model server with /predict + /metrics
│   ├── policy/                # Policy evaluation
│   ├── adversarial/           # Drift injection, attack simulation
│   └── explainer/             # LLM-generated audit narratives
├── dashboard/                 # War room UI (React)
├── charts/                    # Helm charts for Cerberus itself
├── infra/                     # Local dev setup (kind, values files, setup script)
└── docs/                      # Architecture docs, ADRs, Cypher queries
```
 
### Core Stack
 
| Layer | Technology | Purpose |
|---|---|---|
| Orchestration | Kubernetes + custom operator | Desired state reconciliation |
| Metrics | Prometheus | Time-series metrics backbone |
| ML Monitoring | Evidently | Drift detection, statistical tests |
| Audit Graph | Neo4j | Causal decision lineage |
| Model Registry | MLflow | Experiment tracking, model versioning |
| Pipeline | Prefect | Retraining workflow orchestration |
| Canary | Argo Rollouts | Progressive delivery, automatic rollback |
| Explainability | Ollama (local LLM) | Plain-English decision narratives |
| Dashboard | React + WebSockets | Real-time war room UI |
 
---
 
## Roadmap
 
### Phase 1 — Observability Foundation `[Complete]`
> *The eyes of Cerberus*
 
- [x] Local Kubernetes cluster (kind)
- [x] Prometheus deployment
- [x] Neo4j deployment
- [x] Neo4j audit schema design
- [x] Causal event graph — Cypher queries
- [x] Observability stack reproducible via `infra/setup.sh`
### Phase 2 — Kubernetes Operator `[Complete]`
> *The hands of Cerberus*
 
- [x] Go project scaffold with controller-runtime
- [x] `MLModel` CRD definition
- [x] Reconciliation loop — create and update Deployments
- [x] Operator detecting image changes and auto-updating
- [ ] `MLPolicy` CRD definition
- [ ] Unit + integration tests for reconciler
- [ ] Helm chart for the operator itself
### Phase 3 — ML Lifecycle + Adversarial Engine `[In Progress]`
> *The brain and the stress test*
 
- [x] FastAPI model server with `/predict` and `/metrics` endpoints
- [x] Prometheus metrics — drift score, confidence, latency
- [x] Model server packaged as Docker image
- [x] Operator deploying real model server from MLModel declaration
- [ ] Prometheus scraping model server metrics
- [ ] Chained model topology (Model A → Model B)
- [ ] Evidently drift detection pipeline
- [ ] Adversarial engine — covariate shift injection
- [ ] Adversarial engine — label flip simulation
- [ ] Adversarial engine — sudden concept drift
### Phase 4 — Decision Engine + Autonomous Loop `[Planned]`
> *Cerberus acts on its own*
 
- [ ] Policy engine — threshold evaluation
- [ ] Cascade analyzer — dependency graph traversal
- [ ] Autonomous retraining trigger
- [ ] Canary deployment via Argo Rollouts
- [ ] Automatic promotion / rollback logic
- [ ] Ollama integration — plain-English decision narratives
- [ ] Full audit trail written to Neo4j on every action
- [ ] End-to-end autonomous loop demo
### Phase 5 — War Room Dashboard `[Planned]`
> *The face of Cerberus*
 
- [ ] React project scaffold
- [ ] WebSocket connection for real-time updates
- [ ] Live model topology map (health pulsing on nodes)
- [ ] Live event feed
- [ ] Decision timeline with audit trail viewer
- [ ] Chaos injection UI (the button that breaks things)
- [ ] Model dependency graph visualization
### Phase 6 — Production Hardening `[Planned]`
> *From portfolio project to real product*
 
- [ ] RBAC — least-privilege service accounts
- [ ] Secrets management (External Secrets Operator)
- [ ] TLS everywhere (cert-manager)
- [ ] Horizontal pod autoscaling for model servers
- [ ] Backup and restore for Neo4j audit graph
- [ ] Helm chart — single-command install of entire platform
- [ ] GitHub Actions CI — lint, test, build, push
- [ ] Semantic versioning + automated releases
- [ ] Full documentation site
### Phase 7 — Open Source Launch `[Planned]`
> *Ship it*
 
- [ ] Architecture decision records (ADRs) in `/docs`
- [ ] Contribution guide
- [ ] Local quickstart — working in under 10 minutes
- [ ] Demo video
- [ ] Blog post — "How I built an autonomous ML ops platform"
- [ ] ProductHunt / HackerNews launch
---
 
## Getting Started
 
### Prerequisites
 
- Docker
- kind
- kubectl
- helm
- Go 1.22+
- Python 3.11+
### Quickstart
 
```bash
git clone https://github.com/jawad-glitch/Cerberus.git
cd Cerberus
 
# Start Docker, spin up cluster, deploy core stack
bash infra/setup.sh
```
 
> **Note:** Full quickstart guide coming at Phase 3 completion. Currently in active development.
 
---
 
## Contributing
 
Cerberus is being built in public. Contributions, issues, and feedback are welcome.
 
See `CONTRIBUTING.md` (coming soon) for guidelines.
 
---
 
## License
 
MIT © Jawad
 
---
 
*"If it's not automated, it's not done."*