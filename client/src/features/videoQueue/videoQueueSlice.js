import { createSlice } from '@reduxjs/toolkit'

export const videoQueueSlice = createSlice({
    name: 'videoQueue',
    initialState: {
        queue: [],
        currPlayingIndex: 0
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
            // if currPlayingIndex is -1 then this is the first video in queue
            if (state.currPlayingIndex == -1) {
                state.currPlayingIndex = action.payload
                state.queue[action.payload].active = true
            } else {
                state.queue[state.currPlayingIndex].active = false
                state.queue[action.payload].active = true
                state.currPlayingIndex = action.payload
            }
        },
        emptyVideos: state => {
            state.queue.length = 0
        }
    }
})

export const selectVideoQueue = state => state.videoQueue
export const { pushVideo, removeVideo, emptyVideos, setVideoActive } = videoQueueSlice.actions
export default videoQueueSlice.reducer