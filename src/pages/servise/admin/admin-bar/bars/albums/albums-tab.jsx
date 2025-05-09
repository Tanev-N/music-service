import { AlmubList } from "@/components/album/album";
import { Button } from "@/components/button/button";
import { createAlbum } from "@/components/album/album-api";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { useContext } from "react";
import { Input } from "@/components/input/input";
import { useState } from "react";
import { CreateAlbumModal } from "./album-modal/album-modal";
const AlbumsTab = () => {
  const { user } = useContext(AuthContext);
  console.log(user.token)
  const [isModalOpen, setIsModalOpen] = useState(false);

  const openCreateAlbumWindow = () => {
    setIsModalOpen(true);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
  };
  return (
    <>
      <AlmubList />
      <Button
        type="submit"
        onClick={openCreateAlbumWindow}
        text="Создать новый альбом"
      />
      {isModalOpen && (
        <CreateAlbumModal onClose={handleModalClose} token={user.token} />
      )}
    </>
  );
};

export { AlbumsTab };
