import { useNavigate } from "react-router-dom";
import styles from "./navbar.module.css";
import { Button } from "@/components/button/button";
const Navbar = () => {
  const navigate = useNavigate();

  return (
    <nav className={styles.nav}>
      <div className={styles.logo__container}>
        <img src="/music-icon4.png"className={styles.logo__image}></img>
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
