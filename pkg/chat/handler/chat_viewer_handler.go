package chat_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const chatViewerHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Evolution GO — Chat Viewer</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#111b21;color:#e9edef;height:100vh;display:flex;flex-direction:column}
#header{background:#202c33;padding:12px 20px;display:flex;align-items:center;gap:12px;border-bottom:1px solid #2a3942}
#header h1{font-size:16px;font-weight:600;color:#e9edef}
#header span{font-size:12px;color:#8696a0}
#app{display:flex;flex:1;overflow:hidden}
#sidebar{width:320px;background:#111b21;border-right:1px solid #2a3942;display:flex;flex-direction:column}
#sidebar-header{padding:12px 16px;background:#202c33;border-bottom:1px solid #2a3942}
#sidebar-header select{width:100%;background:#2a3942;color:#e9edef;border:none;padding:8px 10px;border-radius:6px;font-size:13px;cursor:pointer}
#sidebar-header select option{background:#202c33}
#chat-list{flex:1;overflow-y:auto}
.chat-item{display:flex;align-items:center;padding:12px 16px;cursor:pointer;border-bottom:1px solid #2a3942;gap:10px;transition:background .15s}
.chat-item:hover,.chat-item.active{background:#2a3942}
.chat-avatar{width:42px;height:42px;border-radius:50%;background:#00a884;display:flex;align-items:center;justify-content:center;font-size:16px;font-weight:600;color:#fff;flex-shrink:0}
.chat-info{flex:1;min-width:0}
.chat-name{font-size:14px;font-weight:500;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.chat-preview{font-size:12px;color:#8696a0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;margin-top:2px}
.chat-time{font-size:11px;color:#8696a0;flex-shrink:0;align-self:flex-start;margin-top:2px}
#main{flex:1;display:flex;flex-direction:column;background:#0b141a}
#chat-header{background:#202c33;padding:12px 16px;display:flex;align-items:center;gap:10px;border-bottom:1px solid #2a3942}
#chat-header-avatar{width:38px;height:38px;border-radius:50%;background:#00a884;display:flex;align-items:center;justify-content:center;font-size:15px;font-weight:600;color:#fff}
#chat-header-name{font-size:15px;font-weight:500}
#messages{flex:1;overflow-y:auto;padding:20px 60px;display:flex;flex-direction:column;gap:4px;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3C/svg%3E")}
.msg{max-width:65%;padding:7px 12px 6px;border-radius:8px;font-size:14px;line-height:1.45;position:relative;word-wrap:break-word}
.msg.out{background:#005c4b;align-self:flex-end;border-bottom-right-radius:2px}
.msg.in{background:#202c33;align-self:flex-start;border-bottom-left-radius:2px}
.msg-meta{display:flex;align-items:center;gap:4px;justify-content:flex-end;margin-top:3px}
.msg-time{font-size:11px;color:#8696a0}
.msg-type{font-size:11px;color:#8696a0;font-style:italic}
.msg-sender{font-size:11px;color:#00a884;font-weight:600;margin-bottom:3px}
.date-divider{text-align:center;margin:12px 0}
.date-divider span{background:#182229;color:#8696a0;font-size:12px;padding:4px 10px;border-radius:8px}
#empty{display:flex;flex:1;align-items:center;justify-content:center;flex-direction:column;gap:12px;color:#8696a0}
#empty svg{opacity:.3}
.loading{text-align:center;padding:20px;color:#8696a0;font-size:13px}
#token-bar{padding:8px 16px;background:#1f2c33;display:flex;gap:8px;align-items:center;border-bottom:1px solid #2a3942}
#token-bar input{flex:1;background:#2a3942;border:none;padding:7px 10px;border-radius:6px;color:#e9edef;font-size:13px}
#token-bar button{background:#00a884;border:none;color:#fff;padding:7px 14px;border-radius:6px;cursor:pointer;font-size:13px;font-weight:600}
#token-bar button:hover{background:#017561}
.badge{background:#00a884;color:#fff;border-radius:10px;font-size:11px;padding:1px 6px;font-weight:600}
</style>
</head>
<body>
<div id="header">
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
  <h1>Evolution GO &mdash; Chat Viewer</h1>
  <span id="header-info">Selecione uma instância e informe o token</span>
</div>
<div id="token-bar">
  <input id="token-input" type="password" placeholder="Token da instância (apikey)" />
  <button onclick="loadChats()">Carregar Chats</button>
</div>
<div id="app">
  <div id="sidebar">
    <div id="chat-list"><div class="loading">Informe o token e ID da instância acima e clique em Carregar Chats.</div></div>
  </div>
  <div id="main">
    <div id="empty">
      <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
      <span>Selecione uma conversa</span>
    </div>
  </div>
</div>
<script>
let currentChat = null;
let currentToken = '';

function getInitials(jid) {
  const num = jid.split('@')[0].split(':')[0];
  return num.slice(-2);
}

function formatTime(ts) {
  const d = new Date(ts);
  return d.toLocaleTimeString('pt-BR', {hour:'2-digit', minute:'2-digit'});
}

function formatDate(ts) {
  const d = new Date(ts);
  const today = new Date();
  const yesterday = new Date(today); yesterday.setDate(today.getDate()-1);
  if (d.toDateString() === today.toDateString()) return 'Hoje';
  if (d.toDateString() === yesterday.toDateString()) return 'Ontem';
  return d.toLocaleDateString('pt-BR');
}

async function loadChats() {
  currentToken = document.getElementById('token-input').value.trim();
  if (!currentToken) { alert('Informe o token da instância.'); return; }

  document.getElementById('header-info').textContent = 'Carregando...';
  document.getElementById('chat-list').innerHTML = '<div class="loading">Carregando chats...</div>';

  try {
    const res = await fetch('/message/chats?limit=100', {
      headers: { 'apikey': currentToken }
    });
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const json = await res.json();
    document.getElementById('header-info').textContent = (json.data || []).length + ' chats carregados';
    renderChatList(json.data || []);
  } catch(e) {
    document.getElementById('chat-list').innerHTML = '<div class="loading" style="color:#ef4444">Erro: ' + e.message + '</div>';
  }
}

function renderChatList(chats) {
  if (!chats || chats.length === 0) {
    document.getElementById('chat-list').innerHTML = '<div class="loading">Nenhum chat encontrado. Verifique se DATABASE_SAVE_MESSAGES=true.</div>';
    return;
  }
  const el = document.getElementById('chat-list');
  el.innerHTML = '';
  chats.forEach(c => {
    const div = document.createElement('div');
    div.className = 'chat-item';
    div.dataset.jid = c.chat;
    const initials = getInitials(c.chat);
    const preview = c.content || ('[' + (c.messageType || 'mensagem') + ']');
    div.innerHTML = '<div class="chat-avatar">'+initials+'</div>' +
      '<div class="chat-info"><div class="chat-name">'+c.chat.split('@')[0]+'</div>' +
      '<div class="chat-preview">'+(c.fromMe?'Você: ':'')+escHtml(preview)+'</div></div>' +
      '<div class="chat-time">'+formatTime(c.timestamp)+'</div>';
    div.onclick = () => openChat(c.chat, div);
    el.appendChild(div);
  });
}

async function openChat(jid, el) {
  document.querySelectorAll('.chat-item').forEach(i => i.classList.remove('active'));
  el.classList.add('active');
  currentChat = jid;

  const main = document.getElementById('main');
  main.innerHTML = '<div id="chat-header"><div id="chat-header-avatar">'+getInitials(jid)+'</div>' +
    '<div><div id="chat-header-name">'+jid.split('@')[0]+'</div></div></div>' +
    '<div id="messages"><div class="loading">Carregando mensagens...</div></div>';

  try {
    const res = await fetch('/message/chat?chat='+encodeURIComponent(jid)+'&limit=100', {
      headers: { 'apikey': currentToken }
    });
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const json = await res.json();
    renderMessages(json.data || []);
  } catch(e) {
    document.getElementById('messages').innerHTML = '<div class="loading" style="color:#ef4444">Erro: ' + e.message + '</div>';
  }
}

function renderMessages(msgs) {
  const container = document.getElementById('messages');
  if (!msgs || msgs.length === 0) {
    container.innerHTML = '<div class="loading">Nenhuma mensagem encontrada.</div>';
    return;
  }
  // messages come newest first, reverse to show oldest first
  const sorted = [...msgs].reverse();
  container.innerHTML = '';
  let lastDate = '';

  sorted.forEach(m => {
    const d = formatDate(m.timestamp);
    if (d !== lastDate) {
      lastDate = d;
      const divider = document.createElement('div');
      divider.className = 'date-divider';
      divider.innerHTML = '<span>'+d+'</span>';
      container.appendChild(divider);
    }

    const div = document.createElement('div');
    div.className = 'msg ' + (m.fromMe ? 'out' : 'in');

    let inner = '';
    if (!m.fromMe) {
      inner += '<div class="msg-sender">'+m.source.split('@')[0]+'</div>';
    }
    const content = m.content || ('[' + (m.messageType || 'mensagem') + ']');
    inner += escHtml(content);
    inner += '<div class="msg-meta"><span class="msg-time">'+formatTime(m.timestamp)+'</span>';
    if (m.messageType && m.messageType !== 'text') {
      inner += '<span class="msg-type">'+m.messageType+'</span>';
    }
    inner += '</div>';

    div.innerHTML = inner;
    container.appendChild(div);
  });
  container.scrollTop = container.scrollHeight;
}

function escHtml(t) {
  return String(t).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g,'<br>');
}

// Allow Enter key to trigger load
document.addEventListener('keydown', e => {
  if (e.key === 'Enter' && document.activeElement.id === 'token-input') {
    loadChats();
  }
});
</script>
</body>
</html>`

// ChatViewer serves the chat viewer HTML page.
func ChatViewer(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, chatViewerHTML)
}
