import React, { Component } from 'react';

export default class Sandbox extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: []
        }
    }
    
    componentDidMount() {
        // connect((msg) => {
        //     console.log("New Message");
        //     this.setState({chatHistory: [...this.state.chatHistory, JSON.parse(msg.data)]})
        //     console.log(this.state);
        // });
    }

    send() {
        console.log("hello");
    }

    render() {
        return (
            <div>
                <ul>
                    {this.state.chatHistory.map((msg, index) => (
                        <li key={index}>{msg.body}</li>
                    ))}
                </ul>
                <button onClick={this.send}>Send</button>
            </div>
        )
    }
}