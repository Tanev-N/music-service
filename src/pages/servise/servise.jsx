import { useContext } from "react"
import { AuthContext } from "@features/auth-provider/auth-provider"
import { AdminPage } from "./admin/admin-page"
import { UserPage } from "./user/user-page"
const ServisePage = () => {
    const { user } = useContext(AuthContext)

    return (
        <>
            {user.permission == "admin" ? <AdminPage/> : <UserPage/>}
        </>
    )
}

export {ServisePage}