package org.paasify.cfvault.example.integration;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.vault.authentication.ClientAuthentication;
import org.springframework.vault.support.VaultResponse;
import org.springframework.vault.support.VaultToken;
import org.springframework.web.client.RestClientException;
import org.springframework.web.client.RestOperations;

public class CloudFoundryAuthentication implements ClientAuthentication {

	private static final Log logger = LogFactory.getLog(CloudFoundryAuthentication.class);

	private final CloudFoundryAuthenticationOptions options;

	private final RestOperations restOperations;

	public CloudFoundryAuthentication(CloudFoundryAuthenticationOptions options, RestOperations restOperations) {
		this.options = options;
		this.restOperations = restOperations;
	}

	@Override
	public VaultToken login() {
		Map<String, String> login = new HashMap<String, String>();

		try {
			login.put("certificate", encodeFileToBase64Binary(options.getCertificateSupplier().get()));
			login.put("key", encodeFileToBase64Binary(options.getKeySupplier().get()));
		} catch (IOException ioe) {
			throw new RuntimeException("Failed to load CloudFoundry instance certificate/key", ioe);
		}

		try {
			VaultResponse response = restOperations.postForObject("auth/{mount}/login", login, VaultResponse.class,
					options.getPath());

			logger.debug("Login successful using CloudFoundry authentication");

			return LoginTokenUtil.from(response.getAuth());
		} catch (RestClientException e) {
			throw VaultLoginException.create("CloudFoundry", e);
		}
	}

	private String encodeFileToBase64Binary(String file) throws IOException {
		return java.util.Base64.getEncoder().encodeToString(file.getBytes());
	}
}