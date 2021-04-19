import store from '../../store'

export default class Synchronizer {
    constructor(seekTo, player, offset=0.5) {
        this.offset = offset;
        this.seekTo = seekTo;
        this.player = player;
        this.interval = null;
    }

    start() {
        let self = this;
        let a = function() {
            let elapsed = timerElapsed() / 1000;
            let currTime = self.player.getCurrentTime();
            let diff = Math.abs(elapsed - currTime);
            // console.log(`timer: ${elapsed.toFixed(2)}\nvideo: ${currTime.toFixed(2)}\ndiff:${diff.toFixed(2)}`)
            if (diff > self.offset) {
                // console.log('NOT IN SYNC')
                let elapsed = timerElapsed()
                let str = elapsed / 1000
                // console.log(`seekTo ${str.toFixed(2)}`)
                self.seekTo(elapsed);
            }
        }

        if (this.interval == null) {
            this.interval = setInterval(a, 1000)
        }
    }

    stop() {
        clearInterval(this.interval);
        this.interval = null;
    }
}

function timerElapsed() {
    const state = store.getState()
    const { start, progress, stop } = state.playerClock
    const now = new Date()

    if (stop) {
        return progress
    } else {
        const elapsed = progress + (now.getTime() - start)
        return elapsed
    }
}