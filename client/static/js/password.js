document.addEventListener("DOMContentLoaded", function () {
  var togglePassword = document.getElementById("toggle-password");
  var passwordInput = document.getElementById("password");

  if (togglePassword && passwordInput) {
    togglePassword.addEventListener("click", function (event) {
      event.preventDefault();
      if (passwordInput.getAttribute("type") === "password") {
        passwordInput.setAttribute("type", "text");
        togglePassword.classList.add("view");
      } else {
        passwordInput.setAttribute("type", "password");
        togglePassword.classList.remove("view");
      }
    });
  }
});
