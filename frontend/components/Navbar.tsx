import { StatusCodes } from "http-status-codes"
import Link from "next/link"
import { useEffect, useState } from "react"

export default function Navbar() {

    async function logout() {
        const response = await fetch('/api/logout', {
            method: 'POST'
        })
        response.text().then(() => {
            window.location.reload()
        })

    }


    const [username, setUsername] = useState('')
    useEffect(() => {
        async function getSessionDetails() {
            const response = await fetch('/api/session')
            if (response.status == StatusCodes.OK) {
                response.json().then((data) => {
                    setUsername(data.username)
                })
            }

        }
        getSessionDetails()
    }, [])

    return (
        <>
            <header>
                <Link href="/">Home</Link>
                {username.length > 0 &&
                    <>
                        <span>{username}</span>
                        <button onClick={logout}>Logout</button>
                    </>
                }
                {username.length == 0 && <Link href='/login'>Login</Link>}
                <Link href="/session">Session</Link>
            </header>
        </>
    )
}