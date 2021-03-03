// Component that represents a video including it's duration, title and url
import React, { useContext, useEffect } from 'react'
import '../../styles/flex.css'
import socketContext from '../../context/socket'

export default function Room(props) {
    return (
        <div className="d-flex flex-col">
            <div>
                <h3 style={{margin: 0}}>{props.name}</h3>
                <p>{props.id}</p>
            </div>
        </div>
    )
}