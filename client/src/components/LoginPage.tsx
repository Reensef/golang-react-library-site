import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // Отправляем данные на сервер
    try {
      const response = await fetch("http://localhost:5000/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("token", data.token); // Сохраняем токен
        navigate("/dashboard"); // Перенаправляем на защищенную страницу
      } else {
        alert("Неверный логин или пароль");
      }
    } catch (error) {
      console.error("Ошибка авторизации:", error);
      alert("Ошибка соединения с сервером");
    }
  };

  return (
    <div>
      <h1>Авторизация</h1>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Логин"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button type="submit">Войти</button>
      </form>
    </div>
  );
};

export default LoginPage;
