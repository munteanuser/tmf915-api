/* ── State ─────────────────────────────────────────────────────────────────── */
const S = {
  apiBase: localStorage.getItem('aim_base') || 'http://localhost:8080/tmf-api/AiM/v4',
  apiKey:  localStorage.getItem('aim_key')  || '',
  theme:   localStorage.getItem('aim_theme') || 'dark',
  metrics: JSON.parse(localStorage.getItem('aim_metrics') || '[]'), // [{ts,method,path,status,ms,bytes}]
  counts:  {},
};

/* ── Resource definitions ──────────────────────────────────────────────────── */
const RESOURCES = {
  aiModelSpecification: {
    label: 'AI Model Specification', plural: 'AI Model Specifications',
    icon: '📋', color: '#388bfd',
    cols: ['name','version','state','modelType','framework'],
    fields: [
      {k:'name',       label:'Name',        type:'text', required:true},
      {k:'description',label:'Description', type:'textarea'},
      {k:'version',    label:'Version',     type:'text', placeholder:'1.0'},
      {k:'state',      label:'State',       type:'select', options:['active','inactive','draft']},
      {k:'category',   label:'Category',    type:'text'},
      {k:'modelType',  label:'Model Type',  type:'text', placeholder:'generative'},
      {k:'framework',  label:'Framework',   type:'text', placeholder:'transformer'},
      {k:'isBundle',   label:'Is Bundle',   type:'select', options:['false','true']},
    ],
  },
  aiModel: {
    label: 'AI Model', plural: 'AI Models',
    icon: '🤖', color: '#3fb950',
    cols: ['name','state','modelType','vendor','version'],
    fields: [
      {k:'name',       label:'Name',        type:'text', required:true},
      {k:'description',label:'Description', type:'textarea'},
      {k:'state',      label:'State',       type:'select', options:['active','inactive','draft','retired']},
      {k:'modelType',  label:'Model Type',  type:'text'},
      {k:'framework',  label:'Framework',   type:'text'},
      {k:'vendor',     label:'Vendor',      type:'text'},
      {k:'version',    label:'Version',     type:'text'},
    ],
  },
  aiContractSpecification: {
    label: 'AI Contract Specification', plural: 'AI Contract Specifications',
    icon: '📄', color: '#bc8cff',
    cols: ['name','version','state'],
    fields: [
      {k:'name',       label:'Name',        type:'text', required:true},
      {k:'description',label:'Description', type:'textarea'},
      {k:'version',    label:'Version',     type:'text'},
      {k:'state',      label:'State',       type:'select', options:['active','inactive','draft']},
    ],
  },
  aiContract: {
    label: 'AI Contract', plural: 'AI Contracts',
    icon: '📑', color: '#d29922',
    cols: ['name','state','documentNumber'],
    fields: [
      {k:'name',          label:'Name',            type:'text', required:true},
      {k:'description',   label:'Description',     type:'textarea'},
      {k:'state',         label:'State',           type:'select', options:['active','inactive','draft','terminated']},
      {k:'documentNumber',label:'Document Number', type:'text'},
    ],
  },
  aiContractViolation: {
    label: 'AI Contract Violation', plural: 'AI Contract Violations',
    icon: '⚠️', color: '#f85149',
    cols: ['name','violationType','severity'],
    readOnly: true,
    fields: [
      {k:'name',          label:'Name',           type:'text', required:true},
      {k:'description',   label:'Description',    type:'textarea'},
      {k:'violationType', label:'Violation Type', type:'text'},
      {k:'severity',      label:'Severity',       type:'select', options:['minor','major','critical']},
    ],
  },
  alarm: {
    label: 'Alarm', plural: 'Alarms',
    icon: '🔔', color: '#f85149',
    cols: ['alarmType','perceivedSeverity','state','specificProblem'],
    fields: [
      {k:'alarmType',        label:'Alarm Type',         type:'text', required:true},
      {k:'perceivedSeverity',label:'Perceived Severity', type:'select', options:['critical','major','minor','warning','indeterminate','cleared']},
      {k:'state',            label:'State',              type:'select', options:['raised','cleared','updated']},
      {k:'alarmDetails',     label:'Alarm Details',      type:'textarea'},
      {k:'specificProblem',  label:'Specific Problem',   type:'text'},
    ],
  },
  rule: {
    label: 'Rule', plural: 'Rules',
    icon: '⚙️', color: '#58a6ff',
    cols: ['name','ruleType','state','expression'],
    fields: [
      {k:'name',       label:'Name',       type:'text', required:true},
      {k:'description',label:'Description',type:'textarea'},
      {k:'ruleType',   label:'Rule Type',  type:'text'},
      {k:'state',      label:'State',      type:'select', options:['active','inactive','draft']},
      {k:'expression', label:'Expression', type:'text', placeholder:'accuracy < 0.85'},
    ],
  },
  monitor: {
    label: 'Monitor', plural: 'Monitors',
    icon: '📈', color: '#3fb950',
    cols: ['name','state','monitorType'],
    fields: [
      {k:'name',        label:'Name',         type:'text', required:true},
      {k:'description', label:'Description',  type:'textarea'},
      {k:'state',       label:'State',        type:'select', options:['active','suspended','inactive']},
      {k:'monitorType', label:'Monitor Type', type:'text'},
    ],
  },
  topic: {
    label: 'Topic', plural: 'Topics',
    icon: '💬', color: '#bc8cff',
    cols: ['name','state'],
    fields: [
      {k:'name',       label:'Name',        type:'text', required:true},
      {k:'description',label:'Description', type:'textarea'},
      {k:'state',      label:'State',       type:'select', options:['active','inactive']},
    ],
  },
  hub: {
    label: 'Hub', plural: 'Hubs',
    icon: '🌐', color: '#388bfd',
    cols: ['callback','query'],
    noEdit: true,
    fields: [
      {k:'callback', label:'Callback URL', type:'text', required:true, placeholder:'https://your-server.example.com/notify'},
      {k:'query',    label:'Query Filter', type:'text', placeholder:'eventType=AiModelCreateEvent'},
    ],
  },
};

const TEMPLATES = [
  { resource:'aiModel', name:'GPT-4 Production', desc:'Large language model deployment template', body:{name:'GPT-4',description:'OpenAI GPT-4 large language model',state:'active',modelType:'generative',framework:'transformer',vendor:'OpenAI',version:'4.0'} },
  { resource:'aiModel', name:'Llama 3 Open Source', desc:'Meta open-source LLM template', body:{name:'Llama-3-70B',description:'Meta Llama 3 70B open-source model',state:'active',modelType:'generative',framework:'transformer',vendor:'Meta',version:'3.0'} },
  { resource:'aiModelSpecification', name:'Standard LLM Spec', desc:'Baseline spec for large language models', body:{name:'Standard LLM Specification',description:'Specification for large language models',version:'1.0',state:'active',category:'LLM',modelType:'generative',framework:'transformer',isBundle:false} },
  { resource:'aiContract', name:'Enterprise SLA', desc:'90-day enterprise contract with SLA', body:{name:'Enterprise AI Contract',description:'Enterprise-grade AI service contract',state:'active',documentNumber:'CTR-ENTERPRISE-001'} },
  { resource:'alarm', name:'Accuracy Degradation', desc:'Alarm for model accuracy drop', body:{alarmType:'qualityOfServiceAlarm',perceivedSeverity:'critical',state:'raised',specificProblem:'accuracy-degradation',alarmDetails:'Model accuracy dropped below 85% threshold'} },
  { resource:'rule', name:'Accuracy Threshold', desc:'Trigger when accuracy < 85%', body:{name:'Accuracy Threshold Rule',description:'Trigger alarm when model accuracy drops below 85%',ruleType:'threshold',state:'active',expression:'accuracy < 0.85'} },
  { resource:'rule', name:'Latency SLA', desc:'Trigger when latency > 500ms', body:{name:'Latency SLA Rule',description:'Trigger when P95 latency exceeds 500ms',ruleType:'threshold',state:'active',expression:'p95_latency_ms > 500'} },
  { resource:'monitor', name:'Accuracy Monitor', desc:'Monitor model accuracy metrics', body:{name:'Model Accuracy Monitor',description:'Monitors model accuracy in production',state:'active',monitorType:'performance'} },
  { resource:'hub', name:'Slack Webhook', desc:'Notify Slack on any event', body:{callback:'https://hooks.slack.com/services/YOUR/WEBHOOK/URL',query:''} },
  { resource:'aiContractSpecification', name:'Standard SLA Spec', desc:'Reusable SLA contract specification', body:{name:'Standard SLA Specification',description:'Reusable specification for AI service level agreements',version:'1.0',state:'active'} },
];

/* ── API client ────────────────────────────────────────────────────────────── */
async function api(method, path, body) {
  const url = S.apiBase + path;
  const t0 = Date.now();
  let res, text = '';
  try {
    res = await fetch(url, {
      method,
      headers: { 'X-API-Key': S.apiKey, 'Content-Type': 'application/json' },
      body: body ? JSON.stringify(body) : undefined,
    });
    text = await res.text();
  } catch(e) {
    recordMetric(method, path, 0, Date.now()-t0, 0);
    throw new Error('Network error: ' + e.message);
  }
  recordMetric(method, path, res.status, Date.now()-t0, text.length);
  let json = null;
  try { json = JSON.parse(text); } catch(_) {}
  return { ok: res.status >= 200 && res.status < 300, status: res.status, body: json, text };
}

function recordMetric(method, path, status, ms, bytes) {
  S.metrics.push({ ts: Date.now(), method, path, status, ms, bytes });
  if (S.metrics.length > 500) S.metrics = S.metrics.slice(-500);
  localStorage.setItem('aim_metrics', JSON.stringify(S.metrics));
}

async function apiList(resource) {
  const r = await api('GET', `/${resource}?limit=100`);
  if (!r.ok) throw new Error(`${r.status}: ${r.text}`);
  return r.body || [];
}
async function apiGet(resource, id) {
  const r = await api('GET', `/${resource}/${id}`);
  if (!r.ok) throw new Error(`${r.status}`);
  return r.body;
}
async function apiCreate(resource, body) {
  const r = await api('POST', `/${resource}`, body);
  if (!r.ok) throw new Error(r.body?.message || r.text);
  return r.body;
}
async function apiPatch(resource, id, body) {
  const r = await api('PATCH', `/${resource}/${id}`, body);
  if (!r.ok) throw new Error(r.body?.message || r.text);
  return r.body;
}
async function apiDelete(resource, id) {
  const r = await api('DELETE', `/${resource}/${id}`);
  if (!r.ok) throw new Error(r.body?.message || r.text);
}

/* ── Toast ─────────────────────────────────────────────────────────────────── */
function toast(msg, type='success') {
  const el = document.createElement('div');
  el.className = `toast toast-${type}`;
  el.innerHTML = `<div class="toast-dot"></div><span>${esc(msg)}</span>`;
  document.getElementById('toast-container').appendChild(el);
  setTimeout(() => el.remove(), 3500);
}

/* ── Modal ─────────────────────────────────────────────────────────────────── */
function openModal(title, bodyHTML, footerHTML) {
  document.getElementById('modal-title').textContent = title;
  document.getElementById('modal-body').innerHTML = bodyHTML;
  document.getElementById('modal-footer').innerHTML = footerHTML || '';
  document.getElementById('modal-overlay').classList.remove('hidden');
}
function closeModal() { document.getElementById('modal-overlay').classList.add('hidden'); }

/* ── Router ─────────────────────────────────────────────────────────────────── */
const routes = {};
function on(path, fn) { routes[path] = fn; }

async function route() {
  const hash = location.hash.slice(1) || '/';
  document.querySelectorAll('.nav-item').forEach(el => {
    el.classList.toggle('active', el.dataset.route === hash || (hash.startsWith(el.dataset.route+'/') && el.dataset.route !== '/'));
    if (el.dataset.route === '/' && hash === '/') el.classList.add('active');
  });
  const content = document.getElementById('page-content');
  content.innerHTML = '<div class="loading-spinner">Loading…</div>';
  const parts = hash.split('/').filter(Boolean);
  const key = '/' + (parts[0] || '');
  const fn = routes[key];
  if (fn) {
    try { await fn(parts); }
    catch(e) { content.innerHTML = `<div class="empty-state"><div class="empty-state-icon">⚠️</div><h3>Error</h3><p>${esc(e.message)}</p></div>`; }
  } else {
    content.innerHTML = `<div class="empty-state"><h3>404</h3><p>Page not found</p></div>`;
  }
}

/* ── Helpers ────────────────────────────────────────────────────────────────── */
function esc(s) { return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;'); }
function nav(r) { location.hash = r; }
function fmtTime(ts) {
  const d = new Date(ts), now = new Date();
  const diff = (now - d) / 1000;
  if (diff < 60) return `${Math.floor(diff)}s ago`;
  if (diff < 3600) return `${Math.floor(diff/60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff/3600)}h ago`;
  return d.toLocaleDateString();
}
function fmtStatus(s) {
  if (!s) return '<span class="activity-status status-4xx">—</span>';
  const cls = s >= 500 ? 'status-5xx' : s >= 400 ? 'status-4xx' : 'status-2xx';
  return `<span class="activity-status ${cls}">${s}</span>`;
}
function stateBadge(state) {
  const m = {active:'badge-active', inactive:'badge-inactive', raised:'badge-raised', cleared:'badge-cleared', draft:'badge-warning', terminated:'badge-inactive'};
  return `<span class="badge ${m[state]||'badge-unknown'}">${esc(state||'—')}</span>`;
}

/* ── Health check ───────────────────────────────────────────────────────────── */
async function checkHealth() {
  const el = document.getElementById('health-badge');
  if (!el) return;
  try {
    const r = await fetch(S.apiBase.replace('/tmf-api/AiM/v4','') + '/health');
    if (r.ok) { el.className = 'badge badge-ok'; el.textContent = '● online'; }
    else throw new Error();
  } catch(_) { el.className = 'badge badge-error'; el.textContent = '● offline'; }
}

/* ════════════════════════════════════════════════════════════════════════════
   PAGES
   ════════════════════════════════════════════════════════════════════════════ */

/* ── Dashboard ──────────────────────────────────────────────────────────────── */
on('/', async () => {
  const content = document.getElementById('page-content');
  content.innerHTML = '<div class="loading-spinner">Loading dashboard…</div>';

  const counts = {};
  await Promise.all(Object.keys(RESOURCES).map(async r => {
    try { const d = await apiList(r); counts[r] = d.length; S.counts[r] = d.length; }
    catch(_) { counts[r] = '?'; }
  }));

  const totalCalls = S.metrics.length;
  const totalBytes = S.metrics.reduce((a,m) => a + (m.bytes||0), 0);
  const tokenEst = Math.round(totalBytes / 4);
  const avgMs = S.metrics.length ? Math.round(S.metrics.reduce((a,m)=>a+(m.ms||0),0)/S.metrics.length) : 0;
  const errRate = S.metrics.length ? Math.round(S.metrics.filter(m=>m.status>=400).length/S.metrics.length*100) : 0;

  const recent = [...S.metrics].reverse().slice(0,10);

  const barData = buildHourlyBars();

  content.innerHTML = `
<div class="page-header">
  <div><div class="page-title">Dashboard</div><div class="page-subtitle">TMF915 AI Management Suite overview</div></div>
  <div class="page-actions">
    <button class="btn btn-secondary" onclick="location.hash='#/settings'">⚙ Settings</button>
    <button class="btn btn-primary" onclick="refreshDashboard()">↺ Refresh</button>
  </div>
</div>

<div class="stats-grid">
  ${Object.entries(RESOURCES).map(([k,r]) => `
  <div class="stat-card" onclick="nav('/${k}')">
    <div class="stat-card-icon">${r.icon}</div>
    <div class="stat-card-value">${counts[k]??'—'}</div>
    <div class="stat-card-label">${r.plural}</div>
  </div>`).join('')}
</div>

<div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:16px;margin-bottom:24px">
  <div class="panel">
    <div class="panel-header"><span class="panel-title">API Calls</span></div>
    <div class="panel-body" style="text-align:center">
      <div style="font-size:36px;font-weight:700;color:var(--accent)">${totalCalls}</div>
      <div style="font-size:12px;color:var(--text-secondary)">total requests this session</div>
    </div>
  </div>
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Avg Latency</span></div>
    <div class="panel-body" style="text-align:center">
      <div style="font-size:36px;font-weight:700;color:var(--success)">${avgMs}<span style="font-size:16px;font-weight:400">ms</span></div>
      <div style="font-size:12px;color:var(--text-secondary)">average response time</div>
    </div>
  </div>
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Error Rate</span></div>
    <div class="panel-body" style="text-align:center">
      <div style="font-size:36px;font-weight:700;color:${errRate>5?'var(--danger)':'var(--success)'}">${errRate}<span style="font-size:16px;font-weight:400">%</span></div>
      <div style="font-size:12px;color:var(--text-secondary)">4xx / 5xx responses</div>
    </div>
  </div>
</div>

<div style="display:grid;grid-template-columns:2fr 1fr;gap:16px;margin-bottom:24px">
  <div class="panel">
    <div class="panel-header">
      <span class="panel-title">Request Volume (last 24h)</span>
      <div class="method-pills">
        <span class="method-pill method-GET">GET</span>
        <span class="method-pill method-POST">POST</span>
        <span class="method-pill method-PATCH">PATCH</span>
        <span class="method-pill method-DELETE">DELETE</span>
      </div>
    </div>
    <div class="panel-body">
      <div class="bar-chart" id="bar-chart">
        ${barData.map(b => {
          const pct = barData.max ? Math.max(4, Math.round(b.count/barData.max*100)) : 4;
          return `<div class="bar-chart-bar" style="height:${pct}%"><div class="bar-tooltip">${b.label}: ${b.count}</div></div>`;
        }).join('')}
      </div>
      <div class="bar-chart-labels">
        ${barData.map(b => `<div class="bar-chart-label">${b.label}</div>`).join('')}
      </div>
    </div>
  </div>
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Token Estimation</span></div>
    <div class="panel-body">
      <div style="margin-bottom:16px">
        <div style="font-size:11px;color:var(--text-muted);margin-bottom:2px">DATA TRANSFERRED</div>
        <div style="font-size:20px;font-weight:700">${fmtBytes(totalBytes)}</div>
      </div>
      <div style="margin-bottom:16px">
        <div style="font-size:11px;color:var(--text-muted);margin-bottom:2px">ESTIMATED TOKENS</div>
        <div style="font-size:20px;font-weight:700;color:var(--purple)">${tokenEst.toLocaleString()}</div>
        <div style="font-size:11px;color:var(--text-secondary)">~4 chars/token</div>
      </div>
      <div>
        <div style="font-size:11px;color:var(--text-muted);margin-bottom:2px">GPT-4o COST EST.</div>
        <div style="font-size:20px;font-weight:700;color:var(--warning)">$${(tokenEst/1000*0.005).toFixed(4)}</div>
        <div style="font-size:11px;color:var(--text-secondary)">at $0.005 / 1K tokens</div>
      </div>
    </div>
  </div>
</div>

<div class="panel">
  <div class="panel-header">
    <span class="panel-title">Recent Activity</span>
    <button class="btn btn-ghost btn-sm" onclick="S.metrics=[];localStorage.removeItem('aim_metrics');route()">Clear</button>
  </div>
  <div class="activity-feed">
    ${recent.length ? recent.map(m => `
    <div class="activity-item">
      <span class="activity-time">${fmtTime(m.ts)}</span>
      <span class="method-pill method-${m.method}">${m.method}</span>
      <span class="activity-resource">${esc(m.path)}</span>
      <span style="color:var(--text-muted);font-size:11px">${m.ms}ms</span>
      ${fmtStatus(m.status)}
    </div>`).join('') : '<div style="padding:24px;text-align:center;color:var(--text-muted);font-size:13px">No API calls recorded yet</div>'}
  </div>
</div>`;
});

function refreshDashboard() { route(); }

function buildHourlyBars() {
  const now = Date.now();
  const bars = Array.from({length:24},(_,i) => {
    const h = new Date(now - (23-i)*3600000);
    return { label: h.getHours()+':00', count: 0, hour: h.getHours() };
  });
  S.metrics.forEach(m => {
    const hAgo = Math.floor((now - m.ts)/3600000);
    if (hAgo < 24) bars[23 - hAgo].count++;
  });
  bars.max = Math.max(...bars.map(b=>b.count), 1);
  return bars;
}

function fmtBytes(b) {
  if (b < 1024) return b + ' B';
  if (b < 1048576) return (b/1024).toFixed(1) + ' KB';
  return (b/1048576).toFixed(2) + ' MB';
}

/* ── Generic resource list ───────────────────────────────────────────────────── */
Object.keys(RESOURCES).forEach(resource => {
  on('/'+resource, async (parts) => {
    if (parts[1]) { await renderDetail(resource, parts[1]); return; }
    await renderList(resource);
  });
});

async function renderList(resource) {
  const R = RESOURCES[resource];
  const content = document.getElementById('page-content');
  let items = [], err = null;
  try { items = await apiList(resource); }
  catch(e) { err = e.message; }

  const filterInput = document.getElementById('search-'+resource);
  const q = filterInput ? filterInput.value.toLowerCase() : '';
  const filtered = q ? items.filter(i => JSON.stringify(i).toLowerCase().includes(q)) : items;

  content.innerHTML = `
<div class="page-header">
  <div>
    <div class="page-title">${R.icon} ${R.plural}</div>
    <div class="page-subtitle">${filtered.length} of ${items.length} records</div>
  </div>
  <div class="page-actions">
    ${!R.readOnly ? `<button class="btn btn-primary" onclick="openCreateModal('${resource}')">+ New ${R.label}</button>` : ''}
  </div>
</div>
${err ? `<div style="color:var(--danger);padding:16px;background:var(--danger-subtle);border-radius:var(--radius);margin-bottom:16px">${esc(err)}</div>` : ''}
<div class="panel">
  <div class="toolbar">
    <input type="text" class="toolbar-search" id="search-${resource}" placeholder="Filter ${R.plural.toLowerCase()}…" value="${esc(q)}" oninput="renderList('${resource}')"/>
    <div class="toolbar-right">
      <span style="font-size:12px;color:var(--text-muted)">${filtered.length} items</span>
    </div>
  </div>
  <div class="panel-body-flush">
    ${filtered.length === 0 ? `
    <div class="empty-state">
      <div class="empty-state-icon">${R.icon}</div>
      <h3>No ${R.plural} yet</h3>
      <p>Create your first ${R.label.toLowerCase()} to get started.</p>
      ${!R.readOnly ? `<button class="btn btn-primary" onclick="openCreateModal('${resource}')">+ New ${R.label}</button>` : ''}
    </div>` : `
    <table class="data-table">
      <thead><tr>
        <th>ID</th>
        ${R.cols.map(c=>`<th>${esc(c)}</th>`).join('')}
        <th>Actions</th>
      </tr></thead>
      <tbody>
        ${filtered.map(item => `
        <tr>
          <td class="id-cell" title="${esc(item.id)}">${esc((item.id||'').substring(0,8))}…</td>
          ${R.cols.map(c => {
            const v = item[c];
            if (c==='state') return `<td>${stateBadge(v)}</td>`;
            return `<td title="${esc(v)}">${esc(v)||'<span style="color:var(--text-muted)">—</span>'}</td>`;
          }).join('')}
          <td>
            <div class="actions">
              <button class="btn btn-ghost btn-sm" onclick="openDetailModal('${resource}','${esc(item.id)}')">View</button>
              ${!R.readOnly && !R.noEdit ? `<button class="btn btn-ghost btn-sm" onclick="openEditModal('${resource}','${esc(item.id)}')">Edit</button>` : ''}
              ${!R.readOnly ? `<button class="btn btn-danger btn-sm" onclick="confirmDelete('${resource}','${esc(item.id)}','${esc(item.name||item.id)}')">Delete</button>` : ''}
            </div>
          </td>
        </tr>`).join('')}
      </tbody>
    </table>`}
  </div>
</div>`;
}

/* ── Create modal ────────────────────────────────────────────────────────────── */
function openCreateModal(resource, prefill={}) {
  const R = RESOURCES[resource];
  openModal(`New ${R.label}`, buildForm(resource, prefill), `
    <button class="btn btn-secondary" onclick="closeModal()">Cancel</button>
    <button class="btn btn-primary" onclick="submitCreate('${resource}')">Create</button>`);
}

async function submitCreate(resource) {
  const body = readForm(resource);
  try {
    await apiCreate(resource, body);
    closeModal();
    toast(`${RESOURCES[resource].label} created`, 'success');
    await renderList(resource);
  } catch(e) { toast(e.message, 'error'); }
}

/* ── Edit modal ──────────────────────────────────────────────────────────────── */
async function openEditModal(resource, id) {
  const R = RESOURCES[resource];
  try {
    const item = await apiGet(resource, id);
    openModal(`Edit ${R.label}`, buildForm(resource, item, id), `
      <button class="btn btn-secondary" onclick="closeModal()">Cancel</button>
      <button class="btn btn-primary" onclick="submitEdit('${resource}','${id}')">Save Changes</button>`);
  } catch(e) { toast(e.message, 'error'); }
}

async function submitEdit(resource, id) {
  const body = readForm(resource);
  try {
    await apiPatch(resource, id, body);
    closeModal();
    toast(`${RESOURCES[resource].label} updated`, 'success');
    await renderList(resource);
  } catch(e) { toast(e.message, 'error'); }
}

/* ── Detail modal ────────────────────────────────────────────────────────────── */
async function openDetailModal(resource, id) {
  const R = RESOURCES[resource];
  try {
    const item = await apiGet(resource, id);
    openModal(`${R.label} Detail`,
      `<div class="json-block">${esc(JSON.stringify(item, null, 2))}</div>`,
      `<button class="btn btn-secondary" onclick="closeModal()">Close</button>
       ${!R.readOnly && !R.noEdit ? `<button class="btn btn-primary" onclick="closeModal();openEditModal('${resource}','${id}')">Edit</button>` : ''}`
    );
  } catch(e) { toast(e.message, 'error'); }
}

/* ── Delete confirm ──────────────────────────────────────────────────────────── */
function confirmDelete(resource, id, name) {
  openModal(`Delete ${RESOURCES[resource].label}`,
    `<p style="color:var(--text-secondary)">Are you sure you want to delete <strong>${esc(name)}</strong>? This cannot be undone.</p>`,
    `<button class="btn btn-secondary" onclick="closeModal()">Cancel</button>
     <button class="btn btn-danger" onclick="doDelete('${resource}','${id}')">Delete</button>`);
}
async function doDelete(resource, id) {
  try {
    await apiDelete(resource, id);
    closeModal();
    toast(`${RESOURCES[resource].label} deleted`, 'success');
    await renderList(resource);
  } catch(e) { toast(e.message, 'error'); }
}

/* ── Form builder ────────────────────────────────────────────────────────────── */
function buildForm(resource, data={}, id=null) {
  const R = RESOURCES[resource];
  return `<form id="res-form" onsubmit="return false">
    ${id ? `<input type="hidden" name="_id" value="${esc(id)}"/>` : ''}
    ${R.fields.map(f => {
      const v = data[f.k] ?? '';
      if (f.type === 'textarea') return `
        <div class="form-group">
          <label class="form-label">${esc(f.label)}${f.required?' <span style="color:var(--danger)">*</span>':''}</label>
          <textarea class="form-control" name="${f.k}" rows="3">${esc(v)}</textarea>
        </div>`;
      if (f.type === 'select') return `
        <div class="form-group">
          <label class="form-label">${esc(f.label)}</label>
          <select class="form-control" name="${f.k}">
            ${f.options.map(o=>`<option value="${o}"${String(v)===String(o)?' selected':''}>${o}</option>`).join('')}
          </select>
        </div>`;
      return `
        <div class="form-group">
          <label class="form-label">${esc(f.label)}${f.required?' <span style="color:var(--danger)">*</span>':''}</label>
          <input class="form-control" type="text" name="${f.k}" value="${esc(v)}" placeholder="${esc(f.placeholder||'')}"/>
        </div>`;
    }).join('')}
  </form>`;
}

function readForm(resource) {
  const form = document.getElementById('res-form');
  const data = {};
  const R = RESOURCES[resource];
  R.fields.forEach(f => {
    const el = form.elements[f.k];
    if (!el) return;
    const v = el.value.trim();
    if (v !== '') data[f.k] = v;
  });
  return data;
}

/* ── Detail page ─────────────────────────────────────────────────────────────── */
async function renderDetail(resource, id) {
  const R = RESOURCES[resource];
  const content = document.getElementById('page-content');
  const item = await apiGet(resource, id);
  content.innerHTML = `
<div class="page-header">
  <div>
    <div style="font-size:12px;color:var(--text-secondary);margin-bottom:4px">
      <a href="#/${resource}">${R.plural}</a> / ${esc(id.substring(0,8))}…
    </div>
    <div class="page-title">${R.icon} ${esc(item.name || id)}</div>
  </div>
  <div class="page-actions">
    ${!R.readOnly && !R.noEdit ? `<button class="btn btn-secondary" onclick="openEditModal('${resource}','${id}')">Edit</button>` : ''}
    ${!R.readOnly ? `<button class="btn btn-danger" onclick="confirmDelete('${resource}','${id}','${esc(item.name||id)}')">Delete</button>` : ''}
  </div>
</div>
<div class="panel">
  <div class="panel-body">
    <div class="json-block">${esc(JSON.stringify(item, null, 2))}</div>
  </div>
</div>`;
}

/* ── Templates ───────────────────────────────────────────────────────────────── */
on('/templates', async () => {
  const content = document.getElementById('page-content');
  const custom = JSON.parse(localStorage.getItem('aim_templates') || '[]');
  const all = [...TEMPLATES, ...custom];
  content.innerHTML = `
<div class="page-header">
  <div><div class="page-title">📐 Templates</div><div class="page-subtitle">Pre-built request templates for common scenarios</div></div>
</div>
<div class="templates-grid">
  ${all.map((t,i) => {
    const R = RESOURCES[t.resource];
    return `
    <div class="template-card">
      <div class="template-card-header">
        <span style="font-size:18px">${R?.icon||'📄'}</span>
        <span class="template-card-tag">${esc(R?.label||t.resource)}</span>
      </div>
      <div class="template-card-name">${esc(t.name)}</div>
      <div class="template-card-desc">${esc(t.desc)}</div>
      <div class="template-card-actions">
        <button class="btn btn-primary btn-sm" onclick="useTemplate(${i})">Use Template</button>
        <button class="btn btn-ghost btn-sm" onclick="previewTemplate(${i})">Preview</button>
      </div>
    </div>`;
  }).join('')}
  <div class="template-card" style="border-style:dashed;cursor:pointer" onclick="openCreateTemplateModal()">
    <div style="text-align:center;padding:16px;color:var(--text-muted)">
      <div style="font-size:32px;margin-bottom:8px">+</div>
      <div style="font-weight:600">Custom Template</div>
      <div style="font-size:12px;margin-top:4px">Save your own request body</div>
    </div>
  </div>
</div>`;
});

function useTemplate(idx) {
  const all = [...TEMPLATES, ...JSON.parse(localStorage.getItem('aim_templates')||'[]')];
  const t = all[idx];
  if (!t) return;
  openCreateModal(t.resource, t.body);
}
function previewTemplate(idx) {
  const all = [...TEMPLATES, ...JSON.parse(localStorage.getItem('aim_templates')||'[]')];
  const t = all[idx];
  openModal(`Template: ${t.name}`, `<div class="json-block">${esc(JSON.stringify(t.body, null, 2))}</div>`,
    `<button class="btn btn-secondary" onclick="closeModal()">Close</button>
     <button class="btn btn-primary" onclick="closeModal();useTemplate(${idx})">Use</button>`);
}
function openCreateTemplateModal() {
  const resOptions = Object.entries(RESOURCES).map(([k,r])=>`<option value="${k}">${r.label}</option>`).join('');
  openModal('Save Custom Template', `
    <div class="form-group"><label class="form-label">Resource</label>
      <select class="form-control" id="tpl-resource">${resOptions}</select></div>
    <div class="form-group"><label class="form-label">Template Name</label>
      <input class="form-control" id="tpl-name" placeholder="My Template"/></div>
    <div class="form-group"><label class="form-label">Description</label>
      <input class="form-control" id="tpl-desc" placeholder="What this template does"/></div>
    <div class="form-group"><label class="form-label">JSON Body</label>
      <textarea class="form-control" id="tpl-body" rows="6" placeholder='{"name":"..."}'></textarea></div>`,
    `<button class="btn btn-secondary" onclick="closeModal()">Cancel</button>
     <button class="btn btn-primary" onclick="saveCustomTemplate()">Save</button>`);
}
function saveCustomTemplate() {
  try {
    const body = JSON.parse(document.getElementById('tpl-body').value);
    const t = { resource: document.getElementById('tpl-resource').value, name: document.getElementById('tpl-name').value, desc: document.getElementById('tpl-desc').value, body };
    const arr = JSON.parse(localStorage.getItem('aim_templates')||'[]');
    arr.push(t);
    localStorage.setItem('aim_templates', JSON.stringify(arr));
    closeModal();
    toast('Template saved', 'success');
    route();
  } catch(e) { toast('Invalid JSON body: '+e.message, 'error'); }
}

/* ── Analytics ───────────────────────────────────────────────────────────────── */
on('/analytics', async () => {
  const content = document.getElementById('page-content');
  const methods = {GET:0,POST:0,PATCH:0,DELETE:0};
  const statuses = {'2xx':0,'4xx':0,'5xx':0};
  const byResource = {};
  S.metrics.forEach(m => {
    methods[m.method] = (methods[m.method]||0)+1;
    if (m.status>=500) statuses['5xx']++;
    else if (m.status>=400) statuses['4xx']++;
    else statuses['2xx']++;
    const seg = m.path.split('/')[1]||'other';
    byResource[seg] = (byResource[seg]||0)+1;
  });
  const totalCalls = S.metrics.length;
  const totalBytes = S.metrics.reduce((a,m)=>a+(m.bytes||0),0);
  const tokenEst = Math.round(totalBytes/4);
  const msArr = S.metrics.map(m=>m.ms).filter(Boolean).sort((a,b)=>a-b);
  const p50 = msArr[Math.floor(msArr.length*.5)]||0;
  const p95 = msArr[Math.floor(msArr.length*.95)]||0;
  const p99 = msArr[Math.floor(msArr.length*.99)]||0;

  content.innerHTML = `
<div class="page-header">
  <div><div class="page-title">📊 Analytics</div><div class="page-subtitle">API usage statistics and performance metrics</div></div>
  <div class="page-actions">
    <button class="btn btn-secondary" onclick="exportMetrics()">⬇ Export CSV</button>
    <button class="btn btn-ghost btn-sm" onclick="if(confirm('Clear all metrics?')){S.metrics=[];localStorage.removeItem('aim_metrics');route()}">Clear</button>
  </div>
</div>

<div style="display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin-bottom:24px">
  ${[['Total Calls',totalCalls,'var(--accent)'],['Data Transferred',fmtBytes(totalBytes),'var(--success)'],['Est. Tokens',tokenEst.toLocaleString(),'var(--purple)'],['Session Cost','$'+(tokenEst/1000*0.005).toFixed(4),'var(--warning)']].map(([l,v,c])=>`
  <div class="panel"><div class="panel-body" style="text-align:center">
    <div style="font-size:28px;font-weight:700;color:${c}">${v}</div>
    <div style="font-size:12px;color:var(--text-secondary);margin-top:4px">${l}</div>
  </div></div>`).join('')}
</div>

<div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:24px">
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Request Methods</span></div>
    <div class="panel-body">
      ${Object.entries(methods).map(([m,c])=>`
      <div style="display:flex;align-items:center;gap:12px;margin-bottom:10px">
        <span class="method-pill method-${m}" style="min-width:56px;text-align:center">${m}</span>
        <div style="flex:1;background:var(--bg-inset);border-radius:4px;height:8px;overflow:hidden">
          <div style="height:100%;background:var(--accent);width:${totalCalls?Math.round(c/totalCalls*100):0}%;border-radius:4px"></div>
        </div>
        <span style="font-size:13px;min-width:32px;text-align:right">${c}</span>
      </div>`).join('')}
    </div>
  </div>
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Response Status Distribution</span></div>
    <div class="panel-body">
      ${[['2xx Success','var(--success)',statuses['2xx']],['4xx Client Error','var(--warning)',statuses['4xx']],['5xx Server Error','var(--danger)',statuses['5xx']]].map(([l,c,n])=>`
      <div style="display:flex;align-items:center;gap:12px;margin-bottom:10px">
        <span style="font-size:12px;color:${c};min-width:100px">${l}</span>
        <div style="flex:1;background:var(--bg-inset);border-radius:4px;height:8px;overflow:hidden">
          <div style="height:100%;background:${c};width:${totalCalls?Math.round(n/totalCalls*100):0}%;border-radius:4px"></div>
        </div>
        <span style="font-size:13px;min-width:32px;text-align:right">${n}</span>
      </div>`).join('')}
    </div>
  </div>
</div>

<div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:24px">
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Latency Percentiles</span></div>
    <div class="panel-body">
      ${[['P50 (median)', p50],['P95', p95],['P99', p99]].map(([l,v])=>`
      <div style="display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid var(--border-muted)">
        <span style="color:var(--text-secondary);font-size:13px">${l}</span>
        <span style="font-weight:700">${v} ms</span>
      </div>`).join('')}
    </div>
  </div>
  <div class="panel">
    <div class="panel-header"><span class="panel-title">Calls by Resource</span></div>
    <div class="panel-body">
      ${Object.entries(byResource).sort((a,b)=>b[1]-a[1]).slice(0,6).map(([r,c])=>`
      <div style="display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid var(--border-muted)">
        <span style="color:var(--text-secondary);font-size:13px;font-family:monospace">${esc(r)}</span>
        <span style="font-weight:700">${c}</span>
      </div>`).join('')}
    </div>
  </div>
</div>

<div class="panel">
  <div class="panel-header"><span class="panel-title">Request Log</span></div>
  <div class="activity-feed">
    ${[...S.metrics].reverse().slice(0,50).map(m=>`
    <div class="activity-item">
      <span class="activity-time">${fmtTime(m.ts)}</span>
      <span class="method-pill method-${m.method}">${m.method}</span>
      <span class="activity-resource">${esc(m.path)}</span>
      <span style="color:var(--text-muted);font-size:11px">${m.ms}ms · ${fmtBytes(m.bytes)}</span>
      ${fmtStatus(m.status)}
    </div>`).join('') || '<div style="padding:24px;text-align:center;color:var(--text-muted)">No requests yet</div>'}
  </div>
</div>`;
});

function exportMetrics() {
  const rows = ['timestamp,method,path,status,ms,bytes'];
  S.metrics.forEach(m => rows.push(`${new Date(m.ts).toISOString()},${m.method},${m.path},${m.status},${m.ms},${m.bytes}`));
  const blob = new Blob([rows.join('\n')], {type:'text/csv'});
  const a = document.createElement('a'); a.href = URL.createObjectURL(blob); a.download = 'tmf915-metrics.csv'; a.click();
}

/* ── Onboarding ──────────────────────────────────────────────────────────────── */
on('/onboarding', async () => {
  const step = parseInt(localStorage.getItem('aim_ob_step')||'0');
  const content = document.getElementById('page-content');
  const steps = ['Welcome','Configure','Test','First Resource','Done'];

  const stepBar = steps.map((s,i) => {
    const cls = i < step ? 'done' : i === step ? 'active' : '';
    return `<li class="wizard-step ${cls}"><span class="wizard-step-num">${i<step?'✓':i+1}</span><span>${s}</span></li>`;
  }).join('');

  let body = '';
  if (step === 0) body = `
    <h2>Welcome to TMF915 AiM Suite 👋</h2>
    <p>This management console gives you full control over your AI Management API. You can manage AI models, contracts, monitors, alarms, and more — all from one place.</p>
    <ul class="feature-list">
      <li>Manage all 12 TMF915 resource types</li>
      <li>Real-time analytics and token cost estimation</li>
      <li>Pre-built templates for common scenarios</li>
      <li>Webhook hub management</li>
      <li>Export metrics as CSV</li>
    </ul>
    <div class="form-actions">
      <button class="btn btn-primary" onclick="obNext()">Get Started →</button>
    </div>`;
  else if (step === 1) body = `
    <h2>Configure API Connection</h2>
    <p>Enter your TMF915 API endpoint and API key. These are saved to your browser's localStorage.</p>
    <div class="form-group">
      <label class="form-label">API Base URL</label>
      <div class="form-hint">Default: http://localhost:8080/tmf-api/AiM/v4</div>
      <input class="form-control" id="ob-base" value="${esc(S.apiBase)}" placeholder="http://localhost:8080/tmf-api/AiM/v4"/>
    </div>
    <div class="form-group">
      <label class="form-label">API Key (X-API-Key)</label>
      <input class="form-control" id="ob-key" type="text" value="${esc(S.apiKey)}" placeholder="your-api-key"/>
    </div>
    <div class="form-actions">
      <button class="btn btn-secondary" onclick="obPrev()">← Back</button>
      <button class="btn btn-primary" onclick="obSaveCreds()">Save & Continue →</button>
    </div>`;
  else if (step === 2) {
    let healthStatus = '…testing…';
    body = `
    <h2>Test Connection</h2>
    <p>Let's verify your API is reachable and your key works.</p>
    <div id="ob-test-result" style="padding:16px;background:var(--bg-inset);border-radius:var(--radius);font-family:monospace;font-size:13px;margin-bottom:16px">Testing…</div>
    <div class="form-actions">
      <button class="btn btn-secondary" onclick="obPrev()">← Back</button>
      <button class="btn btn-primary" id="ob-continue" disabled onclick="obNext()">Continue →</button>
    </div>`;
    setTimeout(async () => {
      const el = document.getElementById('ob-test-result');
      if (!el) return;
      try {
        const base = S.apiBase.replace('/tmf-api/AiM/v4','');
        const hr = await fetch(base+'/health');
        const hok = hr.ok;
        const ar = await api('GET', '/aiModel');
        el.innerHTML = hok
          ? `<span style="color:var(--success)">✓ API is online</span><br><span style="color:var(--success)">✓ API key accepted (${ar.status})</span><br><span style="color:var(--text-secondary)">Connected to ${esc(S.apiBase)}</span>`
          : `<span style="color:var(--danger)">✗ Health check failed</span>`;
        const cont = document.getElementById('ob-continue');
        if (cont && hok && ar.ok) cont.disabled = false;
        else if (cont) cont.disabled = false;
      } catch(e) {
        if(el) el.innerHTML = `<span style="color:var(--danger)">✗ Connection failed: ${esc(e.message)}</span><br><span style="color:var(--text-secondary)">Make sure the server is running and your config is correct.</span>`;
        const cont = document.getElementById('ob-continue');
        if (cont) cont.disabled = false;
      }
    }, 300);
  } else if (step === 3) body = `
    <h2>Create Your First AI Model</h2>
    <p>Let's create an AI model record to make sure everything works end-to-end.</p>
    <div class="form-group">
      <label class="form-label">Name <span style="color:var(--danger)">*</span></label>
      <input class="form-control" id="ob-model-name" value="My First AI Model" placeholder="Model name"/>
    </div>
    <div class="form-group">
      <label class="form-label">Description</label>
      <textarea class="form-control" id="ob-model-desc" rows="2">Created during onboarding</textarea>
    </div>
    <div class="form-group">
      <label class="form-label">State</label>
      <select class="form-control" id="ob-model-state">
        <option value="active">active</option><option value="draft">draft</option>
      </select>
    </div>
    <div class="form-actions">
      <button class="btn btn-secondary" onclick="obPrev()">← Back</button>
      <button class="btn btn-ghost" onclick="obNext()">Skip</button>
      <button class="btn btn-primary" onclick="obCreateModel()">Create & Continue →</button>
    </div>`;
  else body = `
    <h2>You're all set! 🎉</h2>
    <p>Your TMF915 AI Management console is ready. Here's what you can do next:</p>
    <ul class="feature-list">
      <li>Browse and manage all AI resources from the sidebar</li>
      <li>Use Templates for quick resource creation</li>
      <li>Monitor usage in Analytics</li>
      <li>Register webhooks via Hubs</li>
    </ul>
    <div class="form-actions">
      <button class="btn btn-primary" onclick="localStorage.setItem('aim_ob_step','0');location.hash='/'">Go to Dashboard →</button>
    </div>`;

  content.innerHTML = `
<div style="max-width:700px;margin:0 auto">
  <div class="page-title" style="margin-bottom:24px">Onboarding</div>
  <ul class="wizard-steps">${stepBar}</ul>
  <div class="wizard-card">${body}</div>
</div>`;
  if (step === 2) {} // test runs via setTimeout above
});

function obNext() { localStorage.setItem('aim_ob_step', String(parseInt(localStorage.getItem('aim_ob_step')||'0')+1)); route(); }
function obPrev() { localStorage.setItem('aim_ob_step', String(Math.max(0,parseInt(localStorage.getItem('aim_ob_step')||'0')-1))); route(); }
function obSaveCreds() {
  const base = document.getElementById('ob-base').value.trim().replace(/\/$/,'');
  const key  = document.getElementById('ob-key').value.trim();
  if (!base) { toast('API Base URL is required','error'); return; }
  S.apiBase = base; S.apiKey = key;
  localStorage.setItem('aim_base', base); localStorage.setItem('aim_key', key);
  obNext();
}
async function obCreateModel() {
  const name = document.getElementById('ob-model-name').value.trim();
  const desc = document.getElementById('ob-model-desc').value.trim();
  const state = document.getElementById('ob-model-state').value;
  try {
    await apiCreate('aiModel', {name, description:desc, state});
    toast('AI Model created!','success');
    obNext();
  } catch(e) { toast(e.message,'error'); }
}

/* ── Settings ────────────────────────────────────────────────────────────────── */
on('/settings', () => {
  const content = document.getElementById('page-content');
  content.innerHTML = `
<div class="page-header">
  <div><div class="page-title">⚙ Settings</div><div class="page-subtitle">Configure your management console</div></div>
</div>

<div class="settings-section">
  <h3>API Connection</h3>
  <div class="settings-row">
    <div class="settings-row-label"><strong>Base URL</strong><span>TMF915 API endpoint</span></div>
    <div class="settings-row-control">
      <input class="form-control" id="set-base" value="${esc(S.apiBase)}" style="max-width:480px"/>
    </div>
  </div>
  <div class="settings-row">
    <div class="settings-row-label"><strong>API Key</strong><span>X-API-Key header value</span></div>
    <div class="settings-row-control">
      <input class="form-control" id="set-key" type="text" value="${esc(S.apiKey)}" style="max-width:480px"/>
    </div>
  </div>
  <div class="settings-row">
    <div class="settings-row-label"></div>
    <div class="settings-row-control">
      <button class="btn btn-primary" onclick="saveSettings()">Save Connection</button>
      <button class="btn btn-secondary" onclick="testConnection()" style="margin-left:8px">Test Connection</button>
    </div>
  </div>
  <div id="conn-result" style="margin-top:8px"></div>
</div>

<div class="settings-section">
  <h3>Appearance</h3>
  <div class="settings-row">
    <div class="settings-row-label"><strong>Theme</strong><span>Light or dark color scheme</span></div>
    <div class="settings-row-control">
      <select class="form-control" id="set-theme" style="max-width:160px" onchange="applyTheme(this.value)">
        <option value="dark" ${S.theme==='dark'?'selected':''}>Dark</option>
        <option value="light" ${S.theme==='light'?'selected':''}>Light</option>
      </select>
    </div>
  </div>
</div>

<div class="settings-section">
  <h3>Data</h3>
  <div class="settings-row">
    <div class="settings-row-label"><strong>Metrics</strong><span>${S.metrics.length} recorded requests</span></div>
    <div class="settings-row-control">
      <button class="btn btn-danger" onclick="if(confirm('Clear all metrics?')){S.metrics=[];localStorage.removeItem('aim_metrics');toast('Metrics cleared','success');route()}">Clear Metrics</button>
    </div>
  </div>
  <div class="settings-row">
    <div class="settings-row-label"><strong>Custom Templates</strong><span>${JSON.parse(localStorage.getItem('aim_templates')||'[]').length} saved</span></div>
    <div class="settings-row-control">
      <button class="btn btn-danger" onclick="if(confirm('Delete all custom templates?')){localStorage.removeItem('aim_templates');toast('Templates cleared','success');route()}">Clear Templates</button>
    </div>
  </div>
  <div class="settings-row">
    <div class="settings-row-label"><strong>Onboarding</strong><span>Reset the onboarding wizard</span></div>
    <div class="settings-row-control">
      <button class="btn btn-secondary" onclick="localStorage.setItem('aim_ob_step','0');location.hash='/onboarding'">Reset Onboarding</button>
    </div>
  </div>
</div>`;
});

function saveSettings() {
  const base = document.getElementById('set-base').value.trim().replace(/\/$/,'');
  const key  = document.getElementById('set-key').value.trim();
  S.apiBase = base; S.apiKey = key;
  localStorage.setItem('aim_base', base); localStorage.setItem('aim_key', key);
  toast('Settings saved','success');
}
async function testConnection() {
  const el = document.getElementById('conn-result');
  el.innerHTML = '<span style="color:var(--text-muted)">Testing…</span>';
  try {
    const base = document.getElementById('set-base').value.trim().replace(/\/$/,'');
    const key  = document.getElementById('set-key').value.trim();
    const hr = await fetch(base.replace('/tmf-api/AiM/v4','')+'/health');
    const ar = await fetch(base+'/aiModel', {headers:{'X-API-Key':key}});
    el.innerHTML = hr.ok && ar.ok
      ? `<span style="color:var(--success)">✓ Connected successfully</span>`
      : `<span style="color:var(--danger)">✗ Health: ${hr.status} | Auth: ${ar.status}</span>`;
  } catch(e) { el.innerHTML = `<span style="color:var(--danger)">✗ ${esc(e.message)}</span>`; }
}

/* ── Theme ───────────────────────────────────────────────────────────────────── */
function applyTheme(t) {
  S.theme = t;
  localStorage.setItem('aim_theme', t);
  document.documentElement.setAttribute('data-theme', t);
  document.getElementById('theme-icon-dark').style.display = t==='light'?'none':'inline';
  document.getElementById('theme-icon-light').style.display = t==='light'?'inline':'none';
  document.getElementById('theme-label').textContent = t==='light'?'Dark mode':'Light mode';
}

/* ── Bootstrap ───────────────────────────────────────────────────────────────── */
document.addEventListener('DOMContentLoaded', () => {
  applyTheme(S.theme);

  document.getElementById('theme-toggle').addEventListener('click', () => {
    applyTheme(S.theme === 'dark' ? 'light' : 'dark');
  });

  document.querySelectorAll('.nav-item[data-route]').forEach(el => {
    el.addEventListener('click', e => { e.preventDefault(); location.hash = el.dataset.route; });
  });

  document.getElementById('modal-close').addEventListener('click', closeModal);
  document.getElementById('modal-overlay').addEventListener('click', e => {
    if (e.target === document.getElementById('modal-overlay')) closeModal();
  });

  window.addEventListener('hashchange', route);
  checkHealth();
  setInterval(checkHealth, 30000);

  if (!S.apiKey && !localStorage.getItem('aim_ob_step')) {
    location.hash = '/onboarding';
  } else {
    route();
  }
});
