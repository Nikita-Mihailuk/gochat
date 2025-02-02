async function loginUser() {
  const email = document.getElementById("login-email").value.trim();
  const password = document.getElementById("login-password").value.trim();

  try {
    const response = await fetch(`${API_V1}login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    if (!response.ok) return handleError(response);

    const userRole = await IdentifyUser();

    if (userRole === "admin") {
      window.location.href = "admin.html";
    } else {
      window.location.href = "home.html";
    }
  } catch (err) {
    console.error("Ошибка авторизации:", err);
  }
}
