// JavaScript handlers for the todo app
// perform actions via fetch and update the UI without reloading

document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector(".todo-form");
    if (form) {
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
});
