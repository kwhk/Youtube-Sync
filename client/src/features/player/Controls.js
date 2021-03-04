import React, { useState, useEffect, useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux'
import socketContext from '../../context/socket'
import { selectCurrVideo, setCurrVideoPlaybackStatus } from '../currVideo/currVideoSlice'
import VideoTimer from '../timer/Timer';
import Synchronizer from '../sync/Synchronizer';
import './YoutubePlayer.css';
import ProgressBar from './ProgressBar'

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
        console.log(`Video seeked to ${ms / 1000} sec`)
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
        <div className="d-flex flex-row justify-content-space-evenly" style={{marginTop: "5px"}}>
            { isPlaying ?
                <button onClick={pauseVideoEmit}>pause</button>
                :
                <button onClick={playVideoEmit}>playy</button>
            }
            <ProgressBar player={props.player} isPlaying={isPlaying} seekToEmit={seekToEmit}/>
        </div>
    )
}