package report

// CoverageHTMLTemplate is the standalone HTML template for the coverage dashboard.
const CoverageHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>req42-tracer: Coverage Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
               background: #f5f5f5; color: #333; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                  color: white; padding: 24px 32px; }
        .header h1 { font-size: 24px; margin-bottom: 4px; }
        .header p  { opacity: 0.85; font-size: 14px; }
        .container { max-width: 1200px; margin: 24px auto; padding: 0 24px; }
        .summary-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px,1fr));
                         gap: 16px; margin-bottom: 24px; }
        .card { background: white; border-radius: 6px; padding: 20px;
                box-shadow: 0 1px 4px rgba(0,0,0,.08); border-left: 4px solid #667eea; }
        .card.good    { border-left-color: #27ae60; }
        .card.warning { border-left-color: #f39c12; }
        .card.danger  { border-left-color: #e74c3c; }
        .card-label { font-size: 11px; text-transform: uppercase; color: #888; margin-bottom: 6px; }
        .card-value { font-size: 28px; font-weight: 700; }
        .controls { display: flex; gap: 12px; margin-bottom: 16px; flex-wrap: wrap; align-items: center; }
        .controls select, .controls input {
            padding: 7px 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 13px; }
        table { width: 100%; border-collapse: collapse; background: white;
                border-radius: 6px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,.08); }
        th { background: #f0f0f0; padding: 10px 12px; text-align: left; font-size: 12px;
             text-transform: uppercase; color: #666; cursor: pointer; white-space: nowrap; }
        th:hover { background: #e8e8e8; }
        td { padding: 9px 12px; border-bottom: 1px solid #f0f0f0; font-size: 13px; }
        tr:last-child td { border-bottom: none; }
        tr:hover td { background: #fafafa; }
        .bar-wrap { width: 120px; background: #eee; border-radius: 3px; height: 8px; display: inline-block; }
        .bar-fill  { height: 8px; border-radius: 3px; }
        .good    .bar-fill, .pct-good    { background: #27ae60; color: #27ae60; }
        .warning .bar-fill, .pct-warning { background: #f39c12; color: #f39c12; }
        .danger  .bar-fill, .pct-danger  { background: #e74c3c; color: #e74c3c; }
        .arch-badge { background: #7ed321; color: #fff; border-radius: 3px;
                      padding: 1px 6px; font-size: 11px; white-space: nowrap; }
        .no-arch { color: #ccc; font-size: 11px; }
        tfoot td { font-weight: 700; background: #f8f9fa; border-top: 2px solid #ddd; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Coverage Dashboard</h1>
        <p>Statement-level test coverage mapped to architecture components</p>
    </div>
    <div class="container">
        <div class="summary-cards" id="summary-cards"></div>
        <div class="controls">
            <select id="filter-level">
                <option value="">All levels</option>
                <option value="danger">Danger (&lt;70%)</option>
                <option value="warning">Warning (70–80%)</option>
                <option value="good">Good (≥80%)</option>
            </select>
            <input id="search" type="text" placeholder="Search package or arch…">
            <span id="row-count" style="font-size:13px;color:#888"></span>
        </div>
        <table>
            <thead>
                <tr>
                    <th onclick="sortBy('arch_id')">Arch Component ↕</th>
                    <th onclick="sortBy('package')">Package ↕</th>
                    <th onclick="sortBy('statements')" style="text-align:right">Statements ↕</th>
                    <th onclick="sortBy('covered')" style="text-align:right">Covered ↕</th>
                    <th onclick="sortBy('pct')">Coverage ↕</th>
                </tr>
            </thead>
            <tbody id="tbody"></tbody>
            <tfoot id="tfoot"></tfoot>
        </table>
    </div>
<script>
    const data = <!--COVERAGE_DATA_JSON-->;
    let sortCol = 'pct', sortAsc = true;

    function sortBy(col) {
        sortAsc = sortCol === col ? !sortAsc : true;
        sortCol = col;
        render();
    }

    function render() {
        const level  = document.getElementById('filter-level').value;
        const search = document.getElementById('search').value.toLowerCase();

        let rows = (data.rows || []).filter(r => {
            if (level  && r.level !== level) return false;
            if (search && !r.package.toLowerCase().includes(search) &&
                         !r.arch_id.toLowerCase().includes(search) &&
                         !r.arch_title.toLowerCase().includes(search)) return false;
            return true;
        });

        rows.sort((a, b) => {
            const va = a[sortCol] ?? '', vb = b[sortCol] ?? '';
            if (typeof va === 'number') return sortAsc ? va - vb : vb - va;
            return sortAsc ? String(va).localeCompare(String(vb)) : String(vb).localeCompare(String(va));
        });

        document.getElementById('row-count').textContent = rows.length + ' package(s)';

        document.getElementById('tbody').innerHTML = rows.map(r => {
            const pctStr = r.pct.toFixed(1) + '%';
            const bar = '<div class="bar-wrap"><div class="bar-fill ' + r.level + '" style="width:' + Math.min(r.pct,100).toFixed(1) + '%"></div></div>';
            const arch = r.arch_id
                ? '<span class="arch-badge">' + esc(r.arch_id) + '</span>'
                : '<span class="no-arch">—</span>';
            return '<tr class="' + r.level + '">' +
                '<td>' + arch + (r.arch_title ? ' <small style="color:#888">' + esc(r.arch_title) + '</small>' : '') + '</td>' +
                '<td style="font-family:monospace">' + esc(r.package) + '</td>' +
                '<td style="text-align:right">' + r.statements + '</td>' +
                '<td style="text-align:right">' + r.covered + '</td>' +
                '<td><span class="pct-' + r.level + '" style="font-weight:600;margin-right:8px">' + pctStr + '</span>' + bar + '</td>' +
                '</tr>';
        }).join('');

        const totStmts = rows.reduce((s,r) => s + r.statements, 0);
        const totCov   = rows.reduce((s,r) => s + r.covered, 0);
        const totPct   = totStmts > 0 ? (totCov / totStmts * 100).toFixed(1) : '0.0';
        document.getElementById('tfoot').innerHTML =
            '<tr><td colspan="2"><strong>Total (' + rows.length + ' packages)</strong></td>' +
            '<td style="text-align:right"><strong>' + totStmts + '</strong></td>' +
            '<td style="text-align:right"><strong>' + totCov + '</strong></td>' +
            '<td><strong>' + totPct + '%</strong></td></tr>';
    }

    function renderCards() {
        const lvlCount = { good: 0, warning: 0, danger: 0 };
        (data.rows || []).forEach(r => { if (lvlCount[r.level] !== undefined) lvlCount[r.level]++; });
        document.getElementById('summary-cards').innerHTML =
            card(data.overall_level, 'Overall Coverage', data.overall_pct.toFixed(1) + '%') +
            card('good',    'Packages ≥80%',  lvlCount.good) +
            card('warning', 'Packages 70–80%', lvlCount.warning) +
            card('danger',  'Packages <70%',  lvlCount.danger) +
            card('', 'Total Statements', data.total_stmts) +
            card('', 'Total Covered',    data.total_covered);
    }

    function card(level, label, value) {
        return '<div class="card ' + level + '"><div class="card-label">' + label + '</div>' +
               '<div class="card-value">' + value + '</div></div>';
    }

    function esc(s) {
        return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
    }

    document.getElementById('filter-level').addEventListener('change', render);
    document.getElementById('search').addEventListener('input', render);

    renderCards();
    render();
</script>
</body>
</html>
`
