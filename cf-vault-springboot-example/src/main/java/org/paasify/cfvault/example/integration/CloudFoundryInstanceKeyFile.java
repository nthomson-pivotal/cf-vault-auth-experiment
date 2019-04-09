package org.paasify.cfvault.example.integration;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;

import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.util.Assert;
import org.springframework.util.StreamUtils;
import org.springframework.vault.VaultException;

public class CloudFoundryInstanceKeyFile implements CloudFoundryStringSupplier {

    /**
     * Default path to the instance certificate file.
     */
    public static final String DEFAULT_CF_INSTANCE_CERTIFICATE_FILE = "/etc/cf-instance-credentials/instance.key";

    private byte[] certificate;

    public CloudFoundryInstanceKeyFile() {
        this(DEFAULT_CF_INSTANCE_CERTIFICATE_FILE);
    }

    public CloudFoundryInstanceKeyFile(String path) {
        this(new FileSystemResource(path));
    }

    public CloudFoundryInstanceKeyFile(File file) {
        this(new FileSystemResource(file));
    }

    public CloudFoundryInstanceKeyFile(Resource resource) {

		Assert.isTrue(resource.exists(),
				() -> String.format("Resource %s does not exist", resource));

		try {
			this.certificate = readCertificate(resource);
		}
		catch (IOException e) {
			throw new VaultException(String.format(
					"CF instance key retrieval from %s failed", resource), e);
		}
	}

	@Override
	public String get() {
		return new String(certificate, StandardCharsets.US_ASCII);
	}

	protected static byte[] readCertificate(Resource resource) throws IOException {

		Assert.notNull(resource, "Resource must not be null");

		try (InputStream is = resource.getInputStream()) {
			return StreamUtils.copyToByteArray(is);
		}
	}
}