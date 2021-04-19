import { useContext, useEffect, useState } from 'react'
import { useDispatch } from 'react-redux'
import SocketContext from '../context/socket'
import { setConnectedUsers } from '../features/connectedUsers/connectedUsersSlice'
import { setCurrVideo } from '../features/currVideo/currVideoSlice' 
import { setVideoQueue }  from '../features/videoQueue/videoQueueSlice'

export default function useJoinRoom(roomID) {
	console.log('useJoinRoom mounted')
    const { socket } = useContext(SocketContext)
    const [render, setRender] = useState(false)
	const dispatch = useDispatch()

	useEffect(() => {
		socket.emit('join-room', roomID)
		socket.on('room-welcome', data => {
			setRender(true)
			console.log('received room-welcome', data)
			dispatch(setConnectedUsers(data.connectedUsers))
			dispatch(setVideoQueue(data.videoQueue, data.currVideo.index))

			// If index is -1, this means there are no videos in queue
			// so not possible to have current video playing.
			if (data.currVideo.index != -1) {
				dispatch(setCurrVideo(data.currVideo))
			}
		})
		return () => {console.log('leave room'); socket.emit('leave-room')}
	}, [])

	console.log('render val in useJoinRoom:', render)
    return render
}