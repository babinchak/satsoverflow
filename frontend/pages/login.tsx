import { useState } from "react"


export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    async function submitLogin() {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                'email': email,
                'password': password
            })
        })
    }

    return (
        <>
            <div>
                <label>Email:</label>
                <input onChange={(e) => setEmail(e.target.value)} />
                <label>Password:</label>
                <input onChange={(e) => setPassword(e.target.value)} />
                <button onClick={submitLogin}>Submit</button>
            </div>
        </>
    )
}