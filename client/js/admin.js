async function fetchAdminData(endpoint, method = "GET", body = null) {
    const options = {
        method,
        headers: { "Content-Type": "application/json" },
    };
    if (body) options.body = JSON.stringify(body);

    const response = await fetchWithAuth(`${API_V1}admins/${endpoint}`, options);
    if (!response.ok) {
        alert(`Ошибка: ${response.statusText}`);
        return null;
    }
    return response;
}

function closeUpdateDialog() {
    document.getElementById("update-profile-modal").style.display = "none";
    document.getElementById("update-room-modal").style.display = "none";
}

function goToHome() {
    window.location.href = "home.html";
}

async function loadUsers() {
    const response = await fetchWithAuth(`${API_V1}admins/users`);
    const users = await response.json();
    if (!users) return;
    
    let content = `<table class="admin-table">
        <tr><th>ID</th><th>Email</th><th>Имя</th><th>Роль</th><th>Фото</th><th>Действия</th></tr>`;
    
    users.forEach(user => {
        let photo_url = "assets/default.jpg";
        if (user.photo_url !== ""){
            photo_url = `${API_BASE}${user.photo_url}`;
        }
        content += `<tr>
            <td>${user.id}</td>
            <td>${user.email}</td>
            <td>${user.name}</td>
            <td>${user.role}</td>
            <td><img src="${photo_url}" width="50" height="50"></td>
            <td>
                <button class="edit-btn" onclick="editUser(${user.id})">Редактировать</button>
                <button class="delete-btn" onclick="deleteUser(${user.id})">Удалить</button>
            </td>
        </tr>`;
    });

    content += "</table>";
    document.getElementById("content").innerHTML = content;
}

async function loadRooms() {
    const response = await fetchWithAuth(`${API_V1}admins/rooms`);
    const rooms = await response.json();
    if (!rooms) return;

    let content = `<table class="admin-table">
        <tr><th>ID</th><th>Название</th><th>Действия</th></tr>`;

    rooms.forEach(room => {
        content += `<tr>
            <td>${room.id}</td>
            <td>${room.name}</td>
            <td>
                <button class="edit-btn" onclick="editRoom(${room.id})">Редактировать</button>
                <button class="delete-btn" onclick="deleteRoom(${room.id})">Удалить</button>
            </td>
        </tr>`;
    });

    content += "</table>";
    document.getElementById("content").innerHTML = content;
}

async function loadSessions() {
    const response = await fetchWithAuth(`${API_V1}admins/sessions`);
    const sessions = await response.json();
    if (!sessions) return;

    let content = `<table class="admin-table">
        <tr><th>ID</th><th>User ID</th><th>Refresh-токен</th><th>Действия</th></tr>`;

    sessions.forEach(session => {
        content += `<tr>
            <td>${session.id}</td>
            <td>${session.user_id}</td>
            <td>${session.refresh_token}</td>
            <td>
                <button class="delete-btn" onclick="deleteSession(${session.id})">Удалить</button>
            </td>
        </tr>`;
    });

    content += "</table>";
    document.getElementById("content").innerHTML = content;
}

async function deleteUser(userId) {
    if (!confirm("Удалить пользователя?")) return;

    const response = await fetchWithAuth(`${API_V1}admins/users/${userId}`, {method: "DELETE"});
    if (response.ok) {
        await loadRooms();
    }
}

async function deleteRoom(roomId) {
    if (!confirm("Удалить комнату?")) return;
    
    const response = await fetchWithAuth(`${API_V1}admins/rooms/${roomId}`, {method: "DELETE"});
    if (response.ok) {
        await loadRooms();
    }
}

async function deleteSession(sessionId) {
    if (!confirm("Удалить сессию?")) return;
    
    const response = await fetchWithAuth(`${API_V1}admins/sessions/${sessionId}`, {method: "DELETE"});
    if (response.ok) {
        await loadSessions();
    }
}


function editUser(userId) {
    editingUserId = userId;
    document.getElementById("update-profile-modal").style.display = "flex";
}

function editRoom(roomId) {
    editingRoomId = roomId;
    document.getElementById("update-room-modal").style.display = "flex";
}

async function saveUpdatedUser() {
    if (!editingUserId) return;

    const newName = document.getElementById("update-name").value;
    const fileInput = document.getElementById("profile-photo-upload");
    const formData = new FormData();

    if (newName) formData.append("name", newName);
    if (fileInput.files.length > 0) formData.append("photo", fileInput.files[0]);

    const response = await fetchWithAuth(`${API_V1}admins/users/${editingUserId}`, {
        method: "PATCH",
        body: formData,
    });

    if (response.ok) {
        alert("Пользователь обновлен");
        await loadUsers();
        closeUpdateDialog();
    } else {
        handleError(response);
    }
}

async function saveUpdatedRoomr() {
    if (!editingRoomId) return;

    const newName = document.getElementById("update-name-room").value;
    if (!newName) return alert("Введите новое имя комнаты");

    const response = await fetchWithAuth(`${API_V1}admins/rooms/${editingRoomId}`, { method: 'PATCH', body: JSON.stringify({name: newName}), });

    if (response) {
        alert("Комната обновлена");
        await loadRooms();
        closeUpdateDialog();
    }
}
