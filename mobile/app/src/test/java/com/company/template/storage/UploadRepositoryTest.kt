package com.company.template.storage

import kotlinx.coroutines.test.runTest
import okhttp3.OkHttpClient
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class UploadRepositoryTest {

    private lateinit var backendServer: MockWebServer
    private lateinit var r2Server: MockWebServer
    private lateinit var repository: R2UploadRepository

    @Before
    fun setUp() {
        backendServer = MockWebServer()
        r2Server = MockWebServer()
        backendServer.start()
        r2Server.start()
        repository = R2UploadRepository(
            backendBaseUrl = backendServer.url("").toString().trimEnd('/'),
            httpClient = OkHttpClient(),
        )
    }

    @After
    fun tearDown() {
        backendServer.shutdown()
        r2Server.shutdown()
    }

    @Test
    fun `upload returns public URL on success`() = runTest {
        val r2Url = r2Server.url("/test-key").toString()
        backendServer.enqueue(
            MockResponse()
                .setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setBody("""{"upload_url":"$r2Url","public_url":"https://pub.r2.dev/test-key"}"""),
        )
        r2Server.enqueue(MockResponse().setResponseCode(200))

        val result = repository.upload(
            filename = "test.jpg",
            contentType = "image/jpeg",
            fileBytes = byteArrayOf(1, 2, 3),
            idToken = "fake-token",
        )

        assertTrue(result.isSuccess)
        assertEquals("https://pub.r2.dev/test-key", result.getOrNull())

        val presignReq = backendServer.takeRequest()
        assertEquals("POST", presignReq.method)
        assertEquals("/api/v1/storage/presign", presignReq.path)
        assertEquals("Bearer fake-token", presignReq.getHeader("Authorization"))

        val r2Req = r2Server.takeRequest()
        assertEquals("PUT", r2Req.method)
        assertEquals("image/jpeg", r2Req.getHeader("Content-Type"))
    }

    @Test
    fun `upload returns failure when presign fails`() = runTest {
        backendServer.enqueue(MockResponse().setResponseCode(401))

        val result = repository.upload(
            filename = "test.jpg",
            contentType = "image/jpeg",
            fileBytes = byteArrayOf(),
            idToken = "bad-token",
        )

        assertTrue(result.isFailure)
    }

    @Test
    fun `upload returns failure when R2 PUT fails`() = runTest {
        val r2Url = r2Server.url("/test-key").toString()
        backendServer.enqueue(
            MockResponse()
                .setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setBody("""{"upload_url":"$r2Url","public_url":"https://pub.r2.dev/test-key"}"""),
        )
        r2Server.enqueue(MockResponse().setResponseCode(403))

        val result = repository.upload(
            filename = "test.jpg",
            contentType = "image/jpeg",
            fileBytes = byteArrayOf(1, 2, 3),
            idToken = "token",
        )

        assertTrue(result.isFailure)
    }
}
