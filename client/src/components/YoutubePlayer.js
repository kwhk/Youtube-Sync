import React, { Component } from 'react';
import socketIo from "socket.io-client";

import Youtube from 'react-youtube';
import './YoutubePlayer.css';

export default class YoutubePlayer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      player: null,
      currVideoTime: 0,
      socket: null
    }

    this.opts = {
      height: '720',
      width: '1280',
      playerVars: {
        // https://developers.google.com/youtube/player_parameters
        controls: 0,
        disablekb: 0,
        modestbranding: 1
      },
    }
  }

  componentDidMount() {
    let socket = socketIo(process.env.REACT_APP_API_HOST);
    this.setState({socket}, () => {
      console.log('socket connected');
    });

    socket.on('syncTime', data => {
      console.log(data);
    });
  }

  componentWillUnmount() {
    this.state.socket.disconnect();
  }
  
  playVideo = () => {
    console.log('playing');
    this.state.player.playVideo();
  }
  
  pauseVideo = () => {
    console.log('pause');
    this.state.player.pauseVideo();
  }

  onVideoReady = (e) => {
    console.log('onReady...');
    this.setState({player: e.target});
  }

  adjustProgress = (e) => {
    this.setState({currVideoTime: e.target.value});
  }

  seekToTimestamp = (e) => {
    console.log(e.target.value);
    this.state.socket.emit('syncTime', e.target.value);
  }

  

  render() {
    return (
      <div>
        <Youtube videoId="D1PvIWdJ8xo" opts={this.opts} onReady={this.onVideoReady}/>
        <div className="d-flex flex-row justify-content-center width-100" style={{marginTop: "5px"}}>
          <div>
            <button onClick={this.playVideo}>play</button>
            <button onClick={this.pauseVideo}>pause</button>
          </div>
          <div>
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoTime} onInput={this.adjustProgress} onMouseUp={this.seekToTimestamp} step="0.5"/>
          </div>
        </div>
      </div>
    )
  }


}
