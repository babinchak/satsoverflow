import { StatusCodes } from "http-status-codes"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import Modal from "../components/Modal"
import Navbar from "../components/Navbar"


export default function Profile() {
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('')
    const [createdDate, setCreatedDate] = useState('')
    const [balance, setBalance] = useState(0);
    const [sats, setSats] = useState('')
    const [openModal, setOpenModal] = useState(false)
    const [invoiceHash, setInvoiceHash] = useState('')
    const [paymentRequest, setPaymentRequest] = useState('')
    const router = useRouter();

    useEffect(() => {
        async function getProfileDetails() {
            const response = await fetch('/api/profile')
            if (response.status == StatusCodes.OK) {
                response.json().then((data) => {
                    setUsername(data.username)
                    setEmail(data.email)
                    setCreatedDate(data.createdDate)
                    setBalance(data.balance)
                })
            } else {
                router.push('/login')
            }

        }
        getProfileDetails()
    }, [])

    async function addFunds() {
        const body = { "sats": parseInt(sats) }
        const resp = await fetch('/api/deposit', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
        resp.json().then((data) => {
            setInvoiceHash(data.payment_request)
            setOpenModal(true)
            waitInvoicePaid(data.hash)
        })
    }

    async function waitInvoicePaid(hash: string) {
        const response = await fetch('/api/waitInvoicePaid?' + new URLSearchParams({ hash: hash }), {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        })
        response.json().then((data) => {
            console.log("status: ", data.status)
            window.location.reload()
        })
    }

    async function withdrawalFunds() {
        const body = { "payment_request": paymentRequest }
        const response = await fetch('/api/withdrawal', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
        if (response.status == StatusCodes.OK) {
            response.json().then((data) => {
                console.log("status: ", data.status)
                window.location.reload()
            })
        }
    }

    return (
        <>
            <Navbar />
            {openModal && <Modal hash={invoiceHash} setCloseModal={setOpenModal} setCloseController={() => { console.log("Close controller") }} />}
            <div>Session username: {username}</div>
            <div>Email: {email}</div>
            <div>Created: {createdDate}</div>
            <div>Balance: {balance}</div>
            <button onClick={addFunds}>Add Funds</button>
            <input type="number" min={0} defaultValue={0} onChange={(e) => { setSats(e.target.value) }}></input>
            <button onClick={withdrawalFunds}>Withdrawal Funds</button>
            <input placeholder="paste invoice" onChange={(e) => { setPaymentRequest(e.target.value) }}></input>
        </>
    )
}