import React, { useContext, useEffect } from 'react'
import SocketContext from '../../context/socket'
import { useHistory, Redirect } from 'react-router-dom'
import Button from '../../components/Button'
import VideoInput from '../videoQueue/VideoInput'
import ThemeToggle from './Toggle'

export default function Navbar() {
    const socket = useContext(SocketContext)
    let history = useHistory()

    useEffect(() => {
        socket.on('create-room', data => {
            console.log(data);
            history.push('/room/' + data.id)
        })
    }, [])

    const createRoom = () => {
        socket.emit('create-room')
    }

    return (
        <nav className="xs:px-2 sm:px-5 md:px-8 lg:px-10 bg-primary">
            <div>
                <div className="relative flex items-center justify-between h-16">
                    {/* <div className="absolute inset-y-0 left-0 flex items-center sm:hidden"> */}
                        {/* <button type="button" className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white" aria-controls="mobile-menu" aria-expanded="false">
                            <span className="sr-only">Open main menu</span>
                            <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
                            </svg>
                            <svg className="hidden h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button> */}
                    {/* </div> */}
                    <div className="flex items-center justify-center sm:items-stretch sm:justify-start">
                        <div className="flex-shrink-0 flex items-center">
                            <img className="block lg:hidden h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-mark-indigo-500.svg" alt="Workflow"/>
                            <img className="hidden lg:block h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-logo-indigo-500-mark-white-text.svg" alt="Workflow"/>
                        </div>
                    </div>
                    { window.location.pathname.split('/')[1] == "room"
                    ?  
                    <div className="w-1/2">
                        <VideoInput/>
                    </div>
                    : null
                    }
                    <div className="flex items-center pr-2 sm:ml-6 sm:pr-0">
                        <Button onClick={createRoom} bgColor="bg-highlight" bgHoverColor="bg-white" textHoverColor="text-black" iconHoverColor="text-black" text="New Room" icon="add-outline"/>
                    </div>
                </div>
            </div>
        </nav>
    )
}