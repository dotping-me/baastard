const app = document.querySelector(".app");
const toggleSQLEditor = document.querySelector("#toggle-workspace");

toggleSQLEditor.addEventListener("click", () => {
    app.classList.toggle("workspace-hidden");
});