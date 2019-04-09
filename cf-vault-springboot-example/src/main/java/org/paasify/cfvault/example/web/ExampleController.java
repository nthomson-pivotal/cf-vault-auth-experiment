package org.paasify.cfvault.example.web;

import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.vault.core.VaultTemplate;
import org.springframework.vault.support.Versioned;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("/")
public class ExampleController {

	@Autowired
	private VaultTemplate vaultTemplate;

	private final static String ENCRYPTED_KEY_NAME = "github.oauth2.key";

	@GetMapping()
	public String basic() {
		Versioned<Map<String, Object>> response = vaultTemplate.opsForVersionedKeyValue("kv").get("github");

		if (response.getData().containsKey(ENCRYPTED_KEY_NAME)) {
			String value = (String) response.getData().get(ENCRYPTED_KEY_NAME);

			return "Encrypted value: " + value;
		}

		return "No value found for key";
	}
}