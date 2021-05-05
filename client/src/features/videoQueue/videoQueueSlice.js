import { createSlice } from '@reduxjs/toolkit'
import { getYoutubeVideo } from './utils'

export const videoQueueSlice = createSlice({
    name: 'videoQueue',
    initialState: {
        queue: [],
        currPlayingIndex: -1
    },
    reducers: {
        // add video to end of queue
        pushVideo: (state, action) => {
            state.queue.push(action.payload)
        },
        removeVideo: (state, action) => {
            state.queue = state.queue.filter((_, index) => {
                if (index === action.payload) {
                    if (index === state.currPlayingIndex) {
                        // remove currPlayingIndex or else index error will occur
                        state.currPlayingIndex = -1
                    }
                    return false
                }
                return true
            })
        },
        setVideoActive: (state, action) => {
            // if state.currPlayingIndex is -1 then this is the first video in queue
            if (state.currPlayingIndex == -1) {
                state.currPlayingIndex = action.payload
                state.queue[action.payload].active = true
            } else {
                state.queue[state.currPlayingIndex].active = false
                state.queue[action.payload].active = true
                state.currPlayingIndex = action.payload
            }
        },
        setVideoQueue: (state, action) => {
            state.queue = action.payload
        },
        emptyVideos: state => {
            state.queue.length = 0
        }
    }
})

export function pushVideo(url) {
    return async function pushVideoThunk(dispatch, getState) {
        let videoInfo = await getYoutubeVideo(url)
        videoInfo.active = false
        dispatch(videoQueueSlice.actions.pushVideo(videoInfo))
    }
}

export function setVideoQueue(queue, currPlayingIndex) {
    return async function setVideoQueueThunk(dispatch, getState) {
        let newQueue = []
        for (const video of queue) {
            let videoInfo = await getYoutubeVideo(video.url)
            videoInfo.active = false
            newQueue.push(videoInfo)
        }
        dispatch(videoQueueSlice.actions.setVideoQueue(newQueue))
        
        // If index is -1, this means there are no videos in queue
        // so not possible to have current video playing.
        if (currPlayingIndex != -1) {
            dispatch(videoQueueSlice.actions.setVideoActive(currPlayingIndex))
        }
    }
}

export const selectVideoQueue = state => state.videoQueue
export const { removeVideo, emptyVideos, setVideoActive } = videoQueueSlice.actions
export default videoQueueSlice.reducer