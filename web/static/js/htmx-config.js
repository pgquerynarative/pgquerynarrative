// HTMX Configuration
// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
	if (document.body) {
		document.body.addEventListener('htmx:beforeRequest', function(evt) {
			// Show loading indicator
			const loading = document.getElementById('loading');
			if (loading) {
				loading.style.display = 'block';
			}
		});

		document.body.addEventListener('htmx:afterRequest', function(evt) {
			// Hide loading indicator
			const loading = document.getElementById('loading');
			if (loading) {
				loading.style.display = 'none';
			}
		});

		document.body.addEventListener('htmx:responseError', function(evt) {
			// Display error message
			const target = evt.detail.target;
			if (target) {
				let errorMsg = 'Request failed';
				try {
					const response = JSON.parse(evt.detail.xhr.responseText);
					errorMsg = response.message || response.name || errorMsg;
				} catch (e) {
					errorMsg = evt.detail.xhr.responseText || errorMsg;
				}
				const errorHtml = '<div class="error-message"><strong>Error:</strong> ' + errorMsg + '</div>';
				target.innerHTML = errorHtml;
			}
		});
	}
});
