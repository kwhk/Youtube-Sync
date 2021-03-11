import YoutubeAPI from '../../api/youtube'
import moment from 'moment'

export async function getYoutubeVideo(url) {
    const res = await YoutubeAPI.get("/videos", {
        params: {
            id: url
        }
    })

    const data = res.data.items[0]

    const videoInfo = {
        title: data.snippet.title,
        thumbnail: data.snippet.thumbnails.medium,
        duration: moment.duration(data.contentDetails.duration).asMilliseconds(),
        channelTitle: data.snippet.channelTitle,
        url: url
    }

    return videoInfo
}

export function youtubeParser(url){
    var regExp = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*/;
    var match = url.match(regExp);
    return (match&&match[7].length==11)? match[7] : false;
}

export function msToTime(duration) {
    var seconds = Math.floor((duration / 1000) % 60),
        minutes = Math.floor((duration / (1000 * 60)) % 60),
        hours = Math.floor((duration / (1000 * 60 * 60)) % 24);
    
    hours = (hours < 10) ? "0" + hours : hours;
    minutes = (minutes < 10) ? "0" + minutes : minutes;
    seconds = (seconds < 10) ? "0" + seconds : seconds;

    if (hours == 0) {
        return minutes + ":" + seconds;
    }
    return hours + ":" + minutes + ":" + seconds;
}