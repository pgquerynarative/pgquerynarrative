(function () {
  'use strict';

  var resultChart = null;

  function isNumeric(v) {
    if (v === null || v === undefined) return false;
    return typeof v === 'number' && !Number.isNaN(v);
  }

  function getLabelAndValueColumns(data) {
    var cols = data.columns || [];
    var rows = data.rows || [];
    if (cols.length === 0 || rows.length === 0) return null;
    var labelIdx = 0;
    var valueIdx = 1;
    if (cols.length === 1) {
      if (rows[0] && isNumeric(rows[0][0])) {
        var vals = rows.map(function (r) { return r[0]; });
        var lbls = vals.map(function (_, i) { return 'Row ' + (i + 1); });
        return { labels: lbls, values: vals, valueLabel: cols[0] };
      }
      return { labels: rows.map(function (r) { return String(r[0] != null ? r[0] : ''); }), values: [], valueLabel: '' };
    }
    for (var i = 0; i < cols.length; i++) {
      var sample = rows[0] && rows[0][i];
      if (isNumeric(sample)) {
        valueIdx = i;
        labelIdx = i === 0 ? 1 : 0;
        break;
      }
    }
    var labels = rows.map(function (r) { return String(r[labelIdx] != null ? r[labelIdx] : ''); });
    var values = rows.map(function (r) {
      var v = r[valueIdx];
      return isNumeric(v) ? v : (parseFloat(v) || 0);
    });
    return { labels: labels, values: values, valueLabel: cols[valueIdx] || 'Value' };
  }

  function drawChart(canvas, type, data) {
    var parsed = getLabelAndValueColumns(data);
    if (!parsed || (type !== 'table' && parsed.values.length === 0)) return;

    if (resultChart) {
      resultChart.destroy();
      resultChart = null;
    }

    var ctx = canvas.getContext('2d');
    var cfg = {
      type: type === 'bar' ? 'bar' : type === 'line' ? 'line' : type === 'pie' ? 'pie' : 'bar',
      data: {
        labels: parsed.labels,
        datasets: [{
          label: parsed.valueLabel,
          data: parsed.values,
          backgroundColor: type === 'pie'
            ? ['#3498db', '#2ecc71', '#e74c3c', '#f39c12', '#9b59b6', '#1abc9c', '#34495e', '#e67e22']
            : 'rgba(52, 152, 219, 0.6)',
          borderColor: 'rgba(52, 152, 219, 1)',
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: { display: type === 'pie' }
        },
        scales: type !== 'pie' ? {
          y: { beginAtZero: true }
        } : {}
      }
    };

    if (type === 'pie' && parsed.labels.length > 12) {
      cfg.data.labels = parsed.labels.slice(0, 12).concat(['Other']);
      var otherSum = parsed.values.slice(12).reduce(function (a, b) { return a + b; }, 0);
      cfg.data.datasets[0].data = parsed.values.slice(0, 12).concat([otherSum]);
    }

    resultChart = new Chart(ctx, cfg);
  }

  function initChartsInContainer(container) {
    if (!container || !window.Chart) return;
    var area = container.querySelector('.chart-area');
    if (!area) return;
    var raw = area.getAttribute('data-chart-data');
    if (!raw) return;
    var data;
    try {
      data = JSON.parse(raw);
    } catch (e) {
      return;
    }
    var canvas = area.querySelector('#result-chart');
    var select = area.querySelector('#chart-type-select');
    if (!canvas || !select) return;

    function updateChart() {
      var type = select.value;
      if (type) drawChart(canvas, type, data);
      else if (resultChart) {
        resultChart.destroy();
        resultChart = null;
      }
    }

    select.addEventListener('change', updateChart);

    area.querySelectorAll('.chart-type-btn').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var t = btn.getAttribute('data-chart-type');
        if (t && t !== 'table') {
          select.value = t;
          updateChart();
        }
      });
    });
  }

  function onHtmxSettle(ev) {
    var target = ev.detail && ev.detail.target ? ev.detail.target : ev.target;
    if (!target) return;
    var resultsEl = target.id === 'results' ? target : (target.querySelector && target.querySelector('#results'));
    if (resultsEl) initChartsInContainer(resultsEl);
  }

  document.addEventListener('htmx:afterSettle', onHtmxSettle);
})();
