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

    socket.on('seekTo', sec => {
      this.seekToSec(sec);
    });

    socket.on('play', () => {
      console.log('SERVER SAYS PLAY!')
      this.playVideo();
    })

    socket.on('pause', () => {
      console.log('SERVER SAYS PAUSE!');
      this.pauseVideo();
    })
  }

  componentWillUnmount() {
    this.state.socket.disconnect();
  }
  
  playVideoEmit = () => {
    if (this.playVideo()) {
      this.state.socket.emit('play');
    }
  }

  playVideo = () => {
    if (this.state.player) {
      this.state.player.playVideo();
      this.lockProgressBar();
      return true;
    }
    return false;
  }

  pauseVideoEmit = () => {
    if (this.pauseVideo()) {
      this.state.socket.emit('pause');
    }
  }
  
  pauseVideo = () => {
    if (this.state.player) {
      this.state.player.pauseVideo();
      this.unlockProgressBar();
      return true;
    }
    return false;
  }

  lockProgressBar = () => {
    if (!this.calcProgressInterval) {
      this.calcProgressInterval = setInterval(() => this.calculateProgress(), 100);
    }
  }

  unlockProgressBar = () => {
    if (this.calcProgressInterval) {
      clearInterval(this.calcProgressInterval);
      this.calcProgressInterval = null;
    }
  }  

  onVideoReady = (e) => {
    this.setState({player: e.target});
  }

  changeProgress = (e) => {
    if (this.state.player) {
      this.setState({currVideoPercent: e.target.value});
    }
  }

  percentToSec = (percent) => {
    return (percent / 100) * this.state.player.getDuration();
  }

  seekToSec = (sec) => {
    if (this.state.player) {
      this.state.player.seekTo(sec)
      return true;
    }
    return false;
  }

  seekToSecEmit = (sec) => {
    if (this.seekToSec(sec)) {
      this.state.socket.emit('seekTo', sec);
    }
  }

  calculateProgress = () => {
    if (this.state.player) {
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
            <button onClick={this.playVideoEmit}>play</button>
            <button onClick={this.pauseVideoEmit}>pause</button>
          </div>
          <div>
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoPercent} onChange={this.changeProgress} onMouseUp={(e) => {this.seekToSecEmit(this.percentToSec(e.target.value))}} onMouseDown={this.unlockProgressBar} step="0.1"/>
          </div>
        </div>
      </div>
    )
  }


}
