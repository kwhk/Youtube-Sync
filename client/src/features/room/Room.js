import React, { useEffect, useState } from 'react'
import YoutubePlayer from '../player/YoutubePlayer'
import VideoInput from '../videoQueue/VideoInput'
import VideoQueue from '../videoQueue/VideoQueue'
import { useDispatch } from 'react-redux'
import { setElapsed, setUrl, setPlaybackStatus } from '../currVideo/currVideoSlice' 
import Socket from '../../api/socket'
import SocketContext from '../../context/socket'
import '../../styles/flex.css';
import Ping from '../../api/ping';

export default function Room() {
	const [socket, setSocket] = useState(null)
	const dispatch = useDispatch()
	
	useEffect(() => {
		async function connect() {
			const s = new Socket()
			await s.connect()
			setSocket(s)
		}

		connect()

		return async function closeSocket() {
			await socket.disconnect()
		}
	}, []);

	if (socket != null) {
		socket.on('join', data => {
			dispatch(setUrl(data.videoUrl))
			dispatch(setElapsed(data.videoElapsed))
			dispatch(setPlaybackStatus(data.videoIsPlaying))
			socket.clientID = data.clientID;
			socket.roomID = data.roomID;
			new Ping(socket);
			console.log(data)
		})

		return (
			<SocketContext.Provider value={socket}>
				<div className="d-flex flex-row flex-wrap">
					<YoutubePlayer/>
					<div className="d-flex flex-col">
						<VideoInput/>
						<VideoQueue/>
					</div>
				</div>
			</SocketContext.Provider>
		)
	} else {
		return (
			<div>Not connected</div>
		)
	}
}