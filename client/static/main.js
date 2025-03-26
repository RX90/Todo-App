let createList = document.getElementById("new-list");
let menu = document.querySelector(".menu");
let details = document.getElementById("details");
let title = document.getElementById("title");
let listIdCounter = 0;

let plusList = document.querySelector(".plus-icon");

let loginButton = document.getElementById("login-button-signin");
let registerButton = document.getElementById("login-button-signup");
let logoutButton = document.getElementById("login-button-logout");

let infoText = document.getElementById("info-text");

let panelSignin = document.querySelector(".background-sign-in");
let panelSignUp = document.querySelector(".background-sign-up");

let signinSendData = document.getElementById("signin-button");
let signupSendData = document.getElementById("signup-button");

let usernameLogin = document.querySelector(".signin-name-input");
let passwordLogin = document.querySelector(".signin-password-input");

let passwordRegister = document.querySelector(".signup-password-input");

let moveToRegistr = document.getElementById("move-to-registr");

let activePanel = null;

let signinViewSvg = document.querySelector(".signin-view-svg");
let signinHiddenSvg = document.querySelector(".signin-hide-svg");

let signupViewSvg = document.querySelector(".signup-view-svg");
let signupHiddenSvg = document.querySelector(".signup-hide-svg");

// Скрыть/Показать пароль у логина
signinViewSvg.addEventListener("click", function () {
  passwordLogin.type = "text";
  signinViewSvg.style.display = "none";
  signinHiddenSvg.style.display = "block";
});

signinHiddenSvg.addEventListener("click", function () {
  passwordLogin.type = "password";
  signinHiddenSvg.style.display = "none";
  signinViewSvg.style.display = "block";
});

//Скрыть/Показать пароль у регистрации
signupViewSvg.addEventListener("click", function () {
  passwordRegister.type = "text";
  signupViewSvg.style.display = "none";
  signupHiddenSvg.style.display = "block";
});

signupHiddenSvg.addEventListener("click", function () {
  passwordRegister.type = "password";
  signupHiddenSvg.style.display = "none";
  signupViewSvg.style.display = "block";
});

document.addEventListener("keydown", function (event) {
  if (event.key === "Escape") {
    hiddenPopupSignin();
    hiddenPopupSignUp();
  }
});

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
  loginButton.style.display = "none";
  logoutButton.style.display = "block";

  logoutButton.addEventListener("click", async function () {
    await logout();

    location.reload();
  });
}

//Логинизация
function showPopupSignin() {
  panelSignin.style.display = "block";
  usernameLogin.value = "";
  passwordLogin.value = "";

  usernameLogin.style.outline = "";
  passwordLogin.style.outline = "";

  signinError.textContent = "";
  signinError2.textContent = "";
}

function hiddenPopupSignin() {
  panelSignin.style.display = "none";
}

loginButton.addEventListener("click", function () {
  showPopupSignin();
});

signinSendData.addEventListener("click", async function () {
  const user = usernameLogin.value.trim();
  const pass = passwordLogin.value.trim();

  const success = await signIn(user, pass);
  if (success) {
    console.log("Вход прошел успешно");
    hiddenPopupSignin();

    clearRenderedLists();
    const lists = await getAllLists();
    if (lists && Array.isArray(lists)) {
      lists.forEach((list) => renderSingleList(list.id, list.title));
    }
  } else {
    console.error("Ошибка входа! Проверьте логин и пароль.");
  }
});

moveToRegistr.addEventListener("click", function () {
  hiddenPopupSignin();

  showPopupSignUp();
});

//Регистрация
function showPopupSignUp() {
  panelSignUp.style.display = "block";
}

function hiddenPopupSignUp() {
  panelSignUp.style.display = "none";
}

registerButton.addEventListener("click", function () {
  showPopupSignUp();
});

signupSendData.addEventListener("click", async function () {
  const user = usernameRegister.value;
  const pass = passwordRegister.value;

  // let letterCheck = document.getElementById("letter-check");
  // let numberCheck = document.getElementById("number-check");
  // let lengthCheck = document.getElementById("length-check");

  let letterLabel = document.getElementById("letter-label");
  let letterCheckbox = document.getElementById("checkbox-letter");

  let numberLabel = document.getElementById("number-label");
  let numberCheckbox = document.getElementById("checkbox-number");

  let lengthLabel = document.getElementById("length-label");
  let lengthCheckbox = document.getElementById("checkbox-length");

  let errorMessage = document.getElementById("error-message");

  let isValid = true;

  errorUser.textContent = "";
  errorMessage.textContent = "";
  passwordRegister.style.outline = "";

  if (/[а-яА-ЯёЁ]/.test(pass)) {
    console.log("Пароль содержит русские буквы!");
    errorMessage.textContent =
      "Пароль должен содержать только английские буквы"; /*!import*/
    isValid = false;
    passwordRegister.style.outline = "3px solid red";
  }

  if (user.length < 3 || user.length > 32) {
    console.log("Имя от 3 до 32 символов");
    errorUser.textContent = "Логин должен содежать минимум 3 символа";
    isValid = false;
    usernameRegister.style.outline = "3px solid red";
  }

  if (/[а-яА-ЯёЁ]/.test(user)) {
    console.log("Логин содержит русские буквы!");
    errorUser.textContent =
      "Логин должен содержать только латинские символы, цифры, дефисы (-) и нижние подчёркивания (_)";
    isValid = false;
    usernameRegister.style.outline = "3px solid red";
  }

  if (!isValid) {
    return; // Если есть ошибки, прекращаем выполнение
  }

  if (/[\d]/.test(pass)) {
    numberLabel.style.color = "white";
    numberCheckbox.src = "../src/img/violet-checkbox.svg";
  } else {
    numberLabel.style.color = "red";
    errorMessage.textContent = "Добавьте еще одну цифру";
    isValid = false;
    passwordRegister.style.outline = "3px solid red";
  }

  if (/[a-z]/.test(pass) && /[A-Z]/.test(pass)) {
    letterLabel.style.color = "white";
    letterCheckbox.src = "../src/img/violet-checkbox.svg";
  } else {
    letterLabel.style.color = "red";
    errorMessage.textContent = "Забыл про большую и маленькую букву";
    isValid = false;
    passwordRegister.style.outline = "3px solid red";
  }

  if (pass.length >= 8 && pass.length <= 32) {
    lengthLabel.style.color = "white";
    lengthCheckbox.src = "../src/img/violet-checkbox.svg";
  } else {
    lengthLabel.style.color = "red";
    errorMessage.textContent = "Пароль должен содежать минимум 8 символов";
    isValid = false;
    passwordRegister.style.outline = "3px solid red";
  }

  if (!isValid) {
    if (!errorMessage.textContent) {
      errorMessage.textContent = "";
    }
    if (!errorUser.textContent) {
      errorUser.textContent = "";
    }
    return;
  }

  const success = await signUp(user, pass);
  if (success) {
    console.log("Регистрация прошла успешно");
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const signInSuccess = await signIn(user, pass);
    if (signInSuccess) {
      console.log("Вход выполнен");
      hiddenPopupSignUp();
    }

    clearRenderedLists();
    const lists = await getAllLists();
    if (lists && Array.isArray(lists)) {
      lists.forEach((list) => renderSingleList(list.id, list.title));
    }
  } else {
    console.error("Ошибка входа! Проверьте логин и пароль.");
  }
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
      addList.setSelectionRange(addList.value.length, addList.value.length);

      addList.addEventListener("keydown", async function (event) {
        if (event.key === "Enter") {
          const newTitleList = addList.value.trim();
          const maxLength = 32;

          const existingList = [...menu.querySelectorAll(".add-list")]
            .filter((list) => list !== addList)
            .map((list) => list.value.trim().toLowerCase())
            .includes(newTitleList.toLowerCase());

          if (existingList) {
            addList.value = title;
            addList.disabled = true;
            event.preventDefault();
            dotsPanel.remove();
            return;
          }

          if (newTitleList !== "" && newTitleList !== title) {
            if (newTitleList.length > maxLength) {
              newTitleList = newTitleList.substring(0, maxLength);
              addList.value = newTitleList;
              console.log("Больше 32 символов");
            }
            await EditList(listid, newTitleList);
            title = newTitleList;
          }
          addList.disabled = true;
          dotsPanel.remove();
        }
      });
    });

    const place = dots.getBoundingClientRect();
    dotsPanel.style.top = `${place.top + window.scrollY - 100}px`;
    dotsPanel.style.left = `${place.left + window.scrollX + 30}px`;

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
  const maxLength = 65;

  const editTask = document.createElement("img");
  editTask.src = "/src/img/violet-edit.svg";
  editTask.classList.add("edit-task");

  editTask.addEventListener("click", async function (event) {
    titleTask.disabled = false; // Даем возможность редактировать
    titleTask.focus(); // Фокус на поле ввода

    titleTask.addEventListener("input", function () {
      // Ограничение по длине названия задачи
      if (titleTask.value.length > maxLength) {
        titleTask.value = titleTask.value.substring(0, maxLength);
      }
    });

    titleTask.addEventListener("keydown", async function (event) {
      const newTitle = titleTask.value.trim();

      if (event.key === "Enter" && newTitle !== "") {
        const existingTask = [...taskList.querySelectorAll(".title-task")]
          .filter((input) => input !== titleTask)
          .map((input) => input.value.trim().toLowerCase())
          .includes(newTitle.toLowerCase());

        if (existingTask) {
          titleTask.value = task.title;
          titleTask.disabled = true;
          event.preventDefault();
          return;
        }

        const listId = details.getAttribute("data-id");
        const taskId = Number(menuTask.getAttribute("data-task-id"));

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
      menuTask.remove();
      setTimeout(() => {
        if (completedList.children.length === 0) {
          document.querySelector(".title-completed").style.display = "none";
        }
      }, 0);
    } else {
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
    let maxList = 10;

    const currentListsCount = document.querySelectorAll(".menu-item").length;
    if (currentListsCount > maxList) {
      console.log("Достигнуто максимальное количество листов");
      return;
    }

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
        createList.value = "";
      }
    } catch (error) {
      console.error("Ошибка при создании листа:", error);
      createList.preventDefault();
    }
    renderSingleList(newTitleList);
  }
});

plusList.addEventListener("click", async function () {
  const title = createList.value.trim();

  let maxList = 10;

  const currentListsCount = document.querySelectorAll(".menu-item").length;
  if (currentListsCount > maxList) {
    console.log("Достигнуто максимальное количество листов");
    createList.preventDefault();
    return;
  }

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
      createList.value = "";
    }
  } catch (error) {
    console.error("Ошибка при создании листа:", error);
    createList.preventDefault();
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
      taskInput.value = "";
    }
  } catch (error) {
    console.error("Ошибка при создании задачи:", error);
    taskInput.preventDefault();
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

//Фикс бага с дюпом листов
function clearRenderedLists() {
  const menu = document.querySelector(".menu");
  if (menu) {
    const newListInputContainer =
      menu.querySelector("#new-list")?.parentElement;

    menu.innerHTML = "";

    if (newListInputContainer) {
      menu.appendChild(newListInputContainer);
    }
  }

  const details = document.getElementById("details");
  if (details) details.style.display = "none";
}

document.addEventListener("DOMContentLoaded", async () => {
  if (localStorage.getItem("accessToken")) {
    console.log("Токен найден, загружаем списки...");

    clearRenderedLists();
    const lists = await getAllLists();
    if (lists && Array.isArray(lists)) {
      lists.forEach((list) => {
        renderSingleList(list.id, list.title);
        getAllTasks(list.id);
      });
    }
  } else {
    console.log("Токен не найден, показываем окно входа.");
  }
});
