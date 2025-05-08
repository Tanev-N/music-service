import { useState } from "react";
import { Input } from "@components/input/input";
import { Button } from "@components/button/button";
import styles from "@pages/login/login.module.css";

const RegisterPage = () => {
    const [form, setForm] = useState({
        login: "",
        password: "",
        repeatPassword: "",
    });

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        // Здесь обработка регистрации
        console.log(form);
    };

    return (
        <main className={styles.main}>
            <form className={styles.form} onSubmit={handleSubmit}>
                <h3>Регистрация</h3>
                <Input
                    text="Логин"
                    name="login"
                    value={form.login}
                    onChange={handleChange}
                />
                <Input
                    text="Пароль"
                    type="password"
                    name="password"
                    value={form.password}
                    onChange={handleChange}
                />
                <Input
                    text="Повторите пароль"
                    type="password"
                    name="repeatPassword"
                    value={form.repeatPassword}
                    onChange={handleChange}
                />
                <Button text="Зарегистрироваться" type="submit" />
            </form>
        </main>
    );
};

export { RegisterPage };
