// Функция для получения всех бригад и отображения их на странице
async function fetchBrigades() {
    try {
        const response = await fetch(`${apiUrl}/brigades/`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                ...getAuthHeaders()
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

// Основной код для работы с DOM
document.addEventListener("DOMContentLoaded", function () {
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

    fetchBrigades(); // Загружаем бригады при загрузке страницы
});
