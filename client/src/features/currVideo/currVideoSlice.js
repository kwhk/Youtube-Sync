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
        setElapsed: (state, action) => {
            state.elapsed = action.payload
        },
        setUrl: (state, action) => {
            state.url = action.payload
        },
        setPlaybackStatus: (state, action) => {
            state.isPlaying = action.payload
        },
        setVideo: (state, action) => {
            const video = action.payload
            state.isPlaying = video.isPlaying
            state.elapsed = video.elapsed
            state.url = video.url
        }
    }
})

export const selectCurrVideo = state => state.currVideo
export const { setPlaybackStatus, setElapsed, setUrl, setVideo } = currVideoSlice.actions
export default currVideoSlice.reducer