import { Routes, Route } from "react-router-dom";
import HomePage from "@pages/home/home";
import { LoginPage } from "@/pages/login/login";
import { RegisterPage } from "@/pages/register/register";
const MusicRouter = () => {
  const routes = [
    { path: "/", element: <HomePage /> },
    { path: "/login", element: <LoginPage /> },
    { path: "/register", element: <RegisterPage /> }
  ];
  return (
    <Routes>
      {routes.map((route) => (
        <Route key={route.path} {...route} />
      ))}
    </Routes>
  );
};

export default MusicRouter;
