import styles from "./album.module.css";
import { getAllAlbums, getAlbum, deleteAlbum } from "./album-api";
import { useEffect, useState } from "react";
import { AlbumDetailsModal } from "./AlbumDetailsModal/album-details-modal";

const AlmubList = () => {
    const [albums, setAlbums] = useState([]);
    const [selectedAlbumId, setSelectedAlbumId] = useState(null);

    useEffect(() => {
        const fetchAlbums = async () => {
            const response = await getAllAlbums();
            if (response.ok) {
                const data = await response.json();
                setAlbums(data);
            } else {
                setAlbums([]);
            }
        };
        fetchAlbums();
    }, []);

    return (
        <div className={styles.album_container}>
            <h2 className={styles.album_section_title}>Музыкальные альбомы</h2>
            <div className={styles.album_list}>
                {albums && albums.length > 0 ? (
                    albums.map(album => (
                        <AlbumCard
                            key={album.id}
                            id={album.id}
                            title={album.title}
                            artist={album.artist}
                            release_date={album.release_date}
                            cover_url={album.cover_url}
                            onSelect={setSelectedAlbumId}
                        />
                    ))
                ) : (
                    <div className={styles.empty_state}>
                        <p>Пока нет доступных альбомов</p>
                        <p>Добавьте новые альбомы, чтобы они появились здесь</p>
                    </div>
                )}
            </div>
            {selectedAlbumId && (
                <AlbumDetailsModal 
                    albumId={selectedAlbumId} 
                    onClose={() => setSelectedAlbumId(null)}
                    onAlbumDeleted={(deletedId) => {
                        setAlbums(prev => prev.filter(album => album.id !== deletedId));
                        setSelectedAlbumId(null);
                    }}
                />
            )}
        </div>
    );
};

const AlbumCard = ({ id, title, artist, release_date, cover_url, onSelect }) => {
    return (
        <div className={styles.album_card} onClick={() => onSelect(id)}>
            <img src={cover_url} alt={title} className={styles.album_cover} />
            <div className={styles.album_info}>
                <h3 className={styles.album_title}>{title}</h3>
                <p className={styles.album_artist}>{artist}</p>
                <p className={styles.album_date}>{release_date}</p>
            </div>
        </div>
    );
};

const Almub = ({ id }) => {};

export { Almub, AlmubList };