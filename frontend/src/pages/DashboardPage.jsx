import { useState, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import api from "../api/axios";

function DashboardPage() {
  const [locations, setLocations] = useState([]);
  const [currentData, setCurrentData] = useState(null);
  const [predictionData, setPredictionData] = useState(null);
  const [loading, setLoading] = useState(true);

  // State for our new location form
  const [form, setForm] = useState({ name: "", latitude: "", longitude: "" });

  const { logout } = useAuth();

  // This function fetches all the user's saved locations
  const fetchLocations = async () => {
    try {
      setLoading(true);
      const response = await api.get("/api/locations"); // Uses our auth token
      setLocations(response.data || []); // Ensure it's an array

      // If they have locations, fetch weather for the first one
      if (response.data && response.data.length > 0) {
        handleFetchWeather(response.data[0]);
      } else {
        setLoading(false);
      }
    } catch (error) {
      console.error("Error fetching locations:", error);
      setLoading(false);
    }
  };

  // This function fetches weather/prediction for a specific location
  const handleFetchWeather = async (location) => {
    try {
      setLoading(true);
      setCurrentData(null);
      setPredictionData(null);

      //pass the lat/lon to our weather endpoints
      const [currentRes, predictionRes] = await Promise.all([
        api.get("/api/current-weather"), // TODO: Update backend to use lat/lon
        api.get("/api/prediction"), // TODO: Update backend to use lat/lon
      ]);

      setCurrentData(currentRes.data);
      setPredictionData(predictionRes.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
    setLoading(false);
  };

  // Fetch locations when the page loads
  useEffect(() => {
    fetchLocations();
  }, []); //Runs once

  // FORM HANDLERS
  const handleFormChange = (e) => {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  };

  const handleAddLocation = async (e) => {
    e.preventDefault();
    try {
      // Send the new location to our protected backend route
      const response = await api.post("/api/locations", {
        name: form.name,
        latitude: parseFloat(form.latitude),
        longitude: parseFloat(form.longitude),
      });

      // Add the new location to our state and clear the form
      setLocations([...locations, response.data]);
      setForm({ name: "", latitude: "", longitude: "" });
    } catch (error) {
      console.error("Error adding location:", error);
    }
  };

  const handleDeleteLocation = async (id) => {
    if (window.confirm("Are you sure you want to delete this location?")) {
      try {
        await api.delete(`/api/locations/${id}`); // Uses our auth token!
        // Filter out the deleted location from our state
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

      {/*ADD LOCATION FORM*/}
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

      {/*LOCATIONS LIST*/}
      <h2>My Locations</h2>
      {locations.length === 0 && <p>You have no saved locations.</p>}
      <ul>
        {locations.map((loc) => (
          <li key={loc.id}>
            {loc.name} ({loc.latitude}, {loc.longitude})
            {/* We'll add the weather fetch later */}
            {/* <button onClick={() => handleFetchWeather(loc)}>Get Weather</button> */}
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

      {/*WEATHER DISPLAY*/}
      <h2>Weather Data (Hard-coded for now)</h2>
      {loading && <p>Loading data...</p>}

      <h3>Current Air Quality (from Go/OpenWeather)</h3>
      {currentData && <pre>{JSON.stringify(currentData, null, 2)}</pre>}

      <h3>AI Prediction (from Go - Python)</h3>
      {predictionData && <pre>{JSON.stringify(predictionData, null, 2)}</pre>}
    </div>
  );
}

export default DashboardPage;
