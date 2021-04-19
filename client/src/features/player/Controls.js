import React, { useState, useEffect, useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux'
import socketContext from '../../context/socket'
import { selectCurrVideo, setCurrVideoPlaybackStatus } from '../currVideo/currVideoSlice'
import Synchronizer from '../sync/Synchronizer';
import './YoutubePlayer.css';
import ProgressBar from './ProgressBar'
import { toggleTheatre } from './playerSizeSlice'
import { pauseClock, playClock, seekToClock, resetClock } from '../clock/playerClockSlice'
import { selectVideoQueue } from '../videoQueue/videoQueueSlice'
import { msToMinutesAndSeconds } from '../../utils/convertTime'

/*

BUGS :

sync isn't great esp when new client joins
clients disconnect from: socket more frequently now (caused by TCP error FIND THE ISSUE)
seek control is buggy

*/

export default function Controls(props) {
    const { isPlaying, elapsed, duration } = useSelector(selectCurrVideo)
    const { currPlayingIndex } = useSelector(selectVideoQueue)
    const { socket } = useContext(socketContext)
    const dispatch = useDispatch()
    const seekTo = (ms) => {
        props.player.seekTo(ms / 1000);
        dispatch(seekToClock({ms: ms, now: (new Date()).getTime()}))
    }
    
    const [sync] = useState(new Synchronizer(seekTo, props.player))
   
    // Reset player clock when there is new current video.
    useEffect(() => {
        dispatch(resetClock())
    }, [currPlayingIndex])

    
    const playVideoEmit = () => {
        let currTimeMs = Math.floor(props.player.getCurrentTime() * 1000);
        dispatch(seekToClock({ms: currTimeMs, now: (new Date()).getTime()}))
        socket.broadcast('play-video', currTimeMs);
    }

    const playVideo = () => {
        console.log('PLAY')
        props.player.playVideo();
        props.player.unMute();
        dispatch(playClock((new Date()).getTime()))
        sync.start();
        dispatch(setCurrVideoPlaybackStatus(true))
    }

    const pauseVideoEmit = () => {
        let currTimeMs = Math.floor(props.player.getCurrentTime() * 1000);
        dispatch(seekToClock({ms: currTimeMs, now: (new Date()).getTime()}))
        socket.broadcast('pause-video', currTimeMs);
    }
    
    const pauseVideo = () => {
        console.log('PAUSE')
        props.player.pauseVideo();
        dispatch(pauseClock((new Date()).getTime()))
        sync.stop();
        dispatch(setCurrVideoPlaybackStatus(false))
    }

    const seekToEmit = (ms) => {
        socket.broadcast('seekto-video', ms);
    }
    
    const toggleTheatreSize = () => {
        dispatch(toggleTheatre())
    }
    
    useEffect(() => {
        socket.on('seekto-video', ms => {
            seekTo(ms);
        });

        socket.on('play-video', ms => {
            playVideo();
            seekTo(ms);
        })

        socket.on('pause-video', ms => {
            pauseVideo();
            seekTo(ms);
        })

        if (isPlaying) {
            playVideo()
            seekTo(elapsed)
        }
    }, [])

    return (
        <div className="flex flex-col">
            <ProgressBar player={props.player} isPlaying={isPlaying} seekToEmit={seekToEmit}/>
            <div className="transition-colors duration-500 bg-secondary lg-rounded-b-xl p-3 sm:px-3 lg:px-5 xl:px-8 grid grid-cols-13 items-center">
                <div className="col-start-1 flex flex-row items-center justify-start text-white text-sm">
                    {msToMinutesAndSeconds(elapsed)} / {msToMinutesAndSeconds(duration)}
                </div>
                <div className="col-start-7 flex flex-row items-center justify-center">
                    { isPlaying ?
                        <ion-icon name="pause-sharp" class="text-white text-3xl visible cursor-pointer" onClick={pauseVideoEmit}></ion-icon>
                        :
                        <ion-icon name="play-sharp" class="text-white text-3xl visible cursor-pointer" onClick={playVideoEmit}></ion-icon>
                    } 
                </div>
                <div className="text-gray-700 flex flex-row col-start-13 justify-end" onClick={toggleTheatreSize}>
                    <ion-icon class="text-2xl visible cursor-pointer" name="tv-outline"></ion-icon>
                </div>
            </div>
        </div>

    )
}