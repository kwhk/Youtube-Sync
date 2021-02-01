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
    <Room />
  );
}

export default App;
