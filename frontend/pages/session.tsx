// import { StatusCodes } from "http-status-codes"
// import { useEffect, useState } from "react"


// export default function SessionDetails() {
//     const [username, setUsername] = useState('')
//     const [email, setEmail] = useState('')
//     const [createdDate, setCreatedDate] = useState('')

//     useEffect(() => {
//         async function getSessionDetails() {
//             const response = await fetch('/api/session')
//             if (response.status == StatusCodes.OK) {
//                 response.json().then((data) => {
//                     setUsername(data.username)
//                     setEmail(data.email)
//                     setCreatedDate(data.createdDate)
//                 })
//             }

//         }
//         getSessionDetails()
//     }, [])

//     return (
//         <>
//             <div>Session username: {username}</div>
//             <div>Email: {email}</div>
//             <div>Created: {createdDate}</div>
//         </>
//     )
// }