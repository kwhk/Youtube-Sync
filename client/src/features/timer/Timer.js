export default class VideoTimer {
    constructor(start, progress) {
        this.start = start;

        // Time elapsed in ms
        this.progress = progress;
        this.stop = true;
        this.interval = null;
    }

    pause() {
        let now = new Date();
        this.progress = now.getTime() - this.start.getTime() + this.progress;
        this.start = now;
        this.stop = true;
        
        // clearInterval(this.interval)
        // this.interval = null;
        return this
    }

    play() {
        this.start = new Date();
        this.stop = false;
        // if (this.interval == null) {
        //     this.interval = setInterval(() => {console.log(self.elapsed())}, 1000)
        // }
        return this
    }

    // Returns time elapsed since video started in ms
    elapsed() {
        let now = new Date();

        if (this.stop) {
            return this.progress
        } else {
            let elapsed = this.progress + (now.getTime() - this.start.getTime());
            return elapsed;
        }
    }

    seekTo(ms) {
        this.progress = ms;
        // console.log(`Time seeked to ${this.progress / 1000} seconds`);
        this.start = new Date();
        return this;
    }
}