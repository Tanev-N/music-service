import {domain} from '@api/api'


const Register = async (login, password) => {
    return fetch(`${domain}/users`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ login: login, password: password }),
        credentials: "include" 
    });
};

export {Register}