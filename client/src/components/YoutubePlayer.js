import React, { Component } from 'react';
import socketIo from "socket.io-client";
import Youtube from 'react-youtube';

import './YoutubePlayer.css';

export default class YoutubePlayer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      player: null,
      currVideoPercent: 0,
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
      console.log('PLAY');
      this.state.player.playVideo();
      this.lockProgressBar();
    }
  }
  
  pauseVideo = () => {
    if (this.state.player) {
      console.log('PAUSE');
      this.state.player.pauseVideo();
      this.unlockProgressBar();
    }
  }

  lockProgressBar = () => {
    if (!this.calcProgressInterval) {
    console.log('\tlocking progress bar...');
      this.calcProgressInterval = setInterval(() => this.calculateProgress(), 100);
    }
  }

  unlockProgressBar = () => {
    if (this.calcProgressInterval) {
      console.log('\tunlocking progress bar...');
      clearInterval(this.calcProgressInterval);
      this.calcProgressInterval = null;
    }
  }  

  onVideoReady = (e) => {
    this.setState({player: e.target});
  }

  changeProgress = (e) => {
    if (this.state.player) {
      console.log('changing progress...');
      this.setState({currVideoPercent: e.target.value});
    }
  }

  seekToTimestamp = (e) => {
    if (this.state.player) {
      let sec = (e.target.value / 100) * this.state.player.getDuration();
      this.state.socket.emit('syncTime', sec);
      this.state.player.seekTo(sec)
      console.log('\tseeked to ' + sec);
      this.setState({currVideoPercent: e.target.value});
    }
  }

  calculateProgress = () => {
    if (this.state.player) {
      console.log('calculate progress...');
      let percent = this.state.player.getCurrentTime() / this.state.player.getDuration() * 100;
      this.setState({ currVideoPercent: Math.round(10 * percent) / 10 });
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
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoPercent} onChange={this.changeProgress} onMouseUp={(e) => {console.log('MOUSE UP'); this.seekToTimestamp(e)}} onMouseDown={() => {console.log('MOUSE DOWN'); this.unlockProgressBar()}} step="0.1"/>
          </div>
        </div>
      </div>
    )
  }


}
