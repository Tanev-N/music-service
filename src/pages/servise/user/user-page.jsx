import { Button } from "@/components/button/button";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { useContext } from "react";
import styles from "../admin/admin-page.module.css";
import stylesUser from "./user-page.module.css";
import { History } from "@/components/history/history";
import { Search } from "@/components/search/search";
import { PlaylistContainer } from "@/components/playlist/playlist";
const UserPage = () => {
  const { user, logout } = useContext(AuthContext);
  const servises = [<History />, <Search />, <PlaylistContainer />];
  return (
    <main className={styles.main}>
      <section className={styles.main__section}>
        <div className={stylesUser.logo__container}>
          <img src="/music-icon4.png" className={stylesUser.logo__image}></img>
          <span className={stylesUser.logo__name}>Sonix</span>
        </div>
        <Button type="delete" text="Выйти" size="small" onClick={logout} />
      </section>
      <div className={stylesUser.music__servises}>
        {servises.map((el, index) => (
          <section key={index} className={stylesUser.content__section}>
            {el}
          </section>
        ))}
      </div>
    </main>
  );
};

export { UserPage };
