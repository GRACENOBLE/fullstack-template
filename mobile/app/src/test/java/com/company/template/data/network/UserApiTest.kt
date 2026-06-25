package com.company.template.data.network

import kotlinx.coroutines.test.runTest
import okhttp3.OkHttpClient
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class UserApiTest {

    private lateinit var server: MockWebServer

    // A plain OkHttpClient with no auth interceptor — tests don't need Firebase
    private val testClient = OkHttpClient()

    @Before
    fun setUp() {
        server = MockWebServer()
        server.start()
    }

    @After
    fun tearDown() {
        server.shutdown()
    }

    @Test
    fun `getMe returns UserProfile on successful response`() = runTest {
        server.enqueue(
            MockResponse()
                .setResponseCode(200)
                .setBody("""{"data":{"uid":"u1","email":"a@b.com"}}""")
        )

        val result = UserApi.getMe(
            baseUrl = server.url("/").toString().trimEnd('/'),
            client = testClient,
        )

        assertTrue(result.isSuccess)
        val profile = result.getOrThrow()
        assertEquals("u1", profile.uid)
        assertEquals("a@b.com", profile.email)
    }

    @Test
    fun `getMe returns failure with backend message on error response`() = runTest {
        server.enqueue(
            MockResponse()
                .setResponseCode(401)
                .setBody("""{"error":{"code":"UNAUTHENTICATED","message":"no token"}}""")
        )

        val result = UserApi.getMe(
            baseUrl = server.url("/").toString().trimEnd('/'),
            client = testClient,
        )

        assertTrue(result.isFailure)
        assertEquals("no token", result.exceptionOrNull()?.message)
    }
}
