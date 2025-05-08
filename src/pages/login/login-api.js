import {domain} from '@api/api'


const Login = async (login, password) => {
    return fetch(`${domain}/user/auth`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ login, password }),
        credentials: "include" 
    });
};

export {Login}