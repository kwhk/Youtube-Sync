import React, { useEffect, useState } from 'react'
import {
	BrowserRouter as Router,
	Switch,
	Route,
	Link
} from 'react-router-dom'
import './App.css';
import SocketContext from './context/socket'
import RoomList from './features/roomList/RoomList'
import Socket from './api/socket'
import Room from './features/room/Room'

function App() {
	const [socket, setSocket] = useState(null)
	
	useEffect(() => {
		async function connect() {
			const s = new Socket()
			await s.connect()
			setSocket(s)
		}

		connect()

		return async function closeSocket() {
			await socket.disconnect()
		}
	}, []);

  if (socket != null) {
    return (
		<SocketContext.Provider value={socket}>
			<Router>
				<Switch>
					<Route exact path="/">
						<RoomList />
					</Route>
					<Route exact path="/room/:id" render={(props) => <Room {...props}/>}/>
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
