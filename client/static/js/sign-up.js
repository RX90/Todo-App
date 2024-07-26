document
  .getElementById("signup-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const name = document.getElementById("name").value;
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const data = {
      name: name,
      username: username,
      password: password,
    };

    try {
      const response = await fetch("/auth/sign-up", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      const result = await response.json();

      if (response.ok) {
        document.getElementById("message").textContent = "Registration successful!";
        document.getElementById("message").style.color = "green";
        document.getElementById("name").value = "";
        document.getElementById("username").value = "";
        document.getElementById("password").value = "";
      } else {
        document.getElementById("message").textContent =
          "Registration failed: " + result.message;
      }
    } catch (error) {
      document.getElementById("message").textContent =
        "An error occurred: " + error.message;
    }
  });

function show_hide_password(target) {
  var input = document.getElementById("password");
  if (input.getAttribute("type") == "password") {
    target.classList.add("view");
    input.setAttribute("type", "text");
  } else {
    target.classList.remove("view");
    input.setAttribute("type", "password");
  }
  return false;
}
