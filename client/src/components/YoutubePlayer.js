import React, { Component, useState, useEffect } from 'react';
import socketIo from "socket.io-client";

import Youtube from 'react-youtube';

export default function YoutubePlayer() {
  const opts = {
    height: '390',
    width: '640',
    playerVars: {
      // https://developers.google.com/youtube/player_parameters
      controls: 0,
      disablekb: 0,
      modestbranding: 1
    },
  }

  const [player, setPlayer] = useState(null);

  const play = () => {
    console.log('playing');
    player.playVideo();
  }
  
  const pause = () => {
    console.log('pause');
    player.pauseVideo();
  }

  const onReady = (e) => {
    console.log('onReady...');
    setPlayer(e.target);
  }

  useEffect(() => {
    console.log('socket connected');
    const socket = socketIo(process.env.REACT_APP_API_HOST);
    return () => socket.disconnect();
  });

  return (
    <div>
      <Youtube videoId="D1PvIWdJ8xo" opts={opts} onReady={onReady}/>
      <div>
        <button onClick={play}>play</button>
        <button onClick={pause}>pause</button>
      </div>
    </div>
  )
}