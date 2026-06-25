package com.company.template.websocket

import okhttp3.Request
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString

/**
 * Test double for [WebSocketFactory]. Captures listener and created sockets so tests
 * can drive WebSocket callbacks without real network I/O.
 */
class FakeWebSocketFactory : WebSocketFactory {
    val sockets = mutableListOf<FakeWebSocket>()
    var lastListener: WebSocketListener? = null

    val lastSocket: FakeWebSocket get() = sockets.last()

    override fun newWebSocket(
        request: Request,
        listener: WebSocketListener,
    ): WebSocket {
        lastListener = listener
        return FakeWebSocket(request).also { sockets.add(it) }
    }
}

class FakeWebSocket(
    private val req: Request = Request.Builder().url("ws://localhost/ws").build(),
) : WebSocket {
    val sentMessages = mutableListOf<String>()
    var closed = false

    override fun request(): Request = req

    override fun queueSize(): Long = 0L

    override fun send(text: String): Boolean {
        sentMessages.add(text)
        return true
    }

    override fun send(bytes: ByteString): Boolean = false

    override fun close(
        code: Int,
        reason: String?,
    ): Boolean {
        closed = true
        return true
    }

    override fun cancel() {
        closed = true
    }
}
