document
  .getElementById("signin-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const data = {
      username: username,
      password: password,
    };

    try {
      const response = await fetch("/auth/sign-in", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      const result = await response.json();

      if (response.ok) {
        document.getElementById("message").textContent = "Sign In successful!";
        document.getElementById("message").style.color = "green";
        document.getElementById("username").value = "";
        document.getElementById("password").value = "";
        localStorage.setItem("bearerToken", result.token);
      } else {
        document.getElementById("message").textContent =
          "Sign In failed: " + result.message;
      }
    } catch (error) {
      document.getElementById("message").textContent =
        "An error occurred: " + error.message;
    }
  });

function show_hide_password(target) {
  var input = document.getElementById("password-input");
  if (input.getAttribute("type") === "password") {
    target.classList.add("view");
    input.setAttribute("type", "text");
  } else {
    target.classList.remove("view");
    input.setAttribute("type", "password");
  }
  return false;
}

// fetch("/куда-то там", {
//   headers: {
//     Authorization: "Bearer " + localStorage.getItem("bearerToken"),
//   },
// });
