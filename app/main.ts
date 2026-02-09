const socket = new WebSocket('ws://localhost:3000/ws');

socket.onopen = () => {
    console.log('Success!');
    socket.send('client!!! xD')
}

socket.onclose = (e) => {
    console.log('Socket closed connection', e);
}

socket.onerror = (err) => {
    console.log('Socket err', err);
}