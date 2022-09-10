import { useState } from "react"
import { useRouter } from "next/router"
import { StatusCodes } from "http-status-codes"

export default function Register() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const router = useRouter();

    async function submitRegister() {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                'email': email,
                'password': password
            })
        })

        if (response.status == StatusCodes.OK) {
            router.push("/login")
        }
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