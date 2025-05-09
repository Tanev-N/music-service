import { domain } from "@api/api";

const createAlbum = async (
  title,
  artist,
  release_date,
  cover_url,
  authToken
) => {
  return fetch(`${domain}/albums`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
    body: JSON.stringify({ title, artist, release_date, cover_url }),
    credentials: "include",
  });
};

const deleteAlbum = async (id, authToken) => {
  return fetch(`${domain}/albums/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
    credentials: "include",
  });
};

const getAlbum = async (id) => {
  return fetch(`${domain}/albums/${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });
};

const getAllAlbums = async () => {
  return fetch(`${domain}/albums`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });
};

export { getAllAlbums, createAlbum, getAlbum, deleteAlbum };
