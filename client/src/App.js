import React, { useEffect, useState } from 'react'
import {
	BrowserRouter as Router,
	Switch,
	Route
} from 'react-router-dom'
import './App.css';
import SocketContext from './context/socket'
import HomePage from './features/home/HomePage'
import Socket from './api/socket'
import Room from './features/room/Room'
import useScript from './utils/useScript'
import Navbar from './features/navbar/Navbar'
import Ping from './api/ping'

function App() {
	const [socket, setSocket] = useState(null)
	// const history = createBrowserHistory()
	
	useScript('https://unpkg.com/ionicons@5.4.0/dist/ionicons.js')
	useEffect(() => {
		async function connect() {
			const s = new Socket()
			await s.connect()
			new Ping(s)
			setSocket(s)
		}

		connect()

		return async function closeSocket() {
			if (socket != null) {
				console.log(socket)
				await socket.disconnect()
			}
		}
	}, []);


	if (socket != null) {
		return (
			<SocketContext.Provider value={socket}>
				<Router>
					<div className="w-full">
						<Navbar/>
					</div>
					<Switch>
						<Route exact path="/">
							<HomePage />
						</Route>
						{/* <Route exact path="/room/:id" render={(props) => <Room {...props}/>}/> */}
						<Route exact path="/room/:id">
							<Room />
						</Route>
					</Switch>
				</Router>
			</SocketContext.Provider>
		);
	} else {
		return (
			<div>Not connected</div>
		)
	}
}

export default App;