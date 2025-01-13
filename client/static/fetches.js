let userid = document.getElementById("userid");
let username = document.getElementById("window-input-username");
let password = document.getElementById("window-input-password");
let isDone = false;
const accessToken = "";
localStorage.setItem("accessToken", accessToken);
const apiURL = "http://localhost:8000/api/auth/sign-up";

const User = {
  Id: "",
  Username: username.value,
  Password: password.value,
};

// const UserList = {
//   Id: "",
//   UserId: userid,
//   ListId: "",
// };

// const ListsTAsks = {
//   Id: "",
//   ListId: "",
//   TaskID: "",
// };

// async function fetchAccessToken() {
//   try {
//     const response = await fetch(apiURL, {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify(User),
//     });

//     if (!response.ok) {
//       throw new Error(`Ошибка: ${response.statusText}`);
//     }

//     const data = await response.json();
//     const accessToken = data.token;
//     localStorage.setItem("accessToken", accessToken);

//     console.log("Access token сохранен:", accessToken);
//   } catch (error) {
//     console.error("Ошибка получения токена:", error);
//   }
// }

async function sendUserData() {
  try {
    const response = await fetch(
      "http://192.168.71.111:8000/api/auth/sign-in",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(User), //Обьект с данными о пользователе
      }
    );

    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(errorResponse.message || "Ошибка отправки пользователя");
    }
    const result = await response.json();
    console.log("Данные пользователя отправлены:", result);
    alert("Пользователь успешно отправлен!");
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось отправить пользователя: " + error.message);
  }
}

async function sendList() {
  const ListData = {
    id: "",
    Title: createList.value,
  };
  try {
    const response = await fetch("/api/lists/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(ListData), //Обьект с данными о листе
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(errorResponse.message || "Ошибка отправки листов");
    }
    const result = await response.json();
    console.log("Данные о листе отправлены", result);
    renderSingleList(result);

    alert("Данные о листе отправлены");
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось отправить данные листа " + error.message);
  }
}

async function getAllList() {
  try {
    const response = await fetch("/api/lists/", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(
        console.log(errorResponse.message || "Ошибка получения данных листов")
      );
    }
    const result = await response.json();
    console.log("Листы получены:", result);
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось получить листы: " + error.message);
  }
}

async function sendTask() {
  const Task = {
    Id: "",
    Title: taskInput.textContent,
    Done: isDone,
  };
  try {
    const response = await fetch("/api/lists/tasks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(Task), //Обьект с данными о задаче
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(errorResponse.message || "Ошибка отправки задачи");
    }
    const result = await response.json();
    console.log("Данные о задаче отправлены", result);
    alert("Данные о задаче отправлены");
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось отправить данные задачи " + error);
  }
}

async function getAllTask() {
  try {
    const response = await fetch("http://192.168.71.111:8000/api/tasks", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      throw new Error(errorResponse.message || "Ошибка получения данных задач");
    }
    const result = await response.json();
    console.log("Задачи получены:", result);

    result.forEach((task) => {
      renderSingleTask(task); // Рендерим каждую задачу по очереди
    });

    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось получить задачи: " + error.message);
  }
}
