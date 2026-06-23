package com.company.template.storage

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody

interface UploadRepository {
    suspend fun upload(
        filename: String,
        contentType: String,
        fileBytes: ByteArray,
        idToken: String,
    ): Result<String>  // returns public URL
}

class R2UploadRepository(
    private val backendBaseUrl: String,
    private val httpClient: OkHttpClient = OkHttpClient(),
) : UploadRepository {

    @Serializable
    private data class PresignRequest(
        val filename: String,
        @SerialName("content_type") val contentType: String,
    )

    @Serializable
    private data class PresignResponse(
        @SerialName("upload_url") val uploadUrl: String,
        @SerialName("public_url") val publicUrl: String,
    )

    private val json = Json { ignoreUnknownKeys = true }

    override suspend fun upload(
        filename: String,
        contentType: String,
        fileBytes: ByteArray,
        idToken: String,
    ): Result<String> = withContext(Dispatchers.IO) {
        runCatching {
            val presignResponse = presign(filename, contentType, idToken)
            uploadToR2(presignResponse.uploadUrl, fileBytes, contentType)
            presignResponse.publicUrl
        }
    }

    private fun presign(filename: String, contentType: String, idToken: String): PresignResponse {
        val payload = json.encodeToString(PresignRequest.serializer(), PresignRequest(filename, contentType))
        val request = Request.Builder()
            .url("$backendBaseUrl/api/v1/storage/presign")
            .post(payload.toRequestBody("application/json".toMediaType()))
            .header("Authorization", "Bearer $idToken")
            .build()
        val response = httpClient.newCall(request).execute()
        check(response.isSuccessful) { "presign failed: ${response.code}" }
        val body = checkNotNull(response.body?.string()) { "presign: empty body" }
        return json.decodeFromString(PresignResponse.serializer(), body)
    }

    private fun uploadToR2(uploadUrl: String, fileBytes: ByteArray, contentType: String) {
        val request = Request.Builder()
            .url(uploadUrl)
            .put(fileBytes.toRequestBody(contentType.toMediaType()))
            .build()
        val response = httpClient.newCall(request).execute()
        check(response.isSuccessful) { "R2 upload failed: ${response.code}" }
    }
}
