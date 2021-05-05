import axios from 'axios'

const getRooms = async () => {
    const res = await axios.get('/api/db/rooms')
    return res.data
}

export { getRooms }