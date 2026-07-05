package ru.bonifero.crypto.crypto;

import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import ru.bonifero.crypto.exception.CryptoServiceException;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.Provider;
import java.security.Security;
import java.util.Base64;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class CryptoServiceTest {

    private static CryptoService cryptoService;

    @BeforeAll
    static void setUp() throws Exception {
        Provider provider = new BouncyCastleProvider();
        Security.addProvider(provider);

        KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA", provider);
        generator.initialize(2048);
        KeyPair keyPair = generator.generateKeyPair();

        cryptoService = new CryptoService(keyPair.getPublic(), keyPair.getPrivate(), provider);
    }

    @Test
    void encryptThenDecrypt_returnsOriginalMessage() {
        String message = "Тестовое сообщение для шифрования 123!";

        HybridEncryptionResult encrypted = cryptoService.encrypt(message);
        assertNotNull(encrypted.encryptedKey());
        assertNotNull(encrypted.iv());
        assertNotNull(encrypted.cipherText());
        assertNotEquals(message, encrypted.cipherText());

        String decrypted = cryptoService.decrypt(encrypted);
        assertEquals(message, decrypted);
    }

    @Test
    void encrypt_producesDifferentCipherTextForSameMessage() {
        String message = "Repeated message";

        HybridEncryptionResult first = cryptoService.encrypt(message);
        HybridEncryptionResult second = cryptoService.encrypt(message);

        assertNotEquals(first.iv(), second.iv(), "IV must be random per call");
        assertNotEquals(first.cipherText(), second.cipherText(),
                "different IV/AES key must yield different cipher text");
    }

    @Test
    void decrypt_rejectsTamperedCipherText() {
        HybridEncryptionResult encrypted = cryptoService.encrypt("original message");

        byte[] tampered = Base64.getDecoder().decode(encrypted.cipherText());
        tampered[0] ^= 0x01;
        HybridEncryptionResult withTamperedCipherText = new HybridEncryptionResult(
                encrypted.encryptedKey(), encrypted.iv(), Base64.getEncoder().encodeToString(tampered));

        assertThrows(CryptoServiceException.class, () -> cryptoService.decrypt(withTamperedCipherText));
    }

    @Test
    void decrypt_rejectsInvalidBase64() {
        HybridEncryptionResult invalid = new HybridEncryptionResult("not-base64!!", "not-base64!!", "not-base64!!");

        assertThrows(CryptoServiceException.class, () -> cryptoService.decrypt(invalid));
    }

    @Test
    void sha256Hash_isDeterministic() {
        String message = "Hello, world!";

        HashResult first = cryptoService.hash(message, HashAlgorithm.SHA_256);
        HashResult second = cryptoService.hash(message, HashAlgorithm.SHA_256);

        assertEquals("SHA-256", first.algorithm());
        assertEquals(first.hashHex(), second.hashHex());
        assertEquals(64, first.hashHex().length());
    }

    @Test
    void pbkdf2Hash_usesRandomSaltAndIsVerifiable() {
        String message = "Hello, world!";

        HashResult first = cryptoService.hash(message, HashAlgorithm.PBKDF2);
        HashResult second = cryptoService.hash(message, HashAlgorithm.PBKDF2);

        assertNotNull(first.saltHex());
        assertNotEquals(first.saltHex(), second.saltHex(), "salt must be random per call");
        assertNotEquals(first.hashHex(), second.hashHex(), "different salt must yield different hash");
        assertEquals(100_000, first.iterations());
    }
}
