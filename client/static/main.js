let createList = document.getElementById("new-list");
let menu = document.querySelector(".menu");
let details = document.getElementById("details");
let title = document.getElementById("title");
let listIdCounter = 0;
let popupWindow = document.getElementById("background");
let popupButton = document.getElementById("window-button");

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
  hiddenPopup();
  console.log("Имя: ", username.value, "Пароль:", password.value);
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
  });
}

function renderSingleTask(task) {
  const menuTask = document.createElement("div");
  menuTask.classList.add("menu-task");

  const circleIcon = document.createElement("img");
  circleIcon.src = "/src/img/circle.svg";
  circleIcon.classList.add("circle-icon");

  const titleTask = document.createElement("span");
  titleTask.textContent = taskTitle;
  titleTask.classList.add("title-task");

  if (task.done) {
    titleTask.classList.add("task-done");
  }

  menuTask.appendChild(circleIcon);
  menuTask.appendChild(titleTask);

  panel.appendChild(menuTask);
}

//Создание листов
createList.addEventListener("keydown", async function (event) {
  if (event.key === "Enter" && createList.value.trim() !== "") {
    const title = createList.value.trim();

    try {
      const newList = await sendList(title);

      if (newList) {
        const menuItem = document.createElement("div");
        menuItem.classList.add("menu-item");
        menuItem.setAttribute("data-id", newList.id);

        const imgList = document.createElement("img");
        imgList.classList.add("icon-list");
        imgList.setAttribute("src", "/src/img/list.svg");

        const addList = document.createElement("p");
        addList.classList.add("add-list");
        addList.textContent = newList.Title;

        menuItem.appendChild(imgList);
        menuItem.appendChild(addList);

        menu.appendChild(menuItem);

        openPanel(newList.id, newList.Title);

        menuItem.addEventListener("click", function () {
          const listId = menuItem.getAttribute("data-id");
          const listName = menuItem.querySelector(".add-list").textContent;
          openPanel(listId, listName);
        });
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
}

const taskInput = document.getElementById("task-input");
const taskButton = document.getElementById("task-button");
const panel = document.querySelector(".panel");

//Создание задач
function createTask() {
  const menuTask = document.createElement("div");
  menuTask.classList.add("menu-task");

  const circleIcon = document.createElement("img");
  circleIcon.src = "/src/img/circle.svg";
  circleIcon.classList.add("circle-icon");

  const titleTask = document.createElement("span");
  titleTask.textContent = taskInput.value;
  titleTask.classList.add("title-task");

  menuTask.appendChild(circleIcon);
  menuTask.appendChild(titleTask);

  panel.appendChild(menuTask);
}

taskInput.addEventListener("keydown", async function (event) {
  if (event.key === "Enter" && taskInput.value.trim() !== "") {
    const taskTitle = createTask.value.trim();

    try {
      const newTask = await sendTask();

      if (newTask) {
        createTask();
      }
    } catch (error) {
      console.error("Ошибка при создании задачи:", error);
      alert("Не удалось создать задачу: " + error.message);
    }
  }
});

//Рендер листов
document.addEventListener("DOMContentLoaded", async () => {
  const lists = await getAllList();
  if (lists) renderLists(lists);
});
