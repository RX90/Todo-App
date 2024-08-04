document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("signup-form");
  const messageDiv = document.getElementById("message");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const name = document.getElementById("name").value;
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
      const signUpResponse = await fetch("/auth/sign-up", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, username, password }),
      });

      const signUpData = await signUpResponse.json();

      if (signUpResponse.status === 200) {
        try {
          const signInResponse = await fetch("/auth/sign-in", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password }),
          });

          const signInData = await signInResponse.json();

          if (signInResponse.status === 200) {
            const token = signInData.token;

            localStorage.setItem("accessToken", token);

            window.location.href = "/";
          } else {
            messageDiv.textContent = "Sign-in failed: " + signInData.message;
            messageDiv.style.color = "red";
          }
        } catch (signInError) {
          messageDiv.textContent = "Sign-in failed: " + signInError.message;
          messageDiv.style.color = "red";
        }
      } else {
        messageDiv.textContent = "Sign-up failed: " + signUpData.message;
        messageDiv.style.color = "red";
      }
    } catch (signUpError) {
      messageDiv.textContent = "Sign-up failed: " + signUpError.message;
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
