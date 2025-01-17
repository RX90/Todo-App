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

//Проверка существует ли у пользователя токен, если нет, мы отправляем его регаться или логиниться
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

function showPopup() {
  popupWindow.style.display = "block";
}

function hiddenPopup() {
  popupWindow.style.display = "none";
}

popupButton.addEventListener("click", function () {
  if (popupButton.textContent === "Sign Up!") {
    console.log("Регистрация пользователя:", username.value);
    signUp(username.value, password.value); // Вызываем функцию регистрации
    signIn(username.value, password.value);
  } else {
    console.log("Вход пользователя:", username.value);
    signIn(username.value, password.value); // Вызываем функцию входа
  }
  if (error || error.status === 401) {
    popupButton.disabled = true;
  } else {
    hiddenPopup();
  }
});

createAccount.addEventListener("click", function () {
  pupopTitle.textContent = "Sign-Up";
  infoText.textContent = "Please fill in the fields to create an account";
  popupButton.textContent = "Sign Up!";
  createAccount.style.display = "none";
});

function renderSingleList(list) {
  const menuItem = document.createElement("div");
  menuItem.classList.add("menu-item");
  menuItem.setAttribute("data-id", list.id);

  const imgList = document.createElement("img");
  imgList.classList.add("icon-list");
  imgList.setAttribute("src", "/src/img/list.svg");

  const addList = document.createElement("p");
  addList.classList.add("add-list");
  addList.textContent = list.title;

  menuItem.appendChild(imgList);
  menuItem.appendChild(addList);
  menu.appendChild(menuItem);

  menuItem.addEventListener("click", () => {
    openPanel(list.id, list.title);
    console.log(list.id, list.title);
  });
}

function renderSingleTask(task) {
  const taskList = document.querySelector(".task-list");
  const menuTask = document.createElement("div");
  menuTask.classList.add("menu-task");

  const circleIcon = document.createElement("img");
  circleIcon.src = "/src/img/circle.svg";
  circleIcon.classList.add("circle-icon");

  const titleTask = document.createElement("span");
  titleTask.textContent = task.title;
  titleTask.classList.add("title-task");

  menuTask.appendChild(circleIcon);
  menuTask.appendChild(titleTask);

  taskList.appendChild(menuTask);
}

//Создание листов
createList.addEventListener("keydown", async function (event) {
  if (event.key === "Enter" && createList.value.trim() !== "") {
    const title = createList.value.trim();

    try {
      const newList = await sendList(title);

      if (newList) {
        renderSingleList(newList);
        openPanel(newList.id, title);
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
  taskList.innerHTML = ""; // Очищаем только задачи

  getAllTasks(listId);
}

const taskInput = document.getElementById("task-input");
const taskButton = document.getElementById("task-button");

taskInput.addEventListener("keydown", async function (event) {
  if (event.key === "Enter" && taskInput.value.trim() !== "") {
    const taskTitle = taskInput.value.trim();
    const listId = details.getAttribute("data-id");
    try {
      const newTaskId = await sendTask(listId, taskTitle);
      console.log("Title:", taskTitle);

      if (newTaskId) {
        const newTask = {
          id: newTaskId,
          title: taskTitle,
        };

        renderSingleTask(newTask);
        taskInput.value = "";
      }
    } catch (error) {
      console.error("Ошибка при создании задачи:", error);
      alert("Не удалось создать задачу: " + error.message);
    }
  }
});

document.addEventListener("DOMContentLoaded", async () => {
  try {
    const lists = await getAllLists();
    if (lists && Array.isArray(lists)) {
      lists.forEach(renderSingleList);
    } else {
      console.error("Не удалось загрузить листы.");
    }

    const tasks = await getAllTasks(listId);
    if (tasks && Array.isArray(tasks)) {
      tasks.forEach(renderSingleTask);
    } else {
      console.error("Не удалось загрузить задачи.");
    }
  } catch (error) {
    console.error("Ошибка при загрузке данных:", error);
  }
});
