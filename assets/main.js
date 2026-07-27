let credentials = null;

const loginView = document.getElementById("loginView");
const dashboardView = document.getElementById("dashboardView");
const loginForm = document.getElementById("loginForm");
const whitelistForm = document.getElementById("whitelistForm");
const playerList = document.getElementById("playerList");
const playerCount = document.getElementById("playerCount");
const statusMessage = document.getElementById("statusMessage");
const loginError = document.getElementById("loginError");

function headers(json = false) {
    const result = { "X-Api-Username": credentials.username, "X-Api-Key": credentials.password };
    if (json) result["Content-Type"] = "application/json";
    return result;
}

async function request(path, options = {}) {
    const response = await fetch(path, { ...options, headers: { ...headers(Boolean(options.body)), ...options.headers } });
    const body = await response.json();
    if (!response.ok || body.success !== true) throw new Error(body.message || "Request failed");
    return body;
}

function showStatus(message, type = "success") {
    statusMessage.textContent = message;
    statusMessage.className = `status-message ${type}`;
    statusMessage.hidden = false;
}

function renderPlayers(players) {
    playerList.replaceChildren();
    playerCount.textContent = players.length;
    if (players.length === 0) {
        const empty = document.createElement("p");
        empty.className = "empty-state";
        empty.textContent = "No players are currently whitelisted.";
        playerList.append(empty);
        return;
    }
    for (const username of players) {
        const row = document.createElement("div");
        row.className = "player-row";
        const name = document.createElement("span");
        name.textContent = username;
        const remove = document.createElement("button");
        remove.type = "button";
        remove.className = "remove-button";
        remove.textContent = "Remove";
        remove.addEventListener("click", () => removePlayer(username, remove));
        row.append(name, remove);
        playerList.append(row);
    }
}

async function loadPlayers() {
    playerList.setAttribute("aria-busy", "true");
    try {
        const response = await request("/api/whitelist");
        renderPlayers(response.players || []);
    } catch (error) {
        renderPlayers([]);
        showStatus(error.message, "error");
    } finally {
        playerList.removeAttribute("aria-busy");
    }
}

async function removePlayer(username, button) {
    button.disabled = true;
    try {
        const response = await request("/api/whitelist/remove", { method: "POST", body: JSON.stringify({ username }) });
        showStatus(response.message || `${username} was removed.`);
        await loadPlayers();
    } catch (error) {
        showStatus(error.message, "error");
        button.disabled = false;
    }
}

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const form = new FormData(loginForm);
    credentials = { username: form.get("username"), password: form.get("password") };
    loginError.hidden = true;
    try {
        await request("/api/validate", { method: "POST" });
        loginView.hidden = true;
        dashboardView.hidden = false;
        await loadPlayers();
    } catch (error) {
        credentials = null;
        loginError.textContent = error.message;
        loginError.hidden = false;
    }
});

whitelistForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const input = document.getElementById("playerUsername");
    const username = input.value.trim();
    if (!/^[A-Za-z0-9_]{3,16}$/.test(username)) {
        showStatus("Use a Minecraft Java username: 3–16 letters, numbers, or underscores.", "error");
        return;
    }
    try {
        const response = await request("/api/whitelist/add", { method: "POST", body: JSON.stringify({ username }) });
        showStatus(response.message || `${username} was added.`);
        whitelistForm.reset();
        await loadPlayers();
    } catch (error) {
        showStatus(error.message, "error");
    }
});

document.getElementById("refreshButton").addEventListener("click", loadPlayers);
document.getElementById("signOutButton").addEventListener("click", () => {
    credentials = null;
    dashboardView.hidden = true;
    loginView.hidden = false;
    loginForm.reset();
    document.getElementById("loginUsername").focus();
});
