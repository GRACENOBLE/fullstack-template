package com.company.template.fcm

import kotlinx.serialization.json.Json
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class FcmRegistrationPayloadTest {
    private val json = Json { encodeDefaults = true }

    @Test
    fun serializes_token_and_platform() {
        val payload = FcmRegistrationPayload(token = "test-fcm-token-123")
        val encoded = json.encodeToString(FcmRegistrationPayload.serializer(), payload)

        assertTrue("JSON must contain token field", encoded.contains("\"token\""))
        assertTrue("JSON must contain platform field", encoded.contains("\"platform\""))
        assertTrue("JSON must contain token value", encoded.contains("test-fcm-token-123"))
        assertTrue("JSON must contain android platform", encoded.contains("\"android\""))
    }

    @Test
    fun default_platform_is_android() {
        val payload = FcmRegistrationPayload(token = "tok")
        assertEquals("android", payload.platform)
    }

    @Test
    fun platform_can_be_overridden() {
        val payload = FcmRegistrationPayload(token = "tok", platform = "ios")
        assertEquals("ios", payload.platform)
    }

    @Test
    fun roundtrip_serialization() {
        val original = FcmRegistrationPayload(token = "round-trip-token", platform = "android")
        val encoded = json.encodeToString(FcmRegistrationPayload.serializer(), original)
        val decoded = json.decodeFromString(FcmRegistrationPayload.serializer(), encoded)
        assertEquals(original, decoded)
    }
}
