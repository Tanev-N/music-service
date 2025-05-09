import "./App.css";
import MusicRouter from "./features/router/router";
import { AuthProvider } from "./features/auth-provider/auth-provider";
import Navbar from "./features/navbar/navbar";
import { useContext } from "react";
import { AuthContext } from "./features/auth-provider/auth-provider";

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider> 
  );
}

function AppContent() {
  const { isAuthenticated } = useContext(AuthContext);
  return (
    <>
      {!isAuthenticated && <Navbar/>}
      <MusicRouter />
    </>
  );
}

export default App;
