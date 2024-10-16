const apiUrl = "http://localhost:8080/api";
let lastNotifications = [];

// Проверка наличия токена
function getAuthHeaders() {
    const token = localStorage.getItem("token");
    return token ? { "Authorization": `Bearer ${token}` } : {};
}

// Функция для получения всех бригад и отображения их на странице
async function fetchBrigades() {
    try {
        const response = await fetch(`${apiUrl}/brigades/`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                ...getAuthHeaders() // Добавляем токен в запрос
            }
        });
        if (!response.ok) throw new Error('Ошибка при получении бригад: ' + response.statusText);

        const data = await response.json();
        const brigadeList = document.getElementById('brigade-list');
        if (!brigadeList) return; // Прекращаем выполнение, если элемента нет

        brigadeList.innerHTML = '';
        data.length === 0
            ? brigadeList.innerHTML = '<li>Нет доступных бригад.</li>'
            : data.forEach(brigade => {
                const listItem = document.createElement('li');
                listItem.textContent = `ID: ${brigade.id}, Бригада: ${brigade.name}, Статус: ${brigade.status}`;
                brigadeList.appendChild(listItem);
            });
    } catch (error) {
        console.error(error);
    }
}

// Функция для создания бронирования
async function createBooking() {
    if (!elementExists('bookingOutput')) return; // Прекращаем выполнение, если элемента нет

    const slotId = document.getElementById("slotId").value;
    const date = document.getElementById("date").value;
    const status = document.getElementById("status").value;

    try {
        const response = await fetch(`${apiUrl}/bookings/`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...getAuthHeaders()
            },
            body: JSON.stringify({ slot_id: slotId, date, status })
        });
        if (!response.ok) throw new Error(`Ошибка создания бронирования: ${response.statusText}`);
        const result = await response.json();
        document.getElementById("bookingOutput").textContent = JSON.stringify(result, null, 2);
    } catch (error) {
        console.error(error);
        document.getElementById("bookingOutput").textContent = "Ошибка при создании бронирования.";
    }
}

// Функции для уведомлений
async function fetchNotifications() {
    try {
        const response = await fetch(`${apiUrl}/notifications`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });

        if (!response.ok) throw new Error(`Ошибка получения уведомлений: ${response.statusText}`);
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

// Основной код для работы с DOM
document.addEventListener("DOMContentLoaded", function () {
    if (elementExists('create-brigade-form')) {
        document.getElementById('create-brigade-form').addEventListener('submit', function (event) {
            event.preventDefault();
            const brigadeName = document.getElementById('brigade-name').value;
            const brigadeStatus = document.getElementById('brigade-status').value;

            fetch(`${apiUrl}/brigades/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem("token")}`
                },
                body: JSON.stringify({ name: brigadeName, status: brigadeStatus })
            })
                .then(response => {
                    if (!response.ok) throw new Error('Ошибка создания бригады');
                    return response.json();
                })
                .then(data => {
                    alert('Бригада создана: ' + data.id);
                    fetchBrigades();
                    document.getElementById('create-brigade-form').reset();
                })
                .catch(error => {
                    console.error(error);
                    alert('Ошибка при создании бригады');
                });
        });
    }

    if (elementExists('create-task-form')) {
        document.getElementById('create-task-form').addEventListener('submit', async function (event) {
            event.preventDefault();
            const taskBrigadeId = parseInt(document.getElementById('task-brigade-id').value);
            const taskDescription = document.getElementById('task-description').value;
            const taskAssignedAt = new Date(document.getElementById('task-assignedat').value + "T00:00:00Z").toISOString();
            const taskStatus = document.getElementById('task-status').value;

            try {
                const response = await fetch(`${apiUrl}/tasks/`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem("token")}`
                    },
                    body: JSON.stringify({
                        brigade_id: taskBrigadeId,
                        description: taskDescription,
                        assigned_at: taskAssignedAt,
                        status: taskStatus
                    })
                });

                if (!response.ok) throw new Error('Ошибка создания задачи');
                const data = await response.json();
                alert('Задача создана с ID: ' + data.id);
                document.getElementById('create-task-form').reset();
            } catch (error) {
                console.error('Ошибка при создании задачи:', error);
                alert('Ошибка при создании задачи: ' + error.message);
            }
        });
    }
});

// Загрузка уведомлений сразу при загрузке страницы
window.onload = async () => {
    await loadAndDisplayNotifications();
    setInterval(loadAndDisplayNotifications, 5000);
};
