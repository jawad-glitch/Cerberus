import time
import random
import numpy as np
from fastapi import FastAPI
from pydantic import BaseModel
from prometheus_client import Gauge, generate_latest, CONTENT_TYPE_LATEST
from fastapi.responses import Response
from sklearn.ensemble import RandomForestClassifier

app = FastAPI(title="Cerberus Model Server")

def train_model():
    np.random.seed(42)
    X = np.random.randn(1000, 4)  
    y = (X[:, 0] + X[:, 1] > 1).astype(int) 
    clf = RandomForestClassifier(n_estimators=10)
    clf.fit(X, y)
    return clf

model = train_model()

DRIFT_SCORE = Gauge("model_drift_score", "Current drift score of the model")
CONFIDENCE = Gauge("model_prediction_confidence", "Average prediction confidence")
LATENCY = Gauge("model_inference_latency_ms", "Inference latency in milliseconds")

drift_counter = 0

class PredictRequest(BaseModel):
    amount: float
    frequency: float
    location_score: float
    time_of_day: float

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/predict")
def predict(req: PredictRequest):
    global drift_counter
    start = time.time()

    features = np.array([[req.amount, req.frequency, req.location_score, req.time_of_day]])

    drift_counter += 1
    drift_noise = min(drift_counter * 0.001, 1.0)
    features += drift_noise * np.random.randn(*features.shape)

    proba = model.predict_proba(features)[0]
    confidence = float(max(proba))
    prediction = int(np.argmax(proba))

    latency_ms = (time.time() - start) * 1000

    drift_score = min(drift_noise * 2, 1.0)
    DRIFT_SCORE.set(drift_score)
    CONFIDENCE.set(confidence)
    LATENCY.set(latency_ms)

    return {
        "prediction": prediction,
        "confidence": confidence,
        "drift_score": drift_score,
        "is_fraud": bool(prediction)
    }

@app.get("/metrics")
def metrics():
    return Response(generate_latest(), media_type=CONTENT_TYPE_LATEST)