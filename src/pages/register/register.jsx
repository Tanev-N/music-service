import { useState, useContext } from "react";
import { Input } from "@components/input/input";
import { Button } from "@components/button/button";
import styles from "@pages/login/login.module.css";
import { Register } from "./register-api";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { useNavigate } from "react-router-dom";

const RegisterPage = () => {
  const [form, setForm] = useState({
    login: "",
    password: "",
    repeatPassword: "",
  });
  const [errors, setErrors] = useState({
    login: "",
    password: "",
    repeatPassword: "",
  });
  const navigate = useNavigate();
  const { login } = useContext(AuthContext);
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });

    if (name === "login") {
      setErrors({
        ...errors,
        login: value.length < 8 ? "Логин должен быть не менее 8 символов" : "",
      });
    }
    if (name === "password") {
      setErrors({
        ...errors,
        password:
          value.length < 8 ? "Пароль должен быть не менее 8 символов" : "",
      });
    }
    if (name === "repeatPassword") {
      if (form.password != form.repeatPassword) {
        setErrors({ ...errors, repeatPassword: "Пароли не совпадают" });
      } else {
        setErrors({ ...errors, repeatPassword: "" });
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (
      errors.login == "" &&
      errors.password == "" &&
      errors.repeatPassword == "" &&
      form.login != "" &&
      form.password != "" &&
      form.repeatPassword != ""
    ) {
      const response = await Register(form.login, form.password);
      if (response.ok) {
        
        navigate("/login");
      } else {
        setErrors({
          ...errors,
          repeatPassword: "Такой пользователь уже существует",
        });
      }
    }
  };
  const toLogin = (e) => {
    e.preventDefault();
    navigate("/login");
  }
  return (
    <main className={styles.main}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <h3>Registration</h3>
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
        <Button text="Есть аккаунт" onClick={toLogin}/>

      </form>
    </main>
  );
};

export { RegisterPage };
