const apiUrl = "http://localhost:8080/api"; // Настройте URL API по мере необходимости
let lastNotifications = []; // Массив для отслеживания последних уведомлений

// Функция для получения всех бригад и отображения их на странице
async function fetchBrigades() {
    try {
        const response = await fetch(`${apiUrl}/brigades/`);
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        const data = await response.json();
        console.log('Бригады получены:', data); // Отладочный вывод
        const brigadeList = document.getElementById('brigade-list');
        brigadeList.innerHTML = ''; // Очистить список перед добавлением новых элементов

        if (data.length === 0) {
            brigadeList.innerHTML = '<li>Нет доступных бригад.</li>'; // Сообщение, если нет бригад
        } else {
            data.forEach(brigade => {
                const listItem = document.createElement('li');
                // Обновляем текст, чтобы включить ID
                listItem.textContent = `ID: ${brigade.id}, Бригада: ${brigade.name}, Статус: ${brigade.status}`;
                brigadeList.appendChild(listItem);
            });
        }
    } catch (error) {
        console.error('Ошибка при получении бригад:', error);
    }
}

// Функция для создания бронирования
async function createBooking() {
    const slotId = document.getElementById("slotId").value;
    const date = document.getElementById("date").value;
    const status = document.getElementById("status").value;

    try {
        const response = await fetch(`${apiUrl}/bookings/`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ slot_id: slotId, date: date, status: status })
        });

        if (!response.ok) {
            throw new Error(`Error creating booking: ${response.statusText}`);
        }

        const result = await response.json();
        document.getElementById("bookingOutput").textContent = JSON.stringify(result, null, 2);
    } catch (error) {
        console.error(error);
        document.getElementById("bookingOutput").textContent = "Ошибка при создании бронирования.";
    }
}


// Функция для получения уведомлений
async function fetchNotifications() {
    try {
        const response = await fetch(`${apiUrl}/notifications`, {
            method: "GET",
            headers: { "Content-Type": "application/json" }
        });

        if (!response.ok) {
            throw new Error(`Error fetching notifications: ${response.statusText}`);
        }

        const notifications = await response.json();
        return notifications; // Возвращаем уведомления для дальнейшей обработки
    } catch (error) {
        console.error(error);
        return []; // Возвращаем пустой массив в случае ошибки
    }
}

// Функция для отображения уведомлений
function displayNotifications(notifications) {
    const notificationsDiv = document.getElementById("notificationsOutput");
    notificationsDiv.innerHTML = ""; // Очищаем предыдущие уведомления

    notifications.forEach(notification => {
        const notificationElement = document.createElement("div");
        notificationElement.classList.add("notification");
        notificationElement.textContent = notification;

        // CSS для плавного появления
        notificationElement.style.opacity = 0; // Начинаем невидимо
        notificationsDiv.appendChild(notificationElement);

        // Эффект плавного появления
        setTimeout(() => {
            notificationElement.style.opacity = 1; // Плавно появляется
        }, 50); // Короткая задержка для рендеринга
    });
}

// Функция для загрузки и отображения уведомлений
async function loadAndDisplayNotifications() {
    const notifications = await fetchNotifications();

    // Показать новые уведомления, если они изменились
    if (JSON.stringify(notifications) !== JSON.stringify(lastNotifications)) {
        lastNotifications = notifications; // Обновляем последние уведомления
        displayNotifications(notifications); // Отображаем уведомления
    }
}

// Основной код для работы с DOM
document.addEventListener("DOMContentLoaded", function() {
    // Обработчики событий для форм создания бригад
    document.getElementById('create-brigade-form').addEventListener('submit', function(event) {
        event.preventDefault();
        const brigadeName = document.getElementById('brigade-name').value;
        const brigadeStatus = document.getElementById('brigade-status').value;

        // Отправка данных о бригаде на сервер
        fetch(`${apiUrl}/brigades/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: brigadeName, status: brigadeStatus })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            alert('Бригада создана: ' + data.id);
            fetchBrigades(); // Обновить список бригад
            document.getElementById('create-brigade-form').reset();
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Ошибка при создании бригады');
        });
    });

    // Обработчики событий для форм создания задач
    document.getElementById('create-task-form').addEventListener('submit', async function(event) {
        event.preventDefault();
    
        const taskBrigadeId = parseInt(document.getElementById('task-brigade-id').value);
        const taskDescription = document.getElementById('task-description').value;
        const assignedDateInput = document.getElementById('task-assignedat').value; // YYYY-MM-DD
        const taskAssignedAt = new Date(assignedDateInput + "T00:00:00Z").toISOString(); // Преобразование в ISO формат
        const taskStatus = document.getElementById('task-status').value;
    
        // Логируем данные перед отправкой
        console.log('Данные перед отправкой:', {
            brigade_id: taskBrigadeId,
            description: taskDescription,
            assigned_at: taskAssignedAt,
            status: taskStatus
        });
    
        try {
            const response = await fetch(`http://localhost:8080/api/tasks/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    brigade_id: taskBrigadeId,
                    description: taskDescription,
                    assigned_at: taskAssignedAt,
                    status: taskStatus
                })
            });
    
            if (!response.ok) {
                const errorDetails = await response.json();
                throw new Error('Ошибка сервера: ' + JSON.stringify(errorDetails));
            }
    
            const data = await response.json();
            alert('Задача создана с ID: ' + data.id);
            document.getElementById('create-task-form').reset();
        } catch (error) {
            console.error('Ошибка при создании задачи:', error);
            alert('Ошибка при создании задачи: ' + error.message);
        }
    });
});
async function createBrigade(brigadeData) {
    try {
        const response = await fetch(`${apiUrl}/brigades/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(brigadeData),
        });

        if (!response.ok) {
            throw new Error('Ошибка создания бригады: ' + response.statusText);
        }

        const data = await response.json();
        console.log('Бригада создана:', data);
    } catch (error) {
        console.error('Ошибка при создании бригады:', error);
    }
}


// Загрузка уведомлений сразу при загрузке страницы
window.onload = async () => {
    await loadAndDisplayNotifications(); // Получаем и отображаем уведомления при загрузке
    setInterval(loadAndDisplayNotifications, 5000); // Обновление уведомлений каждые 5 секунд
};
