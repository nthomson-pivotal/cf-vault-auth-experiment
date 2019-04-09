package org.paasify.cfvault.example.integration;

import java.util.function.Supplier;

@FunctionalInterface
public interface CloudFoundryStringSupplier extends Supplier<String> {

	/**
	 * Get a JWT for Kubernetes authentication.
	 *
	 * @return the Kubernetes Service Account JWT.
	 */
	@Override
	String get();
}