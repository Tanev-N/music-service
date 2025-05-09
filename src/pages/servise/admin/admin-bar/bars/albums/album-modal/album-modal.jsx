import { useState } from 'react';
import { Input } from '@/components/input/input';
import { Button } from '@/components/button/button';
import { createAlbum } from '@/components/album/album-api';
import styles from "./album-modal.module.css";

const CreateAlbumModal = ({ onClose, token, onSuccess }) => {
    const [formData, setFormData] = useState({
        title: '',
        artist: '',
        release_date: '',
        cover_url: ''
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const formattedReleaseDate = formData.release_date
            ? new Date(formData.release_date).toISOString()
            : new Date().toISOString();
        await createAlbum(formData.title, formData.artist, formattedReleaseDate, formData.cover_url, token);
        onSuccess?.();
        onClose();
    };

    return (
        <div className={styles['modal-overlay']}>
            <div className={styles['modal-content']}>
                <h2 className={styles['h2']}>Создать новый альбом</h2>
                <form onSubmit={handleSubmit} className={styles['form2']}>
                    <Input
                        type="text"
                        name="title"
                        value={formData.title}
                        onChange={handleChange}
                        placeholder="Название альбома"
                        required
                    />
                    <Input
                        type="text"
                        name="artist"
                        value={formData.artist}
                        onChange={handleChange}
                        placeholder="Исполнитель"
                        required
                    />
                    <Input
                        type="date"
                        name="release_date"
                        value={formData.release_date}
                        onChange={handleChange}
                        required
                    />
                    <Input
                        type="url"
                        name="cover_url"
                        value={formData.cover_url}
                        onChange={handleChange}
                        placeholder="URL обложки"
                        required
                    />
                    <div className={styles['modal-buttons']}>
                        <Button type="submit" text="Создать" />
                        <Button type="button" onClick={onClose} text="Отмена" />
                    </div>
                </form>
            </div>
        </div>
    );
};

export { CreateAlbumModal };