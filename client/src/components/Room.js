import React, { Component } from 'react';
import socketIo from 'socket.io-client';
import YoutubePlayer from './YoutubePlayer';

export default class Room extends Component {
  constructor(props) {
    super(props);
    this.state = {
      socket: null,
      currPlayingUrl: null
    }
  }

  componentDidMount() {
    const socket = socketIo(process.env.REACT_APP_API_HOST);

    this.setState({socket}, () => {
      console.log('socket connected');
    });
  }

  componentWillUnmount() {
    this.state.socket.disconnect();
  }

  render() {
    if (this.state.socket) {
      return (
        <YoutubePlayer currPlaying="D1PvIWdJ8xo" socket={this.state.socket}/>
      )
    } else {
      return (
        <div>Not connected</div>
      )
    }
  }
}