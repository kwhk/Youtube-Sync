import React, { Component, useState, useEffect } from 'react';
import socketIo from "socket.io-client";

import Youtube from 'react-youtube';
import './YoutubePlayer.css';

export default class YoutubePlayer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      player: null,
      currTime: 0,
      socket: null
    }

    this.opts = {
      height: '390',
      width: '640',
      playerVars: {
        // https://developers.google.com/youtube/player_parameters
        controls: 0,
        disablekb: 0,
        modestbranding: 1
      },
    }
  }

  componentDidMount() {
    this.setState({socket: socketIo(process.env.REACT_APP_API_HOST)}, () => {
      console.log('socket connected');
    });
  }

  componentWillUnmount() {
    this.state.socket.disconnect();
  }
  
  play = () => {
    console.log('playing');
    this.state.player.playVideo();
  }
  
  pause = () => {
    console.log('pause');
    this.state.player.pauseVideo();
  }

  onReady = (e) => {
    console.log('onReady...');
    this.setState({player: e.target});
  }

  changeTime = (e) => {
    this.setState({currTime: e.target.value});
  }

  seek = (e) => {
    console.log(e.target.value);
  }

  render() {
    return (
      <div>
        <Youtube videoId="D1PvIWdJ8xo" opts={this.opts} onReady={this.onReady}/>
        <div className="d-flex flex-row justify-content-center width-100" style={{marginTop: "5px"}}>
          <div>
            <button onClick={this.play}>play</button>
            <button onClick={this.pause}>pause</button>
          </div>
          <div>
            <input id="audio-progress-bar" type="range" name="video-seek" min="0" max="100" value={this.state.currTime} onInput={this.changeTime} onMouseUp={this.seek} step="0.5"/>
          </div>
        </div>
      </div>
    )
  }


}
