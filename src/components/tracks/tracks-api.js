import { domain } from "@api/api";

// Поиск треков
const searchTracks = async (q) => {
    return fetch(`${domain}/tracks?q=${encodeURIComponent(q)}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include"
    });
};

// Загрузка нового трека (только для администраторов)
const uploadTrack = async (formData, authToken) => {
    return fetch(`${domain}/tracks`, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${authToken}`
            // Не устанавливаем 'Content-Type' для FormData
        },
        body: formData,
        credentials: "include"
    });
};

// Получить детали трека (публичный эндпоинт)
const getTrackDetails = async (id) => {
    return fetch(`${domain}/tracks/${id}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include"
    });
};

// Удалить трек (только для администраторов)
const deleteTrack = async (id, authToken) => {
    return fetch(`${domain}/tracks/${id}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authToken}`
        },
        credentials: "include"
    });
};

// Прослушать трек (потоковая передача MP3 файла)
const streamTrack = async (id) => {
    return fetch(`${domain}/tracks/${id}/stream`, {
        method: "GET",
        headers: {
            "Accept": "audio/mpeg"
        },
        credentials: "include"
    });
};

export { searchTracks, uploadTrack, getTrackDetails, deleteTrack, streamTrack };
