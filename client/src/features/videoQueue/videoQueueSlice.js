import { createSlice } from '@reduxjs/toolkit'

export const videoQueueSlice = createSlice({
    name: 'videoQueue',
    initialState: {
        queue: []
    },
    // implement reducers to modify FIFO queue
    reducers: {
        // add video to end of queue
        push: (state, action) => {
            state.queue.push(action.payload)
        },
        // remove first video from queue
        pop: state => {
            state.queue.shift()
        },
        empty: state => {
            state.queue.length = 0
        }
    }
})

export const selectVideoQueue = state => state.videoQueue.queue
export const { push, pop, empty } = videoQueueSlice.actions
export default videoQueueSlice.reducer