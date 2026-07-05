package ru.bonifero.crypto.config;

import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Provider;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.security.Security;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.Date;

@Configuration
@EnableConfigurationProperties(KeyStoreProperties.class)
public class KeyStoreConfig {

    private static final String CERT_SUBJECT = "CN=crypto-service, OU=Cybersecurity, O=Bank, C=RU";
    private static final long CERT_VALIDITY_DAYS = 3650;

    @Bean
    public Provider bouncyCastleProvider() {
        Provider provider = new BouncyCastleProvider();
        Security.addProvider(provider);
        return provider;
    }

    @Bean
    public KeyPair rsaKeyPair(KeyStoreProperties properties, Provider bouncyCastleProvider) throws Exception {
        if (properties.getPassword() == null || properties.getPassword().isBlank()) {
            throw new IllegalStateException(
                    "Пароль хранилища ключей не задан. Установите переменную окружения CRYPTO_KEYSTORE_PASSWORD.");
        }

        Path keystorePath = Path.of(properties.getPath());
        char[] password = properties.getPassword().toCharArray();

        if (Files.exists(keystorePath)) {
            return loadKeyPair(keystorePath, password, properties.getAlias());
        }
        return generateAndStoreKeyPair(keystorePath, password, properties, bouncyCastleProvider);
    }

    private KeyPair loadKeyPair(Path keystorePath, char[] password, String alias) throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        try (FileInputStream in = new FileInputStream(keystorePath.toFile())) {
            keyStore.load(in, password);
        }
        PrivateKey privateKey = (PrivateKey) keyStore.getKey(alias, password);
        Certificate certificate = keyStore.getCertificate(alias);
        return new KeyPair(certificate.getPublicKey(), privateKey);
    }

    private KeyPair generateAndStoreKeyPair(Path keystorePath, char[] password,
                                             KeyStoreProperties properties, Provider provider) throws Exception {
        if (keystorePath.getParent() != null) {
            Files.createDirectories(keystorePath.getParent());
        }

        KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA", provider);
        generator.initialize(properties.getKeySize(), new SecureRandom());
        KeyPair keyPair = generator.generateKeyPair();

        X509Certificate certificate = selfSignedCertificate(keyPair, provider);

        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, null);
        keyStore.setKeyEntry(properties.getAlias(), keyPair.getPrivate(), password,
                new Certificate[]{certificate});

        try (FileOutputStream out = new FileOutputStream(keystorePath.toFile())) {
            keyStore.store(out, password);
        }
        return keyPair;
    }

    private X509Certificate selfSignedCertificate(KeyPair keyPair, Provider provider) throws Exception {
        X500Name subject = new X500Name(CERT_SUBJECT);
        BigInteger serial = BigInteger.valueOf(System.currentTimeMillis());
        Date notBefore = new Date();
        Date notAfter = new Date(notBefore.getTime() + CERT_VALIDITY_DAYS * 24L * 60 * 60 * 1000);

        X509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(
                subject, serial, notBefore, notAfter, subject, keyPair.getPublic());

        ContentSigner signer = new JcaContentSignerBuilder("SHA256WithRSA")
                .setProvider(provider)
                .build(keyPair.getPrivate());

        return new JcaX509CertificateConverter()
                .setProvider(provider)
                .getCertificate(builder.build(signer));
    }

    @Bean
    public PublicKey rsaPublicKey(KeyPair rsaKeyPair) {
        return rsaKeyPair.getPublic();
    }

    @Bean
    public PrivateKey rsaPrivateKey(KeyPair rsaKeyPair) {
        return rsaKeyPair.getPrivate();
    }
}
