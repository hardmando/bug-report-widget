// UMD bundle - works as <script src=""> for SaaS customers
(function (window) {
    let apiKey, endpoint, modal, isOpen = false;

    window.BugWidget = {
        init(config) {
            apiKey = config.apiKey;
            endpoint = config.endpoint;

            // Create widget button via Shadow DOM
            const container = document.getElementById('bug-widget-container');
            const shadow = container.attachShadow({ mode: 'open' });

            shadow.innerHTML = `
        <style>
          #widget-btn { 
            position: fixed; bottom: 20px; right: 20px; 
            background: #ff5e5e; color: white; padding: 1rem; 
            border-radius: 8px; cursor: pointer; font-weight: bold;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            border: none; font-size: 14px;
          }
          #widget-btn:hover { background: #e55353; }
        </style>
        <button id="widget-btn">Report Bug</button>
      `;

            // Button click → open modal
            shadow.querySelector('#widget-btn').onclick = () => this.openModal();
            console.log('BugWidget initialized');
        },

        openModal() {
            modal = document.getElementById('bug-modal');
            modal.style.display = 'block';
            isOpen = true;
            document.getElementById('bug-desc').focus();
        },

        closeModal() {
            modal.style.display = 'none';
            isOpen = false;
        },

        async submitBug() {
            const description = document.getElementById('bug-desc').value;
            if (!description.trim()) return;

            // Capture basic context
            const bugData = {
                description,
                url: window.location.href,
                userAgent: navigator.userAgent,
                viewport: { width: window.innerWidth, height: window.innerHeight },
                timestamp: new Date().toISOString(),
                consoleLogs: consoleHistory.slice(-10)  // From monkey patch below
            };

            try {
                // Step 1: Validate API key via auth-service
                const authRes = await fetch('http://localhost:8081/validate-key', {
                    headers: { 'X-API-Key': apiKey }
                });
                if (!authRes.ok) throw new Error('Invalid API key');

                // Step 2: Send to ingestion-service
                const ingestRes = await fetch(`${endpoint}/ingest/bugs`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-API-Key': apiKey
                    },
                    body: JSON.stringify(bugData)
                });

                if (ingestRes.ok) {
                    alert('Bug reported successfully!');
                    this.closeModal();
                } else {
                    alert('Failed to report bug');
                }
            } catch (err) {
                console.error(err);
                alert('Network error');
            }
        }
    };

    // Monkey patch console to capture logs
    const consoleHistory = [];
    const originalLog = console.log;
    const originalError = console.error;
    console.log = (...args) => {
        consoleHistory.push({ type: 'log', args, time: Date.now() });
        originalLog.apply(console, args);
    };
    console.error = (...args) => {
        consoleHistory.push({ type: 'error', args, time: Date.now() });
        originalError.apply(console, args);
    };

})(window);
