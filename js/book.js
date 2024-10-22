// Функция для создания бронирования
async function createBooking() {
    const slotId = document.getElementById("slotId").value;
    const date = document.getElementById("date").value;
    const status = document.getElementById("status").value;

    try {
        const response = await fetch(`${apiUrl}/bookings/`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({ slot_id: slotId, date, status })
        });
        if (!response.ok) throw new Error('Ошибка создания бронирования: ' + response.statusText);
        const result = await response.json();
        document.getElementById("bookingOutput").textContent = JSON.stringify(result, null, 2);
    } catch (error) {
        console.error(error);
        document.getElementById("bookingOutput").textContent = "Ошибка при создании бронирования.";
    }
}

// Привязка обработчика события к кнопке
document.getElementById("bookingButton").addEventListener("click", createBooking);


