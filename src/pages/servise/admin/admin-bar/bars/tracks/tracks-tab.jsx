import { useContext, useState, useEffect } from "react";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { uploadTrack } from "@/components/tracks/tracks-api";
import { getAllAlbums } from "@/components/album/album-api";
import { domain } from "@api/api";
import styles from "./tracks-tab.module.css";
import { Button } from "@/components/button/button";
import { Input } from "@/components/input/input";
import TrackCard  from "@/components/tracks/TrackCard"; 
import { AlbumSelect } from "@/components/album/AlbumSelect";

const TracksTab = () => {
    const { user } = useContext(AuthContext);
    const [formData, setFormData] = useState({
        file: null,
        title: "",
        artist_name: "",
        duration: "",
        album_id: "",
        cover_url: ""
    });
    const [uploadedTrack, setUploadedTrack] = useState(null);
    const [streamUrl, setStreamUrl] = useState("");
    const [uploadStatus, setUploadStatus] = useState(null);
    const [albums, setAlbums] = useState([]);

    useEffect(() => {
        const fetchAlbums = async () => {
            const response = await getAllAlbums();
            if (response.ok) {
                const data = await response.json();
                setAlbums(data);
            }
        };
        fetchAlbums();
    }, []);

    const handleChange = (e) => {
        const { name, value, files } = e.target;
        if (name === "file") {
            const fileObj = files[0];
            setFormData(prev => ({ ...prev, file: fileObj }));
            const url = URL.createObjectURL(fileObj);
            const audio = new Audio(url);
            audio.addEventListener("loadedmetadata", () => {
                const duration = Math.round(audio.duration);
                setFormData(prev => ({ ...prev, duration: duration }));
            });
        } else {
            setFormData(prev => ({ ...prev, [name]: value }));
        }
    };

    const handleUpload = async (e) => {
        e.preventDefault();
        if (!formData.file) {
            alert("Выберите файл трека");
            return;
        }
        const data = new FormData();
        data.append("file", formData.file);
        data.append("title", formData.title);
        data.append("artist_name", formData.artist_name);
        data.append("duration", formData.duration);
        data.append("album_id", formData.album_id);
        if (formData.cover_url) {
            data.append("cover_url", formData.cover_url);
        }
        try {
            setUploadStatus("uploading");
            const response = await uploadTrack(data, user.token);
            if (response.ok) {
                const track = await response.json();
                setUploadedTrack(track);
                setStreamUrl(`${domain}/tracks/${track.ID}/stream`);
                setUploadStatus("success");
            } else {
                setUploadStatus("error");
                console.error("Ошибка загрузки трека");
            }
        } catch (error) {
            setUploadStatus("error");
            console.error("Ошибка при загрузке трека", error);
        }
    };

    return (
        <div className={styles.tracksTab}>
            <h2>Загрузка нового трека</h2>
            <form onSubmit={handleUpload} className={styles.uploadForm}>
                <Input type="file" name="file" accept="audio/mp3" onChange={handleChange} required />
                <Input type="text" name="title" placeholder="Название трека" value={formData.title} onChange={handleChange} required />
                <Input type="text" name="artist_name" placeholder="Имя исполнителя" value={formData.artist_name} onChange={handleChange} required />
                <AlbumSelect albums={albums} value={formData.album_id} onChange={handleChange} placeholder="Выберите альбом" />
                <Input type="url" name="cover_url" placeholder="URL обложки (опционально)" value={formData.cover_url} onChange={handleChange} />
                <Button type="submit" text="Загрузить трек" />
            </form>
            {uploadStatus === "uploading" && <p>Загрузка...</p>}
            {uploadStatus === "error" && <p>Ошибка загрузки трека</p>}
            {uploadedTrack && (
                <div className={styles.uploadedTrack}>
                    <h3>Загруженный трек:</h3>
                    <TrackCard track={uploadedTrack} />
                </div>
            )}
        </div>
    );
};

export { TracksTab };