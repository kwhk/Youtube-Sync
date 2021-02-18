export default class Socket {
    constructor() {
        this.socket = null;
        this.clientID = null;
        this.roomID = null;
    }

    connect() {
        let self = this;

        return new Promise(function(resolve, reject) {
            let server = new WebSocket('ws://localhost:8000/api/ws');
            server.onclose = function(event) {
                let e = JSON.stringify(event, ["message", "arguments", "type", "name"]);
                console.log('Connection closed: ' + e);
            }
            
            server.onerror = function(err) {
                let e = JSON.stringify(err, ["message", "arguments", "type", "name"]);
                console.log('Connection error: ' + e);
                reject(err);
            }
            server.onopen = function() {
                self.socket = server;
                resolve(server)
            }
        });
    }

    disconnect() {
        console.log('I want to disconnect!')
        let self = this;
        return new Promise(function(resolve, reject) {
            self.server.onclose = function(event) {
                let e = JSON.stringify(event, ["message", "arguments", "type", "name"]);
                console.log('Connection closed: ' + e);
                resolve();
            }

            self.server.onerror = function(err) {
                let e = JSON.stringify(err, ["message", "arguments", "type", "name"]);
                console.log('Connection closed: ' + e);
                reject(err);
            }
        })
    }

    on(eventName, cb) {
        this.socket.addEventListener("message", res => {
            let data = JSON.parse(res.data);
            if (data.action === "event" && data.event.name && (data.event.name === eventName)) {
                cb(data.event.data);
            }
        })

        return this;
    }
    
    // data and callback parameters are optional
    // emit sends to all clients except sender
    emit(eventName, data, cb) {
        let target = null

        if (this._includeSender != null) {
            target = {roomID: this.roomID, includeSender: this._includeSender}
        }

        let obj = {action: "event", sourceClientID: this.clientID, target, event: {name: eventName}};

        if (data != null) {
            obj.event.data = data;
        }

        this.socket.send(JSON.stringify(obj));

        this.reset();

        if (cb && typeof(cb) === "function") cb();
    }

    reset() {
        this._includeSender = null;
    }

    // Sends to all clients in room including sender
    in() {
        this._includeSender = true;
        return this;
    }

    // Sends to all clients in room except sender
    to() {
        this._includeSender = false;
        return this;
    }
}