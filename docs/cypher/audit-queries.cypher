// Full causal chain from a drift event
MATCH path = (d:DriftEvent {id: "drift-001"})-[*..6]->(end)
RETURN path

// Cascade impact — find all downstream models at risk
MATCH (m:ModelVersion)-[:FEEDS_INTO*..5]->(downstream)
RETURN downstream.id, downstream.name, downstream.status

// Full audit trail for a model version
MATCH path = (start)-[*..6]->(m:ModelVersion {id: "fraud-detector-v2"})
RETURN path
