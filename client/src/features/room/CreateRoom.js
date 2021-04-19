import { useEffect, useContext } from 'react'
import { useHistory } from 'react-router-dom'
import SocketContext from '../../context/socket'
import withLoading from '../../components/Loading'

export default function CreateRoom(props) {
    // const { setLoading } = props
    const { socket } = useContext(SocketContext)
    let history = useHistory()

    useEffect(() => {
        socket.emit('create-room')
        socket.on('create-room', data => {
            // setLoading(false)
            console.log(data)
            history.push('/room/' + data.id)
        })
    })

    return null
}

// export default withLoading(CreateRoom, 'Creating room...')