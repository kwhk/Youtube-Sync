import React, { useState, useContext } from 'react';
import SocketContext from '../../context/socket';
import { getYoutubeVideo } from './utils';

export default function VideoInput(props) {
    const [id, setId] = useState('');
    const socket = useContext(SocketContext);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (id === "") return;
        const videoInfo = await getYoutubeVideo(id);
        socket.in().emit('addVideoQueue', {
            duration: videoInfo.duration,
            url: id
        })
        setId('');
    }

    const handleChange = (e) => {
        setId(e.target.value);
    }

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Youtube URL
            </label><br></br>
            <input type="text" value={id} onChange={handleChange}></input>
            <input type="submit" value="Add Video"></input>
        </form>
    )
}
