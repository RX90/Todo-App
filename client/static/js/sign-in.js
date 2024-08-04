document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("signin-form");
  const messageDiv = document.getElementById("message");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
      const response = await fetch("/auth/sign-in", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.status === 200) {
        const token = data.token;

        localStorage.setItem("accessToken", token);

        window.location.href = "/";
      } else {
        messageDiv.textContent = "Sign-in failed: " + data.message;
        messageDiv.style.color = "red";
      }
    } catch (error) {
      messageDiv.textContent = "Sign-in failed: " + error.message;
      messageDiv.style.color = "red";
    }
  });
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
