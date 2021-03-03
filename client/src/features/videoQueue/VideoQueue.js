import React, { useContext, useEffect } from 'react';
import { selectVideoQueue } from './videoQueueSlice'
import { useSelector, useDispatch } from 'react-redux'
import { pushVideo, removeVideo, setVideoActive } from './videoQueueSlice'
import { setCurrVideo } from '../currVideo/currVideoSlice'
import socketContext from '../../context/socket'
import Video from './Video'
import { getYoutubeVideo } from './utils'

export default function VideoQueue() {
    const { queue } = useSelector(selectVideoQueue)
    const socket = useContext(socketContext)
    const dispatch = useDispatch()

    useEffect(() => {
        socket.on('add-video-queue', async (data) => {
            let videoInfo = await getYoutubeVideo(data.url)
            videoInfo.active = false
            dispatch(pushVideo(videoInfo))
        })

        socket.on('remove-video-queue', data => {
            dispatch(removeVideo(data.index))
        })

        socket.on('play-video-queue', data => {
            dispatch(setCurrVideo({url: data.url, isPlaying: false, elapsed: 0}))
            dispatch(setVideoActive(data.index))
        })
    }, [])

    return (
        <ul>
            {queue.map((video, index) => <li key={index}><Video {...video} index={index}/></li>)}
        </ul>
    )
}