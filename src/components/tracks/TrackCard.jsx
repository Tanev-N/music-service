import React, { useState, useRef, useEffect, useContext } from "react";
import { domain } from "@api/api";
import {
  listAllGenres,
  assignGenreToTrack,
  getGenresByTrack,
} from "@/components/genre/genre-api";
import { writeTrackOnHistory } from "../history/history-api";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import styles from "./track-card.module.css";
import { Button } from "../button/button";
import GenreCard from "../genre/genreCard";
import { deleteTrack } from "./tracks-api";

const TrackCard = ({ track, onRemoveFromAlbum, isDeletebly }) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const audioRef = useRef(null);
  const { user } = useContext(AuthContext);
  const isAdmin = user.permission === "admin";
  const token = user.token;
  const [visible, setVisible] = useState(true);
  const [trackGenres, setTrackGenres] = useState(
    Array.isArray(track.genres) ? track.genres : []
  );
  const [isGenreDropdownOpen, setIsGenreDropdownOpen] = useState(false);
  const [availableGenres, setAvailableGenres] = useState([]);

  // Fetch genres on mount so they persist
  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const response = await getGenresByTrack(track.ID);
        if (response.ok) {
          const genresData = await response.json();
          setTrackGenres(genresData);
        } else {
          console.error("Ошибка загрузки жанров при перезагрузке");
        }
      } catch (error) {
        console.error("Ошибка загрузки жанров", error);
      }
    };
    fetchGenres();
  }, [track.ID]);

  const togglePlay = async () => {
    if (!audioRef.current) return;
    if (isPlaying) {
      audioRef.current.pause();
      setIsPlaying(false);
    } else {
      audioRef.current.play();
      setIsPlaying(true);

      const response = await writeTrackOnHistory(track.ID, token);
      if (!response.ok) {
        console.error("Ошибка при записи трека в историю");
      }
    }
  };

  const handleTimeUpdate = () => {
    const current = audioRef.current.currentTime;
    const duration = audioRef.current.duration;
    setProgress((current / duration) * 100);
    if (current >= duration) {
      setIsPlaying(false);
    }
  };

  const handleRemoveFromAlbum = async () => {
    try {
      const response = await deleteTrack(track.ID, token);
      if (response.ok) {
        setVisible(false);
        onRemoveFromAlbum && onRemoveFromAlbum(track.ID);
      } else {
        console.error("Ошибка при удалении трека");
      }
    } catch (error) {
      console.error("Ошибка при удалении трека", error);
    }
  };

  const toggleGenreDropdown = async () => {
    if (!isGenreDropdownOpen && availableGenres.length === 0) {
      try {
        const response = await listAllGenres();
        if (response.ok) {
          const genresData = await response.json();
          setAvailableGenres(genresData);
        }
      } catch (error) {
        console.error("Ошибка загрузки жанров", error);
      }
    }
    setIsGenreDropdownOpen((prev) => !prev);
  };

  const handleAddGenre = async (genreId) => {
    try {
      const response = await assignGenreToTrack(track.ID, genreId, token);
      if (response.ok) {
        const addedGenre = availableGenres.find((g) => g.id === genreId);
        if (addedGenre) {
          setTrackGenres((prev) =>
            Array.isArray(prev) ? [...prev, addedGenre] : [addedGenre]
          );
        }
        setIsGenreDropdownOpen(false);
      } else {
        console.error("Ошибка при назначении жанра");
      }
    } catch (error) {
      console.error("Ошибка при назначении жанра", error);
    }
  };

  const handleRemoveGenre = (genreId) => {
    setTrackGenres((prev) => prev.filter((g) => g.id !== genreId));
  };

  if (!visible) return null;

  return (
    <div className={styles.trackCard} style={{ width: "600px" }}>
      <div className={styles.coverContainer} onClick={togglePlay}>
        <img
          src={track.CoverURL}
          alt={track.Title}
          className={styles.trackImage}
        />
        <div className={styles.overlay}>
          <span className={styles.playIcon}>{isPlaying ? "❚❚" : "►"}</span>
        </div>
      </div>
      <div className={styles.trackInfo}>
        <h4 className={styles.trackTitle}>{track.Title}</h4>
        <p className={styles.trackArtist}>{track.ArtistName}</p>
        <p className={styles.trackDuration}>{track.Duration} сек.</p>
        <div className={styles.progressBar}>
          <div
            className={styles.progress}
            style={{ width: `${progress}%` }}
          ></div>
        </div>
        <div className={styles.genreListContainer}>
          {trackGenres &&
            trackGenres.map((genre) => (
              <GenreCard
                key={genre.id}
                name={genre.name}
                id={genre.id}
                idTrack={track.ID}
                deletably={true}
                onRemove={() => handleRemoveGenre(genre.id)}
              />
            ))}
        </div>
      </div>
      {isAdmin && isDeletebly && (
        <div className={styles.trackAdminButtons}>
          <Button
            onClick={handleRemoveFromAlbum}
            type="delete"
            size="small"
            text="Удалить из альбома"
          />
          <div className={styles.genreDropdown}>
            <Button
              onClick={toggleGenreDropdown}
              type="submit"
              size="small"
              text="Добавить жанр"
            />
            {isGenreDropdownOpen && (
              <ul className={styles.genreOptions}>
                {availableGenres.map((genre) => (
                  <li
                    key={genre.id}
                    className={styles.genreOption}
                    onClick={() => handleAddGenre(genre.id)}
                  >
                    <GenreCard
                      name={genre.name}
                      id={genre.id}
                      idTrack={track.ID}
                      deletably={false}
                    />
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
      <audio
        ref={audioRef}
        src={`${domain}/tracks/${track.ID}/stream`}
        onTimeUpdate={handleTimeUpdate}
        onEnded={() => setIsPlaying(false)}
        style={{ display: "none" }}
      />
    </div>
  );
};

export default TrackCard;
