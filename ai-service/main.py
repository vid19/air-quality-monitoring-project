from fastapi import FastAPI

app = FastAPI()

# This is our mock prediction endpoint
# We use @app.post() because in the future,
# we will POST historical data to it for prediction.
@app.post("/predict")
def get_mock_prediction():
    # Return a hard-coded "prediction"
    return {
        "forecast": [
            {"hour": 1, "pm2_5_prediction": 5.2},
            {"hour": 2, "pm2_5_prediction": 5.4},
            {"hour": 3, "pm2_5_prediction": 5.3},
            {"hour": 4, "pm2_5_prediction": 5.1}
        ]
    }