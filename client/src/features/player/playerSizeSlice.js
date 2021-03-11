import { createSlice } from '@reduxjs/toolkit'

export const playerSizeSlice = createSlice({
    name: 'playerSize',
    initialState: {
        theatre: 0
    },
    reducers: {
        toggleTheatre: (state, _) => {
            state.theatre = 1 - state.theatre
        }
    }
})

export const selectPlayerSize = state => state.playerSize
export const { toggleTheatre } = playerSizeSlice.actions
export default playerSizeSlice.reducer