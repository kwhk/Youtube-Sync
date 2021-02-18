import React, { useContext, useEffect } from 'react';
import { selectVideoQueue } from './videoQueueSlice'
import { useSelector, useDispatch } from 'react-redux'
import { push, remove, setActive } from './videoQueueSlice'
import { setVideo } from '../currVideo/currVideoSlice'
import socketContext from '../../context/socket'
import Video from './Video'
import { getYoutubeVideo } from './utils'

/* 

BUGS:
- if there are duplicate videos, selecting another duplicate from the playlist does not refresh the player
- if there are duplicate videos, deleting one of them deletes all other duplicates as well

*/

export default function VideoQueue() {
    const { queue } = useSelector(selectVideoQueue)
    const socket = useContext(socketContext)
    const dispatch = useDispatch()

    useEffect(() => {
        socket.on('addVideoQueue', async (data) => {
            let videoInfo = await getYoutubeVideo(data.url)
            videoInfo.active = false
            dispatch(push(videoInfo))
        })

        socket.on('removeVideoQueue', data => {
            dispatch(remove(data.index))
        })

        socket.on('playVideoQueue', data => {
            dispatch(setActive(data.index))
            dispatch(setVideo({url: data.url, isPlaying: false, elapsed: 0}))
        })
    }, [])

    return (
        <ul>
            {queue.map((video, index) => <li key={index}><Video {...video} index={index}/></li>)}
        </ul>
    )
}