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
    if (this.state.player) {
      this.state.player.playVideo();
      this.unlockProgressBar();
    }
  }
  
  pauseVideo = () => {
    if (this.state.player) {
      this.state.player.pauseVideo();
      this.lockProgressBar();
    }
  }

  unlockProgressBar = () => {
    this.calcProgressInterval = setInterval(() => this.calculateProgress(), 100);
  }

  lockProgressBar = () => {
    clearInterval(this.calcProgressInterval);
  }  

  onVideoReady = (e) => {
    this.setState({player: e.target});
  }

  changeProgress = (e) => {
    if (this.state.player) {
      this.setState({currVideoTime: e.target.value});
    }
  }

  seekToTimestamp = (e) => {
    if (this.state.player) {
      let sec = (e.target.value / 100) * this.state.player.getDuration();
      this.state.socket.emit('syncTime', sec);
      this.state.player.seekTo(sec)
      this.unlockProgressBar();
    }
  }

  calculateProgress = () => {
    if (this.state.player) {
      let percent = this.state.player.getCurrentTime() / this.state.player.getDuration() * 100;
      this.setState({ currVideoTime: Math.round(10 * percent) / 10 });
    }
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
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoTime} onChange={this.changeProgress} onMouseUp={this.seekToTimestamp} onMouseDown={this.lockProgressBar} step="0.1"/>
          </div>
        </div>
      </div>
    )
  }


}
