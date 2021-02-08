import React, { Component } from 'react';
import Youtube from 'react-youtube';
import Ping from '../api/ping';
import VideoTimer from '../features/timer/Timer';
import Synchronizer from '../features/sync/Synchronizer';
import SocketContext from '../context/socket'

import './YoutubePlayer.css';
import './flex.css';

export default class YoutubePlayer extends Component {
  static contextType = SocketContext;
  constructor(props) {
    super(props);
    this.state = {
      player: null,
      currVideoPercent: 0,
      roomID: null,
      // client ID for socket connections
      clientID: null,
      isPlaying: false,
      elapsed: 0,
    }

    this.timer = new VideoTimer(new Date(), 0)
    this.sync = null;
    this.opts = {
      height: '360',
      width: '640',
      playerVars: {
        // https://developers.google.com/youtube/player_parameters
        controls: 1,
        disablekb: 0,
        modestbranding: 1,
        playsinline: 1,
        mute: 1,
        enablejsapi: 1,
        cc_load_policy: 0,
        start: 0
      },
    }
  }

  componentDidMount() {
    let socket = this.context;
    socket.on('join', data => {
      console.log(data)
      this.setState({roomID: data.roomID, clientID: data.clientID, isPlaying: data.videoIsPlaying, elapsed: data.videoElapsed});
      socket.clientID = data.clientID;
      new Ping(socket);
    })

    socket.on('seekTo', ms => {
      this.seekTo(ms);
    });

    socket.on('play', ms => {
      this.playVideo();
      this.seekTo(ms);
    })

    socket.on('pause', ms => {
      this.pauseVideo();
      this.seekTo(ms);
    })
  }

  onVideoReady = (e) => {
    this.setState({player: e.target});
    this.sync = new Synchronizer(this.seekTo, e.target, this.timer);

    if (this.state.isPlaying) {
      this.playVideo();
      this.seekTo(this.state.elapsed);
    }
  }
  
  playVideoEmit = () => {
    let socket = this.context;
    if (this.playVideo()) {
      let currTimeMs = Math.floor(this.state.player.getCurrentTime() * 1000);
      this.timer.seekTo(currTimeMs);
      socket.in(this.state.roomID).emit('play', currTimeMs);
    }
  }

  playVideo = () => {
    let player = this.state.player;

    if (player) {
      console.log('PLAY')
      player.playVideo();
      player.unMute();
      this.timer.play();
      this.sync.start();
      this.lockProgressBar();
      return true;
    }

    return false;
  }

  pauseVideoEmit = () => {
    let socket = this.context;

    if (this.pauseVideo()) {
      let currTimeMs = Math.floor(this.state.player.getCurrentTime() * 1000);
      this.timer.seekTo(currTimeMs);
      socket.in(this.state.roomID).emit('pause', currTimeMs);
    }
  }
  
  pauseVideo = () => {
    let player = this.state.player;
    
    if (player) {
      console.log('PAUSE')
      player.pauseVideo();
      this.timer.pause();
      this.sync.stop();
      this.unlockProgressBar();
      return true;
    }

    return false;
  }

  lockProgressBar = () => {
    if (this.calcProgressInterval == null) {
      this.calcProgressInterval = setInterval(() => this.calculateProgress(), 100);
    }
  }

  unlockProgressBar = () => {
    clearInterval(this.calcProgressInterval);
    this.calcProgressInterval = null;
  }
  changeProgress = (e) => {
    if (this.state.player) {
      this.setState({currVideoPercent: e.target.value});
    }
  }

  percentToMs = (percent) => {
    return (percent / 100) * this.state.player.getDuration() * 1000;
  }

  seekTo = (ms) => {
    if (this.state.player) {
      console.log(`Video seeked to ${ms / 1000} sec`)
      this.state.player.seekTo(ms / 1000);
      this.timer.seekTo(ms);
      return true;
    }
    return false;
  }

  seekToEmit = (ms) => {
    let socket = this.context; 

    if (this.seekTo(ms)) {
      socket.in(this.state.roomID).emit('seekTo', ms);
    }
  }

  calculateProgress = () => {
    if (this.state.player) {
      let percent = this.state.player.getCurrentTime() / this.state.player.getDuration() * 100;
      this.setState({ currVideoPercent: Math.round(10 * percent) / 10 });
    }
  }

  onStateChange = (e) => {
    let state = e.data;
    switch (state) {
      case Youtube.PlayerState.PLAYING:
        this.setState({isPlaying: true});
        break;
      case Youtube.PlayerState.PAUSED:
        this.setState({isPlaying: false});
        break
      default:
        break
    }
  }

  render() {
    return (
      <div className="flex flex-col">
        <Youtube videoId={this.props.currPlaying} opts={this.opts} onReady={this.onVideoReady} onStateChange={this.onStateChange}/>
        <div className="d-flex flex-row justify-content-space-evenly" style={{marginTop: "5px"}}>
          { this.state.isPlaying ?
            <button onClick={this.pauseVideoEmit}>pause</button>
            :
            <button onClick={this.playVideoEmit}>playy</button>
          }
          <div>
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoPercent} onChange={this.changeProgress} onMouseUp={(e) => {this.seekToEmit(this.percentToMs(e.target.value)); this.lockProgressBar();}} onMouseDown={this.unlockProgressBar} step="0.1"/>
          </div>
        </div>
      </div>
    )
  }
}