import { useState, useEffect } from "react";
import axios from "axios";
import { useAuth } from "../context/authContext";

function DashboardPage() {
  const [currentData, setCurrentData] = useState(null);
  const [predictionData, setPredictionData] = useState(null);
  const [loading, setLoading] = useState(true);
  const { logout } = useAuth();

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [currentRes, predictionRes] = await Promise.all([
          axios.get("/api/current-weather"),
          axios.get("/api/prediction"),
        ]);
        setCurrentData(currentRes.data);
        setPredictionData(predictionRes.data);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
      setLoading(false);
    };
    fetchData();
  }, []);

  return (
    <div>
      <button onClick={logout} style={{ float: "right" }}>
        Logout
      </button>
      <h1>Air Quality Monitor</h1>
      {loading && <p>Loading all data...</p>}

      <h2>Current Air Quality (from Go/OpenWeather)</h2>
      {currentData && <pre>{JSON.stringify(currentData, null, 2)}</pre>}

      <h2>AI Prediction (from Go - Python)</h2>
      {predictionData && <pre>{JSON.stringify(predictionData, null, 2)}</pre>}
    </div>
  );
}
export default DashboardPage;
