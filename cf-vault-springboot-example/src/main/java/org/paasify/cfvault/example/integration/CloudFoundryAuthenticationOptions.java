package org.paasify.cfvault.example.integration;

import java.util.function.Supplier;

public class CloudFoundryAuthenticationOptions {

	public static final String DEFAULT_CF_AUTHENTICATION_PATH = "cf";

	/**
	 * Path of the CF authentication backend mount.
	 */
	private final String path;

    private final Supplier<String> certificateSupplier;
    
    private final Supplier<String> keySupplier;

	private CloudFoundryAuthenticationOptions(String path,
			Supplier<String> certificateSupplier, Supplier<String> keySupplier) {

		this.path = path;
        this.certificateSupplier = certificateSupplier;
        this.keySupplier = keySupplier;
	}

	public static CloudFoundryAuthenticationOptionsBuilder builder() {
		return new CloudFoundryAuthenticationOptionsBuilder();
	}

	public String getPath() {
		return path;
	}

	public Supplier<String> getCertificateSupplier() {
		return certificateSupplier;
	}

    public Supplier<String> getKeySupplier() {
		return keySupplier;
	}

	public static class CloudFoundryAuthenticationOptionsBuilder {

		private String path = DEFAULT_CF_AUTHENTICATION_PATH;

        private Supplier<String> certificateSupplier;

        private Supplier<String> keySupplier;
        
		public CloudFoundryAuthenticationOptionsBuilder path(String path) {
			this.path = path;
			return this;
		}
        
		public CloudFoundryAuthenticationOptionsBuilder certificateSupplier(
				Supplier<String> certificateSupplier) {
			this.certificateSupplier = certificateSupplier;
			return this;
        }
        
        public CloudFoundryAuthenticationOptionsBuilder keySupplier(
				Supplier<String> keySupplier) {
			this.keySupplier = keySupplier;
			return this;
		}

		public CloudFoundryAuthenticationOptions build() {
			return new CloudFoundryAuthenticationOptions(path,
					certificateSupplier == null ? new CloudFoundryInstanceCertificateFile()
                            : certificateSupplier,
                    keySupplier == null ? new CloudFoundryInstanceKeyFile()
							: keySupplier);
		}
	}
}