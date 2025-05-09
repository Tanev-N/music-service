import styles from "./admin-page.module.css";
import { useState } from "react";
import { AdminBar } from "./admin-bar/admin-bar";
import { TracksTab } from "./admin-bar/bars/tracks/tracks-tab";
import { GenresTab } from "./admin-bar/bars/genres/genres-tab";
import { AlbumsTab } from "./admin-bar/bars/albums/albums-tab";

const tabs = [
  { name: "Треки", tab: <TracksTab /> },
  { name: "Жанры", tab: <GenresTab /> },
  { name: "Альбомы", tab: <AlbumsTab /> },
];

const AdminPage = () => {
  const [tab, useTab] = useState(null);

  return (
    <main className={styles.main}>
      <section className={styles.main__section}>
        <h1 className={styles.main__section_title}>Панель Администратора</h1>
      </section>
      <AdminBar setTab={useTab} tab={tab} />
      { tab &&
      <section className={styles.content__section}>
        {tab &&
          tabs.map((tab_) => {
            return tab_.name == tab ? tab_.tab : "";
          })}
      </section>
    }
    </main>
  );
};

export { AdminPage, tabs };
