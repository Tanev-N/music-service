import React, { useState, useEffect, useContext } from "react";
import styles from "./playlist.module.css";
import { getUserPlaylists, createPlaylist } from "./playlists-api";
import { Button } from "../button/button";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { Input } from "../input/input";
import { PlaylistDetailsModal } from "./PlaylistDetailsModal/playlist-details-modal";
// Displays a single playlist card
const PlaylistCard = ({ playlist, onSelect }) => {
  return (
    <div className={styles.playlistCard} onClick={() => onSelect(playlist.ID)}>
      {/* Display cover if provided */}
      {playlist && playlist.CoverURL && (
        <img
          src={playlist.CoverURL}
          alt={playlist.Name}
          className={styles.playlistCover} // Add corresponding CSS
        />
      )}
      <div className={styles.playlistInfo}>
        <h3 className={styles.playlistTitle}>{playlist.Name}</h3>
        <p className={styles.playlistDescription}>{playlist.Description}</p>
      </div>
    </div>
  );
};

// Modal form for creating a new playlist
const CreatePlaylistModal = ({ onClose, onCreate }) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [cover, setCover] = useState(""); // New state for cover URL

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (name.trim() === "") return;
    await onCreate(name, description, cover); // Pass cover to onCreate
    onClose();
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalContent}>
        <h2>Создать плейлист</h2>
        <form onSubmit={handleSubmit} className={styles.createForm}>
          <Input
            placeholder="Название"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          {/* New input for cover URL */}
          <Input
            placeholder="Обложка (URL)"
            value={cover}
            onChange={(e) => setCover(e.target.value)}
          />
          <Input
            placeholder="Описание"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <div className={styles.modalButtons} style={{ gap: "16px" }}>
            <Button type="submit" text="Создать" size="small" />
            <Button type="delete" text="Отмена" size="small" onClick={onClose} />
          </div>
        </form>
      </div>
    </div>
  );
};

// Container for showing the list of playlists and a button to create a new playlist
const PlaylistContainer = () => {
  const [playlists, setPlaylists] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [selectedPlaylistId, setSelectedPlaylistId] = useState(null);
  const { user } = useContext(AuthContext);
  const token = user.token;

  useEffect(() => {
    const fetchPlaylists = async () => {
      const response = await getUserPlaylists(token);
      if (response.ok) {
        const data = await response.json();
        setPlaylists(data);
      } else {
        console.error("Ошибка получения плейлистов");
        setPlaylists([]);
      }
    };
    fetchPlaylists();
  }, [token]);

  const handleCreatePlaylist = async (name, description, cover) => {
    const response = await createPlaylist(name, description, token, cover);
    if (response.ok) {
      const newPlaylist = await response.json();
      setPlaylists((prev) =>
        Array.isArray(prev) ? [...prev, newPlaylist] : [newPlaylist]
      );
    } else {
      console.error("Ошибка создания плейлиста");
    }
  };

  return (
    <div className={styles.playlistContainer}>
      <div className={styles.header}>
        <h2>Плейлисты</h2>
        <Button
          text="Создать плейлист"
          onClick={() => setShowModal(true)}
          type="submit"
          size="small"
        />
      </div>
      <div className={styles.playlistList}>
        {playlists && playlists.length > 0 ? (
          playlists.map((pl) => (
            <PlaylistCard
              key={pl.ID}
              playlist={pl}
              onSelect={(id) => setSelectedPlaylistId(id)}
            />
          ))
        ) : (
          <span className={styles.emptyState}>Плейлисты отсутствуют</span>
        )}
      </div>
      {showModal && (
        <CreatePlaylistModal
          onClose={() => setShowModal(false)}
          onCreate={handleCreatePlaylist}
        />
      )}
      {selectedPlaylistId && (
        <PlaylistDetailsModal
          playlistId={selectedPlaylistId}
          onClose={() => setSelectedPlaylistId(null)}
          onPlaylistDeleted={(deletedId) => {
            setPlaylists((prev) => prev.filter((pl) => pl.ID !== deletedId));
            setSelectedPlaylistId(null);
          }}
        />
      )}
    </div>
  );
};

export { PlaylistContainer };
