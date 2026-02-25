// JavaScript handlers for the todo app
// perform actions via fetch and update the UI without reloading

document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector(".todo-form");
    if (form) {
        const inputField = form.querySelector("input[name='content']");
        if (inputField) {
            // ensure futuristic styling via JS in case CSS caching is an issue
            inputField.classList.add('futuristic-input');
        }

        form.addEventListener("submit", async (e) => {
            e.preventDefault();
            const input = form.querySelector("input[name='content']");
            const content = input && input.value.trim();
            if (!content) {
                return; // nothing to add
            }
            try {
                await fetch("/add", {
                    method: "POST",
                    headers: { "Content-Type": "application/x-www-form-urlencoded" },
                    body: new URLSearchParams({ content }),
                });
            } catch (err) {
                console.error("add error", err);
            }
            // clear input and reload list (simple approach)
            input.value = "";
            location.reload();
        });
    }

    // toggle and delete buttons
    document.querySelectorAll(".btn-check").forEach(btn => {
        btn.addEventListener("click", async () => {
            const id = btn.dataset.id;
            if (!id) return;
            btn.disabled = true;
            try {
                await fetch("/toggle", {
                    method: "POST",
                    headers: { "Content-Type": "application/x-www-form-urlencoded" },
                    body: new URLSearchParams({ id }),
                });
            } catch (err) {
                console.error("toggle error", err);
            }
            location.reload();
        });
    });

    document.querySelectorAll(".btn-delete").forEach(btn => {
        btn.addEventListener("click", async () => {
            const id = btn.dataset.id;
            if (!id) return;
            btn.disabled = true;
            try {
                await fetch("/delete", {
                    method: "POST",
                    headers: { "Content-Type": "application/x-www-form-urlencoded" },
                    body: new URLSearchParams({ id }),
                });
            } catch (err) {
                console.error("delete error", err);
            }
            location.reload();
        });
    });

    // edit buttons
    document.querySelectorAll(".btn-edit").forEach(btn => {
        btn.addEventListener("click", async () => {
            const id = btn.dataset.id;
            if (!id) return;
            const currentText = btn.closest('li').querySelector('.task-text').textContent;
            const newText = prompt("Modifier la mission:", currentText);
            if (newText === null || newText.trim() === "") return;
            btn.disabled = true;
            try {
                await fetch("/edit", {
                    method: "POST",
                    headers: { "Content-Type": "application/x-www-form-urlencoded" },
                    body: new URLSearchParams({ id, content: newText.trim() }),
                });
            } catch (err) {
                console.error("edit error", err);
            }
            location.reload();
        });
    });
});
