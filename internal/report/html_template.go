package report

// HTMLTemplate is the embedded HTML template for graph visualization with D3.js.
const HTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>req42-tracer: Traceability Report</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: #f5f5f5;
            color: #333;
        }

        .container {
            display: flex;
            height: 100vh;
        }

        .sidebar {
            width: 250px;
            background: #2c3e50;
            color: white;
            padding: 20px;
            overflow-y: auto;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }

        .sidebar h2 {
            font-size: 18px;
            margin-bottom: 20px;
            text-transform: uppercase;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }

        .sidebar h3 {
            font-size: 14px;
            margin-top: 15px;
            margin-bottom: 10px;
            color: #ecf0f1;
        }

        .legend {
            display: flex;
            flex-direction: column;
            gap: 8px;
            margin-bottom: 20px;
        }

        .legend-item {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 13px;
        }

        .legend-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
        }

        .legend-dot.requirement { background: #4a90e2; }
        .legend-dot.arch { background: #7ed321; }
        .legend-dot.test-spec { background: #f5a623; }
        .legend-dot.test-code { background: #f5a623; opacity: 0.7; }
        .legend-dot.test-result { background: #999; }

        .stats {
            background: rgba(255,255,255,0.1);
            border-radius: 4px;
            padding: 12px;
            font-size: 12px;
            margin-bottom: 20px;
        }

        .stats-item {
            display: flex;
            justify-content: space-between;
            margin: 4px 0;
        }

        .controls {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }

        .control-group {
            display: flex;
            flex-direction: column;
            gap: 6px;
        }

        .control-group label {
            font-size: 12px;
            text-transform: uppercase;
            color: #bdc3c7;
        }

        button {
            background: #3498db;
            color: white;
            border: none;
            padding: 8px 12px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 12px;
            text-transform: uppercase;
            transition: background 0.3s;
        }

        button:hover {
            background: #2980b9;
        }

        button.active {
            background: #e74c3c;
        }

        .main {
            flex: 1;
            display: flex;
            flex-direction: column;
        }

        .header {
            background: white;
            padding: 20px;
            border-bottom: 1px solid #e0e0e0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }

        .header h1 {
            font-size: 24px;
            margin-bottom: 5px;
        }

        .header p {
            color: #666;
            font-size: 14px;
        }

        .tabs {
            display: flex;
            gap: 0;
            border-bottom: 2px solid #e0e0e0;
            margin-top: 15px;
        }

        .tab-button {
            background: transparent;
            border: none;
            padding: 12px 24px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            color: #666;
            border-bottom: 3px solid transparent;
            transition: all 0.3s;
        }

        .tab-button:hover {
            color: #333;
        }

        .tab-button.active {
            color: #3498db;
            border-bottom-color: #3498db;
        }

        .tab-content {
            display: none;
            flex: 1;
            background: #fff;
            position: relative;
            overflow: auto;
        }

        .tab-content.active {
            display: flex;
        }

        #graph {
            flex: 1;
            background: #fff;
            position: relative;
        }

        #matrix {
            flex: 1;
            flex-direction: column;
        }

        .matrix-controls {
            padding: 20px;
            background: #f8f9fa;
            border-bottom: 1px solid #e0e0e0;
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
            align-items: center;
        }

        .matrix-controls input,
        .matrix-controls select {
            padding: 8px 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 13px;
        }

        .matrix-table-wrapper {
            flex: 1;
            overflow: auto;
            padding: 20px;
        }

        .matrix-table {
            border-collapse: collapse;
            font-size: 12px;
            background: white;
        }

        .matrix-table th,
        .matrix-table td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: center;
        }

        .matrix-table th {
            background: #f0f0f0;
            font-weight: 600;
            position: sticky;
            top: 0;
        }

        .matrix-table td.req-id {
            text-align: left;
            background: #f8f9fa;
            font-weight: 500;
            min-width: 150px;
        }

        .matrix-cell {
            cursor: pointer;
            transition: background 0.2s;
        }

        .matrix-cell.covered {
            background: #d4edda;
            color: #155724;
        }

        .matrix-cell.missing {
            background: #f8d7da;
            color: #721c24;
        }

        .matrix-cell.stale {
            background: #fff3cd;
            color: #856404;
        }

        .matrix-stats {
            padding: 20px;
            background: #f8f9fa;
            border-top: 1px solid #e0e0e0;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }

        .stat-item {
            background: white;
            padding: 12px;
            border-radius: 4px;
            border-left: 4px solid #3498db;
        }

        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
        }

        .stat-value {
            font-size: 18px;
            font-weight: bold;
            margin-top: 4px;
        }

        svg {
            width: 100%;
            height: 100%;
        }

        .node {
            stroke: #fff;
            stroke-width: 2px;
            cursor: pointer;
        }

        .node:hover {
            stroke: #000;
            stroke-width: 3px;
        }

        .node.selected {
            stroke: #f5a623;
            stroke-width: 3px;
            filter: drop-shadow(0 0 6px rgba(245, 166, 35, 0.8));
        }

        .node-label {
            font-size: 11px;
            pointer-events: none;
            text-anchor: middle;
            dominant-baseline: central;
            font-weight: 500;
        }

        .link {
            stroke: #bbb;
            stroke-width: 1.5px;
            stroke-dasharray: 0;
            opacity: 0.6;
        }

        .link.faded {
            opacity: 0.1;
        }

        .link-label {
            font-size: 10px;
            pointer-events: none;
            fill: #666;
            background: white;
        }

        .tooltip {
            position: absolute;
            background: rgba(0, 0, 0, 0.9);
            color: white;
            padding: 12px;
            border-radius: 4px;
            font-size: 12px;
            pointer-events: none;
            z-index: 1000;
            max-width: 300px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            display: none;
        }

        .tooltip.visible {
            display: block;
        }

        .tooltip-title {
            font-weight: bold;
            margin-bottom: 6px;
            word-break: break-word;
        }

        .tooltip-content {
            font-size: 11px;
            line-height: 1.5;
        }

        .tooltip-meta {
            margin-top: 6px;
            padding-top: 6px;
            border-top: 1px solid rgba(255,255,255,0.2);
            font-style: italic;
        }

        .zoom-controls {
            position: absolute;
            top: 20px;
            right: 20px;
            display: flex;
            gap: 8px;
            z-index: 100;
        }

        .zoom-btn {
            background: white;
            color: #333;
            border: 1px solid #ddd;
            width: 36px;
            height: 36px;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 4px;
            cursor: pointer;
            font-size: 18px;
            font-weight: bold;
            transition: all 0.2s;
        }

        .zoom-btn:hover {
            background: #f0f0f0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        @media (max-width: 1024px) {
            .sidebar {
                width: 200px;
            }
        }

        @media (max-width: 768px) {
            .container {
                flex-direction: column;
            }

            .sidebar {
                width: 100%;
                max-height: 200px;
            }
        }

        /* ASPICE Dashboard */
        #aspice {
            flex: 1;
            flex-direction: column;
            overflow-y: auto;
            background: #f8f9fa;
        }

        .aspice-overview {
            background: white;
            border-bottom: 1px solid #e0e0e0;
            padding: 24px 30px;
        }

        .aspice-overview h2 {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 16px;
            color: #333;
        }

        .overall-coverage {
            display: flex;
            align-items: center;
            gap: 24px;
        }

        .overall-pct {
            font-size: 42px;
            font-weight: bold;
            min-width: 90px;
        }

        .overall-pct.good { color: #27ae60; }
        .overall-pct.warning { color: #f39c12; }
        .overall-pct.danger { color: #e74c3c; }

        .overall-bar-wrap { flex: 1; }

        .overall-bar-label {
            font-size: 12px;
            color: #888;
            text-transform: uppercase;
            margin-bottom: 6px;
        }

        .overall-bar-track {
            height: 10px;
            background: #e0e0e0;
            border-radius: 5px;
            overflow: hidden;
        }

        .overall-bar-fill {
            height: 100%;
            border-radius: 5px;
            transition: width 0.6s ease;
        }

        .overall-bar-fill.good { background: #27ae60; }
        .overall-bar-fill.warning { background: #f39c12; }
        .overall-bar-fill.danger { background: #e74c3c; }

        .aspice-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
            gap: 16px;
            padding: 20px;
        }

        .process-card {
            background: white;
            border-radius: 8px;
            border: 1px solid #e0e0e0;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }

        .process-card-header {
            padding: 14px 16px;
            display: flex;
            align-items: center;
            gap: 12px;
            border-left: 4px solid #3498db;
        }

        .process-card-header.good { border-left-color: #27ae60; }
        .process-card-header.warning { border-left-color: #f39c12; }
        .process-card-header.danger { border-left-color: #e74c3c; }

        .process-id-badge {
            background: #2c3e50;
            color: white;
            padding: 3px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
            white-space: nowrap;
        }

        .process-card-name {
            flex: 1;
            font-weight: 600;
            font-size: 13px;
            color: #333;
        }

        .process-card-pct {
            font-size: 16px;
            font-weight: bold;
        }

        .process-card-pct.good { color: #27ae60; }
        .process-card-pct.warning { color: #f39c12; }
        .process-card-pct.danger { color: #e74c3c; }

        .process-progress {
            height: 4px;
            background: #f0f0f0;
        }

        .process-progress-fill {
            height: 100%;
            transition: width 0.5s ease;
        }

        .process-progress-fill.good { background: #27ae60; }
        .process-progress-fill.warning { background: #f39c12; }
        .process-progress-fill.danger { background: #e74c3c; }

        .bp-list { padding: 4px 0; }

        .bp-item {
            padding: 8px 16px;
            border-bottom: 1px solid #f5f5f5;
        }

        .bp-item:last-child { border-bottom: none; }

        .bp-row {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .bp-id {
            font-size: 11px;
            font-weight: bold;
            color: #888;
            min-width: 72px;
        }

        .bp-title {
            flex: 1;
            font-size: 12px;
            color: #444;
        }

        .bp-badge {
            font-size: 11px;
            font-weight: bold;
            padding: 2px 6px;
            border-radius: 3px;
        }

        .bp-badge.good { background: #d5f5e3; color: #1e8449; }
        .bp-badge.warning { background: #fef9e7; color: #b7950b; }
        .bp-badge.danger { background: #fdedec; color: #c0392b; }

        .bp-gaps {
            margin: 4px 0 0 80px;
            padding: 0;
            list-style: none;
        }

        .bp-gaps li {
            font-size: 11px;
            color: #e74c3c;
            padding: 1px 0;
        }

        .bp-gaps li::before {
            content: "→ ";
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="sidebar">
            <h2>Traceability Graph</h2>

            <div class="legend">
                <h3>Node Types</h3>
                <div class="legend-item">
                    <div class="legend-dot requirement"></div>
                    <span>Requirement</span>
                </div>
                <div class="legend-item">
                    <div class="legend-dot arch"></div>
                    <span>Architecture</span>
                </div>
                <div class="legend-item">
                    <div class="legend-dot test-spec"></div>
                    <span>Test Spec</span>
                </div>
                <div class="legend-item">
                    <div class="legend-dot test-code"></div>
                    <span>Test Code</span>
                </div>
                <div class="legend-item">
                    <div class="legend-dot test-result"></div>
                    <span>Test Result</span>
                </div>
            </div>

            <div class="legend">
                <h3>Link Types</h3>
                <div class="legend-item">→ satisfies</div>
                <div class="legend-item">→ implements</div>
                <div class="legend-item">→ verifies</div>
                <div class="legend-item">→ derives</div>
                <div class="legend-item">→ covers</div>
            </div>

            <div class="stats" id="stats">
                <div class="stats-item">
                    <span>Requirements:</span>
                    <span id="stat-reqs">0</span>
                </div>
                <div class="stats-item">
                    <span>Architecture:</span>
                    <span id="stat-arch">0</span>
                </div>
                <div class="stats-item">
                    <span>Test Specs:</span>
                    <span id="stat-specs">0</span>
                </div>
                <div class="stats-item">
                    <span>Test Results:</span>
                    <span id="stat-results">0</span>
                </div>
                <div class="stats-item">
                    <span>Total Nodes:</span>
                    <span id="stat-total">0</span>
                </div>
                <div class="stats-item">
                    <span>Total Links:</span>
                    <span id="stat-links">0</span>
                </div>
            </div>

            <div class="controls">
                <div class="control-group">
                    <label>Filters</label>
                    <button id="btn-reset">Reset View</button>
                </div>
                <div class="control-group">
                    <label>Node Types</label>
                    <button id="btn-filter-all" class="active">All</button>
                    <button id="btn-filter-req">Requirements</button>
                    <button id="btn-filter-arch">Architecture</button>
                    <button id="btn-filter-tests">Tests</button>
                </div>
                <div class="control-group">
                    <label>Layout</label>
                    <button id="btn-layout-force" class="active">Force</button>
                    <button id="btn-layout-radial">Radial</button>
                </div>
            </div>
        </div>

        <div class="main">
            <div class="header">
                <div>
                    <h1>Traceability Report</h1>
                    <p>Interactive dependency graph and traceability matrix</p>
                </div>
                <div class="tabs">
                    <button class="tab-button active" onclick="switchTab('graph')">Graph View</button>
                    <button class="tab-button" onclick="switchTab('matrix')">Matrix View</button>
                    <button class="tab-button" onclick="switchTab('aspice')">ASPICE Dashboard</button>
                </div>
            </div>

            <div id="graph" class="tab-content active">
                <div class="zoom-controls">
                    <button class="zoom-btn" id="btn-zoom-in">+</button>
                    <button class="zoom-btn" id="btn-zoom-out">−</button>
                    <button class="zoom-btn" id="btn-zoom-fit">⊙</button>
                </div>
                <div class="tooltip" id="tooltip"></div>
            </div>

            <div id="matrix" class="tab-content">
                <div class="matrix-controls">
                    <input type="text" id="search-req" placeholder="Search requirements...">
                    <select id="filter-priority">
                        <option value="">All Priorities</option>
                        <option value="high">High</option>
                        <option value="medium">Medium</option>
                        <option value="low">Low</option>
                    </select>
                    <select id="filter-status">
                        <option value="">All Status</option>
                        <option value="approved">Approved</option>
                        <option value="draft">Draft</option>
                        <option value="deprecated">Deprecated</option>
                    </select>
                    <button onclick="exportMatrixCSV()" style="margin-left: auto;">📥 Export CSV</button>
                </div>
                <div class="matrix-table-wrapper" id="matrix-table-container">
                    <!-- Matrix table will be generated by JavaScript -->
                </div>
                <div class="matrix-stats">
                    <div class="stat-item">
                        <div class="stat-label">Total Requirements</div>
                        <div class="stat-value" id="matrix-total">0</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Covered</div>
                        <div class="stat-value" id="matrix-covered">0</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Missing</div>
                        <div class="stat-value" id="matrix-missing">0</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Coverage</div>
                        <div class="stat-value" id="matrix-coverage">0%</div>
                    </div>
                </div>
            </div>

            <div id="aspice" class="tab-content">
                <div class="aspice-overview">
                    <h2>ASPICE PAM 4.0 Compliance</h2>
                    <div class="overall-coverage">
                        <div class="overall-pct" id="aspice-overall-pct">-</div>
                        <div class="overall-bar-wrap">
                            <div class="overall-bar-label">Overall Coverage</div>
                            <div class="overall-bar-track">
                                <div class="overall-bar-fill" id="aspice-overall-fill" style="width:0%"></div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="aspice-grid" id="aspice-grid"></div>
            </div>
        </div>
    </div>

    <script>
        // Graph data will be injected here: <!--GRAPH_DATA-->
        const graphData = <!--GRAPH_DATA_JSON-->;
        const matrixData = <!--MATRIX_DATA_JSON-->;
        const aspiceData = <!--ASPICE_DATA_JSON-->;

        // Global variables for filtering
        let globalNode = null;
        let globalLink = null;
        let globalLabel = null;

        function initializeGraph(data) {
            // Update statistics
            updateStats(data);

            // Set up SVG
            const container = document.getElementById('graph');
            const width = container.clientWidth;
            const height = container.clientHeight;

            const svg = d3.select('#graph').append('svg');

            // Create groups for zooming
            const g = svg.append('g');
            const linkGroup = g.append('g').attr('class', 'links');
            const nodeGroup = g.append('g').attr('class', 'nodes');
            const labelGroup = g.append('g').attr('class', 'labels');

            // Color mapping
            const colorMap = {
                'requirement': '#4a90e2',
                'arch': '#7ed321',
                'test-spec': '#f5a623',
                'test-code': '#f5a623',
                'test-result': '#999'
            };

            // Simulation
            const simulation = d3.forceSimulation(data.nodes)
                .force('link', d3.forceLink(data.edges)
                    .id(d => d.id)
                    .distance(d => 100 + (d.value * 20)))
                .force('charge', d3.forceManyBody().strength(-300))
                .force('center', d3.forceCenter(width / 2, height / 2))
                .force('collision', d3.forceCollide(25));

            // Create links (store globally for filtering)
            globalLink = linkGroup.selectAll('line')
                .data(data.edges)
                .enter()
                .append('line')
                .attr('class', 'link')
                .attr('stroke-width', d => Math.sqrt(d.value))
                .on('mouseover', function() {
                    d3.select(this).classed('active', true);
                })
                .on('mouseout', function() {
                    d3.select(this).classed('active', false);
                });

            // Create nodes (store globally for filtering)
            globalNode = nodeGroup.selectAll('circle')
                .data(data.nodes)
                .enter()
                .append('circle')
                .attr('class', 'node')
                .attr('r', 12)
                .attr('fill', d => colorMap[d.type] || '#999')
                .call(drag(simulation))
                .on('mouseover', showTooltip)
                .on('mouseout', hideTooltip)
                .on('click', function(event, d) {
                    event.stopPropagation();
                    highlightConnected(d, data);
                });

            // Create labels (store globally for filtering)
            globalLabel = labelGroup.selectAll('text')
                .data(data.nodes)
                .enter()
                .append('text')
                .attr('class', 'node-label')
                .attr('dy', '0.31em')
                .text(d => d.id)
                .style('font-size', '11px');

            // Zoom behavior
            const zoom = d3.zoom()
                .on('zoom', event => {
                    g.attr('transform', event.transform);
                });

            svg.call(zoom);

            // Zoom buttons
            document.getElementById('btn-zoom-in').addEventListener('click', () => {
                svg.transition().duration(750).call(zoom.scaleBy, 1.3);
            });

            document.getElementById('btn-zoom-out').addEventListener('click', () => {
                svg.transition().duration(750).call(zoom.scaleBy, 0.7);
            });

            document.getElementById('btn-zoom-fit').addEventListener('click', () => {
                const bounds = g.node().getBBox();
                const fullWidth = width;
                const fullHeight = height;
                const midX = bounds.x + bounds.width / 2;
                const midY = bounds.y + bounds.height / 2;

                if (bounds.width > 0 && bounds.height > 0) {
                    const scale = 0.8 / Math.max(bounds.width / fullWidth, bounds.height / fullHeight);
                    const translate = [fullWidth / 2 - scale * midX, fullHeight / 2 - scale * midY];

                    svg.transition()
                        .duration(750)
                        .call(zoom.transform, d3.zoomIdentity.translate(translate[0], translate[1]).scale(scale));
                }
            });

            // Update simulation
            simulation.on('tick', () => {
                globalLink
                    .attr('x1', d => d.source.x)
                    .attr('y1', d => d.source.y)
                    .attr('x2', d => d.target.x)
                    .attr('y2', d => d.target.y);

                globalNode
                    .attr('cx', d => d.x)
                    .attr('cy', d => d.y);

                globalLabel
                    .attr('x', d => d.x)
                    .attr('y', d => d.y);
            });

            // Reset view button
            document.getElementById('btn-reset').addEventListener('click', () => {
                globalNode.classed('selected', false);
                globalLink.classed('faded', false);
                globalNode.classed('faded', false);
            });

            // Filter buttons
            document.getElementById('btn-filter-all').addEventListener('click', () => {
                updateFilterButtons('all');
                filterByType('all');
            });

            document.getElementById('btn-filter-req').addEventListener('click', () => {
                updateFilterButtons('requirement');
                filterByType('requirement');
            });

            document.getElementById('btn-filter-arch').addEventListener('click', () => {
                updateFilterButtons('arch');
                filterByType('arch');
            });

            document.getElementById('btn-filter-tests').addEventListener('click', () => {
                updateFilterButtons('tests');
                filterByType('test-spec', 'test-code', 'test-result');
            });

            function filterByType(...types) {
                const typeMap = {};
                types.forEach(t => typeMap[t] = true);

                globalNode.style('display', d => typeMap[d.type] ? 'block' : 'none');
                globalLabel.style('display', d => typeMap[d.type] ? 'block' : 'none');
                globalLink.style('display', d => typeMap[d.source.type] && typeMap[d.target.type] ? 'line' : 'none');
            }

            function updateFilterButtons(active) {
                document.getElementById('btn-filter-all').classList.remove('active');
                document.getElementById('btn-filter-req').classList.remove('active');
                document.getElementById('btn-filter-arch').classList.remove('active');
                document.getElementById('btn-filter-tests').classList.remove('active');

                if (active === 'all') document.getElementById('btn-filter-all').classList.add('active');
                else if (active === 'requirement') document.getElementById('btn-filter-req').classList.add('active');
                else if (active === 'arch') document.getElementById('btn-filter-arch').classList.add('active');
                else if (active === 'tests') document.getElementById('btn-filter-tests').classList.add('active');
            }

            function showTooltip(event, d) {
                const tooltip = document.getElementById('tooltip');
                let metaHTML = '';
                Object.entries(d.metadata || {}).forEach(([k, v]) => {
                    metaHTML += '<strong>' + k + ':</strong> ' + v + '<br>';
                });

                tooltip.innerHTML =
                    '<div class="tooltip-title">' + d.label + '</div>' +
                    '<div class="tooltip-content">' +
                    '<strong>Type:</strong> ' + d.type + '<br>' +
                    '<strong>ID:</strong> ' + d.id +
                    '</div>' +
                    '<div class="tooltip-meta">' + metaHTML + '</div>';

                tooltip.classList.add('visible');
                tooltip.style.left = (event.pageX + 10) + 'px';
                tooltip.style.top = (event.pageY + 10) + 'px';
            }

            function hideTooltip() {
                document.getElementById('tooltip').classList.remove('visible');
            }

            function highlightConnected(node, data) {
                const connectedIds = new Set([node.id]);

                // Find all connected nodes
                data.edges.forEach(edge => {
                    if (edge.source.id === node.id) connectedIds.add(edge.target.id);
                    if (edge.target.id === node.id) connectedIds.add(edge.source.id);
                });

                nodeGroup.selectAll('.node')
                    .classed('selected', d => d.id === node.id)
                    .classed('faded', d => !connectedIds.has(d.id));

                linkGroup.selectAll('.link')
                    .classed('faded', d => !(connectedIds.has(d.source.id) && connectedIds.has(d.target.id)));
            }

            function drag(simulation) {
                function dragstarted(event) {
                    if (!event.active) simulation.alphaTarget(0.3).restart();
                    event.subject.fx = event.subject.x;
                    event.subject.fy = event.subject.y;
                }

                function dragged(event) {
                    event.subject.fx = event.x;
                    event.subject.fy = event.y;
                }

                function dragended(event) {
                    if (!event.active) simulation.alphaTarget(0);
                    event.subject.fx = null;
                    event.subject.fy = null;
                }

                return d3.drag()
                    .on('start', dragstarted)
                    .on('drag', dragged)
                    .on('end', dragended);
            }
        }

        function updateStats(data) {
            if (!data || !data.nodes) {
                console.error('No data or nodes in updateStats:', data);
                return;
            }

            const stats = {
                requirement: 0,
                arch: 0,
                'test-spec': 0,
                'test-code': 0,
                'test-result': 0
            };

            data.nodes.forEach(node => {
                if (node.type in stats) {
                    stats[node.type]++;
                }
            });

            const reqEl = document.getElementById('stat-reqs');
            const archEl = document.getElementById('stat-arch');
            const specsEl = document.getElementById('stat-specs');
            const resultsEl = document.getElementById('stat-results');
            const totalEl = document.getElementById('stat-total');
            const linksEl = document.getElementById('stat-links');

            if (reqEl) reqEl.textContent = stats.requirement;
            if (archEl) archEl.textContent = stats.arch;
            if (specsEl) specsEl.textContent = stats['test-spec'];
            if (resultsEl) resultsEl.textContent = stats['test-result'];
            if (totalEl) totalEl.textContent = data.nodes.length;
            if (linksEl) linksEl.textContent = data.edges.length;
        }

        // Tab switching
        function switchTab(tabName) {
            document.querySelectorAll('.tab-content').forEach(el => {
                el.classList.remove('active');
            });
            document.querySelectorAll('.tab-button').forEach(el => {
                el.classList.remove('active');
            });

            if (tabName === 'graph') {
                document.getElementById('graph').classList.add('active');
                document.querySelector('button[onclick="switchTab(\'graph\')"]').classList.add('active');
            } else if (tabName === 'matrix') {
                document.getElementById('matrix').classList.add('active');
                document.querySelector('button[onclick="switchTab(\'matrix\')"]').classList.add('active');
                renderMatrix();
            } else if (tabName === 'aspice') {
                document.getElementById('aspice').classList.add('active');
                document.querySelector('button[onclick="switchTab(\'aspice\')"]').classList.add('active');
                renderASPICE();
            }
        }

        function coverageClass(pct) {
            if (pct >= 80) return 'good';
            if (pct >= 50) return 'warning';
            return 'danger';
        }

        function renderASPICE() {
            if (!aspiceData) return;

            var overall = aspiceData.overall || 0;
            var overallPct = Math.round(overall);
            var cls = coverageClass(overall);

            var pctEl = document.getElementById('aspice-overall-pct');
            var fillEl = document.getElementById('aspice-overall-fill');
            if (pctEl) {
                pctEl.textContent = overallPct + '%';
                pctEl.className = 'overall-pct ' + cls;
            }
            if (fillEl) {
                fillEl.style.width = overallPct + '%';
                fillEl.className = 'overall-bar-fill ' + cls;
            }

            var grid = document.getElementById('aspice-grid');
            if (!grid || !aspiceData.processes) return;

            var html = '';
            aspiceData.processes.forEach(function(proc) {
                var covPct = Math.round(proc.coverage || 0);
                var pcls = coverageClass(proc.coverage || 0);

                var bpHTML = '';
                if (proc.bps) {
                    proc.bps.forEach(function(bp) {
                        var bpPct = Math.round(bp.coverage || 0);
                        var bcls = coverageClass(bp.coverage || 0);
                        var gapsHTML = '';
                        if (bp.gaps && bp.gaps.length > 0) {
                            gapsHTML = '<ul class="bp-gaps">';
                            bp.gaps.forEach(function(g) {
                                gapsHTML += '<li>' + g + '</li>';
                            });
                            gapsHTML += '</ul>';
                        }
                        bpHTML += '<div class="bp-item">' +
                            '<div class="bp-row">' +
                            '<span class="bp-id">' + bp.id + '</span>' +
                            '<span class="bp-title">' + bp.title + '</span>' +
                            '<span class="bp-badge ' + bcls + '">' + bpPct + '%</span>' +
                            '</div>' +
                            gapsHTML +
                            '</div>';
                    });
                }

                html += '<div class="process-card">' +
                    '<div class="process-card-header ' + pcls + '">' +
                    '<span class="process-id-badge">' + proc.id + '</span>' +
                    '<span class="process-card-name">' + proc.name + '</span>' +
                    '<span class="process-card-pct ' + pcls + '">' + covPct + '%</span>' +
                    '</div>' +
                    '<div class="process-progress">' +
                    '<div class="process-progress-fill ' + pcls + '" style="width:' + covPct + '%"></div>' +
                    '</div>' +
                    '<div class="bp-list">' + bpHTML + '</div>' +
                    '</div>';
            });

            grid.innerHTML = html || '<p style="padding:20px;color:#999;">No ASPICE process data available.</p>';
        }

        let currentMatrixData = matrixData;

        function renderMatrix() {
            const container = document.getElementById('matrix-table-container');
            const priority = document.getElementById('filter-priority').value;
            const status = document.getElementById('filter-status').value;
            const search = document.getElementById('search-req').value.toLowerCase();

            let filtered = matrixData.rows.filter(row => {
                const matchesPriority = !priority || row.Priority === priority;
                const matchesStatus = !status || row.Status === status;
                const matchesSearch = !search || row.RequirementID.toLowerCase().includes(search) || row.Title.toLowerCase().includes(search);
                return matchesPriority && matchesStatus && matchesSearch;
            });

            currentMatrixData = { ...matrixData, rows: filtered };

            let html = '<table class="matrix-table"><thead><tr><th>Requirement</th><th>Priority</th><th>Status</th>';
            matrixData.columns.forEach(col => {
                html += '<th title="' + col.Title + '">' + col.ID + '</th>';
            });
            html += '</tr></thead><tbody>';

            filtered.forEach(row => {
                html += '<tr>';
                html += '<td class="req-id" title="' + row.Title + '">' + row.RequirementID + '</td>';
                html += '<td>' + row.Priority + '</td>';
                html += '<td>' + row.Status + '</td>';

                matrixData.columns.forEach(col => {
                    const cell = row.Cells[col.ID];
                    const status = cell ? cell.Status : 'missing';
                    const symbol = status === 'covered' ? '✓' : (status === 'stale' ? '⚠' : '✗');
                    html += '<td class="matrix-cell ' + status + '" title="' + (cell ? cell.Evidence : 'Not covered') + '">' + symbol + '</td>';
                });

                html += '</tr>';
            });

            html += '</tbody></table>';
            container.innerHTML = html;

            // Update statistics
            const stats = matrixData.statistics;
            document.getElementById('matrix-total').textContent = stats.TotalRequirements;
            document.getElementById('matrix-covered').textContent = stats.CoveredRequirements;
            document.getElementById('matrix-missing').textContent = stats.MissingRequirements;
            document.getElementById('matrix-coverage').textContent = Math.round(stats.CoveragePercentage) + '%';
        }

        function exportMatrixCSV() {
            let csv = 'Requirement,Priority,Status';
            matrixData.columns.forEach(col => {
                csv += ',' + col.ID;
            });
            csv += '\n';

            currentMatrixData.rows.forEach(row => {
                csv += row.RequirementID + ',' + row.Priority + ',' + row.Status;
                matrixData.columns.forEach(col => {
                    const cell = row.Cells[col.ID];
                    const symbol = cell && cell.Status === 'covered' ? '✓' : (cell && cell.Status === 'stale' ? '⚠' : '✗');
                    csv += ',' + symbol;
                });
                csv += '\n';
            });

            const blob = new Blob([csv], { type: 'text/csv' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'traceability-matrix.csv';
            a.click();
            window.URL.revokeObjectURL(url);
        }

        // Wire up filter controls
        document.getElementById('filter-priority').addEventListener('change', renderMatrix);
        document.getElementById('filter-status').addEventListener('change', renderMatrix);
        document.getElementById('search-req').addEventListener('input', renderMatrix);

        // Initialize graph visualization after all functions are defined
        document.addEventListener('DOMContentLoaded', function() {
            if (graphData && graphData.nodes) {
                initializeGraph(graphData);
            }
        });

        // Also initialize immediately if DOM is already loaded
        if (document.readyState !== 'loading') {
            if (graphData && graphData.nodes) {
                initializeGraph(graphData);
            }
        }
    </script>
</body>
</html>
`
