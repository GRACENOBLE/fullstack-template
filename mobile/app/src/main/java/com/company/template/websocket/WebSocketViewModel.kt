package com.company.template.websocket

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update

data class WsState(
    val isConnected: Boolean = false,
    val lastMessage: WsEnvelope? = null,
)

class WebSocketViewModel(
    private val manager: WebSocketManager,
) : ViewModel() {
    private val _state = MutableStateFlow(WsState())
    val state: StateFlow<WsState> = _state.asStateFlow()

    init {
        manager.onOpen = { _state.update { it.copy(isConnected = true) } }
        manager.onClose = { _state.update { it.copy(isConnected = false) } }
        manager.onMessage = { envelope ->
            _state.update { it.copy(lastMessage = envelope) }
        }
    }

    fun connect(token: String? = null) {
        manager.connect(token)
    }

    fun disconnect() {
        manager.disconnect()
    }

    fun send(envelope: WsEnvelope) {
        manager.send(envelope)
    }

    public override fun onCleared() {
        manager.disconnect()
    }
}
