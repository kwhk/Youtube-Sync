import React from 'react'
import { Link } from 'react-router-dom'
import Button from '../../components/Button'
import VideoInput from '../videoQueue/VideoInput'

export default function Navbar() {
    return (
        <nav className="xs:px-2 sm:px-5 md:px-8 lg:px-10 bg-primary">
            <div>
                <div className="relative flex items-center justify-between h-16">
                    <div className="flex items-center justify-center sm:items-stretch sm:justify-start">
                        <a href="/">
                            <div className="flex-shrink-0 flex items-center">
                                <img className="block lg:hidden h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-mark-indigo-500.svg" alt="Workflow"/>
                                <img className="hidden lg:block h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-logo-indigo-500-mark-white-text.svg" alt="Workflow"/>
                            </div>
                        </a>
                    </div>
                    { window.location.pathname.split('/')[1] === "room" ? 
                    <div className="w-1/2"><VideoInput/></div> : null }
                    { window.location.pathname.split('/')[1] !== "room" ? 
                    <div className="flex items-center pr-2 sm:ml-6 sm:pr-0">
                        <Link to="/room/create">
                            <Button bgColor="bg-highlight" bgHoverColor="bg-white" textHoverColor="text-black" iconHoverColor="text-black" text="New Room" icon="add-outline"/>
                        </Link>
                    </div>
                    : <div className="flex items-center pr-2 sm:ml-6 sm:pr-0"></div> }
                </div>
            </div>
        </nav>
    )
}