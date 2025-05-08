const API_BASE = "http://localhost:8080/"
const API_V1 = `${API_BASE}api/v1/`


document.addEventListener("DOMContentLoaded", async () => {
  const currentPath = window.location.pathname;

  if (currentPath.endsWith("index.html") || currentPath.endsWith("register.html")) {
    return;
  }
  try {
    const userRole = await IdentifyUser();
    if (!userRole) {
      window.location.href = "index.html";
      return;
    } else {
      
      if (currentPath.endsWith("home.html")) {
        document.getElementById("logout-button").addEventListener("click", logout);
        if (userRole === "admin") {
          createAdminButton();
        }
        loadRooms();
        loadProfile();
      } else if (currentPath.endsWith("chat-room.html")) {
        const roomID = localStorage.getItem("currentRoomID");
        if (roomID) joinRoom(roomID);
      } else if (currentPath.endsWith("admin.html") && userRole === "admin") {
        document.getElementById("tab-users").addEventListener("click", loadUsers);
        document.getElementById("tab-rooms").addEventListener("click", loadRooms);
        document.getElementById("tab-sessions").addEventListener("click", loadSessions);

        loadUsers();
      } else {
        await logout();
      }
    }
  } catch (err) {
    console.error("Ошибка проверки авторизации", err);
    window.location.href = "index.html";
  }
});

async function IdentifyUser() {
  const response = await fetchWithAuth(`${API_V1}users/auth/check`);
  if (response.status !== 200) return null;
  
  data = await response.json();
  const userRole = data.role;
  return userRole
}

function preventBackNavigation() {
  history.replaceState(null, null, window.location.href);
  window.addEventListener("popstate", function () {
    history.pushState(null, null, window.location.href);
  });
}

async function fetchWithAuth(url, options = {}) {
  let response = await fetch(url, { 
    ...options, 
    credentials: "include",
    headers: { "Cache-Control": "no-cache, no-store, must-revalidate" }
  });
  if (response.status === 401) {
    const refreshResponse = await fetch(`${API_V1}auth/refresh`, { method: 'POST', credentials: 'include' });
    if (refreshResponse.ok) {
      response = await fetch(url, { 
        ...options, 
        credentials: "include",
        headers: { "Cache-Control": "no-cache, no-store, must-revalidate" }
      });
    }
  }
  return response;
}

function handleError(response) {
  response.json().then(error => alert(error.error));
}

async function logout() {
  await fetchWithAuth(`${API_V1}users/logout`, { method: 'DELETE'});
  localStorage.clear();
  window.location.href = "index.html";
}

function createAdminButton() {
  const button = document.createElement("button");
  button.textContent = "Панель администратора";
  button.classList.add("admin-button");
  button.onclick = function() {
    window.location.href = "admin.html"; 
  };
    
  document.body.appendChild(button);
}