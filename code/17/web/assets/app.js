const loginForm = document.getElementById("login-form");
const chatForm = document.getElementById("chat-form");
const apiHostInput = document.getElementById("api-host");
const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");
const messageInput = document.getElementById("message");
const sendButton = document.getElementById("send-btn");
const chatLog = document.getElementById("chat-log");
const statusBadge = document.getElementById("status-badge");
const loginHint = document.getElementById("login-hint");
const chatHint = document.getElementById("chat-hint");

let token = "";
let busy = false;

if (!apiHostInput.value) {
  apiHostInput.value = window.location.origin;
}

lockChat(true);
updateStatus("idle", "未登录");

loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  if (busy) return;

  const host = normalizeHost(apiHostInput.value);
  if (!host) {
    setHint(loginHint, "请填写正确的服务地址。", true);
    return;
  }

  busy = true;
  updateStatus("busy", "登录中…");
  setHint(loginHint, "正在请求 token…");
  lockChat(true);

  try {
    const resp = await fetch(`${host}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: usernameInput.value.trim(),
        password: passwordInput.value,
      }),
    });
    if (!resp.ok) {
      const msg = await resp.text();
      throw new Error(`登录失败 (${resp.status}): ${msg || "unknown error"}`);
    }
    const data = await resp.json();
    token = data.token;
    updateStatus("ok", "已登录");
    setHint(loginHint, "获取 token 成功，开始聊天吧！");
    lockChat(false);
    messageInput.focus();
  } catch (err) {
    console.error(err);
    setHint(loginHint, err.message, true);
    updateStatus("idle", "未登录");
  } finally {
    busy = false;
  }
});

chatForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  if (!token) {
    setHint(chatHint, "请先登录获取 token。", true);
    return;
  }
  const message = messageInput.value.trim();
  if (!message || busy) return;

  appendMessage("user", message);
  messageInput.value = "";
  setHint(chatHint, "等待模型回应…");
  updateStatus("busy", "对话中…");
  lockChat(true);
  busy = true;

  try {
    await streamChat(message);
    updateStatus("ok", "已登录");
    setHint(chatHint, "完成。");
  } catch (err) {
    console.error(err);
    appendMessage("error", err.message);
    setHint(chatHint, err.message, true);
    updateStatus("idle", "需要重新登录");
    token = "";
  } finally {
    busy = false;
    lockChat(false);
    messageInput.focus();
  }
});

messageInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    chatForm.requestSubmit();
  }
});

function normalizeHost(input) {
  const trimmed = input.trim().replace(/\/+$/, "");
  if (!trimmed.startsWith("http")) {
    return "";
  }
  return trimmed;
}

function updateStatus(state, text) {
  statusBadge.textContent = text;
  statusBadge.className = "status-badge";
  if (state === "ok") {
    statusBadge.classList.add("status-ok");
  } else if (state === "busy") {
    statusBadge.classList.add("status-busy");
  } else {
    statusBadge.classList.add("status-idle");
  }
}

function lockChat(disabled) {
  messageInput.disabled = disabled;
  sendButton.disabled = disabled;
}

function setHint(node, text, isError = false) {
  node.textContent = text;
  node.style.color = isError ? "#b4341d" : "var(--muted)";
}

function appendMessage(role, text) {
  const container = document.createElement("div");
  container.className = `message ${role}`;

  const label = document.createElement("div");
  label.className = "role";
  label.textContent = labelText(role);
  container.appendChild(label);

  const bubble = document.createElement("div");
  bubble.className = "bubble";
  bubble.textContent = text;
  container.appendChild(bubble);

  chatLog.appendChild(container);
  chatLog.scrollTop = chatLog.scrollHeight;
  return bubble;
}

function labelText(role) {
  if (role === "user") return "你";
  if (role === "error") return "错误";
  return "模型";
}

async function streamChat(message) {
  const host = normalizeHost(apiHostInput.value);
  if (!host) {
    throw new Error("服务地址无效");
  }
  const resp = await fetch(`${host}/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ message }),
  });

  if (!resp.ok || !resp.body) {
    const msg = await resp.text();
    throw new Error(`请求失败 (${resp.status}): ${msg || "unknown error"}`);
  }

  const reader = resp.body.getReader();
  const decoder = new TextDecoder("utf-8");
  const target = appendMessage("model", "");
  let buffer = "";
  let currentEvent = "message";

  while (true) {
    const { value, done } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });
    const chunks = buffer.split("\n\n");
    buffer = chunks.pop() || "";
    for (const chunk of chunks) {
      for (const line of chunk.split("\n")) {
        if (line.startsWith("event: ")) {
          currentEvent = line.slice(7).trim() || "message";
          continue;
        }
        if (line.startsWith("data: ")) {
          const payload = line.slice(6);
          if (currentEvent === "error") {
            throw new Error(payload);
          }
          target.textContent += payload;
          chatLog.scrollTop = chatLog.scrollHeight;
        }
      }
      currentEvent = "message";
    }
  }
}
