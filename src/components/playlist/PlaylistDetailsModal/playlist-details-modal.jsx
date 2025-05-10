import React, { useState, useEffect, useContext } from "react";
import styles from "./playlist-details-modal.module.css";
import { getPlaylistWithTracks, deletePlaylist } from "../playlists-api";
import { Button } from "../../button/button";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import TrackCard from "../../tracks/TrackCard";

const PlaylistDetailsModal = ({ playlistId, onClose, onPlaylistDeleted }) => {
  const [playlist, setPlaylist] = useState(null);
  const [loading, setLoading] = useState(true);
  const { user } = useContext(AuthContext);
  const token = user.token;

  useEffect(() => {
    const fetchPlaylist = async () => {
      try {
        const response = await getPlaylistWithTracks(playlistId, token);
        if (response.ok) {
          const data = await response.json();
          // Extract nested playlist and ensure tracks is an array
          const playlistData = { ...data.Playlist, tracks: data.Tracks || [] };
          setPlaylist(playlistData);
        } else {
          console.error("Ошибка загрузки плейлиста");
        }
      } catch (error) {
        console.error("Ошибка при получении плейлиста", error);
      }
      setLoading(false);
    };
    fetchPlaylist();
  }, [playlistId, token]);

  const handleDelete = async () => {
    try {
      const response = await deletePlaylist(playlistId, token);
      if (response.ok) {
        onPlaylistDeleted && onPlaylistDeleted(playlistId);
        onClose();
      } else {
        console.error("Ошибка при удалении плейлиста");
      }
    } catch (error) {
      console.error("Ошибка при удалении плейлиста", error);
    }
  };

  if (loading) {
    return (
      <div className={styles.modalOverlay}>
        <div className={styles.modalContent}>
          <p>Загрузка...</p>
          <Button onClick={onClose} text="Закрыть" size="small" />
        </div>
      </div>
    );
  }

  if (!playlist) {
    return (
      <div className={styles.modalOverlay}>
        <div className={styles.modalContent}>
          <p>Не удалось загрузить данные плейлиста</p>
          <Button onClick={onClose} text="Закрыть" size="small" />
        </div>
      </div>
    );
  }

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h2>{playlist.Name}</h2>
        </div>
        <div className={styles.modalBody}>
          {playlist.CoverURL && (
            <img
              src={playlist.CoverURL}
              alt={playlist.Name}
              className={styles.playlistCover}
            />
          )}
          <div className={styles.playlistDetails}>
            <p>
              <strong>Описание:</strong> {playlist.Description}
            </p>
            <p>
              <strong>Создан:</strong>{" "}
              {new Date(playlist.CreatedDate).toLocaleString()}
            </p>
            <p>
              <strong>Обновлен:</strong> {new Date(playlist.UpdatedAt).toLocaleString()}
            </p>
          </div>
          <div className={styles.tracksSection}>
            <h3>Треки:</h3>
            <div className={styles.tracksList}>
              {playlist.tracks && playlist.tracks.length > 0 ? (
                playlist.tracks.map((track) => (
                  <TrackCard playlistId={playlist.ID} key={track.ID} track={track} isDeletebly={false} isDeleteblyForUser={true} />
                ))
              ) : (
                <p style={{textAlign: "center"}}>Нет треков</p>
              )}
            </div>
          </div>
        </div>
        <div className={styles.modalFooter}>
          <Button
            onClick={handleDelete}
            type="delete"
            text="Удалить плейлист"
            size="small"
          />
          <Button onClick={onClose} text="Закрыть" size="small" />
        </div>
      </div>
    </div>
  );
};

export { PlaylistDetailsModal };