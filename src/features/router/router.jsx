import { Routes, Route, Navigate } from "react-router-dom";
import HomePage from "@pages/home/home";
import { LoginPage } from "@/pages/login/login";
import { RegisterPage } from "@/pages/register/register";
import { ServisePage } from "@/pages/servise/servise";
import { useContext } from "react";
import { AuthContext } from "../auth-provider/auth-provider";

const MusicRouter = () => {
  const { isAuthenticated } = useContext(AuthContext);

  const publicRoutes = [
    { path: "/", element: <HomePage /> },
    { path: "/login", element: <LoginPage /> },
    { path: "/register", element: <RegisterPage /> },
  ];

  const privateRoutes = [{ path: "/servise", element: <ServisePage /> }];

  return (
    <Routes>
      {publicRoutes.map((route) => (
        <Route
          key={route.path}
          path={route.path}
          element={
            isAuthenticated ? <Navigate to="/servise" replace /> : route.element
          }
        />
      ))}
      {privateRoutes.map((route) => (
        <Route
          key={route.path}
          path={route.path}
          element={
            !isAuthenticated ? <Navigate to="/login" replace /> : route.element
          }
        />
      ))}
    </Routes>
  );
};

export default MusicRouter;
