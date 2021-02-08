import { configureStore } from '@reduxjs/toolkit'
import videoQueueReducer from './features/videoQueue/videoQueueSlice';

export default configureStore({
    reducer: {
        videoQueue: videoQueueReducer
    }
})