import React, { useContext, useEffect } from 'react'
import SocketContext from '../context/socket'
import { Redirect } from 'react-router-dom';
import useConnectSocket from '../hooks/useConnectSocket'
import withLoading from './Loading'

export default function ConnectToSocket(props) {
    // const { setLoading } = props
    const socket = useConnectSocket()
    const { setSocket } = useContext(SocketContext)

    useEffect(() => {
        setSocket(socket)
        // setLoading(false)
    }, [socket])

    return socket != null ? <Redirect to={props.match.url}/> : null
    // TODO: Add timeout error message if can't connect
    // Or try reconnecting again
}

// export default withLoading(ConnectToSocket, 'Connecting to server.')