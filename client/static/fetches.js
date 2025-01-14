let userid = document.getElementById("userid");
let username = document.getElementById("window-input-username");
let password = document.getElementById("window-input-password");
let isDone = false;

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

async function fetchAccessToken() {
  try {
    const response = await fetch(apiURL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(User),
    });

    if (!response.ok) {
      throw new Error(`Ошибка: ${response.statusText}`);
    }

    const data = await response.json();
    const accessToken = data.token;
    localStorage.setItem("accessToken", accessToken);

    console.log("Access token сохранен:", accessToken);
  } catch (error) {
    console.error("Ошибка получения токена:", error);
  }
}

async function signUp(username, password) {
  const url = "/api/auth/sign-up";
  const userData = {
    username: username,
    password: password,
  };

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.err || "Unknown error occurred");
    }

    const result = await response.json();
    console.log("User signed up successfully:", result);
  } catch (error) {
    console.error("Error during sign up:", error.message);
  }
}

async function signIn(username, password) {
  const url = "/api/auth/sign-in";
  const userData = {
    username: username,
    password: password,
  };

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.err || "Unknown error occurred");
    }

    const result = await response.json();
    localStorage.setItem("accessToken", result.token);
    console.log("User signed in successfully:", result.token);
  } catch (error) {
    console.error("Error during sign in:", error.message);
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

async function getAllLists() {
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
    Title: taskInput.value,
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

async function getAllTasks() {
  try {
    const response = await fetch("/api/lists/tasks", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
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
