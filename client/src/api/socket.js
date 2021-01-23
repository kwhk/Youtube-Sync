export default class Socket {
    constructor() {
        this.socket = new WebSocket("ws://localhost:8000/api/ws");
        this.clientID = null;
    }

    on(eventName, cb) {
        this.socket.addEventListener("message", res => {
            let data = JSON.parse(res.data);
            if (data.body.eventName && (data.body.eventName === eventName)) {
                cb(data.body.data);
            }
        })

        return this;
    }
    
    // data and callback parameters are optional
    // emit sends to all clients except sender
    emit(eventName, data, cb) {
        let target = {includeSender: false, sourceClientID: this.clientID}

        if (this._targRoomID) {
            target.targRoomID = this._targRoomID;
        } else if (this._targUserID) {
            target.targUserID = this._targUserID;
        }

        let obj = {messageType: "event", target, body: {eventName}};


        if (data) {
            obj.body.data = data;
        }

        this.socket.send(JSON.stringify(obj));

        this._targRoomID = null;
        this._targUserID = null;

        if (cb && typeof(cb) === "function") cb();
    }

    // broadcasts sends to all clients including sender
    broadcast(eventName, data, cb) {
        let target = {includeSender: true}
        
        if (this._targRoomID) {
            target.targRoomID = this._targRoomID;
        } else if (this._targUserID) {
            target.targUserID = this._targUserID;
        }
        
        let obj = {messageType: "event", target, body: {eventName}};

        if (data) {
            obj.body.data = data;
        }

        this.socket.send(JSON.stringify(obj));

        this._targRoomID = null;
        this._targUserID = null;

        if (cb && typeof(cb) === "function") cb();
    }

    // Sends to RoomID
    in(targRoomID) {
        this._targRoomID = targRoomID;
        return this;
    }

    to(targUserID) {
        this._targUserID = targUserID;
        return this;
    }
}