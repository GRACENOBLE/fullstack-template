package com.company.template.websocket

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener

/**
 * Factory abstraction over [OkHttpClient.newWebSocket] — injectable for unit tests.
 */
fun interface WebSocketFactory {
    fun newWebSocket(
        request: Request,
        listener: WebSocketListener,
    ): WebSocket
}

/**
 * Manages a single OkHttp WebSocket connection with exponential-backoff reconnection.
 *
 * Callbacks fire on OkHttp's dispatcher thread. Callers that update UI state must
 * marshal to the main thread (e.g. via StateFlow or viewModelScope).
 *
 * [reconnectScheduler] is injectable so unit tests can drive retries synchronously.
 * The default implementation uses [android.os.Handler] on the main looper.
 */
class WebSocketManager(
    private val serverUrl: String,
    private val factory: WebSocketFactory = defaultFactory(),
    val maxRetries: Int = 10,
    private val reconnectScheduler: (delayMs: Long, action: () -> Unit) -> Unit = { delay, action ->
        android.os.Handler(android.os.Looper.getMainLooper()).postDelayed(action, delay)
    },
) {
    var onOpen: (() -> Unit)? = null
    var onClose: (() -> Unit)? = null
    var onMessage: ((WsEnvelope) -> Unit)? = null
    var onError: ((Throwable) -> Unit)? = null

    private var socket: WebSocket? = null
    private var retryCount = 0
    private var currentToken: String? = null
    var active = false
        private set

    fun connect(token: String?) {
        socket?.close(1000, "reconnecting")
        socket = null
        active = true
        currentToken = token
        retryCount = 0
        openSocket()
    }

    fun disconnect() {
        active = false
        socket?.close(1000, "client disconnect")
        socket = null
    }

    fun send(envelope: WsEnvelope): Boolean {
        val json = Json.encodeToString(envelope)
        return socket?.send(json) ?: false
    }

    private fun openSocket() {
        val urlBuilder = StringBuilder(serverUrl)
        currentToken?.let { urlBuilder.append("?token=").append(it) }

        val request =
            Request
                .Builder()
                .url(urlBuilder.toString())
                .build()

        socket = factory.newWebSocket(request, listener)
    }

    private val listener =
        object : WebSocketListener() {
            override fun onOpen(
                webSocket: WebSocket,
                response: Response,
            ) {
                retryCount = 0
                onOpen?.invoke()
            }

            override fun onMessage(
                webSocket: WebSocket,
                text: String,
            ) {
                parseEnvelope(text)?.let { this@WebSocketManager.onMessage?.invoke(it) }
            }

            override fun onClosed(
                webSocket: WebSocket,
                code: Int,
                reason: String,
            ) {
                onClose?.invoke()
            }

            override fun onFailure(
                webSocket: WebSocket,
                t: Throwable,
                response: Response?,
            ) {
                onError?.invoke(t)
                if (active && retryCount < maxRetries) {
                    val delay = minOf(1000L * (1L shl retryCount), 30_000L)
                    retryCount++
                    reconnectScheduler(delay) { if (active) openSocket() }
                } else {
                    onClose?.invoke()
                }
            }
        }

    private fun parseEnvelope(text: String): WsEnvelope? =
        try {
            Json.decodeFromString<WsEnvelope>(text)
        } catch (_: Exception) {
            null
        }

    companion object {
        private fun defaultFactory(): WebSocketFactory {
            val client = OkHttpClient()
            return WebSocketFactory { request, listener -> client.newWebSocket(request, listener) }
        }
    }
}
