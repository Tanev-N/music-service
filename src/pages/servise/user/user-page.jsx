import { Button } from "@/components/button/button";
import { AuthContext } from "@/features/auth-provider/auth-provider";
import { useContext } from "react";
import styles from "../admin/admin-page.module.css";
import stylesUser from "./user-page.module.css";
import { History } from "@/components/history/history";
import { Search } from "@/components/search/search";
const UserPage = () => {
  const { user, logout } = useContext(AuthContext);
  const servises = [
    <History/>,
    <Search/>
  ]
  return (
    <main className={styles.main}>
      <section className={styles.main__section}>
        <h1 className={styles.main__section_title}>Sonix</h1>
        <Button type="delete" text="Выйти" size="small" onClick={logout}/>
      </section>
      <div className={stylesUser.music__servises}>
        {
            servises.map((el) => {
                return <section key={el} className={styles.content__section}>{el}</section>
            })
        }
        
      </div>
    </main>
  );
};

export { UserPage };
