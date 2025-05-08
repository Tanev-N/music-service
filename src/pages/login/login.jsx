import { Input } from "@components/input/input";
import { Button } from "@components/button/button";
import { Login } from "./login-api";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import styles from "./login.module.css";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useContext } from "react";

const LoginPage = () => {
  const [form, setForm] = useState({ login: "", password: "" });
  const [errors, setErrors] = useState({ password: "" });
  const navigate = useNavigate();
  const { login: authLogin } = useContext(AuthContext);
  const onChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
    setErrors({ password: "" });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (errors.password == "" && form.login != "" && form.password != "") {
      const response = await Login(form.login, form.password);
      if (response.ok) {
        const { id, login, permission } = await response.json();
        authLogin({ id, login, permission });
        navigate("/");
      } else {
        setErrors({
          password: "Неправильный логин или пароль",
        });
      }
    }
  };

  const toRegister = (e) => {
    e.preventDefault();
    navigate("/register");
  };

  return (
    <main className={styles.main}>
      <form className={styles.form}>
        <h3>Log In</h3>
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
          error={errors.password}
          type="password"
          autoComplete="current-password"
          required
          value={form.password}
          onChange={onChange}
        />
        <Button text="Войти" onClick={handleSubmit} type="submit" />
        <Button text="Нет аккаунта" onClick={toRegister} />
      </form>
    </main>
  );
};

export { LoginPage };
