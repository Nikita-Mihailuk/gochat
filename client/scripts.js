const API_BASE = "http://localhost:8080/";
const API_V1 = `${API_BASE}api/v1/`
let socket = null;

document.addEventListener("DOMContentLoaded", () => {
  const currentPath = window.location.pathname;
  const userID = localStorage.getItem("currentUserID");

  if (currentPath.endsWith("index.html")) {
    if (userID) {
      window.location.href = "home.html";
    } else {
      preventBackNavigation();
    }
  } else if (currentPath.endsWith("home.html")) {
    if (!userID) {
      window.location.href = "index.html";
    } else {
      preventBackNavigation();
    }
  } else if (currentPath.endsWith("chat-room.html")) {
    if (!userID) {
      window.location.href = "index.html";
    } else {
      preventBackNavigation();
    }
  }

  if (currentPath.endsWith("home.html") && userID) {
    loadRooms();
    loadProfile();
  }

  if (currentPath.endsWith("chat-room.html")) {
    const roomID = localStorage.getItem('currentRoomID');
    if (roomID) {
      joinRoom(roomID);
    } else {
      console.error("Room ID не найден в localStorage.");
    }
  }
});

function preventBackNavigation() {
  history.replaceState(null, null, window.location.href);
  window.addEventListener("popstate", function () {
    history.pushState(null, null, window.location.href);
  });
}

function logout() {
  localStorage.clear();
  window.location.href = "index.html";
}

function goBackToHome() {
  window.location.href = "home.html";
}

function showTab(tabId) {
  document.querySelectorAll(".tab-content").forEach(tab => tab.style.display = "none");
  document.getElementById(tabId).style.display = "block";
}

async function createRoom() {
  const roomName = document.getElementById("new-room").value;
  if (roomName.length > 30){
    alert("Превышено допустимое количество символов")
    return
  }
  if (roomName) {
    try {
      const response = await fetch(`${API_V1}rooms/`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: roomName }),
      });
      if (response.ok) {
        loadRooms();
      } else {
        handleError(response);
      }
    } catch (err) {
      console.error("Ошибка создания комнаты:", err);
    }
  }
}

async function loadRooms() {
  try {
    const response = await fetch(`${API_V1}rooms/`);
    if (response.ok) {
      const rooms = await response.json();
      const roomList = document.getElementById("room-list");
      roomList.innerHTML = "";
      rooms.forEach(room => {
        const roomItem = document.createElement("li");
        roomItem.innerHTML = `
          <div class="room-group">
            <div class="room-name">
              <span>${room.name}</span>
            </div>
            <button onclick="goToRoom(${room.id})" class="join-room-button">Войти</button>
          </div>
        `;
        roomList.appendChild(roomItem);
      });
    } else {
      alert("Ошибка при загрузке комнат.");
    }
  } catch (err) {
    console.error("Ошибка загрузки комнат:", err);
  }
}

function goToRoom(roomID) {
  localStorage.setItem('currentRoomID', roomID);
  window.location.href = "chat-room.html";
}

async function loginUser() {
  const email = document.getElementById("login-email").value;
  const password = document.getElementById("login-password").value;

  try {
    const response = await fetch(`${API_V1}users/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (response.ok) {
      const data = await response.json();
      const {user_id} = data;

      if (user_id) {
        localStorage.setItem('currentUserID', user_id);
        window.location.href = "home.html";
      }
    } else {
      handleError(response);
    }
  } catch (err) {
    console.error("Ошибка авторизации:", err);
  }
}

async function registerUser() {
  const email = document.getElementById("register-email").value;
  const password = document.getElementById("register-password").value;
  const name = document.getElementById("register-name").value;

  if (email.length > 30 && password.length > 30 && name.length > 30){
    alert("Превышено допустимое количество символов")
    return
  }

  if (email && password && name) {
    try {
      const response = await fetch(`${API_V1}users/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password, name }),
      });
      if (response.ok) {
        window.location.href = "index.html";
      } else {
        handleError(response);
      }
    } catch (err) {
      console.error("Ошибка при регистрации:", err);
    }
  } else {
    alert("Заполните все поля");
  }
}

function openUpdateProfileDialog() {
  document.getElementById("update-profile-modal").style.display = "flex";
}

function closeUpdateProfileDialog() {
  document.getElementById("update-profile-modal").style.display = "none";
}

async function saveUpdatedProfile() {
  const newName = document.getElementById("update-name").value;
  const currentPassword = document.getElementById("current-password").value;
  const newPassword = document.getElementById("new-password").value;
  const profilePhoto = document.getElementById("profile-photo-upload").files[0];
  const userID = localStorage.getItem("currentUserID");

  if (!currentPassword) {
    alert("Пожалуйста введите ваш текущий пароль.");
    return;
  }

  const formData = new FormData();
  formData.append("current_password", currentPassword);
  if (newName) formData.append("name", newName);
  if (newPassword) formData.append("new_password", newPassword);
  if (profilePhoto) formData.append("photo", profilePhoto);

  try {
    const response = await fetch(`${API_V1}users/${userID}`, {
      method: "PATCH",
      body: formData,
    });

    if (response.ok) {
      alert("Профиль успешно изменён!");
      closeUpdateProfileDialog();
      loadProfile();
    } else {
      handleError(response);
    }
  } catch (err) {
    console.error("Ошибка при изменении профиля:", err);
  }
}

async function loadProfile() {
  try {
    const userID = localStorage.getItem("currentUserID");
    const response = await fetch(`${API_V1}users/${userID}`);
    if (response.ok) {
      const profile = await response.json();
      let photo_url = "default.jpg";
      if (profile.photo_url !== ""){
        photo_url = `${API_BASE}${profile.photo_url}`;
      }
      document.getElementById("profile-photo").src = photo_url;
      document.getElementById("profile-email").textContent = `Эл.почта: ${profile.email}`;
      document.getElementById("profile-name-display").textContent = `Имя: ${profile.name}`;
    } else {
      alert("Не удалось загрузить профиль.");
    }
  } catch (err) {
    console.error("Ошибка загрузки профиля:", err);
  }
}

async function joinRoom(roomID) {
  try {
    localStorage.setItem('currentRoomID', roomID);
    setupWebSocket(roomID);

    const response = await fetch(`${API_V1}rooms/${roomID}`);
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

function setupWebSocket(roomID) {
  if (socket) socket.close();
  socket = new WebSocket(`ws://localhost:8080/api/v1/ws/${roomID}?user_id=${localStorage.getItem('currentUserID')}`);

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

  socket.onopen = () => console.log(`WebSocket соединение установлено для комнаты ${roomID}`);
}

function handleError(response) {
  response.json().then(error => alert(error.error));
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
  const messagesDiv = document.getElementById('messages');
  const messageDiv = document.createElement('div');
  messageDiv.classList.add('message');

  

  if (message.user_id == localStorage.getItem('currentUserID')) {
    messageDiv.classList.add('sender');
    messageDiv.innerHTML = `<div class="message-content">${message.content}</div>`;
  } else {
    let photo_url = "default.jpg";
    console.log(message.user_avatar)
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
    let photo_url = "default.jpg";
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
  const messageInput = document.getElementById('new-message');
  const messageContent = messageInput.value.trim();

  if (!messageContent) return;

  const currentUserId = localStorage.getItem('currentUserID');

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
