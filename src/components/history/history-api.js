import { domain } from "@api/api";

// Получить историю юзера
const getHistory = async (authToken) => {
  return fetch(`${domain}/history`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`
    },
    credentials: "include"
  });
};


// Записать трек в историю
const writeTrackOnHistory = async (trackId, authToken) => {
  return fetch(`${domain}/history/tracks/${trackId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
       "Authorization": `Bearer ${authToken}`
    },
    credentials: "include"
  });
};



export { getHistory, writeTrackOnHistory };
