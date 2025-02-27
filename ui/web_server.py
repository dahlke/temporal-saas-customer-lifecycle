import uuid
from dataclasses import dataclass, field
from typing import Dict
from flask import Flask, render_template, request, jsonify, redirect, url_for
from shared.base import get_temporal_client, LifecycleInput, AcceptClaimCodeInput, \
	TEMPORAL_ADDRESS, TEMPORAL_NAMESPACE, TEMPORAL_TASK_QUEUE, ENCRYPT_PAYLOADS

from temporalio.common import TypedSearchAttributes, SearchAttributeKey, \
	SearchAttributePair

app = Flask(__name__)

# Define search attribute keys for workflow search
lifecycle_status_key = SearchAttributeKey.for_text("LifecycleStatus")
TEMPORAL_UI_URL = TEMPORAL_ADDRESS.replace("7233", "8233") if "localhost" in TEMPORAL_ADDRESS \
	else "https://cloud.temporal.io"

# Define the available scenarios
SCENARIOS = {
	"HAPPY_PATH": {
		"title": "Happy Path",
		"description": "This will run the usual scenario with no failures, including signals to re-send claim codes and updates to accept and validate claim codes. The subscription loop will be run within the main workflow."
	},
	"RECOVERABLE_FAILURE": {
		"title": "Recoverable Failure (Bug in Code)",
		"description": "This will cause the workflow to fail in a recoverable way. This means you can comment out the line that causes the error and continue the workflow."
	},
	"NON_RECOVERABLE_FAILURE": {
		"title": "Non Recoverable Failure (Non-Retryable)",
		"description": "This will cause the workflow to fail entirely, by throwing a non-retryable error."
	},
	"FLAKEY_API": {
		"title": "Flakey API (API failures that recover on 5th attempt)",
		"description": "This workflow will fail 5 times at the Create Admin Users activity, simulating a downstream API outage."
	},
	"CHILD_WORKFLOW": {
		"title": "Child Workflow",
		"description": "This will run the Happy Path then spawn a child workflow to manage the subscription lifecycle."
	},
	"NEXUS": {
		"title": "Nexus Workflow",
		"description": "This will run the Happy Path then spawn a workflow in another namespace to manage the subscription lifecycle."
	},
}

# Global variable to store the Temporal client
temporal_client = None

async def _get_singleton_temporal_client():
	global temporal_client
	if temporal_client is None:
		temporal_client = await get_temporal_client()
	return temporal_client

# Define the main route
@app.route("/", methods=["GET", "POST"])
async def main():
	return render_template(
		"index.html",
		scenarios=SCENARIOS,
		temporal_host_url=TEMPORAL_ADDRESS,
		temporal_ui_url=TEMPORAL_UI_URL,
		temporal_namespace=TEMPORAL_NAMESPACE,
		payloads_encrypted=ENCRYPT_PAYLOADS
	)

# Define the run_workflow route
@app.route("/run_workflow", methods=["GET", "POST"])
async def run_workflow():
	# Get the selected scenario and workflow ID from the request arguments
	selected_scenario = request.args.get("scenario", "")
	account_name = request.args.get("accountName", "default")
	wf_id = request.args.get("wfID", "")

	# Create Workflow input
	wf_input = LifecycleInput(
		account_name=account_name,
		emails=["sa@temporal.io", "solutions@temporal.io"],
		price=10.0,
		scenario=selected_scenario,
	)

	# Get the Temporal client
	client = await _get_singleton_temporal_client()

	no_existing_workflow = False
	try:
		# Check if the workflow already exists
		wf_handle = client.get_workflow_handle(wf_id)
		await wf_handle.describe()
	except Exception as e:
		no_existing_workflow = True

	if no_existing_workflow:
		# Start the workflow if it doesn't exist
		await client.start_workflow(
			"LifecycleWorkflow", # Defined in the Go worker
			wf_input,
			id=wf_id,
			task_queue=TEMPORAL_TASK_QUEUE,
			search_attributes=TypedSearchAttributes([
				SearchAttributePair(lifecycle_status_key, ""),
			]),
		)

	return render_template(
		"run_workflow.html",
		wf_id=wf_id,
		selected_scenario=selected_scenario,
		temporal_host_url=TEMPORAL_ADDRESS,
		temporal_ui_url=TEMPORAL_UI_URL,
		temporal_namespace=TEMPORAL_NAMESPACE,
		payloads_encrypted=ENCRYPT_PAYLOADS
	)

# Define the get_progress route
@app.route('/get_progress')
async def get_progress():
	wf_id = request.args.get('wfID', "")
	scenario = request.args.get('scenario', "")

	payload = {}

	try:
		client = await _get_singleton_temporal_client()
		wf_handle = client.get_workflow_handle(wf_id)

		workflow_desc = await wf_handle.describe()
		payload = await wf_handle.query("GetState")

		if workflow_desc.status == 3:
			error_message = "Workflow failed: {wf_id}"
			print(f"Error in get_progress route: {error_message}")
			return redirect(url_for('end_workflow', wfID=wf_id, scenario=scenario))

		payload["temporal_ui_url"] = TEMPORAL_UI_URL
		payload["temporal_namespace"] = TEMPORAL_NAMESPACE

		return jsonify(payload)
	except Exception as e:
		print(e)
		return jsonify(payload)

# Define the ended route
@app.route('/end_workflow')
async def end_workflow():
	wf_id = request.args.get("wfID", "")

	client = await _get_singleton_temporal_client()
	wf_handle = client.get_workflow_handle(wf_id)
	status = await wf_handle.result()

	return render_template(
		"end_workflow.html",
		wf_id=wf_id,
		status=status,
		temporal_host_url=TEMPORAL_ADDRESS,
		temporal_ui_url=TEMPORAL_UI_URL,
		temporal_namespace=TEMPORAL_NAMESPACE,
		payloads_encrypted=ENCRYPT_PAYLOADS
	)

# Define the signal route
@app.route('/signal', methods=["POST"])
async def signal():
	wf_id = request.args.get("wfID", "")
	signal_type = request.json.get("signalType", "")
	scenario = request.args.get("scenario", "")

	print("SIGNAL", wf_id, signal_type)

	try:
		client = await _get_singleton_temporal_client()
		wf_handle = client.get_workflow_handle(wf_id)

		if signal_type == "ResendClaimCodesSignal":
			await wf_handle.signal(signal_type)
		elif signal_type == "CancelSubscriptionSignal":
			await wf_handle.signal(signal_type)

			if scenario == "NEXUS":
				await wf_handle.cancel()
		else:
			raise Exception("Signal type not supported")

	except Exception as e:
		print(f"Error sending signal: {str(e)}")
		return jsonify({"error": str(e)}), 500

	return "Signal received successfully", 200

# Define the update route
@app.route('/update', methods=["POST"])
async def update():
	wf_id = request.args.get("wfID", "")
	claim_code = request.json.get("claim_code", "")

	try:
		client = await _get_singleton_temporal_client()
		wf_handle = client.get_workflow_handle(wf_id)

		claim_code_input = AcceptClaimCodeInput(
			claim_code=claim_code,
		)
		result = await wf_handle.execute_update("AcceptClaimCodeUpdate", claim_code_input)

		return jsonify({"result": result}), 200
	except Exception as e:
		print(f"Error sending update: {str(e)}")
		return jsonify({"result": "Error sending update. Make sure your code is not empty."}), 422

# Run the Flask app
if __name__ == "__main__":
	app.run(debug=True, port=3000)
