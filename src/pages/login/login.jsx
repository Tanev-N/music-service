import { Input } from "@components/input/input";
import { Button } from "@components/button/button";
import { Login } from "./login-api";
import { AuthProvider } from "@/features/auth-provider/auth-provider";
import styles from './login.module.css'
import { useState } from "react";

const LoginPage = () => {
    const [form, setForm] = useState({ login: "", password: "" });

    const onChange = (e) => {
        const { name, value } = e.target;
        setForm((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const onSubmit = async(e) => {
        e.preventDefault();
        const response = await Login(form.login, form.password);
    };

    return (
        <main className={styles.main}>
            <form className={styles.form} onSubmit={onSubmit}>
                <h3>Личный кабинет</h3>
                <Input
                    text="Логин"
                    name="login"
                    label="Логин"
                    type="text"
                    autoComplete="login"
                    required
                    value={form.login}
                    onChange={onChange}
                />
                <Input
                    label="Пароль"
                    text="Пароль"
                    name="password"
                    type="password"
                    autoComplete="current-password"
                    required
                    value={form.password}
                    onChange={onChange}
                />
                <Button text="Войти" type="submit" />
            </form>
        </main>
    );
};

export { LoginPage };
