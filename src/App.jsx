import "./App.css";
import MusicRouter from "./features/router/router";
import { AuthProvider } from "./features/auth-provider/auth-provider";
import Navbar from "./features/navbar/navbar";
function App() {
  return (
    <AuthProvider>
      <Navbar/>
      <MusicRouter />
    </AuthProvider> 
  );
}

export default App;
