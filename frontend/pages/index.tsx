import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import React from 'react'
import { useState, useEffect, useRef } from 'react'
import Modal from '../components/Modal'
import Navbar from '../components/Navbar'
// import styles from '../styles/Home.module.css'

// export async function getServerSideProps() {
//   const response = await fetch('https://localhost:3000/message', {
//     method: 'GET',
//     headers: { 'Content-Type': 'application/json' },
//   })
//   console.log(response.json())
// }
function QuestionLink({ title, bounty, page_id }: { title: string, bounty: number, page_id: number }) {
  return (
    <div className="hover:bg-slate-700">
      <Link href={`/question/${page_id.toString()}`} >
        <a className="grid grid-flow-col">
          <span>{title}</span>
          <span className="justify-self-end">{bounty} sats</span>
        </a>
      </Link>
    </div>

  )
}

const Home: NextPage = () => {
  const [questionTitle, setQuestionTitle] = useState('');
  const [questionBody, setQuestionBody] = useState('');
  const [questionSats, setQuestionSats] = useState('');
  const [questions, setQuestions] = useState<any>([]);
  const [openModal, setOpenModal] = useState(false);
  const [invoiceHash, setInvoiceHash] = useState('');
  const [controller, setController] = useState<AbortController | null>(null)
  // const controller = new AbortController()
  // const controller = useRef
  // const controller = useRef<AbortController | null>(null);
  const c = useRef<WebSocket | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch('/api/questions', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
      })
      response.json().then((data) => {
        console.log("here", data.messages)
        let posts = []
        for (const message of data.messages) {
          if (message == null) { break }
          console.log(message.title)
          posts.push({ title: message.title, bounty: message.bounty, id: message.id })

        }
        setQuestions(posts)
      })
    }
    fetchData()
    setController(new AbortController())
    // c.current = new WebSocket('ws://localhost:8080/api/invoice/ws')
    // c.current.onmessage = (msg) => {
    //   console.log("Received back ", msg.data)
    //   setInvoiceHash(msg.data)
    //   setOpenModal(true)
    //   if (c.current != null) {
    //     c.current.onmessage = (msg) => {
    //       if (msg.data == "Settled") {
    //         console.log("Settled!")
    //       }
    //     }
    //   }
    // }
  }, [])

  async function submitQuestion() {
    console.log("Hi")
    const body = {
      "title": questionTitle,
      "body": questionBody,
      "bounty": parseInt(questionSats)
    }
    console.log(JSON.stringify(body))
    // if (c.current != null) {
    //   c.current.send(JSON.stringify(body))
    // }
    const response = await fetch('/api/question', {
      method: 'POST',
      headers: {
        // Accept: 'application.json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    })
    // console.log(response.json())
    response.json().then((data) => {
      console.log("payrequest: ", data.payment_request)
      setInvoiceHash(data.payment_request)
      setOpenModal(true)
      waitInvoicePaid(data.hash)
    })
  }

  async function waitInvoicePaid(hash: string) {
    if (controller == null) return
    const response = await fetch('/api/waitInvoicePaid?' + new URLSearchParams({ hash: hash }), {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal: controller.signal
    })
    response.json().then((data) => {
      console.log("status: ", data.status)
      window.location.reload()
    })
  }

  function abortInvoice() {
    if (controller != null) {
      controller.abort()
    }
  }

  const handleTitle = (event: React.ChangeEvent<HTMLInputElement>) => {
    // console.log("Changing to ", event.target.value)
    setQuestionTitle(event.target.value)
  }

  const handleBody = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setQuestionBody(event.target.value)
  }

  const handleSats = (event: React.ChangeEvent<HTMLInputElement>) => {
    // console.log(event.target.value)
    setQuestionSats(event.target.value)
  }

  return (
    <>
      <Navbar />
      {openModal && <Modal hash={invoiceHash} setCloseModal={setOpenModal} setCloseController={abortInvoice} />}
      {/* <div className="w-2/5 mx-auto">
        <div className="text-slate-700">stay humble</div>
        <div className="text-slate-700">stack <span className="text-white">sats</span></div>
        <div className="text-white">overflow</div>
      </div> */}

      <div className="border-white border-2 border-dashed w-3/5 min-w-min flex flex-col mx-auto mt-16">
        <Image src="/shss.svg" height={300} width={600} />
        <div>Get your question answered. Lightning Fast.</div>
        <div className="flex flex-nowrap border-red-900 border-2">
          <input type="text" onChange={handleTitle} placeholder="How do I escape vim?" className="basis-2/3 border-2 mr-1 rounded-lg h-12 text-xl"></input>
          <div className="my-auto basis-1/12">Bounty:</div>
          <div className="basis-1/6">
            <div>
              <input type="number" min={0} max={10} step={0.01} defaultValue={0} className="rounded-tl rounded-bl"></input>
              <span> $ </span>
            </div>
            <div>
              <input type="number" min={0} defaultValue={0} onChange={handleSats}></input>
              <span> sats </span>
            </div>
          </div>
          <button className="basis-1/12 hover:bg-slate-700 rounded-lg" onClick={submitQuestion}>Submit & Pay Bounty</button>
        </div>
        <textarea rows={3} cols={60} name="text" onChange={handleBody} placeholder="Question description..."></textarea>

        <div>
          {
            questions.map((post: any) => <QuestionLink key={post.id} title={post.title} bounty={post.bounty} page_id={post.id}></QuestionLink>)
          }
        </div>
      </div>

    </>
  )
}

export default Home
