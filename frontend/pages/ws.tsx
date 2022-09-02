import { useEffect, useRef, useState } from "react"


export default function Ws() {
    const [testMessage, setTestMessage] = useState('')
    // let c: WebSocket
    const c = useRef<WebSocket | null>(null)


    useEffect(() => {
        c.current = new WebSocket('ws://localhost:8080/api/ws')
        c.current.onmessage = (msg) => { console.log("Received back ", msg.data) }
        console.log("Updated websocket")
    }, [])

    function echo() {
        console.log(testMessage)
        if (c.current == null) {
            console.log("c is null")
            return
        }
        // if (typeof (c) == 'undefined') {
        //     c = new WebSocket('ws://localhost:8080/api/ws')
        // }
        c.current.send(testMessage)
    }

    return (
        <>
            <div>Hello</div>
            <input onChange={(e) => { setTestMessage(e.target.value) }}></input>
            <button onClick={echo}>Send</button>
        </>
    )
}