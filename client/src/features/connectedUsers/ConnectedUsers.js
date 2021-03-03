import React, { useState, useContext, useEffect } from 'react'
import SocketContext from '../../context/socket'
import { useSelector, useDispatch } from 'react-redux'
import { selectConnectedUsers, pushUser, removeUser } from './connectedUsersSlice'

export default function ConnectedUsers() {
    const socket = useContext(SocketContext)
    const { users } = useSelector(selectConnectedUsers)
    const dispatch = useDispatch()
    
    useEffect(() => {
        socket.on('join-room', userId => {
            dispatch(pushUser(userId))
        })

        socket.on('leave-room', userId => {
            console.log('leave room fired', userId)
            dispatch(removeUser(userId))
        })
    }, [])

    return (
        <ul>
            {Object.keys(users).map((user, index) => <li key={index}>{user}</li>)}
        </ul>
    )
}