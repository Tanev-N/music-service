import { Link, useLocation } from "react-router";
import styles from "./navbar.module.css";

const Navbar = () => {
    const location = useLocation();

    const navItems = [
        { to: "/", label: "Главная" },
        { to: "/login", label: "Вход" },
        { to: "/register", label: "Регистрация" },
    ];

    return (
        <nav className={styles.nav}>
            {navItems.map((item) => (
                <div
                    key={item.to}
                    className={`${styles.navelement} ${
                        location.pathname === item.to ? styles.active : ""
                    }`}
                >
                    <Link
                        to={item.to}
                        className={styles.link}
                    >
                        {item.label}
                    </Link>
                </div>
            ))}
        </nav>
    );
};

export default Navbar;
