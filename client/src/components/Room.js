import React, { useEffect, useState } from 'react'
import YoutubePlayer from './YoutubePlayer'
import VideoInput from '../features/videoQueue/VideoInput'
import VideoQueue from '../features/videoQueue/VideoQueue'
import Socket from '../api/socket'
import SocketContext from '../context/socket'
import './flex.css';

export default function Room() {
  const [socket, setSocket] = useState(null);
  
  useEffect(() => {
    async function connect() {
      const s = new Socket();
      await s.connect()
      setSocket(s)
    }

    connect()

    return async function closeSocket() {
      await socket.disconnect()
    }
  }, []);

  if (socket != null) {
    return (
      <SocketContext.Provider value={socket}>
        <div className="d-flex flex-row flex-wrap">
          <YoutubePlayer currPlaying="0-q1KafFCLU"/>
          <div className="d-flex flex-col">
            <VideoInput/>
            <VideoQueue/>
          </div>
        </div>
      </SocketContext.Provider>
    )
  } else {
    return (
      <div>Not connected</div>
    )
  }
}