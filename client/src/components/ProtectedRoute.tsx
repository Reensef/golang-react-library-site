import { ReactNode } from "react";
import { Navigate } from "react-router-dom";

const isTokenValid = () => {
  const token = localStorage.getItem("token");
  if (!token) return false;

  try {
    const payload = JSON.parse(atob(token.split(".")[1])); // Раскодируем payload токена
    const currentTime = Math.floor(Date.now() / 1000);
    return payload.exp > currentTime; // Сравниваем время истечения токена
  } catch {
    return false;
  }
};

const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  if (!isTokenValid()) {
    return <Navigate to="/loginpage" replace />;
  }

  return children;
};

export default ProtectedRoute;
