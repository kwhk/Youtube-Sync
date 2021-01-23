import React from 'react';
import logo from './logo.svg';
import './App.css';
import './components/YoutubePlayer';
import Room from './components/Room';
import Sandbox from './components/Sandbox';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom"

function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Player</Link>
            </li>
            <li>
              <Link to="/sandbox">Sandbox</Link>
            </li>
          </ul>
        </nav>

        <Switch>
          <Route path="/sandbox">
            <Sandbox/>
          </Route>
          <Route path="/">
            <Room />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
