// Component that represents a video including it's duration, title and url
import React from 'react'

export default function Room(props) {
    return (
        <div className="flex flex-col">
            <div>
                <h3 style={{margin: 0}}>{props.name}</h3>
                <p>{props.id}</p>
            </div>
        </div>
    )
}