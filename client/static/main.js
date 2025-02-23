let createList = document.getElementById("new-list");
let menu = document.querySelector(".menu");
let details = document.getElementById("details");
let title = document.getElementById("title");
let listIdCounter = 0;
let loginButton = document.getElementById("login-button-signin");
let registerButton = document.getElementById("login-button-signup");
let infoText = document.getElementById("info-text");
let panelSignin = document.querySelector(".background-sign-in");
let signinSendData = document.getElementById("signin-button");
let username = document.querySelector(".signin-name-input");
let password = document.querySelector(".signin-password-input");

let activePanel = null;

if (
  !localStorage.getItem("accessToken") ||
  localStorage.getItem("accessToken").trim() === ""
) {
  console.log("Токен не найден");
  loginButton.addEventListener("click", function () {
    showPopupSignin();
  });
} else {
  console.log("Токен получен", localStorage.getItem("accessToken"));
  hiddenPopupSignin();
}

async function logoutLocalStorage() {
  registerButton.style.display = "none";
  loginButton.textContent = "Выйти";
  loginButton.addEventListener("click", async function () {
    await logout();
    location.reload();
  });
}

function showPopupSignin() {
  panelSignin.style.display = "block";
}

function hiddenPopupSignin() {
  panelSignin.style.display = "none";
}

loginButton.addEventListener("click", function () {
  showPopupSignin();
});

signinSendData.addEventListener("click", async function () {
  const user = username.value;
  const pass = password.value;

  const success = await signIn(user, pass); // Проверяем успешность логина
  if (success) {
    console.log("Логин успешен!");
    hiddenPopupSignin(); // Скрываем форму логина
    await updateLists(); // Загружаем и отображаем листы
  } else {
    console.error("Ошибка входа! Проверьте логин и пароль.");
  }
});

// document.addEventListener("DOMContentLoaded", function () {
function renderSingleList(listid, title) {
  const menuItem = document.createElement("div");
  menuItem.classList.add("menu-item");
  menuItem.setAttribute("data-id", listid);

  const imgList = document.createElement("img");
  imgList.classList.add("icon-list");
  imgList.setAttribute("src", "/src/img/list.svg");

  const addList = document.createElement("input");
  addList.classList.add("add-list");
  addList.value = title;
  addList.disabled = true;

  const dots = document.createElement("img");
  dots.src = "/src/img/dots.svg";
  dots.classList.add("dots");

  dots.addEventListener("click", (event) => {
    event.stopPropagation(); // Останавливаем всплытие
    if (activePanel) {
      activePanel.remove();
    }

    const dotsPanel = document.createElement("div");
    dotsPanel.classList.add("dots-panel");

    const dotsEdit = document.createElement("button");
    dotsEdit.textContent = "Переименовать";
    dotsEdit.classList.add("dots-edit");

    const iconDotsEdit = document.createElement("img");
    iconDotsEdit.src = "/src/img/edit.svg";
    iconDotsEdit.classList.add("dots-delete-icon");

    const dotsDelete = document.createElement("button");
    dotsDelete.textContent = "Удалить";
    dotsDelete.style.color = "red";
    dotsDelete.classList.add("dots-delete");

    const iconDotsDelete = document.createElement("img");
    iconDotsDelete.src = "/src/img/red-delete.svg";
    iconDotsDelete.classList.add("dots-delete-icon");

    dotsEdit.appendChild(iconDotsEdit);
    dotsDelete.appendChild(iconDotsDelete);
    dotsPanel.appendChild(dotsEdit);
    dotsPanel.appendChild(dotsDelete);
    document.body.appendChild(dotsPanel);

    dotsDelete.addEventListener("click", async function () {
      const listId = menuItem.getAttribute("data-id");
      await DeleteList(listId);
      menuItem.remove();
      dotsPanel.remove();
      details.style.display = "none";
    });

    dotsEdit.addEventListener("click", function () {
      addList.disabled = false;
      addList.focus();

      addList.addEventListener("keydown", async function (event) {
        if (event.key === "Enter") {
          const newTitleList = addList.value.trim();
          if (newTitleList !== "" && newTitleList !== title) {
            await EditList(listid, newTitleList);
            title = newTitleList;
          }
          addList.disabled = true;
        }
      });
    });

    const place = dots.getBoundingClientRect();
    dotsPanel.style.top = `${place.top + window.scrollY + 20}px`;
    dotsPanel.style.left = `${place.left + window.scrollX + 10}px`;

    activePanel = dotsPanel;
  });

  menuItem.appendChild(imgList);
  menuItem.appendChild(addList);
  menuItem.appendChild(dots);
  menu.appendChild(menuItem);

  menuItem.addEventListener("click", () => {
    openPanel(listid, title);
    console.log(listid, title);
  });
}
// });

//Сброс панели по клику на экран
document.addEventListener("click", (event) => {
  if (
    activePanel &&
    !activePanel.contains(event.target) &&
    !event.target.classList.contains("dots")
  ) {
    activePanel.remove();
    activePanel = null;
  }
});

function renderSingleTask(task) {
  const taskList = document.querySelector(".task-list");
  const completedList = document.querySelector(".completed");

  const menuTask = document.createElement("div");
  menuTask.classList.add("menu-task");
  menuTask.setAttribute("data-task-id", task.id);

  const circleIcon = document.createElement("img");
  circleIcon.src = task.done ? "/src/img/done.svg" : "/src/img/!done.svg";
  circleIcon.classList.add("circle-icon");

  const titleTask = document.createElement("input");
  titleTask.value = task.title;
  titleTask.classList.add("title-task");
  titleTask.disabled = true;

  const editTask = document.createElement("img");
  editTask.src = "/src/img/violet-edit.svg";
  editTask.classList.add("edit-task");

  editTask.addEventListener("click", async function (event) {
    titleTask.disabled = false;
    titleTask.focus();

    titleTask.addEventListener("keydown", async function (event) {
      if (event.key === "Enter" && titleTask.value.trim() !== "") {
        const listId = details.getAttribute("data-id");
        const taskId = Number(menuTask.getAttribute("data-task-id"));
        const newTitle = titleTask.value.trim();

        await EditTask(taskId, listId, newTitle);

        titleTask.disabled = true;
      }
    });
  });

  const deleteTask = document.createElement("img");
  deleteTask.src = "/src/img/violet-delete.svg";
  deleteTask.classList.add("delete-task");

  deleteTask.addEventListener("click", async function () {
    const listId = details.getAttribute("data-id");
    const taskId = Number(menuTask.getAttribute("data-task-id"));
    const success = await DeleteTask(listId, taskId);

    if (success) {
      menuTask.remove(); // Удаляем задачу из DOM
    } else {
      alert("Ошибка при удалении задачи!");
    }
  });

  menuTask.remove();

  function moveToCorrectPlace() {
    menuTask.remove();
    const listId = details.getAttribute("data-id");
    const titleCompleted = document.querySelector(".title-completed");

    if (task.done) {
      titleCompleted.style.display = "block";
      completedList.appendChild(menuTask);
      titleTask.style.textDecoration = "line-through";
    } else {
      taskList.appendChild(menuTask);
      titleTask.style.textDecoration = "none";
    }

    if (completedList.children.length === 0) {
      titleCompleted.style.display = "none";
    }
  }

  circleIcon.addEventListener("click", async function () {
    const newState = !task.done;
    const listId = details.getAttribute("data-id");
    const taskId = Number(menuTask.getAttribute("data-task-id"));

    try {
      const updatedTask = toggleTaskState(taskId, newState, listId);
      if (updatedTask) {
        task.done = newState;
        circleIcon.src = newState ? "/src/img/done.svg" : "/src/img/!done.svg";
        moveToCorrectPlace();
      } else {
        console.error("Ошибка обновления задачи");
      }
    } catch (error) {
      console.error("Не удалось обновить задачу:", error);
    }
  });

  menuTask.appendChild(circleIcon);
  menuTask.appendChild(editTask);
  menuTask.appendChild(deleteTask);
  menuTask.appendChild(titleTask);

  moveToCorrectPlace();
}

//Создание листов
createList.addEventListener("keydown", async function (event) {
  if (event.key === "Enter" && createList.value.trim() !== "") {
    const title = createList.value.trim();

    if (
      !localStorage.getItem("accessToken") ||
      localStorage.getItem("accessToken").trim() === ""
    ) {
      console.log("Токен не найден");
      showPopupSignin();
    } else {
      console.log("Токен получен", localStorage.getItem("accessToken"));
      hiddenPopupSignin();
    }

    try {
      const listId = await sendList(title);

      if (listId) {
        renderSingleList(listId, title);
        openPanel(listId, title);
      }
    } catch (error) {
      console.error("Ошибка при создании листа:", error);
      alert("Не удалось создать лист: " + error.message);
    }
    createList.value = "";
  }
});

function openPanel(listId, listName) {
  details.style.display = "block";
  details.setAttribute("data-id", listId);
  title.textContent = listName;

  const taskList = document.querySelector(".task-list");
  const completedList = document.querySelector(".completed");

  taskList.innerHTML = "";
  completedList.innerHTML = "";
  // dotsPanel.innerHTML = "";

  const titleCompleted = document.querySelector(".title-completed");
  titleCompleted.style.display = "none";

  getAllTasks(listId);
}

const taskInput = document.getElementById("task-input");
const taskButton = document.getElementById("task-button");

async function createTask() {
  if (taskInput.value.trim() === "") return;

  const taskTitle = taskInput.value.trim();
  const listId = details.getAttribute("data-id");

  if (
    !localStorage.getItem("accessToken") ||
    localStorage.getItem("accessToken").trim() === ""
  ) {
    console.log("Токен не найден");
    showPopup();
    return;
  }

  try {
    const newObject = await sendTask(listId, taskTitle);
    const newTaskId = newObject.task_id;
    console.log("Title:", taskTitle);

    if (newTaskId) {
      const newTask = {
        id: newTaskId,
        title: taskTitle,
        done: false,
      };

      console.log("New TASK:", newTask);
      renderSingleTask(newTask);
      taskInput.value = ""; // Очищаем поле после успешного создания
    }
  } catch (error) {
    console.error("Ошибка при создании задачи:", error);
    alert("Не удалось создать задачу: " + error.message);
  }
}

taskInput.addEventListener("keydown", async function (event) {
  if (event.key === "Enter") {
    await createTask();
  }
});

taskButton.addEventListener("click", async function () {
  await createTask();
});

function update() {
  document.addEventListener("DOMContentLoaded", async () => {
    try {
      const lists = await getAllLists();
      console.log("Получили листы:", lists);
      if (lists && Array.isArray(lists)) {
        lists.forEach((list) => renderSingleList(list.id, list.title));
      } else {
        console.error("Не удалось загрузить листы.");
      }

      const listId = details.getAttribute("data-id"); // Получаем ID текущего списка
      if (listId) {
        const tasks = await getAllTasks(listId);
        if (tasks && Array.isArray(tasks)) {
          tasks.forEach((task) => renderSingleTask(task, listId));
        } else {
          console.error("Не удалось загрузить задачи.");
        }
      }
    } catch (error) {
      console.error("Ошибка при загрузке данных:", error);
    }
  });
}
