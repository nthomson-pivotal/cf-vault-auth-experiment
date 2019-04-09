package org.paasify.cfvault.example.integration.cloud;

import java.net.URI;

import org.paasify.cfvault.example.integration.CloudFoundryAuthentication;
import org.paasify.cfvault.example.integration.CloudFoundryAuthenticationOptions;
import org.springframework.cloud.vault.config.VaultProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.StringUtils;
import org.springframework.vault.authentication.ClientAuthentication;
import org.springframework.vault.client.VaultClients;
import org.springframework.vault.client.VaultEndpoint;
import org.springframework.vault.config.AbstractVaultConfiguration.ClientFactoryWrapper;

@Configuration
public class MyVaultBootstrapConfiguration {
	@Bean
	public ClientAuthentication clientAuthentication(VaultProperties properties, ClientFactoryWrapper clientFactoryWrapper) {
        CloudFoundryAuthenticationOptions options = CloudFoundryAuthenticationOptions.builder().build();
        return new CloudFoundryAuthentication(options, VaultClients.createRestTemplate(createVaultEndpoint(properties), clientFactoryWrapper.getClientHttpRequestFactory()));
	}
	
	static VaultEndpoint createVaultEndpoint(VaultProperties vaultProperties) {
		if (StringUtils.hasText(vaultProperties.getUri())) {
			return VaultEndpoint.from(URI.create(vaultProperties.getUri()));
		}

		VaultEndpoint vaultEndpoint = new VaultEndpoint();
		vaultEndpoint.setHost(vaultProperties.getHost());
		vaultEndpoint.setPort(vaultProperties.getPort());
		vaultEndpoint.setScheme(vaultProperties.getScheme());

		return vaultEndpoint;
	}
}