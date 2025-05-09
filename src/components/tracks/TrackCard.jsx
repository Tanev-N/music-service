import React, { useState, useRef } from "react";
import { domain } from "@api/api";
import styles from "./track-card.module.css";

const TrackCard = ({ track }) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const audioRef = useRef(null);

  const togglePlay = () => {
    if (!audioRef.current) return;
    if (isPlaying) {
      audioRef.current.pause();
      setIsPlaying(false);
    } else {
      audioRef.current.play();
      setIsPlaying(true);
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

  return (
    <div className={styles.trackCard}>
      <div className={styles.coverContainer} onClick={togglePlay}>
        <img src={track.CoverURL} alt={track.Title} className={styles.trackImage} />
        <div className={styles.overlay}>
          <span className={styles.playIcon}>{isPlaying ? "❚❚" : "►"}</span>
        </div>
      </div>
      <div className={styles.trackInfo}>
        <h4 className={styles.trackTitle}>{track.Title}</h4>
        <p className={styles.trackArtist}>{track.ArtistName}</p>
        <p className={styles.trackDuration}>{track.Duration} сек.</p>
        <div className={styles.progressBar}>
          <div className={styles.progress} style={{ width: `${progress}%` }}></div>
        </div>
      </div>
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
