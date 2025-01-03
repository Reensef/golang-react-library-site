import { ReactNode } from "react";
import { Navigate } from "react-router-dom";

const isTokenFrash = () => {
  const token = localStorage.getItem("token");
  if (!token) return false;

  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    const currentTime = Math.floor(Date.now() / 1000);
    return payload.exp > currentTime;
  } catch {
    return false;
  }
};

const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  if (!isTokenFrash()) {
    return <Navigate to="/loginpage" replace />;
  }

  return children;
};

export default ProtectedRoute;
