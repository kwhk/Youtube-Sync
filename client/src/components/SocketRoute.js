import React, { useContext } from 'react';
import { Route } from 'react-router-dom';
import SocketContext from '../context/socket'
import ConnectToSocket from '../components/ConnectToSocket'

const SocketRoute = ({component: Component, ...rest}) => {
    const { socket } = useContext(SocketContext)

    return (
        <Route {...rest} render={props => (
            socket != null ?
                <Component {...props} />
            : <ConnectToSocket {...props} {...rest}/>
        )} />
    );
};

export default SocketRoute;