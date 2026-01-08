document.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById('magic-btn');
    const body = document.getElementById('main-body');

    btn.addEventListener('click', () => {
        body.classList.toggle('active-mode');
        if (body.classList.contains('active-mode')) {
            btn.innerText = "Retourner sur Terre";
        } else {
            btn.innerText = "Activer le Cosmos";
        }
    });
});