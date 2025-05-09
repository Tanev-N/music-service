import styles from "./admin-bar.module.css"
import { Button } from "@/components/button/button"
import { tabs } from "../admin-page"


const AdminBar = ({setTab, tab}) => {
    return (
        <nav className={styles.nav}>
            {tabs.map((tab_) => {
                return <Button key={tab_.name} type={tab == tab_.name ? "submit" : ""} text={tab_.name} onClick={() => {setTab(tab_.name)}}/>
            })}
        </nav>
    )
}

export {AdminBar}