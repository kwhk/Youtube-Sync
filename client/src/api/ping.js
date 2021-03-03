export default class Ping {
    constructor(socket) {
        this.seq = 0;
        this.socket = socket;

        this.initPing();

        socket.on('user-ping', res => {
            if (res.FIN !== 1) {
                this.ping(res);
            }
        })
    }

    initPing() {
        let data = {
            ACK: 0,
            SYN: 1,
            FIN: 0,
            seq: 0,
            ack: 0
        }
        this.socket.emit('user-ping', data);
    }

    ping(res) {
        this.seq += 1;

        let data = {
            ACK: 1,
            SYN: 0,
            FIN: 0,
            seq: this.seq,
            ack: res.seq + 1
        }

        this.socket.emit('user-ping', data)
    }
}

