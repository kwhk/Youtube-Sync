export default class Socket {
    constructor() {
        this.socket = null;
        this.events = {};
    }

    connect() {
        let self = this;

        return new Promise(function(resolve, reject) {
            // Uncomment line below when deploying to kubernetes
            let server = new WebSocket(`ws://${window.location.hostname}/api/ws`);
            // let server = new WebSocket(`ws://localhost:8000/api/ws`)
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
        // Avoid creating duplicate event listeners.
        if (this.events[eventName] != 1) {
            this.events[eventName] = 1
            this.socket.addEventListener("message", res => {
                if (!res.data) {
                    return
                }
                let data = JSON.parse(res.data);
                if (data.action && (data.action === eventName)) {
                    cb(data.data);
                }

            })
        }
        return this;
    }
    
    // emit sends messages to server only
    emit(eventName, data, cb) {
        let obj = {action: eventName, target: null}
        if (data != null) {
            obj.data = data
        }
        this.socket.send(JSON.stringify(obj))
        if (cb && typeof(cb) === "function") cb()
    }

    // broadcast sends messages to room only
    broadcast(eventName, data, cb) {
        let obj = {action: eventName, target: null}
        if (data != null) {
            obj.data = data
        }
        this.socket.send(JSON.stringify(obj))
        if (cb && typeof(cb) === "function") cb()
    }
}