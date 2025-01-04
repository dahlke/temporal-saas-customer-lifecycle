import os
import dataclasses
from dataclasses import dataclass, field
from typing import Dict, Optional, List
from temporalio.client import Client
from temporalio.service import  TLSConfig
from temporalio import converter
from temporalio.runtime import Runtime
from shared.codec import CompressionCodec, EncryptionCodec

# Get the Temporal host URL from environment variable, default to "localhost:7233" if not set
TEMPORAL_ADDRESS = os.environ.get("TEMPORAL_ADDRESS", "localhost:7233")

# Get the mTLS TLS certificate and key file paths from environment variables
TEMPORAL_TLS_CERT = os.environ.get("TEMPORAL_TLS_CERT", "")
TEMPORAL_TLS_KEY = os.environ.get("TEMPORAL_TLS_KEY", "")

# Get the Temporal namespace from environment variable, default to "default" if not set
TEMPORAL_NAMESPACE = os.environ.get("TEMPORAL_NAMESPACE", "default")

# Get the Temporal task queue from environment variable, default to "provision-infra" if not set
TEMPORAL_TASK_QUEUE = os.environ.get("TEMPORAL_TASK_QUEUE", "provision-infra")

# Get the Temporal Cloud API key from environment variable
TEMPORAL_API_KEY = os.environ.get("TEMPORAL_API_KEY", "")

# Determine if payloads should be encrypted based on the value of the "ENCRYPT_PAYLOADS" environment variable
ENCRYPT_PAYLOADS = os.getenv("ENCRYPT_PAYLOADS", 'false').lower() in ('true', '1', 't')


async def get_temporal_client(runtime: Optional[Runtime] = None) -> Client:
	tls_config = False
	data_converter = None

	# If mTLS TLS certificate and key are provided, create a TLSConfig object
	if TEMPORAL_TLS_CERT != "" and TEMPORAL_TLS_KEY != "":
		with open(TEMPORAL_TLS_CERT, "rb") as f:
			client_cert = f.read()

		with open(TEMPORAL_TLS_KEY, "rb") as f:
			client_key = f.read()

		tls_config = TLSConfig(
			client_cert=client_cert,
			client_private_key=client_key,
		)


	if ENCRYPT_PAYLOADS:
		print("Using encryption codec")
		data_converter = dataclasses.replace(
			converter.default(),
			payload_codec=EncryptionCodec(),
			failure_converter_class=converter.DefaultFailureConverterWithEncodedAttributes
		)

	# NOTE: We are using a flag here, since the entire application needs the TEMPORAL_API_KEY
	# to be set, for now.
	if TEMPORAL_API_KEY != "":
		print("Using Cloud API key")
		# Create a Temporal client using the Cloud API key
		client = await Client.connect(
			TEMPORAL_ADDRESS,
			namespace=TEMPORAL_NAMESPACE,
			rpc_metadata={"temporal-namespace": TEMPORAL_NAMESPACE},
			api_key=TEMPORAL_API_KEY,
			data_converter=data_converter,
			tls=True,
		)
	else:
		print("Using MTLS")
		# Create a Temporal client using MTLS
		client: Client = await Client.connect(
			TEMPORAL_ADDRESS,
			namespace=TEMPORAL_NAMESPACE,
			tls=tls_config,
			data_converter=data_converter,
			runtime=runtime
		)

	return client

@dataclass
class AcceptClaimCodeInput:
	claim_code: str

@dataclass
class ClaimCodeStatus:
	email: str
	code: str
	is_claimed: bool

@dataclass
class LifecycleWorkflowState:
	account_name: str
	price: float
	emails: List[str]
	claim_codes: List[ClaimCodeStatus]

@dataclass
class LifecycleWorkflowInput:
	account_name: str
	emails: List[str]
	price: float
	scenario: str
