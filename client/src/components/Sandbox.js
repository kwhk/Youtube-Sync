import React, { Component } from 'react';
import Ping from '../api/ping';

export default class Sandbox extends Component {
    constructor(props) {
        super(props);
        this.state = {
            ping: null
        }
    }

    componentDidMount() {
    }

    
    render() {
        return (
            <div>
            </div>
        )
    }
}