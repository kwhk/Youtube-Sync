import { configureStore } from '@reduxjs/toolkit'
import videoQueueReducer from './features/videoQueue/videoQueueSlice'
import currVideoReducer from './features/currVideo/currVideoSlice'

export default configureStore({
    reducer: {
        videoQueue: videoQueueReducer,
        currVideo: currVideoReducer
    }
})