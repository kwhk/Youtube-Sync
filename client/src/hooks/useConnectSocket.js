import React, { useEffect, useState } from 'react'
import Socket from '../api/socket'
import Ping from '../api/ping'

export default function useConnectSocket() {
    const [socket, setSocket] = useState(null)
	useEffect(() => {
		async function connect() {
			const s = new Socket()
			await s.connect()
			// new Ping(s)
			setSocket(s)
		}

		connect()

		return async function closeSocket() {
			if (socket != null) {
				await socket.disconnect()
			}
		}
	}, []);

    return socket
}