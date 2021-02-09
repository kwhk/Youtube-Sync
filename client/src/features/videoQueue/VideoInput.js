import React, { useState, useContext } from 'react';
import { push } from './videoQueueSlice';
import { useDispatch } from 'react-redux';
import SocketContext from '../../context/socket';
import YoutubeAPI from '../../api/youtube';
import moment from 'moment';

export default function VideoInput(props) {
    const dispatch = useDispatch()
    const [id, setId] = useState('');
    const socket = useContext(SocketContext);

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (id == "") return;

        const res = await YoutubeAPI.get("/videos", {
            params: {
                id: id
            }
        });

        const data = res.data.items[0];
        const duration = moment.duration(data.contentDetails.duration).asMilliseconds();
        
        socket.emit('addVideoQueue', {
            duration: duration,
            url: id
        })

        let videoInfo = {
            title: data.snippet.title,
            thumbnail: data.snippet.thumbnails.standard,
            duration: duration,
            channelTitle: data.snippet.channelTitle
        }

        dispatch(push(videoInfo));
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
