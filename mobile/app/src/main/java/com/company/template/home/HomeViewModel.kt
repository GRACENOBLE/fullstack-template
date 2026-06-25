package com.company.template.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import com.company.template.BuildConfig
import com.company.template.data.network.ApiClient
import com.company.template.data.network.UserApi
import com.company.template.data.network.UserProfile
import com.company.template.ui.state.UiState
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import okhttp3.OkHttpClient

class HomeViewModel(
    private val baseUrl: String,
    private val httpClient: OkHttpClient,
    private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : ViewModel() {
    private val _profileState = MutableStateFlow<UiState<UserProfile>>(UiState.Idle)
    val profileState: StateFlow<UiState<UserProfile>> = _profileState.asStateFlow()

    init {
        fetchProfile()
    }

    fun refresh() {
        fetchProfile()
    }

    private fun fetchProfile() {
        viewModelScope.launch(ioDispatcher) {
            _profileState.value = UiState.Loading
            UserApi
                .getMe(baseUrl = baseUrl, client = httpClient)
                .onSuccess { _profileState.value = UiState.Success(it) }
                .onFailure { _profileState.value = UiState.Error(it.message ?: "Failed to load profile") }
        }
    }

    companion object {
        fun factory(): ViewModelProvider.Factory =
            viewModelFactory {
                initializer {
                    HomeViewModel(
                        baseUrl = BuildConfig.BACKEND_URL,
                        httpClient = ApiClient.httpClient,
                    )
                }
            }
    }
}
