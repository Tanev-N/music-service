import {domain} from '@api/api'


const Login = (login, password) => {
    return fetch(`${domain}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ login, password }),
        credentials: "include" 
    });
};

export {Login}