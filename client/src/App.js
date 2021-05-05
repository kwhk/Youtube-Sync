import React, { useState } from 'react'
import {
	BrowserRouter as Router,
	Switch,
	Route
} from 'react-router-dom'
import './App.css';
import HomePage from './features/home/HomePage'
import useScript from './hooks/useScript'
import Navbar from './features/navbar/Navbar'
import SocketContext from './context/socket'
import SocketRoute from './components/SocketRoute';
import Room from './features/room/Room';
import CreateRoom from './features/room/CreateRoom';

function App() {
	const [socket, setSocket] = useState(null)
	useScript('https://unpkg.com/ionicons@5.4.0/dist/ionicons.js')

	return (
		<SocketContext.Provider value={{socket, setSocket}}>
			<Router>
				<div className="w-full">
					<Navbar/>
				</div>
				<Switch>
					<Route exact path="/">
						<HomePage />
					</Route>
					<SocketRoute component={CreateRoom} exact path="/room/create"/>
					<SocketRoute component={Room} exact path="/room/:id"/>
				</Switch>
			</Router>
		</SocketContext.Provider>
	);
}

export default App;