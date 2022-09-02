import { QRCodeSVG } from 'qrcode.react'

export default function Modal({ hash, open }: { hash: string, open: boolean }) {
    if (!open) return null
    return (
        <>
            <div className="z-10 absolute w-screen h-screen border-white border-2 border-dashed flex flex-col mx-auto backdrop-blur-sm">
                <div>This is a modal</div>
                <QRCodeSVG value={hash} className="h-3/5 w-2/5 mx-auto p-1 bg-white" />
                <div className="break-all">hash: {hash}</div>
            </div>

        </>
    )
}