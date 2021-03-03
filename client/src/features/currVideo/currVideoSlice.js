import { createSlice } from '@reduxjs/toolkit'

export const currVideoSlice = createSlice({
    name: 'currVideo',
    initialState: {
        isPlaying: false,
        elapsed: 0,
        url: ""
    },
    // implement reducers to modify FIFO queue
    reducers: {
        setCurrVideoElapsed: (state, action) => {
            state.elapsed = action.payload
        },
        setCurrVideoUrl: (state, action) => {
            state.url = action.payload
        },
        setCurrVideoPlaybackStatus: (state, action) => {
            state.isPlaying = action.payload
        },
        setCurrVideo: (state, action) => {
            const video = action.payload
            state.isPlaying = video.isPlaying
            state.elapsed = video.elapsed
            state.url = video.url
        }
    }
})

export const selectCurrVideo = state => state.currVideo
export const { setCurrVideoPlaybackStatus, setCurrVideoElapsed, setCurrVideoUrl, setCurrVideo} = currVideoSlice.actions
export default currVideoSlice.reducer