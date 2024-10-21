let lastNotifications = [];

// Функция для получения уведомлений
async function fetchNotifications() {
    try {
        const response = await fetch(`${apiUrl}/notifications`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });

        if (!response.ok) throw new Error('Ошибка получения уведомлений: ' + response.statusText);
        return await response.json();
    } catch (error) {
        console.error(error);
        return [];
    }
}

function displayNotifications(notifications) {
    const notificationsDiv = document.getElementById("notificationsOutput");
    if (!notificationsDiv) return; // Прекращаем выполнение, если элемента нет

    notificationsDiv.innerHTML = ""; // Очищаем предыдущие уведомления

    notifications.forEach(notification => {
        const notificationElement = document.createElement("div");
        notificationElement.classList.add("notification");
        notificationElement.textContent = notification;
        notificationElement.style.opacity = 0;
        notificationsDiv.appendChild(notificationElement);
        setTimeout(() => { notificationElement.style.opacity = 1; }, 50);
    });
}

async function loadAndDisplayNotifications() {
    const notifications = await fetchNotifications();
    if (JSON.stringify(notifications) !== JSON.stringify(lastNotifications)) {
        lastNotifications = notifications;
        displayNotifications(notifications);
    }
}

// Загрузка уведомлений сразу при загрузке страницы
window.onload = async () => {
    await loadAndDisplayNotifications();
    setInterval(loadAndDisplayNotifications, 5000);
};
