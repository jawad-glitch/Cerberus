import time
import requests
from neo4j import GraphDatabase
from datetime import datetime

# Config
PROMETHEUS_URL = "http://localhost:9090"
NEO4J_URL = "bolt://localhost:7687"
NEO4J_USER = "neo4j"
NEO4J_PASSWORD = "cerberus-neo4j"
DRIFT_THRESHOLD = 0.85
CHECK_INTERVAL = 15  # seconds

# Neo4j connection
driver = GraphDatabase.driver(NEO4J_URL, auth=(NEO4J_USER, NEO4J_PASSWORD))

def query_drift_score():
    """Query Prometheus for current drift scores across all model servers."""
    try:
        response = requests.get(
            f"{PROMETHEUS_URL}/api/v1/query",
            params={"query": "model_drift_score"}
        )
        response.raise_for_status()
        data = response.json()
        results = data["data"]["result"]
        models = []
        for result in results:
            instance = result["metric"]["instance"]
            drift_score = float(result["value"][1])
            models.append((instance, drift_score))
        return models
    except Exception as e:
        print(f"[ERROR] Failed to query Prometheus: {e}")
        return []

def write_drift_event(session, instance, drift_score):
    """Write a DriftEvent and PolicyTrigger to the Neo4j audit graph."""
    now = datetime.utcnow().isoformat()
    event_id = f"drift-{instance}-{int(time.time())}"
    policy_id = f"policy-{int(time.time())}"

    session.run("""
        CREATE (d:DriftEvent {
            id: $event_id,
            instance: $instance,
            drift_score: $drift_score,
            detected_at: $now
        })
        CREATE (p:PolicyTrigger {
            id: $policy_id,
            fired_at: $now,
            rule: 'drift_score > 0.85',
            action: 'retrain'
        })
        CREATE (d)-[:TRIGGERED]->(p)
    """, event_id=event_id, instance=instance,
         drift_score=drift_score, now=now,
         policy_id=policy_id)

def decision_loop():
    """Main loop — query metrics, evaluate policy, act."""
    print("[Cerberus] Decision engine started")
    print(f"[Cerberus] Checking every {CHECK_INTERVAL}s | Drift threshold: {DRIFT_THRESHOLD}")

    with driver.session(database="neo4j") as session:
        while True:
            models = query_drift_score()

            for instance, drift_score in models:
                print(f"[{datetime.utcnow().isoformat()}] {instance} → drift={drift_score:.4f}")

                if drift_score > DRIFT_THRESHOLD:
                    print(f"[ALERT] Drift threshold exceeded on {instance}!")
                    print(f"[ACTION] Writing audit event to Neo4j...")
                    write_drift_event(session, instance, drift_score)
                    print(f"[ACTION] Retraining trigger fired for {instance}")
                    print(f"[AUDIT] Model {instance} flagged at drift={drift_score:.4f} — PolicyTrigger created")

            time.sleep(CHECK_INTERVAL)

if __name__ == "__main__":
    decision_loop()