import React from 'react';
import { selectVideoQueue } from './videoQueueSlice'
import { useSelector } from 'react-redux'

export default function VideoQueue() {
    const videoQueue = useSelector(selectVideoQueue)
    return (
        <ul>
            {videoQueue.map((video, index) => <li key={index}>{video.title}</li>)}
        </ul>
    )
}