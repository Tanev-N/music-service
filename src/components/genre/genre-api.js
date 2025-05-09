import { domain } from "@api/api";

// Получить список всех жанров
const listAllGenres = async () => {
  return fetch(`${domain}/genres`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json"
    },
    credentials: "include"
  });
};


// Получить жанры для трека (публичный эндпоинт)
const getGenresByTrack = async (trackId) => {
  return fetch(`${domain}/genres/tracks/${trackId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json"
    },
    credentials: "include"
  });
};

// Назначить жанр треку (только для администраторов)
const assignGenreToTrack = async (trackId, genreId, authToken) => {
  return fetch(`${domain}/genres/tracks/${trackId}/genres`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    body: JSON.stringify({ genre_id: genreId }),
    credentials: "include"
  });
};

// Удалить жанр у трека (только для администраторов)
const removeGenreFromTrack = async (trackId, genreId, authToken) => {
  return fetch(`${domain}/genres/tracks/${trackId}/genres/${genreId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};

// Добавляем функцию createGenre для создания жанра
const createGenre = async (name, authToken) => {
  return fetch(`${domain}/genres`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${authToken}`
    },
    body: JSON.stringify({ name }),
    credentials: "include"
  });
};

export { createGenre, listAllGenres, getGenresByTrack, assignGenreToTrack, removeGenreFromTrack };
