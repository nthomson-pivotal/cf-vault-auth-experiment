package org.paasify.cfvault.example.integration;

import org.springframework.vault.VaultException;
import org.springframework.vault.client.VaultResponses;
import org.springframework.web.client.RestClientResponseException;

/**
 * Niall: HAD TO COPY THIS BECAUSE WHY CAN'T I REFERENCE IT
 * 
 * 
 * 
 * Exception thrown if Vault login fails. The root cause is typically attached as cause.
 *
 * @author Mark Paluch
 * @since 2.1
 */
public class VaultLoginException extends VaultException {

	/**
	 * Create a {@code VaultLoginException} with the specified detail message.
	 *
	 * @param msg the detail message.
	 */
	public VaultLoginException(String msg) {
		super(msg);
	}

	/**
	 * Create a {@code VaultLoginException} with the specified detail message and nested
	 * exception.
	 *
	 * @param msg the detail message.
	 * @param cause the nested exception.
	 */
	public VaultLoginException(String msg, Throwable cause) {
		super(msg, cause);
	}

	/**
	 * Create a {@link VaultLoginException} given {@code authMethod} and a
	 * {@link Throwable cause}.
	 *
	 * @param authMethod must not be {@literal null}.
	 * @param cause must not be {@literal null}.
	 * @return the {@link VaultLoginException}.
	 */
	public static VaultLoginException create(String authMethod, Throwable cause) {

		if (cause instanceof RestClientResponseException) {

			String response = ((RestClientResponseException) cause)
					.getResponseBodyAsString();
			return new VaultLoginException(String.format("Cannot login using %s: %s",
					authMethod, VaultResponses.getError(response)), cause);
		}

		return new VaultLoginException(String.format("Cannot login using %s", cause));
	}
}