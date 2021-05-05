import React from 'react'
import { Link } from 'react-router-dom'
import Button from '../../components/Button'

export default function HomePage() {
    return (
        <div className="h-screen w-screen flex flex-col items-center justify-center">
            <h1 className="text-white text-5xl font-bold">Watch Youtube with friends</h1>
            <div className="mt-5">
                <Link to="/room/create">
                    <Button bgColor="bg-highlight" text="Create a room" fontSize="text-xl"/>
                </Link>
            </div>
        </div>
    )
}