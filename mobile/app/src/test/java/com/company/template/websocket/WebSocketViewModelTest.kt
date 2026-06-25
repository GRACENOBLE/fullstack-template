package com.company.template.websocket

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class WebSocketViewModelTest {
    private lateinit var factory: FakeWebSocketFactory
    private lateinit var manager: WebSocketManager
    private lateinit var viewModel: WebSocketViewModel

    @Before
    fun setUp() {
        Dispatchers.setMain(UnconfinedTestDispatcher())
        factory = FakeWebSocketFactory()
        manager =
            WebSocketManager(
                serverUrl = "ws://localhost:8080/ws",
                factory = factory,
                reconnectScheduler = { _, _ -> },
            )
        viewModel = WebSocketViewModel(manager)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `initial state has isConnected false and no lastMessage`() {
        assertFalse(viewModel.state.value.isConnected)
        assertNull(viewModel.state.value.lastMessage)
    }

    @Test
    fun `connect then open sets isConnected true`() {
        viewModel.connect()
        manager.onOpen!!.invoke()

        assertTrue(viewModel.state.value.isConnected)
    }

    @Test
    fun `close event sets isConnected false`() {
        viewModel.connect()
        manager.onOpen!!.invoke()
        manager.onClose!!.invoke()

        assertFalse(viewModel.state.value.isConnected)
    }

    @Test
    fun `incoming message updates lastMessage`() {
        val envelope = WsEnvelope(type = "ping")
        viewModel.connect()
        manager.onMessage!!.invoke(envelope)

        assertEquals(envelope, viewModel.state.value.lastMessage)
    }

    @Test
    fun `disconnect calls manager disconnect`() {
        viewModel.connect()
        viewModel.disconnect()

        assertFalse(manager.active)
    }

    @Test
    fun `onCleared calls manager disconnect`() {
        viewModel.connect()
        viewModel.onCleared()

        assertFalse(manager.active)
    }
}
