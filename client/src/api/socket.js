export default class Socket {
    constructor() {
        this.socket = new WebSocket("ws://localhost:8000/api/ws");
        this.clientID = null;
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
        let target = {includeSender: false}

        if (this._roomID) {
            target.roomID = this._roomID;
        } else if (this._userID) {
            target.userID = this._userID;
        }

        let obj = {action: "event", sourceClientID: this.clientID, target, event: {name: eventName}};


        if (data) {
            obj.event.data = data;
        }

        this.socket.send(JSON.stringify(obj));

        this.reset();

        if (cb && typeof(cb) === "function") cb();
    }

    reset() {
        this._roomID = null;
        this._userID = null;
    }

    // broadcasts sends to all clients including sender
    broadcast(eventName, data, cb) {
        let target = {includeSender: true}
        
        if (this._roomID) {
            target.roomID = this._roomID;
        } else if (this._userID) {
            target.userID = this._userID;
        }
        
        let obj = {action: "event", sourceClientID: this.clientID, target, event: {name: eventName}};

        if (data) {
            obj.event.data = data;
        }

        this.socket.send(JSON.stringify(obj));

        this.reset();

        if (cb && typeof(cb) === "function") cb();
    }

    // Sends to RoomID
    in(roomID) {
        this._roomID = roomID;
        return this;
    }

    to(userID) {
        this._userID = userID;
        return this;
    }
}