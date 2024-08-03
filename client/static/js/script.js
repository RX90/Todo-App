function validateToken() {
  fetch("/auth/validate", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + token,
    },
  })
    .then((response) => {
      if (response.status === 401) {
        window.location.href = "/auth/sign-in";
      }
    })
    .catch((error) => {
      console.error("Ошибка при проверке токена:", error);
    });
}

const token = localStorage.getItem("token");

validateToken();
