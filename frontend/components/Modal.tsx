import { QRCodeSVG } from 'qrcode.react'
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
// import { faCopy } from '@fortawesome/free-solid-svg-icons'
import ClearIcon from '@mui/icons-material/Clear';
import { useState } from 'react';

export default function Modal({ hash, setCloseModal, setCloseController }: { hash: string, setCloseModal: any, setCloseController: any }) {
    const [clipboardMessage, setClipboardMessage] = useState('')
    function updateClipboard() {
        navigator.clipboard.writeText(hash).then(() => {
            console.log("copied to clipboard");
            setClipboardMessage("Copied to clipboard!")
        }, () => {
            console.log("copy to clipboard failed");
        });
    }

    function closeModal() {
        setCloseModal(false)
        setCloseController()
    }
    // if (!open) return null
    return (
        <>
            <div className="z-10 absolute w-screen h-screen border-white border-2 border-dashed mx-auto backdrop-blur-sm">
                <div className="w-2/5 mx-auto mt-32 backdrop-blur-3xl backdrop-brightness-50 flex flex-col">
                    <button className="self-end" onClick={() => { setCloseModal(false) }}><ClearIcon /></button>
                    <QRCodeSVG value={hash} className="h-3/5 w-2/5 mx-auto p-1 bg-white rounded-lg" />
                    <button className="break-all bg-black rounded-lg p-4 mt-4" onClick={updateClipboard}>{hash}</button>
                    {clipboardMessage.length > 0 && <div className="mx-auto">{clipboardMessage}</div>}
                </div>
            </div>

        </>
    )
}