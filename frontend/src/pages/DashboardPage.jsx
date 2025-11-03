import { useState, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import api from "../api/axios";

function DashboardPage() {
  const [locations, setLocations] = useState([]);
  const [currentData, setCurrentData] = useState(null);
  const [predictionData, setPredictionData] = useState(null);
  const [loading, setLoading] = useState(true);

  // To show which location is currently active
  const [selectedLocation, setSelectedLocation] = useState(null);

  const [form, setForm] = useState({ name: "", latitude: "", longitude: "" });
  const { logout } = useAuth();

  const fetchLocations = async () => {
    try {
      setLoading(true);
      const response = await api.get("/api/locations");
      setLocations(response.data || []);

      if (response.data && response.data.length > 0) {
        // Automatically fetch weather for the first location on load
        handleFetchWeather(response.data[0]);
      } else {
        setLoading(false);
      }
    } catch (error) {
      console.error("Error fetching locations:", error);
      setLoading(false);
    }
  };

  // --- 👇 THIS IS THE KEY UPGRADE ---
  const handleFetchWeather = async (location) => {
    try {
      setLoading(true);
      setSelectedLocation(location); // Set the active location
      setCurrentData(null);
      setPredictionData(null);

      // Pass the lat/lon as query parameters
      const params = {
        lat: location.latitude,
        lon: location.longitude,
      };

      // Run requests in parallel using our dynamic params
      const [currentRes, predictionRes] = await Promise.all([
        api.get("/api/current-weather", { params }),
        api.get("/api/prediction", { params }),
      ]);

      setCurrentData(currentRes.data);
      setPredictionData(predictionRes.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchLocations();
  }, []);

  const handleFormChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleAddLocation = async (e) => {
    e.preventDefault();
    try {
      const response = await api.post("/api/locations", {
        name: form.name,
        latitude: parseFloat(form.latitude),
        longitude: parseFloat(form.longitude),
      });
      setLocations([...locations, response.data]);
      setForm({ name: "", latitude: "", longitude: "" });
    } catch (error) {
      console.error("Error adding location:", error);
    }
  };

  const handleDeleteLocation = async (id) => {
    if (window.confirm("Are you sure you want to delete this location?")) {
      try {
        await api.delete(`/api/locations/${id}`);
        setLocations(locations.filter((loc) => loc.id !== id));
      } catch (error) {
        console.error("Error deleting location:", error);
      }
    }
  };

  return (
    <div>
      <button onClick={logout} style={{ float: "right" }}>
        Logout
      </button>
      <h1>Your Dashboard</h1>

      <hr />

      <h2>Add a New Location</h2>
      <form onSubmit={handleAddLocation}>
        <input
          name="name"
          value={form.name}
          onChange={handleFormChange}
          placeholder="Location Name (e.g., Home)"
          required
        />
        <input
          name="latitude"
          value={form.latitude}
          onChange={handleFormChange}
          placeholder="Latitude (e.g., 40.71)"
          required
        />
        <input
          name="longitude"
          value={form.longitude}
          onChange={handleFormChange}
          placeholder="Longitude (e.g., -74.00)"
          required
        />
        <button type="submit">Add Location</button>
      </form>

      <hr />

      <h2>My Locations</h2>
      {locations.length === 0 && <p>You have no saved locations.</p>}
      <ul>
        {locations.map((loc) => (
          <li key={loc.id}>
            {loc.name} ({loc.latitude}, {loc.longitude})
            {/* --- 👇 ADDED THIS BUTTON --- */}
            <button
              onClick={() => handleFetchWeather(loc)}
              style={{ marginLeft: "10px" }}
            >
              Get Weather
            </button>
            <button
              onClick={() => handleDeleteLocation(loc.id)}
              style={{ color: "red", marginLeft: "10px" }}
            >
              Delete
            </button>
          </li>
        ))}
      </ul>

      <hr />

      {/* --- 👇 DYNAMIC TITLE --- */}
      <h2>
        Weather Data for {selectedLocation ? selectedLocation.name : "..."}
      </h2>
      {loading && <p>Loading data...</p>}

      <h3>Current Air Quality (from Go/OpenWeather)</h3>
      {currentData && <pre>{JSON.stringify(currentData, null, 2)}</pre>}

      <h3>AI Prediction (from Go - Python)</h3>
      {predictionData && <pre>{JSON.stringify(predictionData, null, 2)}</pre>}
    </div>
  );
}

export default DashboardPage;
