import { useState, useEffect } from "react";
import { getAlbum, deleteAlbum } from "../album-api";
import { Button } from "@/components/button/button";
import styles from "./album-details-modal.module.css";
import { useContext } from "react";
import { AuthContext } from "@/features/auth-provider/auth-provider";

const AlbumDetailsModal = ({ albumId, onClose, onAlbumDeleted }) => {
    const [album, setAlbum] = useState(null);
    const [loading, setLoading] = useState(true);
    const {user} = useContext(AuthContext)
    useEffect(() => {
        const fetchAlbum = async () => {
            try {
                const response = await getAlbum(albumId);
                if (response.ok) {
                    const data = await response.json();
                    setAlbum(data);
                } else {
                    console.error("Ошибка загрузки альбома");
                }
            } catch (error) {
                console.error("Ошибка при получении альбома", error);
            }
            setLoading(false);
        };
        fetchAlbum();
    }, [albumId]);

    const handleDelete = async () => {
        try {
            const response = await deleteAlbum(albumId, user.token);
            if (response.ok) {
                if (onAlbumDeleted) onAlbumDeleted(albumId);
                onClose();
            } else {
                console.error("Ошибка удаления альбома");
            }
        } catch (error) {
            console.error("Ошибка при удалении альбома", error);
        }
    };

    if (loading) {
        return (
            <div className={styles.modal_overlay}>
                <div className={styles.modal_content}>
                    <p>Загрузка...</p>
                    <Button onClick={onClose} text="Закрыть" />
                </div>
            </div>
        );
    }

    if (!album) {
        return (
            <div className={styles.modal_overlay}>
                <div className={styles.modal_content}>
                    <p>Не удалось загрузить данные альбома</p>
                    <Button onClick={onClose} text="Закрыть" />
                </div>
            </div>
        );
    }

    return (
        <div className={styles.modal_overlay}>
            <div className={styles.modal_content}>
                <div className={styles.modal_header}>
                    <h2>{album.title}</h2>
                </div>
                <div className={styles.modal_body}>
                    <img src={album.cover_url} alt={album.title} className={styles.album_cover} />
                    <div className={styles.album_details}>
                        <p><strong>Исполнитель:</strong> {album.artist}</p>
                        <p><strong>Дата выпуска:</strong> {new Date(album.release_date).toLocaleDateString()}</p>
                        <p><strong>Создан:</strong> {new Date(album.created_at).toLocaleString()}</p>
                        <p><strong>Обновлен:</strong> {new Date(album.updated_at).toLocaleString()}</p>
                    </div>
                    <div className={styles.tracks_section}>
                        <h3>Треки:</h3>
                        <div className={styles.tracks_container}>
                            {album.tracks && album.tracks.length > 0 ? (
                                <ul className={styles.tracks_list}>
                                    {album.tracks.map(track => (
                                        <li key={track.id} className={styles.track_item}>
                                            {track.title} ({track.duration} сек.)
                                        </li>
                                    ))}
                                </ul>
                            ) : (
                                <p>Нет треков</p>
                            )}
                        </div>
                    </div>
                </div>
                <div className={styles.modal_footer}>
                    <Button onClick={handleDelete} type="delete" text="Удалить альбом" />
                    <Button onClick={onClose} text="Закрыть" />
                </div>
            </div>
        </div>
    );
};

export { AlbumDetailsModal };