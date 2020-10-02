var sockets = {}
sockets.init = function(server) {
  const io = require('socket.io')(server);
  io.on('connection', (socket) => {
    console.log('a user connected');
    socket.on('disconnect', () => {
      console.log('user disconnected');
    })
  });
}

export default sockets