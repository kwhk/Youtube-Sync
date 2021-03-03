import React, { useContext, useEffect, useState } from 'react'
import YoutubePlayer from '../player/YoutubePlayer'
import VideoInput from '../videoQueue/VideoInput'
import VideoQueue from '../videoQueue/VideoQueue'
import { useDispatch } from 'react-redux'
import { setCurrVideoElapsed, setCurrVideoUrl, setCurrVideoPlaybackStatus } from '../currVideo/currVideoSlice' 
import Socket from '../../api/socket'
import SocketContext from '../../context/socket'
import '../../styles/flex.css'
import Ping from '../../api/ping'
import ConnectedUsers from '../connectedUsers/ConnectedUsers'
import { setConnectedUsers } from '../connectedUsers/connectedUsersSlice'

// export default function Room() {
// 	const [socket, setSocket] = useState(null)
// 	const dispatch = useDispatch()
	
// 	useEffect(() => {
// 		async function connect() {
// 			const s = new Socket()
// 			await s.connect()
// 			setSocket(s)
// 		}

// 		connect()

// 		return async function closeSocket() {
// 			await socket.disconnect()
// 		}
// 	}, []);

// 	if (socket != null) {
// 		socket.on('join', data => {
// 			dispatch(setCurrVideoUrl(data.currVideo.url))
// 			dispatch(setCurrVideoElapsed(data.currVideo.elapsed))
// 			dispatch(setCurrVideoPlaybackStatus(data.currVideo.isPlaying))
// 			dispatch(setConnectedUsers(data.connectedUsers))
// 			socket.clientID = data.clientID;
// 			socket.roomID = data.roomID;
// 			new Ping(socket);
// 			console.log(data)
// 		})

// 		return (
// 			<SocketContext.Provider value={socket}>
// 				<div className="d-flex flex-row flex-wrap">
// 					<div className="d-flex flex-col">
// 						<YoutubePlayer/>
// 						<ConnectedUsers/>
// 					</div>
// 					<div className="d-flex flex-col">
// 						<VideoInput/>
// 						<VideoQueue/>
// 					</div>
// 				</div>
// 			</SocketContext.Provider>
// 		)
// 	} else {
// 		return (
// 			<div>Not connected</div>
// 		)
// 	}
// }

export default function Room(props) {
	const socket = useContext(SocketContext)
	const [render, setRender] = useState(false)
	const dispatch = useDispatch()

	useEffect(() => {
		console.log('room mounted')
		socket.emit('join-room', props.match.params.id)
		socket.on('room-welcome', data => {
			dispatch(setCurrVideoUrl(data.currVideo.url))
			dispatch(setCurrVideoElapsed(data.currVideo.elapsed))
			dispatch(setCurrVideoPlaybackStatus(data.currVideo.isPlaying))
			dispatch(setConnectedUsers(data.connectedUsers))
			socket.roomID = data.roomID;
			new Ping(socket);
			setRender(true)
			console.log(data)
		})

		return () => {socket.emit('leave-room')}
	}, [])

	if (render) {
		return (
			<div>
				<div className="d-flex flex-row flex-wrap">
					<div className="d-flex flex-col">
						<YoutubePlayer/>
						<ConnectedUsers/>
					</div>
					<div className="d-flex flex-col">
						<VideoInput/>
						<VideoQueue/>
					</div>
				</div>
			</div>
		)
	} else {
		return (
			<div>Not connected to room</div>
		)
	}
}