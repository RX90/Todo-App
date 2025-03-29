let userid = document.getElementById("userid");
let signinError = document.getElementById("signin-error");
let signinError2 = document.getElementById("signin-error2");

let signinInputName = document.querySelector(".signin-name-input");
let signinInputPassword = document.querySelector(".signin-password-input");

let errorUser = document.getElementById("error-user-message");
let usernameRegister = document.querySelector(".signup-name-input");

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
      console.log("еррор дата", errorData);
      if (
        errorData.err === "can't create user: username is already taken" &&
        response.status === 500
      ) {
        console.log("Такой пользователь уже есть");
        errorUser.textContent = "Пользователь с таким логином уже существует";
        errorUser.style.display = "block";
        usernameRegister.style.outline = "3px solid red";
      }

      throw new Error(errorData.err || "Unknown error occurred");
    }

    const result = await response.json();
    console.log("User signed up successfully:", result);
    return true;
  } catch (error) {
    console.error("Error during sign up:", error.message);
    return false;
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

      if (response.status === 400 || response.status === 401) {
        console.log("Бугагага");
        signinError.textContent = "Некорректно введённые данные";
        signinError.style.visibility = "visible";

        signinError2.textContent = "Некорректно введённые данные";
        signinError2.style.visibility = "visible";

        signinInputName.style.outline = "3px solid red";
        signinInputPassword.style.outline = "3px solid red";

        signinInputName.blur();
        signinInputPassword.blur();
      }
      throw new Error(errorData.err || "Unknown error occurred");
    }

    const result = await response.json();

    localStorage.setItem("accessToken", result.token);
    return true;
  } catch (error) {
    console.error("Error during sign in:", error.message);
    return false;
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка отправки листа");
          return await sendList(title);
        } else {
          showPopupSignin();
        }
      }

      throw new Error(errorResponse.message || "Error send list");
    }

    const result = await response.json();
    console.log("Отправка списков прошла успешно:", result.list_id);
    return result.list_id;
  } catch (error) {
    console.error("Ошибка отправки листа", error.message);
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка отобразить листы");
          return await getAllLists();
        }
      }
      throw new Error(errorResponse.message || "Ошибка получения всех листов");
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка отправки задачи");
          return await sendTask(listId, taskTitle);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Ошибка отправки задачи");
    }

    const result = await response.json();
    console.log("Данные о задаче отправлены", result);
    return result;
  } catch (error) {
    console.error("Ошибка:", error);
  }
}

async function getAllTasks(listId) {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    showPopupSignin();
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
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попвтка принять все задачи");
          return await getAllTasks(listId);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Ошибка получения всех задач");
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
    showPopupSignin();
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
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка изменить название задачи");
          return await toggleTaskState(taskId, isDone, listId);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Ошибка изменения задачи");
    }

    const updatedTask = await response.json();
    console.log("Задача успешно обновлена:", updatedTask);
    return updatedTask;
  } catch (error) {
    console.error("Ошибка в toggleTaskState:", error);
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
    showPopupSignin();
    console.log("Не получилось");
    throw error;
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
        if (errorResponse.err === "tСрок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка выйти из аккаунта");
          return await logout();
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Ошибка при выходе из системы");
    }

    localStorage.removeItem("accessToken");
  } catch (error) {
    console.error("Ошибка при выходе из системы", error);
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка удалить задачу");
          return await DeleteTask(listId, taskId);
        } else {
          showPopupSignin();
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка переименовать задачу");
          return await EditTask(taskId, listId, newTitle);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(
        errorResponse.message || "Ошибка при редактировании задачи"
      );
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.log("Не удалось переименовать задачу");
  }
}

async function DeleteList(listId) {
  try {
    const response = await fetch(`/api/lists/${listId}`, {
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
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка удаления листа");
          return await DeleteList(listId);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Error while delete list");
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.error("Не удалось удалить лист ", error);
  }
}

async function EditList(listId, newTitleList) {
  const updateTitleList = { title: newTitleList };
  try {
    const response = await fetch(`/api/lists/${listId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify(updateTitleList),
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      console.error("Error logout:", errorResponse);
      if (response.status === 401) {
        if (errorResponse.err === "Срок действия токена истек") {
          console.log("Срок действия токена истек, обновление токена...");
          await refreshToken();
          console.log("Повторная попытка переименования листа");
          return await EditList(listId, newTitleList);
        } else {
          showPopupSignin();
        }
      }
      throw new Error(errorResponse.message || "Error while edit list");
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.log("Не удалось переименовать лист");
  }
}
