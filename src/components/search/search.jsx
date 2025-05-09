import React, { useState } from "react";
import { searchTracks } from "../tracks/tracks-api";
import TrackCard from "../tracks/TrackCard";
import { Button } from "../button/button";
import { Input } from "../input/input";
import styles from "./search.module.css";

const Search = () => {
  const [query, setQuery] = useState("");
  const [trackResults, setTrackResults] = useState([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    if (query.trim() === "") {
      setTrackResults([]);
      return;
    }
    setLoading(true);
    const response = await searchTracks(query);
    if (response.ok) {
      const data = await response.json();
      setTrackResults(data);
    } else {
      console.error("Ошибка поиска треков");
    }
    setLoading(false);
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className={styles.searchContainer}>
      <h1 className={styles.title}>Поиск треков</h1>
      <div className={styles.searchBar}>
        <Input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Введите название трека..."
          onKeyPress={handleKeyPress}
        />
        <Button
          text="Искать"
          onClick={handleSearch}
          type="submit"
          size="small"
        />
      </div>
      <div className={styles.resultsContainer}>
        {loading ? (
          <span>Загрузка...</span>
        ) : trackResults && trackResults.length > 0 ? (
          trackResults.map((track) => (
            <TrackCard
              key={track.ID}
              track={track}
              onRemoveFromAlbum={() => {}}
              isDeletebly={false}
            />
          ))
        ) : (
          <span className={styles.noResults}>Ничего не найдено</span>
        )}
      </div>
    </div>
  );
};

export { Search };
