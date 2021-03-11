// Component that represents a video including it's duration, title and url
import React, { useContext } from 'react'
import socketContext from '../../context/socket'
import { msToTime } from './utils'
import Button from '../../components/Button'

export default function Video(props) {
    const socket = useContext(socketContext)

    const handlePlay = () => {
        socket.broadcast('play-video-queue', {url: props.url, index: props.index})
    }

    const handleRemove = () => {
        socket.broadcast('remove-video-queue', {url: props.url, index: props.index})
    }

    return (
        <div className="bg-secondary flex flex-col rounded-xl my-3 p-5">
            { props.active ?
                <div className="flex flex-row items-center mb-2">
                    <div className="h-1 w-1 my-1 rounded-full bg-green-500 z-2 mr-2"></div>
                    <p className="text-gray-300 text-xs font-semibold">NOW PLAYING</p>
                </div>
                :
                null
            }
            <div className="grid grid-cols-5" onClick={handlePlay}>
                <img className="rounded-xl col-span-2 w-full h-auto" src={props.thumbnail.url}/>
                <div className="col-span-3 flex flex-col pl-3">
                    <h3 className="text-white mb-2 xs:text-sm sm:text-sm text-base font-semibold truncate">{props.title}</h3>
                    <p className="text-gray-500 text-sm xs:text-xs sm:text-xs">{props.channelTitle}</p>
                    <p className="text-gray-500 text-sm xs:text-xs sm:text-xs">{msToTime(props.duration)}</p>
                </div>
            </div>
            <div className="flex flex-row mt-5">
                <Button onClick={handleRemove} bgColor="bg-red-600" icon="close-outline" text="Remove"/>
            </div>
        </div>
    )
}