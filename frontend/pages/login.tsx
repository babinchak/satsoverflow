import { StatusCodes } from "http-status-codes";
import { useRouter } from "next/router"
import { useState } from "react"


export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const router = useRouter();

    async function submitLogin() {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                'email': email,
                'password': password
            })
        })

        if (response.status == StatusCodes.OK) {
            router.push("/")
        }
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