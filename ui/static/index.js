function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0,
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

function updateCodesContainer(data) {
    const tbody = document.getElementById('codesTableBody');

    // Clear existing rows
    tbody.innerHTML = '';

    // Add a new row for each email and claim code
    data.forEach(item => {
        const row = document.createElement('tr');

        const emailCell = document.createElement('td');
        emailCell.textContent = item.Email;
        row.appendChild(emailCell);

        const codeCell = document.createElement('td');
        codeCell.textContent = item.Code;
        row.appendChild(codeCell);

        tbody.appendChild(row);
    });
}

var GLOBAL_childWorkflowID = "";

function updateProgress() {
	var urlParams = new URLSearchParams(window.location.search);
	var scenario = urlParams.get("scenario");
	var wfID = urlParams.get("wfID");

	fetch("/get_progress?scenario=" + encodeURIComponent(scenario) + "&wfID=" + encodeURIComponent(wfID))
		.then(response => {
			if (response.ok) {
				return response.json();
			} else {
				throw new Error(`Failed to fetch progress. Status: ${response.status}`);
			}
		})
		.then(data => {
			document.getElementById("errorMessage").innerText = "";
			document.getElementById("progressBar").style.width = data.progress + "%";

			if (scenario === "CHILD_WORKFLOW" && data.child_workflow_id !== "") {
				GLOBAL_childWorkflowID = data.child_workflow_id;

				// Check if the row already exists
				if (!document.getElementById("childWorkflowRow")) {
					const table = document.getElementById("workflowInfoTable");
					const newRow = table.insertRow();
					newRow.id = "childWorkflowRow"; // Set the ID for the new row

					const cell1 = document.createElement('th');
					cell1.innerHTML = "Child Workflow ID";
					newRow.appendChild(cell1);

					const cell2 = newRow.insertCell(1);
					cell2.innerHTML = `<a href="${data.temporal_ui_url}/namespaces/${data.temporal_namespace}/workflows/${GLOBAL_childWorkflowID}" target="_blank">${GLOBAL_childWorkflowID}</a>`;
				}
			}

			var currentStatusEl = document.getElementById("currentStatus");
			if (currentStatusEl != null) {
				// TODO: show the child status
				// TODO: show the nexus status
				currentStatusEl.innerText = data.status;
			}

			if (data.status === "WAITING_FOR_CLAIM_CODES") {
				document.getElementById("resendSignalContainer").style.display = "block";
				document.getElementById("updateContainer").style.display = "block";
				document.getElementById("codesContainer").style.display = "block";

				updateCodesContainer(data.claim_codes);
			} else if (data.status.includes("RENEWED") || (scenario == "CHILD_WORKFLOW" && data.status === "ONBOARDED")) {
				document.getElementById("cancelSignalContainer").style.display = "block";
			} else if (data.status.includes("CODE_NOT_CLAIMED")) {
				window.location.href =
					"/end_workflow?wfID=" + encodeURIComponent(wfID) +
					"&scenario=" + encodeURIComponent(scenario);
			}

			setTimeout(updateProgress, 1000);
		})
		.catch(error => {
			console.error("Error fetching progress:", error.message);
			document.getElementById("errorMessage").innerText = error.message;
			document.getElementById("progressBar").style.backgroundColor = "red";
		});
}

function runWorkflow() {
	var scenario = document.getElementById("scenario").value;
	var accountName = document.getElementById("accountName").value;
	var wfID = `customer-lifecycle-${accountName}-${generateUUID()}`;

	// Redirect to run_workflow page with the selected scenario as a query parameter
	window.location.href =
		"/run_workflow?scenario=" + encodeURIComponent(scenario) +
		"&wfID=" + encodeURIComponent(wfID) +
		"&accountName=" + encodeURIComponent(accountName)
}

function signal(signalType, payload) {
	// Get the wfID from the URL query parameters
	var urlParams = new URLSearchParams(window.location.search);
	var wfID = urlParams.get("wfID");
	if (GLOBAL_childWorkflowID !== "") {
		wfID = GLOBAL_childWorkflowID;
	}
	var scenario = urlParams.get("scenario");

	// Perform AJAX request to the server for signaling
	fetch("/signal?wfID=" + encodeURIComponent(wfID), {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			signalType: signalType,
			payload: payload
		})
	})
	.then(response => {
		if (response.ok) {
			console.log("Signal sent successfully");
			if (signalType === "CancelSubscriptionSignal") {
				window.location.href =
					"/end_workflow?wfID=" + encodeURIComponent(wfID) +
					"&scenario=" + encodeURIComponent(scenario);
			}
		} else {
				console.error("Failed to send signal");

				// Get the signalResult element
				var signalResultEl = document.getElementById("errorMessage");

				// Update the display with the result
				signalResultEl.innerText = "Signal sent failed";
			}
		})
		.catch(error => {
			console.error("Error during signal:", error.message);
		});
}

function update() {
	// Get the wfID from the URL query parameters
	var urlParams = new URLSearchParams(window.location.search);
	var wfID = urlParams.get("wfID");
	var code = document.getElementById("claimCode").value;
	var updateResultEl = document.getElementById("updateResult");
	updateResultEl.style.display = "none";


	// Perform AJAX request to the server for updating
	fetch("/update?wfID=" + encodeURIComponent(wfID), {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			"claim_code": code
		})
	})
	.then(response => {
		if (response.status !== 200) {
			console.error("Failed to send update");
			updateResultEl.style.display = "block";
			updateResultEl.innerText = "Update sent failed, enter a correct claim code and try again."
		} else {
			document.getElementById("updateContainer").style.display = "none";
			document.getElementById("codesContainer").style.display = "none";
			document.getElementById("resendSignalContainer").style.display = "none";
		}
	})
	.catch(error => {
		console.error("Error during update:", error.message);
	});
}

function reloadMainPage() {
	// Redirect to the main page
	window.location.href = "/";
}
