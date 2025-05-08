import { useState } from "react";
import { Input } from "@components/input/input";
import { Button } from "@components/button/button";
import styles from "@pages/login/login.module.css";
import { Register } from "./register-api";
import { AuthProvider } from "@/features/auth-provider/auth-provider";
import { Navigate } from "react-router";

const RegisterPage = () => {
    const [form, setForm] = useState({
        login: "",
        password: "",
        repeatPassword: ""
    });
    const [errors, setErrors] = useState({
        login: "",
        password: "",
        repeatPassword: ""
    })

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit =  async (e) => {
        e.preventDefault();
        if (form.password != form.repeatPassword)
        {
            setErrors({ ...errors, repeatPassword:"Пароли не совпадают" });
        }
        else {
            setErrors({ ...errors, repeatPassword:"" });
            const response = await Register(form.login, form.password);
            if (response.ok) {
                const {login, permission} = response.json()
                // Вот тут нужно использовать Auth провайдер и редиректить на /
            }
            else {

            }
        }
        
    };

    return (
        <main className={styles.main}>
            <form className={styles.form} onSubmit={handleSubmit}>
                <h3>Регистрация</h3>
                <Input
                    text="Логин"
                    name="login"
                    label="Логин"
                    value={form.login}
                    error={errors.login}
                    onChange={handleChange}
                />
                <Input
                    text="Пароль"
                    type="password"
                    name="password"
                    label="Пароль"
                    value={form.password}
                    error={errors.password}
                    onChange={handleChange}
                />
                <Input
                    text="Повторите пароль"
                    type="password"
                    name="repeatPassword"
                    label="Повторите пароль"
                    error={errors.repeatPassword}
                    value={form.repeatPassword}
                    onChange={handleChange}
                />
                <Button text="Зарегистрироваться" type="submit" />
            </form>
        </main>
    );
};

export { RegisterPage };
