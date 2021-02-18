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
        thumbnail: data.snippet.thumbnails.standard,
        duration: moment.duration(data.contentDetails.duration).asMilliseconds(),
        channelTitle: data.snippet.channelTitle,
        url: url
    }

    return videoInfo
}