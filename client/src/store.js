import { configureStore } from '@reduxjs/toolkit'
import videoQueueReducer from './features/videoQueue/videoQueueSlice'
import currVideoReducer from './features/currVideo/currVideoSlice'
import connectedUsersReducer from './features/connectedUsers/connectedUsersSlice'
import playerSizeReducer from './features/player/playerSizeSlice'

export default configureStore({
    reducer: {
        videoQueue: videoQueueReducer,
        currVideo: currVideoReducer,
        connectedUsers: connectedUsersReducer,
        playerSize: playerSizeReducer
    }
})