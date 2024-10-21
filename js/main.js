const apiUrl = "http://localhost:8080/api";

// Проверка наличия токена
function getAuthHeaders() {
    const token = localStorage.getItem("token");
    return token ? { "Authorization": `Bearer ${token}` } : {};
}
async function checkAuthentication() {
    const token = localStorage.getItem("token");
    if (!token) {
        window.location.href = "login.html"; // Перенаправление на страницу входа при отсутствии токена
        return;
    }

    const response = await fetch("http://localhost:8080/auth/userinfo", {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`
        }
    });

    if (!response.ok) {
        localStorage.removeItem("token"); // Удаление токена при ошибке
        window.location.href = "login.html";
    }
}