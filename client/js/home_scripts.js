function showTab(tabId) {
    document.querySelectorAll(".tab-content").forEach(tab => tab.style.display = "none");
    document.getElementById(tabId).style.display = "block";
}
  
async function createRoom() {
    const roomName = document.getElementById("new-room").value.trim();
    if (roomName.length > 30) {
      alert("Превышено допустимое количество символов");
      return;
    }
    if (roomName) {
      try {
        const response = await fetchWithAuth(`${API_V1}rooms/`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: roomName }),
        });
        if (response.ok) loadRooms();
        else handleError(response);
      } catch (err) {
        console.error("Ошибка создания комнаты:", err);
      }
    }
}
  
async function loadRooms() {
    try {
      const response = await fetchWithAuth(`${API_V1}rooms/`);
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
      const response = await fetchWithAuth(`${API_V1}users/`, {
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
    const response = await fetchWithAuth(`${API_V1}users/`);
    if (response.ok) {
      const profile = await response.json();
      let photo_url = "assets/default.jpg";
      if (profile.photo_url !== "") {
        photo_url = `${API_BASE}${profile.photo_url}`;
      }
      document.getElementById("profile-photo").src = photo_url;
      document.getElementById(
        "profile-email"
      ).textContent = `Эл.почта: ${profile.email}`;
      document.getElementById(
        "profile-name-display"
      ).textContent = `Имя: ${profile.name}`;
    } else {
      alert("Не удалось загрузить профиль.");
    }
  } catch (err) {
    console.error("Ошибка загрузки профиля:", err);
  }
}