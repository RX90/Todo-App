let userid = document.getElementById("userid");
let username = document.getElementById("window-input-username");
let password = document.getElementById("window-input-password");
const isDone = false;
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
    console.log("User signed in successfully:", result);
    localStorage.setItem("accessToken", result.token);
    // localStorage.setItem("refreshToken".result.refreshToken)
  } catch (error) {
    console.error("Error during sign in:", error.message);
  }
}

async function sendList(title) {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    showPopup();
    return;
  }
  const url = "/api/lists/";
  const userData = {
    Title: title,
  };

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(userData),
    });

    if (!response.ok) {
      const errorResponse = await response.json();

      if (
        response.status === 401 &&
        errorResponse.message === "token has expired"
      ) {
        console.log("Токен истек, обновление токена...");
        await refreshToken();
        console.log("Повторный запрос...");
        return sendList(title);
      }
    }

    const result = await response.json();
    console.log("Send lists successfully:", result.list_id);
    return result;
  } catch (error) {
    console.error("Error send lists", error.message);
  }
}

async function getAllLists() {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    showPopup();
    return;
  }
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

      if (
        response.status === 401 &&
        errorResponse.message === "token has expired"
      ) {
        console.log("Токен истек, обновление токена...");
        await refreshToken();
        console.log("Повторный запрос...");
        return sendList(title);
      }
    }

    const result = await response.json();
    console.log("Листы получены:", result);
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
  }
}

async function sendTask(listId, taskTitle) {
  const Task = {
    title: taskTitle,
    done: false,
  };

  const accessToken = localStorage.getItem("accessToken");
  if (!accessToken) {
    showPopup();
    return;
  }

  try {
    const response = await fetch(`/api/lists/${listId}/tasks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(Task), //Обьект с данными о задаче
    });

    if (!response.ok) {
      const errorResponse = await response.json();

      if (
        response.status === 401 &&
        errorResponse.message === "token has expired"
      ) {
        console.log("Токен истек, обновление токена...");
        await refreshToken();
        console.log("Повторный запрос...");
        return sendList(title);
      }
    }

    const result = await response.json();
    console.log("Данные о задаче отправлены", result);
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось отправить данные задачи " + error);
  }
}

async function getAllTasks(listId) {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    showPopup();
    return;
  }

  try {
    const response = await fetch(`/api/lists/${listId}/tasks`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });

    if (!response.ok) {
      const errorResponse = await response.json();

      if (
        response.status === 401 &&
        errorResponse.message === "token has expired"
      ) {
        console.log("Токен истек, обновление токена...");
        await refreshToken();
        console.log("Повторный запрос...");
        return sendList(title);
      }
    }

    const result = await response.json();
    console.log("Задачи получены:", result);

    result.forEach((tasks) => {
      renderSingleTask(tasks);
    });

    return result;
  } catch (error) {
    console.error("Ошибка:", error);
    console.log("Не удалось получить задачи: " + error.message);
  }
}

async function toggleTaskState(taskId, isDone, listId) {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    showPopup();
    return;
  }

  console.log("taskId:", taskId);
  console.log("isDone:", isDone);
  console.log("listId:", listId);

  const updatePayload = {
    done: isDone,
    // title: title
  };

  try {
    const response = await fetch(`/api/lists/${listId}/tasks/${taskId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(updatePayload),
    });

    if (!response.ok) {
      const errorResponse = await response.json();

      if (
        response.status === 401 &&
        errorResponse.message === "token has expired"
      ) {
        console.log("Токен истек, обновление токена...");
        await refreshToken();
        console.log("Повторный запрос...");
        return sendList(title);
      }
    }

    const updatedTask = await response.json();
    console.log("Состояние задачи успешно обновлено:", updatedTask);
    return updatedTask;
  } catch (error) {
    console.error("Ошибка:", error);
    alert("Не удалось обновить состояние задачи: " + error.message);
  }
}

async function refreshToken() {
  try {
    const response = await fetch(`/api/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Ошибка сервера:", errorResponse);
      throw new Error(errorResponse.message || "Ошибка отправки accessToken");
    }
    const { accessToken } = await response.json();
    console.log("accessToken обновлен");
    localStorage.setItem("accessToken", accessToken);
    return accessToken;
  } catch (error) {
    showPopup();
    console.log("Не получилось");
  }
}
