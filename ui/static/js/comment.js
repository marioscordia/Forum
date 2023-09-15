document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById("form");
  const textInput = document.getElementById("textInput");

  form.addEventListener("submit", function (event) {
    if (textInput.value.trim() === "") {
      event.preventDefault(); // Prevent the form from submitting
    }
  });
});
