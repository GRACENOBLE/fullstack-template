package com.company.template.fcm

import kotlinx.serialization.Serializable

@Serializable
data class FcmRegistrationPayload(
    val token: String,
    val platform: String = "android",
)
