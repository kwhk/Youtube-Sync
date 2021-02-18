import { createSlice } from '@reduxjs/toolkit'

export const videoQueueSlice = createSlice({
    name: 'videoQueue',
    initialState: {
        queue: [],
        currPlaying: -1
    },
    // implement reducers to modify FIFO queue
    reducers: {
        // add video to end of queue
        push: (state, action) => {
            state.queue.push(action.payload)
        },
        remove: (state, action) => {
            state.queue = state.queue.filter((_, index) => index === action.payload ? false : true)
        },
        setActive: (state, action) => {
            state.queue = state.queue.map((video, index) => index === action.payload ? {...video, active: true} : video)
        },
        empty: state => {
            state.queue.length = 0
        }
    }
})

export const selectVideoQueue = state => state.videoQueue
export const { push, remove, empty, setActive } = videoQueueSlice.actions
export default videoQueueSlice.reducer