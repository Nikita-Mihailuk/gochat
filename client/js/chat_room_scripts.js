let socket = null;

async function joinRoom(roomID) {
    try {
      localStorage.setItem('currentRoomID', roomID);
      await setupWebSocket(roomID);
  
      const response = await fetchWithAuth(`${API_V1}rooms/${roomID}`);
      const messages = await response.json();
      const messagesDiv = document.getElementById('messages');
      if (!messagesDiv) {
        console.error("#messages не найдено в DOM.");
        return;
      }
  
      messagesDiv.innerHTML = '';
      messages.forEach(msg => displayMessage(msg));
  
      messagesDiv.scrollTop = messagesDiv.scrollHeight;
    } catch (err) {
      console.error('Ошибка при подключении к комнате:', err);
    }
}
   
function goBackToHome() {
  localStorage.clear();
  window.location.href = "home.html";
}
  
async function setupWebSocket(roomID) {
    const response = await fetchWithAuth(`${API_V1}ws/`);
    const data = await response.json();
    const userID = data.user_id;
    localStorage.setItem("userID", userID)
  
    if (socket) socket.close();
    socket = new WebSocket(`ws://localhost:8080/api/v1/ws/${roomID}?user_id=${userID}`);
  
    socket.onmessage = function(event) {
      const data = JSON.parse(event.data);
      if (data.type === "participants") {
        updateParticipantsList(data.participants);
      } else if (data.type === "message") {
        displayMessage(data.message);
      } else if (data.type === "notification") {
        displayNotification(data.message);
      }
    };
  
    socket.onerror = error => console.error('Ошибка WebSocket:', error);
  
    socket.onclose = event => {
      console.warn('WebSocket закрыт:', event.reason || 'Соединение закрыто');
      socket = null;
    };
  
}
  

  
function displayNotification(message) {
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message');
    messageDiv.innerHTML = `<div class="notification">${message}</div>`;
    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}
  
function displayMessage(message) {
    console.log(message)
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message');
  
  
    if (message.user_id == localStorage.getItem('userID')) {
      messageDiv.classList.add('sender');
      messageDiv.innerHTML = `<div class="message-content">${message.content}</div>`;
    } else {
      let photo_url = "assets/default.jpg";
      if (message.user_name === ""){

      }
      if (message.user_avatar !== ""){
        photo_url = `${API_BASE}${message.user_avatar}`;
      }
      messageDiv.classList.add('other');
      messageDiv.innerHTML = `
        <div class="message-avatar">
          <img src="${photo_url}" alt="Avatar">
        </div>
        <div class="message-content">
          <strong>${message.user_name}</strong> ${message.content}
        </div>
      `;
    }
  
    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}
  
function updateParticipantsList(participants) {
    const participantsList = document.getElementById("participants-list");
    participantsList.innerHTML = "";
  
    participants.forEach(user => {
      const participantItem = document.createElement("li");
      let photo_url = "assets/default.jpg";
      if (user.photo_url !== ""){
        photo_url = `${API_BASE}${user.photo_url}`;
      }
      participantItem.className = "participant";
      participantItem.innerHTML = `
        <div class="message-avatar">
          <img src="${photo_url}" alt="Avatar">
        </div>
        <span>${user.name}</span>
      `;
      participantsList.appendChild(participantItem);
    });
}
  
function sendMessage() {
    const messageInput = document.getElementById("new-message");
    const messageContent = messageInput.value.trim();
  
    if (!messageContent) return;
  
    const currentUserId = localStorage.getItem('userID');
  
    const message = {
      user_id: parseInt(currentUserId),
      content: messageContent,
    };
  
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
      messageInput.value = '';
    } else {
      console.error('WebSocket не подключен.');
    }
}
  