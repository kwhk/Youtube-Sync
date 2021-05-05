import React, { useContext, useEffect } from 'react'
import SocketContext from '../../context/socket'
import { useSelector, useDispatch } from 'react-redux'
import { selectConnectedUsers, pushUser, removeUser }  from './connectedUsersSlice'
import User from './User'

export default function ConnectedUsers() {
    const { socket } = useContext(SocketContext)
    const { users } = useSelector(selectConnectedUsers)
    const dispatch = useDispatch()
    
    useEffect(() => {
        socket.on('join-room', userId => {
            dispatch(pushUser(userId))
        })

        socket.on('leave-room', userId => {
            dispatch(removeUser(userId))
        })
    }, [])

    return (
        <div className="p-4 bg-gray-900 rounded-xl mb-5">
            <h1 className="text-gray-600 mb-2 font-semibold text-sm">CONNECTED USERS</h1>
            <div className="flex -space-x-1">
                {Object.keys(users).map((user, index) => <User key={index} user={user}/>)}
            </div>
        </div>
    )
}