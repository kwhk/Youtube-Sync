import React, { useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { setCurrVideoElapsed, selectCurrVideo } from '../currVideo/currVideoSlice'

export default function ProgressBar(props) {
    const { duration } = useSelector(selectCurrVideo)
    const [currVideoPercent, setCurrVideoPercent] = useState(.0)
    const [calcProgressInterval, setCalcProgressInterval] = useState(null)
    const dispatch = useDispatch()

    const lockProgressBar = () => {
        if (calcProgressInterval == null) {
            setCalcProgressInterval(setInterval(() => calculateProgress(), 100))
        }
    }
    
    const unlockProgressBar = () => {
        clearInterval(calcProgressInterval);
        setCalcProgressInterval(null);
    }
    
    const changeProgress = (e) => {
        console.log(props.player.getCurrentTime())
        const percent = parseFloat(e.target.value)
        setCurrVideoPercent(percent)
        const elapsed = isNaN(percent) ? 0 : Math.round(percent / 100 * duration)
        dispatch(setCurrVideoElapsed(elapsed))
    }
    
    const percentToMs = (percent) => {
        return (percent / 100) * props.player.getDuration() * 1000;
    }

    const calculateProgress = () => {
        const val = props.player.getCurrentTime() / props.player.getDuration() * 100;
        const percent = Math.round(10 * val) / 10
        setCurrVideoPercent(percent)
        const elapsed = isNaN(percent) ? 0 : Math.round(percent / 100 * duration)
        dispatch(setCurrVideoElapsed(elapsed))
    }
    
    useEffect(() => {
        if (props.isPlaying === true) {
            lockProgressBar()
        } else {
            unlockProgressBar()
        }
    }, [props.isPlaying])
    
    return (
        <input className="slider w-full z-10" type="range" name="video-seek" min="0" max="100" value={currVideoPercent} onChange={changeProgress} onMouseUp={(e) => {props.seekToEmit(percentToMs(e.target.value)); lockProgressBar();}} onMouseDown={unlockProgressBar} step="0.1"/>
    )
}