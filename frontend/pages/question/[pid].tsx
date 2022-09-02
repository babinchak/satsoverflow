import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'

function timeSince(date: number) {
    const now = Date.now()

    var seconds = Math.floor((now - date) / 1000);

    var interval = seconds / 31536000;

    if (interval > 1) {
        return Math.floor(interval) + " years";
    }
    interval = seconds / 2592000;
    if (interval > 1) {
        return Math.floor(interval) + " months";
    }
    interval = seconds / 86400;
    if (interval > 1) {
        return Math.floor(interval) + " days";
    }
    interval = seconds / 3600;
    if (interval > 1) {
        return Math.floor(interval) + " hours";
    }
    interval = seconds / 60;
    if (interval > 1) {
        return Math.floor(interval) + " minutes";
    }
    return Math.floor(seconds) + " seconds";
}

function Answer({ body }: { body: string }) {
    return (
        <div>
            {body}
        </div>
    )
}

export default function Question() {
    const [title, setTitle] = useState('')
    const [body, setBody] = useState('')
    const [bounty, setBounty] = useState(0)
    const [created, setCreated] = useState('')
    const [answerBody, setAnswerBody] = useState('')
    const [answers, setAnswers] = useState<any>([]);
    const [qr, setQr] = useState('1234');
    const router = useRouter()
    const { pid } = router.query


    useEffect(() => {
        if (typeof (pid) != 'string') {
            return
        }
        const fetchData = async () => {
            console.log("pid = ", pid)

            const response = await fetch('/api/question?' + new URLSearchParams({ id: pid }), {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            })
            response.json().then((data) => {
                setTitle(data.title)
                setBody(data.body)
                setBounty(data.bounty)
                const time = Date.parse(data.created)
                const timeElapsed = timeSince(time)
                setCreated(timeElapsed)
            })
        }


        const fetchAnswers = async () => {
            const response = await fetch('/api/answers?' + new URLSearchParams({ question_id: pid }), {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            })
            response.json().then((data) => {
                console.log("Fetching answers", data.answers)
                let bodys = []
                for (const ans of data.answers) {
                    bodys.push({ body: ans.body })

                }
                setAnswers(bodys)
            })
        }
        fetchData()
        fetchAnswers()
    }, [pid])

    async function submitAnswer() {
        console.log("Hi")
        if (typeof (pid) != 'string') {
            console.log('pid not correct type')
            return
        }
        const body = {
            "body": answerBody,
            "question_id": parseInt(pid)
        }
        console.log(JSON.stringify(body))
        const response = await fetch('/api/answer', {
            method: 'POST',
            headers: {
                // Accept: 'application.json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        })
        console.log(response.json())
    }

    return (
        <>
            <div className="mx-auto border-white border-2 w-3/5 min-w-min flex flex-col mt-16">
                <div className="text-3xl">{title}</div>
                <div className="flex flex-row gap-x-16">
                    <span>{created} ago</span>
                    <span className="flex">Bounty: {bounty} sats</span>
                </div>
                <div className="text-lg">{body}</div>
                <div>Answer:</div>
                <textarea onChange={(e) => setAnswerBody(e.target.value)}></textarea>
                <button onClick={submitAnswer}>Submit</button>

                <div>
                    {
                        answers.map((answer: any) => <Answer key={answer.body} body={answer.body}></Answer>)
                    }
                </div>
                <QRCodeSVG value={qr} className="w-3/5 h-3/5" />
                <input onChange={(e) => { setQr(e.target.value) }}></input>
            </div>
        </>
    )
}