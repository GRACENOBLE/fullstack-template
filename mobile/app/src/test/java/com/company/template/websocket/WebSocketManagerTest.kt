package com.company.template.websocket

import okhttp3.Protocol
import okhttp3.Request
import okhttp3.Response
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class WebSocketManagerTest {
    private lateinit var factory: FakeWebSocketFactory
    private val scheduledActions = mutableListOf<() -> Unit>()

    private fun makeManager(maxRetries: Int = 3): WebSocketManager =
        WebSocketManager(
            serverUrl = "ws://localhost:8080/ws",
            factory = factory,
            maxRetries = maxRetries,
            reconnectScheduler = { _, action -> scheduledActions.add(action) },
        )

    private fun fakeResponse(): Response =
        Response
            .Builder()
            .request(Request.Builder().url("ws://localhost:8080/ws").build())
            .protocol(Protocol.HTTP_1_1)
            .code(101)
            .message("Switching Protocols")
            .build()

    @Before
    fun setUp() {
        factory = FakeWebSocketFactory()
        scheduledActions.clear()
    }

    @Test
    fun `onOpen callback is invoked on successful connection`() {
        var opened = false
        val manager = makeManager()
        manager.onOpen = { opened = true }

        manager.connect(token = null)
        factory.lastListener!!.onOpen(factory.lastSocket, fakeResponse())

        assertTrue(opened)
    }

    @Test
    fun `onMessage delivers parsed envelope`() {
        var received: WsEnvelope? = null
        val manager = makeManager()
        manager.onMessage = { received = it }

        manager.connect(token = null)
        factory.lastListener!!.onMessage(
            factory.lastSocket,
            """{"type":"job.completed","payload":{"id":"42"}}""",
        )

        assertEquals("job.completed", received?.type)
    }

    @Test
    fun `onMessage ignores malformed JSON without crashing`() {
        var called = false
        val manager = makeManager()
        manager.onMessage = { called = true }

        manager.connect(token = null)
        factory.lastListener!!.onMessage(factory.lastSocket, "not-json")

        assertFalse(called)
    }

    @Test
    fun `onMessage returns null for missing type field`() {
        var received: WsEnvelope? = null
        val manager = makeManager()
        manager.onMessage = { received = it }

        manager.connect(token = null)
        factory.lastListener!!.onMessage(factory.lastSocket, """{"payload":42}""")

        assertNull(received)
    }

    @Test
    fun `onClose callback fires on graceful close`() {
        var closed = false
        val manager = makeManager()
        manager.onClose = { closed = true }

        manager.connect(token = null)
        factory.lastListener!!.onClosed(factory.lastSocket, 1000, "normal")

        assertTrue(closed)
    }

    @Test
    fun `onFailure schedules retry when retries remain`() {
        val manager = makeManager(maxRetries = 2)
        manager.connect(token = null)

        factory.lastListener!!.onFailure(factory.lastSocket, RuntimeException("refused"), null)

        assertEquals(1, scheduledActions.size)
    }

    @Test
    fun `onFailure does not retry after maxRetries exceeded`() {
        val manager = makeManager(maxRetries = 1)
        manager.connect(token = null)

        factory.lastListener!!.onFailure(factory.lastSocket, RuntimeException("fail"), null)
        scheduledActions.last().invoke()
        factory.lastListener!!.onFailure(factory.lastSocket, RuntimeException("fail"), null)

        // First failure schedules 1 retry; second failure finds retryCount >= maxRetries.
        assertEquals(1, scheduledActions.size)
    }

    @Test
    fun `disconnect prevents scheduled reconnect from opening new socket`() {
        val manager = makeManager(maxRetries = 3)
        manager.connect(token = null)

        factory.lastListener!!.onFailure(factory.lastSocket, RuntimeException("fail"), null)
        assertEquals(1, scheduledActions.size)

        manager.disconnect()
        scheduledActions.first().invoke()

        // After disconnect + action fires, no new socket should be opened.
        assertEquals(1, factory.sockets.size)
    }
}
