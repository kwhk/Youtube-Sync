import React, { useState, useEffect } from 'react'

export default function ProgressBar(props) {
    const [currVideoPercent, setCurrVideoPercent] = useState(.0)
    const [calcProgressInterval, setCalcProgressInterval] = useState(null)

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
        setCurrVideoPercent(e.target.value)
    }
    
    const percentToMs = (percent) => {
        return (percent / 100) * props.player.getDuration() * 1000;
    }

    const calculateProgress = () => {
        let percent = props.player.getCurrentTime() / props.player.getDuration() * 100;
        setCurrVideoPercent(Math.round(10 * percent) / 10)
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