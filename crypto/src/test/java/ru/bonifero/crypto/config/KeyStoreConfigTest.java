package ru.bonifero.crypto.config;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Path;
import java.security.KeyPair;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class KeyStoreConfigTest {

    private final KeyStoreConfig config = new KeyStoreConfig();

    @Test
    void rsaKeyPair_generatesThenReloadsSameKeyFromDisk(@TempDir Path tempDir) throws Exception {
        KeyStoreProperties properties = new KeyStoreProperties();
        properties.setPath(tempDir.resolve("test-keystore.p12").toString());
        properties.setPassword("test-password");
        properties.setAlias("crypto-service");
        properties.setKeySize(2048);

        var provider = config.bouncyCastleProvider();

        KeyPair generated = config.rsaKeyPair(properties, provider);
        assertTrue(tempDir.resolve("test-keystore.p12").toFile().exists(),
                "keystore file should be created on first run");

        KeyPair reloaded = config.rsaKeyPair(properties, provider);

        assertEquals(generated.getPrivate(), reloaded.getPrivate(),
                "second call must load the same key from disk, not generate a new one");
        assertEquals(generated.getPublic(), reloaded.getPublic());
    }

    @Test
    void rsaKeyPair_rejectsMissingPassword(@TempDir Path tempDir) {
        KeyStoreProperties properties = new KeyStoreProperties();
        properties.setPath(tempDir.resolve("test-keystore.p12").toString());
        properties.setPassword("");

        var provider = config.bouncyCastleProvider();

        assertThrows(IllegalStateException.class, () -> config.rsaKeyPair(properties, provider));
    }

    @Test
    void rsaKeyPair_rejectsWrongPasswordOnReload(@TempDir Path tempDir) throws Exception {
        KeyStoreProperties properties = new KeyStoreProperties();
        properties.setPath(tempDir.resolve("test-keystore.p12").toString());
        properties.setPassword("correct-password");
        properties.setAlias("crypto-service");
        properties.setKeySize(2048);

        var provider = config.bouncyCastleProvider();
        config.rsaKeyPair(properties, provider);

        properties.setPassword("wrong-password");
        assertThrows(Exception.class, () -> config.rsaKeyPair(properties, provider));
    }
}
