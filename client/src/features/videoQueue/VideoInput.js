import React, { useState } from 'react';
import { push } from './videoQueueSlice';
import { useDispatch } from 'react-redux';

export default function VideoInput(props) {
    const dispatch = useDispatch()
    const [name, setName] = useState('');

    const handleSubmit = (e) => {
        console.log(name);
        dispatch(push(name))
        setName('')
        e.preventDefault();
    }

    const handleChange = (e) => {
        setName(e.target.value);
    }

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Youtube URL
            </label><br></br>
            <input type="text" value={name} onChange={handleChange}></input>
            <input type="submit" value="Add Video"></input>
        </form>
    )
}
