package com.company.template.home

import com.company.template.data.network.UserProfile
import com.company.template.ui.state.UiState
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import okhttp3.OkHttpClient
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class HomeViewModelTest {
    private lateinit var server: MockWebServer
    private val testDispatcher = UnconfinedTestDispatcher()

    @Before
    fun setUp() {
        server = MockWebServer()
        server.start()
        Dispatchers.setMain(testDispatcher)
    }

    @After
    fun tearDown() {
        server.shutdown()
        Dispatchers.resetMain()
    }

    private fun buildViewModel(): HomeViewModel =
        HomeViewModel(
            baseUrl = server.url("/").toString().trimEnd('/'),
            httpClient = OkHttpClient(),
            ioDispatcher = testDispatcher,
        )

    @Test
    fun `init - success response transitions state to Success with correct data`() =
        runTest {
            server.enqueue(
                MockResponse()
                    .setResponseCode(200)
                    .setHeader("Content-Type", "application/json")
                    .setBody("""{"data":{"uid":"u1","email":"a@b.com","displayName":"Test"}}"""),
            )

            val viewModel = buildViewModel()

            val state = viewModel.profileState.value
            assertTrue("Expected Success but was $state", state is UiState.Success<*>)
            val profile = (state as UiState.Success<*>).data as UserProfile
            assertEquals("u1", profile.uid)
            assertEquals("a@b.com", profile.email)
            assertEquals("Test", profile.displayName)
        }

    @Test
    fun `init - 401 response transitions state to Error`() =
        runTest {
            server.enqueue(
                MockResponse()
                    .setResponseCode(401)
                    .setHeader("Content-Type", "application/json")
                    .setBody("""{"error":{"message":"unauthorized"}}"""),
            )

            val viewModel = buildViewModel()

            val state = viewModel.profileState.value
            assertTrue("Expected Error but was $state", state is UiState.Error)
            assertTrue((state as UiState.Error).message.isNotEmpty())
        }

    @Test
    fun `refresh - re-fetches and updates state`() =
        runTest {
            // First call (from init) — error
            server.enqueue(
                MockResponse()
                    .setResponseCode(401)
                    .setHeader("Content-Type", "application/json")
                    .setBody("""{"error":{"message":"unauthorized"}}"""),
            )
            // Second call (from refresh) — success
            server.enqueue(
                MockResponse()
                    .setResponseCode(200)
                    .setHeader("Content-Type", "application/json")
                    .setBody("""{"data":{"uid":"u2","email":"b@c.com","displayName":"Bob"}}"""),
            )

            val viewModel = buildViewModel()
            assertTrue(viewModel.profileState.value is UiState.Error)

            viewModel.refresh()

            val state = viewModel.profileState.value
            assertTrue("Expected Success but was $state", state is UiState.Success<*>)
            val profile = (state as UiState.Success<*>).data as UserProfile
            assertEquals("u2", profile.uid)
            assertEquals("Bob", profile.displayName)
        }
}
