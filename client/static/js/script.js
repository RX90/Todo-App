document.addEventListener("DOMContentLoaded", function () {
  var addListBtn = document.getElementById("addListBtn");
  var newListTitle = document.getElementById("newListTitle");
  var listsContainer = document.getElementById("listsContainer");

  var modal = document.getElementById("myModal");
  var modalMessage = document.getElementById("modalMessage");
  var span = document.getElementsByClassName("close")[0];

  span.onclick = function () {
    modal.style.display = "none";
  };

  window.onclick = function (event) {
    if (event.target == modal) {
      modal.style.display = "none";
    }
  };

  function showModalMessage(message) {
    modalMessage.textContent = message;
    modal.style.display = "block";
    setTimeout(() => {
      modal.style.display = "none";
    }, 2000);
  }

  async function loadLists() {
    var accessToken = localStorage.getItem("accessToken");

    try {
      const response = await fetch("/api/lists/", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + accessToken,
        },
      });

      if (response.ok) {
        const result = await response.json();
        const lists = result.data;

        listsContainer.innerHTML = "";
        lists.forEach((list) => {
          var listItem = document.createElement("li");
          listItem.textContent = list.title;
          listItem.classList.add("list-item");
          listsContainer.appendChild(listItem);
        });
      } else if (response.status === 401) {
        showModalMessage(
          "Please log in to view your lists. Redirecting to login page..."
        );
        setTimeout(() => {
          window.location.href = "/auth/sign-in";
        }, 2000);
      }
    } catch (error) {
      console.error("There was a problem with the fetch operation:", error);
    }
  }

  loadLists();

  addListBtn.addEventListener("click", async function () {
    var title = newListTitle.value.trim();

    if (title) {
      var accessToken = localStorage.getItem("accessToken");
      var data = {
        title: title,
      };

      try {
        const response = await fetch("/api/lists/", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + accessToken,
          },
          body: JSON.stringify(data),
        });

        newListTitle.value = "";

        if (response.status === 401) {
          showModalMessage(
            "Please log in to add a list. Redirecting to login page..."
          );
          setTimeout(() => {
            window.location.href = "/auth/sign-in";
          }, 2000);
        } else if (response.status === 200) {
          var listItem = document.createElement("li");
          listItem.textContent = title;
          listItem.classList.add("list-item");
          listsContainer.appendChild(listItem);
        }
      } catch (error) {
        console.error("There was a problem with the fetch operation:", error);
      }
    } else {
      showModalMessage("Please fill in both fields.");
    }
  });
});
