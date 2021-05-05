// All functions to grab youtube data through API
import axios from 'axios'
import dotenv from 'dotenv'

dotenv.config()

export default axios.create({
    baseURL: "https://youtube.googleapis.com/youtube/v3",
    params: {
        key: process.env.REACT_APP_YOUTUBE_API_KEY,
        part: 'snippet,contentDetails'
    },
})