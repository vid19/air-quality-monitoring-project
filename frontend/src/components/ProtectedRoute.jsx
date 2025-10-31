import { useAuth } from "../context/authContext";
import { Navigate } from "react-router-dom";

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    // If not logged in, redirect to the login page
    return <Navigate to="/login" replace />;
  }

  // If logged in, show the child component (e.g., the Dashboard)
  return children;
}

export default ProtectedRoute;
