import { useContext } from "react"
import { AuthContext } from "@/features/auth-provider/auth-provider"
const UserPage = () => {
    const {user} = useContext(AuthContext)
    return (
        <main>
            I'm {user.login}!
        </main>
    )
}

export {UserPage}