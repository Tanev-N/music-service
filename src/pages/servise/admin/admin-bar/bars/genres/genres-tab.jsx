import React, { useState, useEffect, useContext } from "react";
import { createGenre, listAllGenres } from "@/components/genre/genre-api";
import { Button } from "@/components/button/button";
import { Input } from "@/components/input/input";
import GenreCard from "@/components/genre/genreCard";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import styles from "./genres-tab.module.css";

const GenresTab = () => {
  const { user } = useContext(AuthContext);
  const [genres, setGenres] = useState([]);
  const [newGenre, setNewGenre] = useState("");
  const [refresh, setRefresh] = useState(0);
  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const response = await listAllGenres();
        if (response.ok) {
          const data = await response.json();
          setGenres(data);
        } else {
          console.error("Ошибка загрузки жанров");
        }
      } catch (error) {
        console.error("Ошибка при получении жанров", error);
      }
    };
    fetchGenres();
  }, [refresh]);

  const handleCreateGenre = async (e) => {
    e.preventDefault();
    if (!newGenre.trim()) return;
    try {
      const response = await createGenre(newGenre, user.token);
      if (response.ok) {
        setNewGenre("");
        setRefresh((prev) => prev + 1);
      } else {
        console.error("Ошибка создания жанра");
      }
    } catch (error) {
      console.error("Ошибка при создании жанра", error);
    }
  };

  return (
    <div className={styles.genresTab}>
      <h2 className={styles.title}>Жанры</h2>
      <div className={styles.genreList}>
        {genres && genres.length > 0 ? (
          genres.map((genre) => (
            <GenreCard
              key={genre.id || genre.ID}
              name={genre.name || genre.Name}
              id={genre.id || genre.ID}
              deletably={false}
            />
          ))
        ) : (
          <p>Жанры отсутствуют</p>
        )}
      </div>
      <form onSubmit={handleCreateGenre} className={styles.genreForm}>
        <Input
          type="text"
          name="genre"
          value={newGenre}
          onChange={(e) => setNewGenre(e.target.value)}
          placeholder="Введите название жанра"
        />
        <Button type="submit" text="Создать жанр" />
      </form>
    </div>
  );
};

export { GenresTab };
