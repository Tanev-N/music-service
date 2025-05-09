import React, { useContext, useState } from "react";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { removeGenreFromTrack } from "./genre-api";
import styles from "./genreCard.module.css";

const GenreCard = ({ name, id, deletably, idTrack, onRemove }) => {
  const { user } = useContext(AuthContext);
  const token = user.token;
  const isAdmin = user.permission === "admin";
  const [visible, setVisible] = useState(true);

  const handleRemove = async () => {
    const response = await removeGenreFromTrack(idTrack, id, token);
    if (response.ok) {
      setVisible(false);
      onRemove && onRemove(id);
    } else {
      console.error("Ошибка при удалении жанра");
    }
  };

  if (!visible) return null;

  return (
    <div className={styles.genreCard}>
      <span className={styles.genreName}>{name}</span>
      {isAdmin && deletably && (
        <span className={styles.deleteIcon} onClick={handleRemove}>
          ×
        </span>
      )}
    </div>
  );
};

export default GenreCard;
