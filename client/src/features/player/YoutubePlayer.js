import React, { useState } from 'react'
import { useSelector } from 'react-redux'
import { selectCurrVideo } from '../currVideo/currVideoSlice'
import Youtube from 'react-youtube'
import Controls from './Controls'

export default function YoutubePlayer() {
	const [player, setPlayer] = useState(null)
	const { url } = useSelector(selectCurrVideo)

	const opts = {
		playerVars: {
			// https://developers.google.com/youtube/player_parameters
			controls: 0,
			disablekb: 0,
			modestbranding: 1,
			playsinline: 1,
			mute: 1,
			enablejsapi: 1,
			cc_load_policy: 0,
			start: 0,
		}
	}

	const onVideoReady = (e) => {
		setPlayer(e.target)
	}
	
	return (
		<div className="flex flex-col rounded-lg overflow-hidden">
			<div id="player-viewer">
				<Youtube videoId={url} opts={opts} onReady={onVideoReady}/>
			</div>
			{ player != null ?
				<Controls player={player}/>
				:
				<h1>Controls disabled because player not ready</h1>
			}
		</div>
	)
}