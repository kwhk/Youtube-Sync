import React, { useState, useEffect, useContext } from 'react'
import SocketContext from '../../context/socket'
import Room from './Room'
import { getRooms } from '../../api/queries'
import { useHistory } from 'react-router-dom'

export default function RoomList() {
    const [rooms, setRooms] = useState([])
    let history = useHistory()
    const { socket } = useContext(SocketContext)

    useEffect(() => {
        async function callDB() {
            let rooms = await getRooms()
            console.log(rooms)
            setRooms(rooms)
        }
        callDB()

        socket.on('create-room', data => {
            console.log(data)
            history.push('/room/' + data.id)
        })
    }, [])

    const createRoom = () => {
        socket.emit('create-room')
    }

    return (
        <div>
            <h1>Room List</h1>
            <a onClick={createRoom}><u>Create Room</u></a>
            <ul>
                {rooms.map((room, index) => <Room key={index} {...room}/>)}
            </ul>
        </div>
    )
}