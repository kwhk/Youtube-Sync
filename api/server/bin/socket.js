var sockets = {}
sockets.init = function(server) {
  const io = require('socket.io')(server);
  io.on('connection', (socket) => {
    console.log('a user connected');

    socket.on('disconnect', () => {
      console.log('user disconnected');
    });

    socket.on('seekTo', sec => {
      socket.broadcast.emit('seekTo', sec);
    });

    socket.on('pause', () => {
      console.log('Broadcast: pause');
      socket.broadcast.emit('pause');
    });

    socket.on('play', () => {
      console.log('Broadcast: play');
      socket.broadcast.emit('play');
    });
  });
}

export default sockets