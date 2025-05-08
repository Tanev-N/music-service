import { useNavigate } from "react-router-dom";
import { useContext } from "react";
import { AuthContext } from "../auth-provider/auth-provider";
import styles from "./navbar.module.css";
import { Button } from "@/components/button/button";
const Navbar = () => {
  const { isAuthenticated, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  return (
    <nav className={styles.nav}>
      <div className={styles.logo__container}>
        <img src="/music-icon.png"className={styles.logo__image}></img>
        <span className={styles.logo__name}>Sonix</span>
      </div>
      <div className={styles.nav__elements_container}>
        <div className={styles.nav__element}>
          <Button text="Главная" onClick={() => {navigate("/")}}/>
        </div>
        <div className={styles.nav__element}>
            <Button text="Войти в сервис" type="submit" onClick={() => {navigate("/login")}}/>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
