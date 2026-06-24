package com.company.template.navigation

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import com.company.template.auth.AuthRepository
import com.company.template.onboarding.OnboardingRepository
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.combine
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch

enum class StartDestination {
    ONBOARDING,
    LOGIN,
    HOME
}

class AppViewModel(
    authRepository: AuthRepository,
    private val onboardingRepository: OnboardingRepository,
) : ViewModel() {

    val startDestination: StateFlow<StartDestination?> =
        combine(
            onboardingRepository.hasSeenOnboarding(),
            authRepository.authStateFlow,
        ) { hasSeen, user ->
            when {
                user != null -> StartDestination.HOME
                !hasSeen -> StartDestination.ONBOARDING
                else -> StartDestination.LOGIN
            }
        }.stateIn(
            scope = viewModelScope,
            started = SharingStarted.WhileSubscribed(5_000),
            initialValue = null,
        )

    fun markOnboardingSeen() {
        viewModelScope.launch {
            runCatching { onboardingRepository.markSeen() }
        }
    }

    companion object {
        fun factory(
            authRepository: AuthRepository,
            onboardingRepository: OnboardingRepository,
        ): ViewModelProvider.Factory = viewModelFactory {
            initializer { AppViewModel(authRepository, onboardingRepository) }
        }
    }
}
