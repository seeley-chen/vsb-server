package middleware

// logViewerHTML 日志查看器页面（内嵌，无外部依赖）
const logViewerHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>VSB Log Viewer</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: 'SF Mono', 'Menlo', 'Consolas', monospace; background: #1e1e2e; color: #cdd6f4; height: 100vh; display: flex; flex-direction: column; }
  .toolbar { background: #181825; padding: 8px 12px; display: flex; gap: 8px; align-items: center; flex-wrap: wrap; border-bottom: 1px solid #313244; }
  .toolbar button, .toolbar select, .toolbar input { background: #313244; color: #cdd6f4; border: 1px solid #45475a; border-radius: 4px; padding: 4px 8px; font-size: 13px; font-family: inherit; cursor: pointer; }
  .toolbar input { width: 200px; }
  .toolbar button:hover { background: #45475a; }
  .toolbar button.active { background: #89b4fa; color: #1e1e2e; }
  .toolbar .sep { width: 1px; height: 20px; background: #45475a; }
  .toolbar .status { margin-left: auto; font-size: 12px; }
  .dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 4px; }
  .dot.on { background: #a6e3a1; }
  .dot.off { background: #f38ba8; }
  .logs { flex: 1; overflow-y: auto; padding: 4px 0; }
  .entry { padding: 3px 12px; border-bottom: 1px solid #11111b; cursor: pointer; font-size: 13px; line-height: 1.5; }
  .entry:hover { background: #313244; }
  .entry .meta { display: flex; gap: 8px; align-items: center; }
  .entry .time { color: #6c7086; font-size: 12px; white-space: nowrap; }
  .entry .lvl { font-weight: bold; width: 44px; text-align: center; border-radius: 3px; padding: 0 4px; }
  .lvl-INFO { background: #45475a; color: #89b4fa; }
  .lvl-WARN { background: #432f2a; color: #fab387; }
  .lvl-ERROR { background: #45283a; color: #f38ba8; }
  .entry .status { width: 32px; text-align: right; font-weight: bold; }
  .status-2 { color: #a6e3a1; }
  .status-4 { color: #fab387; }
  .status-5 { color: #f38ba8; }
  .entry .method { color: #cba6f7; width: 48px; }
  .entry .path { color: #cdd6f4; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .entry .dur { color: #6c7086; font-size: 12px; white-space: nowrap; }
  .entry .rid { color: #585b70; font-size: 12px; white-space: nowrap; }
  .detail { display: none; padding: 8px 12px 8px 60px; background: #11111b; font-size: 12px; }
  .detail.show { display: block; }
  .detail .row { display: flex; gap: 8px; margin: 2px 0; }
  .detail .k { color: #6c7086; width: 80px; flex-shrink: 0; }
  .detail .v { color: #cdd6f4; word-break: break-all; white-space: pre-wrap; }
  .detail .body { background: #181825; padding: 4px 8px; border-radius: 4px; max-height: 200px; overflow-y: auto; }
  .empty { text-align: center; color: #6c7086; padding: 40px; }
</style>
</head>
<body>
<div class="toolbar">
  <button class="filter active" data-lvl="">全部</button>
  <button class="filter" data-lvl="INFO">INFO</button>
  <button class="filter" data-lvl="WARN">WARN</button>
  <button class="filter" data-lvl="ERROR">ERROR</button>
  <div class="sep"></div>
  <input type="text" id="search" placeholder="搜索 path / request_id..." />
  <div class="sep"></div>
  <button id="pause">⏸ 暂停</button>
  <button id="clear">🗑 清空</button>
  <div class="status"><span class="dot off" id="dot"></span><span id="connText">未连接</span> · <span id="count">0</span> 条</div>
</div>
<div class="logs" id="logs"><div class="empty">等待日志...</div></div>
<script>
var all = [], paused = false, curLvl = '', curSearch = '', sse = null;
var logsEl = document.getElementById('logs');
var countEl = document.getElementById('count');

function statusClass(s) { return s >= 500 ? 'status-5' : s >= 400 ? 'status-4' : 'status-2'; }
function statusPrefix(s) { return s >= 500 ? 's5' : s >= 400 ? 's4' : 's2'; }

function matches(e) {
  if (curLvl && e.level !== curLvl) return false;
  if (curSearch) {
    var s = curSearch.toLowerCase();
    if (e.path.toLowerCase().indexOf(s) < 0 && e.request_id.toLowerCase().indexOf(s) < 0) return false;
  }
  return true;
}

function renderEntry(e) {
  var div = document.createElement('div');
  var ts = e.time.replace('T',' ').substring(11,23);
  var html = '<div class="entry">' +
    '<span class="time">' + ts + '</span>' +
    '<span class="lvl lvl-' + e.level + '">' + e.level + '</span>' +
    '<span class="status ' + statusClass(e.status) + '">' + e.status + '</span>' +
    '<span class="method">' + e.method + '</span>' +
    '<span class="path">' + esc(e.path) + (e.query ? '?' + esc(e.query) : '') + '</span>' +
    '<span class="dur">' + e.duration + '</span>' +
    '<span class="rid">' + e.request_id + '</span>' +
    '</div>';
  html += '<div class="detail">' +
    '<div class="row"><span class="k">user_id</span><span class="v">' + esc(e.user_id||'-') + '</span></div>' +
    '<div class="row"><span class="k">ip</span><span class="v">' + esc(e.ip) + '</span></div>' +
    '<div class="row"><span class="k">ua</span><span class="v">' + esc(e.ua) + '</span></div>' +
    '<div class="row"><span class="k">req_body</span><div class="v body">' + esc(e.req_body||'(empty)') + '</div></div>' +
    '<div class="row"><span class="k">resp_body</span><div class="v body">' + esc(e.resp_body||'(empty)') + '</div></div>' +
    '</div>';
  div.innerHTML = html;
  div.querySelector('.entry').onclick = function() {
    div.querySelector('.detail').classList.toggle('show');
  };
  return div;
}

function esc(s) { return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }

function refresh() {
  logsEl.innerHTML = '';
  var shown = 0;
  for (var i = all.length - 1; i >= 0 && shown < 500; i--) {
    if (matches(all[i])) { logsEl.appendChild(renderEntry(all[i])); shown++; }
  }
  if (shown === 0) { logsEl.innerHTML = '<div class="empty">无匹配日志</div>'; }
  countEl.textContent = shown;
}

function addEntry(e) {
  all.push(e);
  if (all.length > 1000) all = all.slice(-1000);
  if (paused) return;
  if (!matches(e)) return;
  if (logsEl.querySelector('.empty')) logsEl.innerHTML = '';
  logsEl.insertBefore(renderEntry(e), logsEl.firstChild);
  // 限制 DOM 节点数
  while (logsEl.children.length > 500) logsEl.removeChild(logsEl.lastChild);
  countEl.textContent = logsEl.children.length;
}

function connect() {
  document.getElementById('dot').className = 'dot off';
  document.getElementById('connText').textContent = '连接中';
  sse = new EventSource('/admin/logs/stream');
  sse.onopen = function() {
    document.getElementById('dot').className = 'dot on';
    document.getElementById('connText').textContent = '已连接';
  };
  sse.onmessage = function(ev) {
    try { addEntry(JSON.parse(ev.data)); } catch(e) {}
  };
  sse.onerror = function() {
    document.getElementById('dot').className = 'dot off';
    document.getElementById('connText').textContent = '已断开，3s 后重连';
    sse.close();
    setTimeout(connect, 3000);
  };
}

document.querySelectorAll('.filter').forEach(function(btn) {
  btn.onclick = function() {
    document.querySelectorAll('.filter').forEach(function(b) { b.classList.remove('active'); });
    btn.classList.add('active');
    curLvl = btn.dataset.lvl;
    refresh();
  };
});

document.getElementById('search').oninput = function(e) {
  curSearch = e.target.value;
  refresh();
};

document.getElementById('pause').onclick = function() {
  paused = !paused;
  this.textContent = paused ? '▶ 继续' : '⏸ 暂停';
};

document.getElementById('clear').onclick = function() {
  all = [];
  logsEl.innerHTML = '<div class="empty">已清空</div>';
  countEl.textContent = '0';
};

connect();
</script>
</body>
</html>`
