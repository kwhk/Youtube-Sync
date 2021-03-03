import { createSlice } from '@reduxjs/toolkit'

export const connectedUsersSlice = createSlice({
    name: 'connectedUsers',
    initialState: {
        users: {}
    },
    // implement reducers to modify FIFO queue
    reducers: {
        // add video to end of queue
        pushUser: (state, action) => {
            state.users[action.payload] = 1
        },
        removeUser: (state, action) => {
            delete state.users[action.payload]
        },
        setConnectedUsers: (state, action) => {
            for (let i = 0; i < action.payload.length; i++) {
                state.users[action.payload[i]] = 1
            }
        }
    }
})

export const selectConnectedUsers = state => state.connectedUsers
export const { pushUser, removeUser, setConnectedUsers } = connectedUsersSlice.actions
export default connectedUsersSlice.reducer