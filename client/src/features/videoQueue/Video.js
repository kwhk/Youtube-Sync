// Component that represents a video including it's duration, title and url
import React, { useContext } from 'react'
import '../../styles/flex.css'
import socketContext from '../../context/socket'

export default function Video(props) {
    const socket = useContext(socketContext)
    // const [bgColor, setBgColor] = useState('rgba(1, 1, 1, 0)')
    let bgColor = 'rgba(1, 1, 1, 0)'

    const handlePlay = () => {
        socket.in().emit('playVideoQueue', {url: props.url, index: props.index})
    }

    const handleRemove = () => {
        socket.in().emit('removeVideoQueue', {url: props.url, index: props.index})
    }
    
    if (props.active) {
        bgColor = 'rgba(1, 1, 1, 0.3)'
    }

    return (
        <div className="d-flex flex-col" style={{backgroundColor: bgColor}}>
            <div onClick={handlePlay}>
                <h3 style={{margin: 0}}>{props.title}</h3>
                <p>
                    {props.duration}
                    <br></br>
                    {props.channelTitle}
                    <br></br>
                    {props.url}
                </p>
            </div>
            <button onClick={handleRemove}>Remove</button>
        </div>
    )
}