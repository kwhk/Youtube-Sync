import React, { Component } from 'react';
import YoutubePlayer from './YoutubePlayer';
import Socket from '../api/socket';

export default class Room extends Component {
  constructor(props) {
    super(props);
    this.state = {
      socket: null,
      currPlayingUrl: null
    }
  }

  componentDidMount() {
    this.setState({socket: new Socket()});
  }

  componentWillUnmount() {
    if (this.state.socket) {
      this.state.socket.close();
    }
  }

  render() {
    // if socket is open and ready to communicate
    if (this.state.socket) {
      return (
        <YoutubePlayer currPlaying="r00ikilDxW4" socket={this.state.socket}/>
      )
    } else {
      return (
        <div>Not connected</div>
      )
    }
  }
}