package org.paasify.cfvault.example.integration;

import java.time.Duration;
import java.util.Map;

import org.springframework.util.Assert;
import org.springframework.vault.authentication.LoginToken;

class LoginTokenUtil {

	/**
	 * Construct a {@link LoginToken} from an auth response.
	 *
	 * @param auth {@link Map} holding a login response.
	 * @return the {@link LoginToken}
	 */
	static LoginToken from(Map<String, Object> auth) {

		Assert.notNull(auth, "Authentication must not be null");

		String token = (String) auth.get("client_token");

		return from(token.toCharArray(), auth);
	}

	/**
	 * Construct a {@link LoginToken} from an auth response.
	 *
	 * @param auth {@link Map} holding a login response.
	 * @return the {@link LoginToken}
	 * @since 2.0
	 */
	static LoginToken from(char[] token, Map<String, Object> auth) {

		Assert.notNull(auth, "Authentication must not be null");

		Boolean renewable = (Boolean) auth.get("renewable");
		Number leaseDuration = (Number) auth.get("lease_duration");

		if (leaseDuration == null) {
			leaseDuration = (Number) auth.get("ttl");
		}

		if (renewable != null && renewable) {
			return LoginToken.renewable(token,
					Duration.ofSeconds(leaseDuration.longValue()));
		}

		if (leaseDuration != null) {
			return LoginToken.of(token, Duration.ofSeconds(leaseDuration.longValue()));
		}

		return LoginToken.of(token);
	}
}