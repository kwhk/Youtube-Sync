import React, { useState, useEffect, useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux'
import socketContext from '../../context/socket'
import { selectCurrVideo, setCurrVideoPlaybackStatus } from '../currVideo/currVideoSlice'
import VideoTimer from '../timer/Timer';
import Synchronizer from '../sync/Synchronizer';
import './YoutubePlayer.css';
import ProgressBar from './ProgressBar'
import { toggleTheatre } from './playerSizeSlice'

/*

BUGS :

sync isn't great esp when new client joins
clients disconnect from: socket more frequently now (caused by TCP error FIND THE ISSUE)
seek control is buggy

*/

export default function Controls(props) {
    const { isPlaying, elapsed } = useSelector(selectCurrVideo)
    const socket = useContext(socketContext)
    const [timer] = useState(new VideoTimer(new Date(), 0))
    const dispatch = useDispatch()
    
    const seekTo = (ms) => {
        // console.log(`Video seeked to ${ms / 1000} sec`)
        props.player.seekTo(ms / 1000);
        timer.seekTo(ms);
    }

    const [sync] = useState(new Synchronizer(seekTo, props.player, timer))
    
    const playVideoEmit = () => {
        // playVideo()
        let currTimeMs = Math.floor(props.player.getCurrentTime() * 1000);
        timer.seekTo(currTimeMs);
        socket.broadcast('play-video', currTimeMs);
    }

    const playVideo = () => {
        console.log('PLAY')
        props.player.playVideo();
        props.player.unMute();
        timer.play();
        sync.start();
        dispatch(setCurrVideoPlaybackStatus(true))
    }

    const pauseVideoEmit = () => {
        // pauseVideo()
        let currTimeMs = Math.floor(props.player.getCurrentTime() * 1000);
        timer.seekTo(currTimeMs);
        socket.broadcast('pause-video', currTimeMs);
    }
    
    const pauseVideo = () => {
        console.log('PAUSE')
        props.player.pauseVideo();
        timer.pause();
        sync.stop();
        dispatch(setCurrVideoPlaybackStatus(false))
    }

    const seekToEmit = (ms) => {
        // seekTo(ms)
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
                <div className="col-start-8 flex flex-row items-center justify-center">
                    { isPlaying ?
                        <ion-icon name="pause-sharp" class="text-white text-3xl visible cursor-pointer" onClick={pauseVideoEmit}></ion-icon>
                        :
                        <ion-icon name="play-sharp" class="text-white text-3xl visible cursor-pointer" onClick={playVideoEmit}></ion-icon>
                    } 
                    <ion-icon class="text-white ml-3 text-base visible cursor-pointer" name="play-skip-forward-sharp"></ion-icon>
                </div>
                <div className="text-gray-700 flex flex-row col-start-13 justify-center" onClick={toggleTheatreSize}>
                    <ion-icon class="text-2xl visible cursor-pointer" name="tv-outline"></ion-icon>
                </div>
            </div>
        </div>

    )
}