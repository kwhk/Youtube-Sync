import React, { useState, useContext } from 'react';
import SocketContext from '../../context/socket';
import { getYoutubeVideo, youtubeParser } from './utils';

export default function VideoInput() {
    const [url, setURL] = useState('');
    const socket = useContext(SocketContext);

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (url === "") return;
        const id = youtubeParser(url)
        const videoInfo = await getYoutubeVideo(id)
        socket.broadcast('add-video-queue', {
            duration: videoInfo.duration,
            url: id
        })
        setURL('')
    }

    const handleChange = (e) => {
        setURL(e.target.value)
    }

    return (
        <form className="relative" onSubmit={handleSubmit}>
            <svg width="20" height="20" fill="gray" className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"><path fillRule="evenodd" clipRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"></path></svg>
            <input value={url} onChange={handleChange} type="text" placeholder="Enter Youtube URL to add to playlist" className="w-full bg-secondary text-sm text-white overflow-ellipsis placeholder-gray-500 rounded-md py-2 pl-10 focus:border-light-blue-500 focus:outline-none focus:ring-1 focus:ring-light-blue-500"/>
        </form>
    )
}
