import { AlmubList } from "@/components/album/album";
import { Button } from "@/components/button/button";
import { CreateAlbumModal } from "./album-modal/album-modal";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { useContext, useState } from "react";

const AlbumsTab = () => {
  const { user } = useContext(AuthContext);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [refresh, setRefresh] = useState(0);

  const openCreateAlbumWindow = () => {
    setIsModalOpen(true);
  };

  const handleAlbumCreated = () => {
    setRefresh(prev => prev + 1);
    setIsModalOpen(false);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <AlmubList refresh={refresh} />
      <Button
        type="submit"
        onClick={openCreateAlbumWindow}
        text="Создать новый альбом"
      />
      {isModalOpen && (
        <CreateAlbumModal
          onClose={handleModalClose}
          onAlbumCreated={handleAlbumCreated}
          token={user.token}
        />
      )}
    </>
  );
};

export { AlbumsTab };
