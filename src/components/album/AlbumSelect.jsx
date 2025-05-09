import React, { useState, useRef, useEffect } from "react";
import styles from "./album-select.module.css";

const AlbumSelect = ({ albums, value, onChange, placeholder = "Выберите альбом" }) => {
  const [isOpen, setIsOpen] = useState(false);
  const selectedAlbum = albums && albums.find(a => a.id === value);
  const dropdownRef = useRef(null);

  const toggleDropdown = () => {
    setIsOpen(prev => !prev);
  };

  const handleSelect = (id) => {
    onChange({ target: { name: "album_id", value: id } });
    setIsOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className={styles.albumSelect} ref={dropdownRef}>
      <div className={styles.selected} onClick={toggleDropdown}>
        {selectedAlbum ? (
          <>
            <img src={selectedAlbum.cover_url} alt={selectedAlbum.title} className={styles.albumCover} />
            <span className={styles.albumTitle}>{selectedAlbum.title}</span>
          </>
        ) : (
          <span className={styles.placeholder}>{placeholder}</span>
        )}
        <span className={styles.arrow}>{isOpen ? "▲" : "▼"}</span>
      </div>
      {isOpen && (
        <ul className={styles.options}>
          {albums && albums.map(album => (
            <li key={album.id} className={styles.option} onClick={() => handleSelect(album.id)}>
              <img src={album.cover_url} alt={album.title} className={styles.optionCover} />
              <span className={styles.optionTitle}>{album.title}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export { AlbumSelect };
