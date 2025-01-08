let createList = document.getElementById("new-list");
let menu = document.querySelector(".menu");
let details = document.getElementById("details");
let detailsText = document.getElementById("details-text");
let listIdCounter = 0;

createList.addEventListener("keydown", function (event) {
  if (event.key === "Enter") {
    console.log("Лист создан:" + createList.value);

    let menuItem = document.createElement("div");
    menuItem.classList.add("menu-item");
    menuItem.setAttribute("data-id", listIdCounter);

    let imgList = document.createElement("img");
    imgList.classList.add("icon-list");
    imgList.setAttribute("src", "/client/src/img/list.svg");

    let addList = document.createElement("p");
    addList.classList.add("add-list");
    addList.textContent = createList.value;

    menuItem.appendChild(imgList);
    menuItem.appendChild(addList);
    menu.appendChild(menuItem);

    openPanel(listIdCounter, createList.value);

    menuItem.addEventListener("click", function () {
      const listId = menuItem.getAttribute("data-id");
      const listName = menuItem.querySelector(".add-list").textContent;
      openPanel(listId, listName);
    });
    createList.value = "";
    listIdCounter++;
    console.log(listIdCounter);
  }
});

function openPanel(listId, listName) {
  details.style.display = "block";
  details.setAttribute("data-id", listId);
  detailsText.textContent = listName;
}
