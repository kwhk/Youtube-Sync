import { createSlice } from '@reduxjs/toolkit'

export const playerClockSlice = createSlice({
    name: 'playerClock',
    initialState: {
        start: new Date().getTime(),
        progress: 0,
        stop: true 
    },
    reducers: {
        pauseClock: (state, action) => {
            const now = action.payload
            state.progress = now - (state.start + state.progress)
            state.start = now
            state.stop = true
        },
        playClock: (state, action) => {
            const now = action.payload
            state.start = now
            state.stop = false
        },
        seekToClock: (state, action) => {
            const {ms, now} = action.payload
            state.progress = ms
            // console.log(`Time seeked to ${ms / 1000} seconds`);
            state.start = now
        },
        resetClock: (state, action) => {
            const now = action.payload
            state.start = now
            state.progress = 0
            state.stop = true
        }
    }
})

export const selectPlayerClock = state => state.playerClock
export const { pauseClock, playClock, seekToClock, resetClock } = playerClockSlice.actions
export default playerClockSlice.reducer