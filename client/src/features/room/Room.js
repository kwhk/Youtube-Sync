import React, { useContext, useEffect, useState } from 'react'
import YoutubePlayer from '../player/YoutubePlayer'
import VideoQueue from '../videoQueue/VideoQueue'
import { useDispatch, useSelector } from 'react-redux'
import { setCurrVideoElapsed, setCurrVideoUrl, setCurrVideoPlaybackStatus } from '../currVideo/currVideoSlice' 
import { setVideoQueue }  from '../videoQueue/videoQueueSlice'
import SocketContext from '../../context/socket'
import ConnectedUsers from '../connectedUsers/ConnectedUsers'
import { setConnectedUsers } from '../connectedUsers/connectedUsersSlice'
import { selectPlayerSize } from '../player/playerSizeSlice'
import { useParams } from 'react-router-dom'

export default function Room(props) {
	const socket = useContext(SocketContext)
	const { theatre } = useSelector(selectPlayerSize)
	const [render, setRender] = useState(false)
	const dispatch = useDispatch()
	const { id } = useParams()

	useEffect(() => {
		console.log('mounting...')
		socket.emit('join-room', id)
		socket.on('room-welcome', data => {
			dispatch(setCurrVideoUrl(data.currVideo.url))
			dispatch(setCurrVideoElapsed(data.currVideo.elapsed))
			dispatch(setCurrVideoPlaybackStatus(data.currVideo.isPlaying))
			dispatch(setConnectedUsers(data.connectedUsers))
			dispatch(setVideoQueue(data.videoQueue))
			setRender(true)
			console.log(data)
		})
		return () => {console.log('leave room'); socket.emit('leave-room')}
	}, [id])

	if (render) {

		if (theatre == 0) {
			return (
				<div className="xs:px-2 sm:px-5 md:px-8 lg:px-10 grid grid-cols-12 gap-10">
					<div className="flex flex-col order-1 xs:col-span-12 sm:col-span-12 md:col-span-8 xl:col-span-9">
						<YoutubePlayer/>
					</div>
					<div className="xs:col-span-12 sm:col-span-6 md:col-span-4 md:order-2 sm:order-3 xl:col-span-3">
						{/* <ConnectedUsers/> */}
						{/* <VideoInput/> */}
						<VideoQueue/>
					</div>
					<div className="xs:col-span-12 md:order-3 sm:order-2 sm:col-span-6">
						<ConnectedUsers/>
					</div>
				</div>
			)
		} else {
			return (
				<div className="grid grid-cols-12 gap-10">
					<div className="col-span-12">
						<YoutubePlayer/>
					</div>
					<div className="xs:col-span-12 sm:col-span-8 lg:col-span-9 xl:col-span-10">
					</div>
					<div className="xs:col-span-12 sm:col-span-4 lg:col-span-3 xl:col-span-2">
						<ConnectedUsers/>
						{/* <VideoInput/> */}
						<VideoQueue/>
					</div>
				</div>
			)

		}
	} else {
		return (
			<div>Not connected to room</div>
		)
	}
}

// function sameRoom(prevProps, nextProps) {
// 	console.log("HIHIHIHIIHII")
// 	console.log(prevProps, nextProps)
// 	console.log(prevProps.match.params.id, nextProps.match.params.id)
// 	return prevProps.match.params.id === nextProps.match.params.id
// }

// export default React.memo(Room, sameRoom)