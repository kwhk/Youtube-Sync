import React, { useEffect, useContext } from 'react'
import { useHistory } from 'react-router-dom'
import Button from '../../components/Button'
import SocketContext from '../../context/socket'

export default function HomePage() {
    const socket = useContext(SocketContext)
    let history = useHistory()

    useEffect(() => {
        socket.on('create-room', data => {
            console.log(data)
            history.push('/room/' + data.id)
        })
    }, [])

    const createRoom = () => {
        socket.emit('create-room')
    }

    return (
        <div className="h-screen w-screen flex flex-col items-center justify-center">
            <h1 className="text-white text-5xl font-bold">Watch Youtube with friends</h1>
            <div className="mt-5">
                <Button onClick={createRoom} bgColor="bg-highlight" text="Create a room" fontSize="text-xl"/>
            </div>
        </div>
    )
}