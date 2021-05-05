import React from 'react'
import YoutubePlayer from '../player/YoutubePlayer'
import VideoQueue from '../videoQueue/VideoQueue'
import { useSelector } from 'react-redux'
import ConnectedUsers from '../connectedUsers/ConnectedUsers'
import { selectPlayerSize } from '../player/playerSizeSlice'
import useJoinRoom from '../../hooks/useJoinRoom'
import { useParams } from 'react-router'

export default function Room(props) {
	const { theatre } = useSelector(selectPlayerSize)
    const { id } = useParams()
	const render = useJoinRoom(id)

	if (render) {
		if (theatre === 0) {
			return (
				<div className="xs:px-2 sm:px-5 md:px-8 lg:px-10 grid grid-cols-12 gap-10">
					<div className="flex flex-col order-1 xs:col-span-12 sm:col-span-12 md:col-span-8 xl:col-span-8">
						<YoutubePlayer/>
					</div>
					<div className="xs:col-span-12 sm:col-span-6 md:col-span-4 md:order-2 sm:order-3 xl:col-span-4">
						<VideoQueue/>
					</div>
					<div className="xs:col-span-12 md:order-3 sm:order-2 sm:col-span-6">
						<ConnectedUsers/>
					</div>
				</div>
			)
		} else {
			return (
				<div className="grid grid-cols-12 gap-10">
					<div className="xs:col-span-12 sm:col-span-12 md:col-span-12 lg:col-span-10 lg:col-start-2">
						<YoutubePlayer/>
					</div>
					<div className="xs:col-span-12 sm:col-span-4 lg:col-span-3 xl:col-span-2">
						<ConnectedUsers/>
						<VideoQueue/>
					</div>
				</div>
			)
		}
	} else {
		return (
			<div className="text-white">Can't connect to room</div>
		)
	}
}
