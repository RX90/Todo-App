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
      console.error("Error send list:", errorResponse);

      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await sendList(title);
        } else {
          showPopup();
        }
      }

      throw new Error(errorResponse.message || "Error send list");
    }

    const result = await response.json();
    console.log("Send lists successfully:", result.list_id);
    return result.list_id;
  } catch (error) {
    console.error("Error send lists", error.message);
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
      console.error("Error updating task:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await getAllLists();
        }
      }
      throw new Error(errorResponse.message || "Error get all lists");
    }

    const result = await response.json();
    console.log("Листы получены:", result);
    await logoutLocalStorage();
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

  try {
    const response = await fetch(`/api/lists/${listId}/tasks/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(Task), //Обьект с данными о задаче
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error updating task:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await sendTask(listId, title);
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error send task");
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
    const response = await fetch(`/api/lists/${listId}/tasks/`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error updating task:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await getAllTasks(listId);
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error get all task");
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

  console.log("AccessToken:", accessToken);
  console.log("taskId:", taskId);
  console.log("isDone:", isDone);
  console.log("listId:", listId);

  const updatePayload = {
    done: isDone,
  };

  try {
    const response = await fetch(`/api/lists/${listId}/tasks/${taskId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      body: JSON.stringify(updatePayload),
    });

    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error updating task:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await sendTask(taskId, listId, isDone);
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error updating task");
    }

    const updatedTask = await response.json();
    console.log("Task updated successfully:", updatedTask);
    return updatedTask;
  } catch (error) {
    console.error("Error in toggleTaskState:", error);
    alert("Failed to update task: " + error.message);
    return null;
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
    const accessToken = await response.json();
    console.log("accessToken обновлен");
    localStorage.setItem("accessToken", accessToken.token);
    return accessToken.token;
  } catch (error) {
    showPopup();
    console.log("Не получилось");
    throw error; // Пробрасываем ошибку дальше
  }
}

async function logout() {
  try {
    const response = await fetch(`/api/auth/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error logout:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying send list");
          await logout();
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error while logout");
    }
    // const result = await response.json();
    localStorage.removeItem("accessToken");
    // return result;
  } catch (error) {
    console.error("Error in logout", error);
    alert("Failed to logout " + error.message);
    localStorage.removeItem("accessToken");
  }
}

async function DeleteTask(listId, taskId) {
  try {
    const response = await fetch(`/api/lists/${listId}/tasks/${taskId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error logout:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying delete task");
          await logout();
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error while delete task");
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.error("Не удалось удалить задачу ", error);
  }
}

async function EditTask(taskId, listId, newTitle) {
  const updateTitle = { title: newTitle };
  try {
    const response = await fetch(`/api/lists/${listId}/tasks/${taskId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(updateTitle),
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error logout:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "token has expired") {
          console.log("Token expired, refreshing token...");
          await refreshToken();
          console.log("Retrying edit task");
          await logout();
        } else {
          showPopup();
        }
      }
      throw new Error(errorResponse.message || "Error while edit task");
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.log("Не удалось переименовать задачу");
  }
}

//Giga5
//GigaMenchik555
