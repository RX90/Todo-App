let createList = document.getElementById("new-list");
let menu = document.querySelector(".menu");
let details = document.getElementById("details");
let title = document.getElementById("title");
let listIdCounter = 0;
let popupWindow = document.getElementById("background");
let popupButton = document.getElementById("window-button");
let createAccount = document.getElementById("create-account");
let pupopTitle = document.getElementById("pupop-title");
let infoText = document.getElementById("info-text");
const panel = document.querySelector(".panel");
const loginButton = document.getElementById("login-button");

if (
  !localStorage.getItem("accessToken") ||
  localStorage.getItem("accessToken").trim() === ""
) {
  console.log("Токен не найден");
  loginButton.addEventListener("click", function () {
    showPopup();
  });
} else {
  console.log("Токен получен", localStorage.getItem("accessToken"));
  hiddenPopup();
}

async function logoutLocalStorage() {
  loginButton.textContent = "Выйти";
  loginButton.addEventListener("click", async function () {
    await logout();
    location.reload();
  });
}

function showPopup() {
  popupWindow.style.display = "block";
}

function hiddenPopup() {
  popupWindow.style.display = "none";
}

popupButton.addEventListener("click", async function () {
  try {
    if (popupButton.textContent === "Sign Up!") {
      console.log("Регистрация пользователя:", username.value);
      await signUp(username.value, password.value);
      await signIn(username.value, password.value);
      location.reload();
    } else {
      console.log("Вход пользователя:", username.value);
      await signIn(username.value, password.value);
      location.reload();
    }
    hiddenPopup();
  } catch (error) {
    console.error("Ошибка:", error.message);
    alert("Не удалось войти или зарегистрироваться: " + error.message);
    popupButton.disabled = true;
  }
});

createAccount.addEventListener("click", function () {
  pupopTitle.textContent = "Sign-Up";
  infoText.textContent = "Please fill in the fields to create an account";
  popupButton.textContent = "Sign Up!";
  createAccount.style.display = "none";
});

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
  addList.readOnly = true;

  const dots = document.createElement("img");
  dots.src = "/src/img/dots.svg";
  dots.classList.add("dots");

  dots.addEventListener("click", (event) => {
    event.stopPropagation(); // Останавливаем всплытие
    const dotsPanel = document.createElement("div");
    dotsPanel.classList.add("dots-panel");

    const dotsEdit = document.createElement("p");
    dotsEdit.textContent = "Edit";
    dotsEdit.classList.add("dots-edit");

    dotsPanel.appendChild(dotsEdit);
    document.body.appendChild(dotsPanel);
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

function renderSingleTask(task) {
  const taskList = document.querySelector(".task-list");
  const completedList = document.querySelector(".completed");

  const menuTask = document.createElement("div");
  menuTask.classList.add("menu-task");
  menuTask.setAttribute("data-task-id", task.id);

  const circleIcon = document.createElement("img");
  circleIcon.src = task.done
    ? "/src/img/color-circle.svg"
    : "/src/img/circle.svg";
  circleIcon.classList.add("circle-icon");

  const titleTask = document.createElement("input");
  titleTask.value = task.title;
  titleTask.classList.add("title-task");
  titleTask.disabled = true;

  const editTask = document.createElement("img");
  editTask.src = "/src/img/edit.svg";
  editTask.classList.add("edit-task");

  editTask.addEventListener("click", async function (event) {
    titleTask.disabled = false;
    titleTask.style.border = "1px solid white";

    titleTask.addEventListener("keydown", async function (event) {
      if (event.key === "Enter" && titleTask.value.trim() !== "") {
        const listId = details.getAttribute("data-id");
        const taskId = Number(menuTask.getAttribute("data-task-id"));
        const newTitle = titleTask.value.trim();

        await EditTask(taskId, listId, newTitle);

        titleTask.disabled = true;
        titleTask.style.border = "none";
      }
    });
  });

  const deleteTask = document.createElement("img");
  deleteTask.src = "/src/img/delete.svg";
  deleteTask.classList.add("delete-task");

  deleteTask.addEventListener("click", async function () {
    const listId = details.getAttribute("data-id");
    const taskId = Number(menuTask.getAttribute("data-task-id"));
    await DeleteTask(listId, taskId);
    menuTask.remove();
  });

  function moveToCorrectPlace() {
    menuTask.remove();
    const listId = details.getAttribute("data-id");

    if (task.done) {
      const titleCompleted = document.querySelector(".title-completed");
      titleCompleted.style.display = "block";
      completedList.appendChild(menuTask);
      titleTask.style.textDecoration = "line-through";
    } else {
      taskList.appendChild(menuTask);
      titleTask.style.textDecoration = "none";
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
        circleIcon.src = newState
          ? "/src/img/color-circle.svg"
          : "/src/img/circle.svg";
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
      showPopup();
    } else {
      console.log("Токен получен", localStorage.getItem("accessToken"));
      hiddenPopup();
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
