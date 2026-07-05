package ru.bonifero.crypto.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI cryptoServiceOpenApi() {
        return new OpenAPI()
                .info(new Info()
                        .title("Crypto Service API")
                        .description("Микросервис шифрования/дешифрования сообщений и расчета хеша (тестовое задание)")
                        .version("1.0.0")
                        .contact(new Contact().name("ru.bonifero")));
    }
}
