let urlUser = "http://192.168.71.111:8080/users";
let userid = document.getElementById("userid");
let username = document.getElementById("username");
let password = document.getElementById("password");

const userData = {
  userid: userid,
  username: username,
  password: password,
};

async function sendUserData() {
  try {
    const response = await fetch(urlUser, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData), //Обьект с данными о пользователе
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(errorResponse.message || "Ошибка отправки пользователя");
    }
    const result = await response.json();
    console.log("Данные пользователя отправлены:", result);
    alert("Пользователь успешно отправлен!");
    return response;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось отправить пользователя: " + error.message);
  }
}

async function getAllList() {
  try {
    const response = await fetch("http://localhost:8080/users", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(
        errorResponse.message || "Ошибка получения данных пользователя"
      );
    }
    const result = await response.json();

    console.log("Пользователь получен:", result);
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось получить пользователя: " + error.message);
  }
}
