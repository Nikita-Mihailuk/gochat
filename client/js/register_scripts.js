async function registerUser() {
    const email = document.getElementById("register-email").value.trim();
    const password = document.getElementById("register-password").value.trim();
    const name = document.getElementById("register-name").value.trim();
  
    if (email.length > 30 || password.length > 30 || name.length > 30){
      alert("Превышено допустимое количество символов")
      return
    }
    
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailPattern.test(email)) {
     alert("Введите корректный email");
     return;
    }

    if (email && password && name) {
      try {
        const response = await fetch(`${API_V1}register`, {
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