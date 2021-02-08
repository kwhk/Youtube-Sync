import { createSlice } from '@reduxjs/toolkit'

export const videoQueueSlice = createSlice({
    name: 'videoQueue',
    initialState: {
        value: []
    },
    // implement reducers to modify FIFO queue
    reducers: {
        // add video to end of queue
        push: (state, action) => {
            state.value.push(action.payload)
        },
        // remove first video from queue
        pop: state => {
            state.value.shift()
        },
        empty: state => {
            state.value.length = 0
        }
    }
})

export const selectVideoQueue = state => state.videoQueue.value
export const { push, pop, empty } = videoQueueSlice.actions
export default videoQueueSlice.reducer