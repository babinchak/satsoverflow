import { useState } from "react"


export default function Register() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    async function submitRegister() {
        const response = await fetch('/api/register', {
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
                <button onClick={submitRegister}>Submit</button>
            </div>
        </>
    )
}