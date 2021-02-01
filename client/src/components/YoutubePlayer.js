import React, { Component } from 'react';
import Youtube from 'react-youtube';
import { poll } from '../api/poll';
import Ping from '../api/ping';

import './YoutubePlayer.css';

export default class YoutubePlayer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      player: null,
      currVideoPercent: 0,
      roomID: null,
      // client ID for socket connections
      clientID: null,
      isPlaying: false
    }

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
        start: 61
      },
    }
  }

  componentDidMount() {
    this.props.socket.on('welcome', data => {
      this.setState({roomID: data.roomID, clientID: data.clientID});
      this.props.socket.clientID = data.clientID;
      new Ping(this.props.socket);
    })

    this.props.socket.on('seekTo', sec => {
      console.log(`SERVER SAYS SEEKTO ${sec} seconds!`);
      this.seekToSec(sec);
    });

    this.props.socket.on('play', () => {
      console.log('SERVER SAYS PLAY!')
      this.playVideo();
    })

    this.props.socket.on('pause', () => {
      console.log('SERVER SAYS PAUSE!');
      this.pauseVideo();
    })
  }
  
  playVideoEmit = () => {
    this.playVideo().then(() => {
      this.props.socket.in(this.state.roomID).emit('play');
    }).catch(err => console.log(err));
  }

  playVideo = () => {
    let self = this;
    let player = this.state.player;

    if (player) {
      player.playVideo();
    }

    return new Promise((resolve, reject) => {
      poll(() => { return player.getPlayerState() === Youtube.PlayerState.PLAYING }, 1500, 1).then(() => {
        player.unMute();
        self.unlockProgressBar();
        return resolve();
      }).catch(err => {console.log(err); return reject('Took too long to play')})
    })
  }

  pauseVideoEmit = () => {
    this.pauseVideo().then(() => {
      this.props.socket.in(this.state.roomID).emit('pause');
    }).catch(err => console.log(err));
  }
  
  pauseVideo = () => {
    let self = this;
    let player = this.state.player;

    if (player) {
      player.pauseVideo();
    }

    return new Promise((resolve, reject) => {
      poll(() => { return player.getPlayerState() === Youtube.PlayerState.PAUSED }, 1500, 1).then(() => {
        self.lockProgressBar();
        return resolve();
      }).catch(err => {console.log(err); return reject('Took too long to pause')})
    })
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
      this.props.socket.in(this.state.roomID).emit('seekTo', sec);
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
    }
  }

  render() {
    return (
      <div className="d-flex flex-col align-items-center">
        <Youtube videoId={this.props.currPlaying} opts={this.opts} onReady={this.onVideoReady} onStateChange={this.onStateChange}/>
        <div className="d-flex flex-row justify-content-center width-100" style={{marginTop: "5px"}}>
          { this.state.isPlaying ?
            <button onClick={this.pauseVideoEmit}>PAUSE</button>
            :
            <button onClick={this.playVideoEmit}>play</button>
          }
          <div>
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currVideoPercent} onChange={this.changeProgress} onMouseUp={(e) => {this.seekToSecEmit(this.percentToSec(e.target.value))}} onMouseDown={this.unlockProgressBar} step="0.1"/>
          </div>
        </div>
      </div>
    )
  }


}
