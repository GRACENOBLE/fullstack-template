package com.company.template.data.network

import com.company.template.BuildConfig
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import okhttp3.OkHttpClient
import okhttp3.Request

@Serializable
data class UserProfile(
    val uid: String,
    val email: String? = null,
    val displayName: String? = null,
)

object UserApi {
    private val json = Json { ignoreUnknownKeys = true }

    suspend fun getMe(
        baseUrl: String = BuildConfig.BACKEND_URL,
        client: OkHttpClient = ApiClient.httpClient,
    ): Result<UserProfile> = runCatching {
        val request = Request.Builder()
            .url("$baseUrl/api/v1/me")
            .get()
            .build()

        client.newCall(request).execute().use { response ->
            val body = response.body?.string() ?: error("empty body")
            if (!response.isSuccessful) {
                val err = json.decodeFromString<ApiErrorResponse>(body)
                error(err.error.message)
            }
            json.decodeFromString<ApiResponse<UserProfile>>(body).data
        }
    }
}
