function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0,
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

function runWorkflow() {
	var selectedScenario = document.getElementById("scenario").value;
	var tfRunID = `customer-lifecycle-${generateUUID()}`;

	// Redirect to run_workflow page with the selected scenario as a query parameter
	window.location.href =
		"/run_workflow?scenario=" + encodeURIComponent(selectedScenario) +
		"&wfID=" + encodeURIComponent(tfRunID)
}

// TODO: clean up all of this code.

function updateProgress() {
	var urlParams = new URLSearchParams(window.location.search);
	var scenario = urlParams.get("scenario");
	var tfRunID = urlParams.get("wfID");

	fetch("/get_progress?wfID=" + encodeURIComponent(tfRunID))
		.then(response => {
			if (response.ok) {
				return response.json();
			} else {
				// If response status is not okay, throw an error
				throw new Error(`Failed to fetch progress. Status: ${response.status}`);
			}
		})
		.then(data => {
			// Update the progress bar
			document.getElementById("errorMessage").innerText = "";
			document.getElementById("progressBar").style.width = data.progress + "%";

			var currentStatusEl = document.getElementById("currentStatus");
			if (currentStatusEl != null) {
				currentStatusEl.innerText = data.status;
			}

			console.log(data);
			if (data.status === "WAITING_FOR_CLAIM_CODES") {
				document.getElementById("signalContainer").style.display = "block";
				document.getElementById("updateContainer").style.display = "block";
			}

			if (data.progress_percent === 100) {
				// Redirect to order confirmation with the tfRunID
				window.location.href =
					"/end_workflow?wfID=" + encodeURIComponent(tfRunID) +
					"&scenario=" + encodeURIComponent(scenario);
			} else {
				// Continue updating progress every second
				setTimeout(updateProgress, 1000);
			}
		})
		.catch(error => {
			// Log the detailed error message to the console
			console.error("Error fetching progress:", error.message);

			// Display the error message in the web browser
			document.getElementById("errorMessage").innerText = error.message;

			// Handle the error by showing a red status bar
			document.getElementById("progressBar").style.backgroundColor = "red";
		});
}

function signal(signalType, payload) {
	// Get the tfRunID from the URL query parameters
	var urlParams = new URLSearchParams(window.location.search);
	var tfRunID = urlParams.get("wfID");

	// Perform AJAX request to the server for signaling
	fetch("/signal?wfID=" + encodeURIComponent(tfRunID), {
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
				if (signalType == "request_continue_as_new") {
					document.getElementById("signalContainer").style.display = "none";
					document.getElementById("updateContainer").style.display = "none";
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

function update(updateType, decision) {
	// Get the tfRunID from the URL query parameters
	var urlParams = new URLSearchParams(window.location.search);
	var tfRunID = urlParams.get("wfID");
	var reason = document.getElementById("reason").value;
	var updateResultEl = document.getElementById("updateResult");
	updateResultEl.style.display = "none";


	// Perform AJAX request to the server for updating
	fetch("/update?wfID=" + encodeURIComponent(tfRunID), {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			decision: decision,
			reason: reason
		})
	})
	.then(response => {
		if (response.status !== 200) {
			console.error("Failed to send update");
			updateResultEl.style.display = "block";
			updateResultEl.innerText = "Update sent failed, enter a reason and try again."
		}
	})
	.catch(error => {
		console.error("Error during update:", error.message);
	});
}

function handleScenarioChange(event) {
	var scenario = event.target.value;
	console.log(scenario);
	// TODO
}

function reloadMainPage() {
	// Redirect to the main page
	window.location.href = "/";
}

function stripAnsi(text) {
  // This regex matches ANSI escape sequences
  const ansiRegex = /\x1b\[[0-9;]*m/g;

  // Replace the ANSI escape codes with an empty string
  return text.replace(ansiRegex, "");
}
