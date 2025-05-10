import { domain } from "@api/api";

// Создание плейлиста
const createPlaylist = async (name, description, authToken, cover_url) => {
  const body = JSON.stringify({ name, description, cover_url });
  return fetch(`${domain}/playlists`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    body,
    credentials: "include"
  });
};

// Получить плейлисты пользователя
const getUserPlaylists = async (authToken) => {
  return fetch(`${domain}/playlists`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};

// Получить плейлист с треками
const getPlaylistWithTracks = async (id, authToken) => {
  return fetch(`${domain}/playlists/${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};



// Удалить плейлист
const deletePlaylist = async (id, authToken) => {
  return fetch(`${domain}/playlists/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};


// Добавить трек в плейлист
const addTrackToPlaylist = async (playlistId, trackId, authToken) => {
  const body = JSON.stringify({ track_id: trackId });
  return fetch(`${domain}/playlists/${playlistId}/tracks`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    body,
    credentials: "include"
  });
};

// Удалить трек из плейлиста
const removeTrackFromPlaylist = async (playlistId, trackId, authToken) => {
  return fetch(`${domain}/playlists/${playlistId}/tracks/${trackId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};

export {
  createPlaylist,
  getUserPlaylists,
  getPlaylistWithTracks,
  deletePlaylist,
  addTrackToPlaylist,
  removeTrackFromPlaylist
};