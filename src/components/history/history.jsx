import { useState, useEffect, useContext } from "react";
import { getHistory } from "./history-api";
import styles from "./history.module.css";
import TrackCard from "../tracks/TrackCard";
import { AuthContext } from "@/features/auth-provider/auth-provider";

const History = () => {
  const [tracks, setTracks] = useState([]);
  const { user } = useContext(AuthContext);

  useEffect(() => {
    const fetchHistory = async () => {
      const response = await getHistory(user.token);
      if (response.ok) {
        const data = await response.json();
        setTracks(data);
      } else {
        console.error("Ошибка загрузки истории");
      }
    };
    fetchHistory();
  }, [user.token]);

// Функция для форматирования даты
const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
};

return (
    <div className={styles.historyContainer}>
        <div className={styles.servise_title}>История</div>
        {tracks && tracks.length > 0 ? (
            tracks.map((item) => (
                <TrackCard
                    key={item.ID}
                    listened={formatDate(item.ListenedAt)}
                    track={item.Track}
                    onRemoveFromAlbum={() => {}}
                    isDeletebly={false}
                    isForUser={true}
                />
            ))
        ) : (
            <span className={styles.noHistory}>В истории еще нет треков!</span>
        )}
    </div>
);
};

export { History };
