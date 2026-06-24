package com.company.template.onboarding

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.launch

class OnboardingViewModel(private val repo: OnboardingRepository) : ViewModel() {

    fun hasSeenOnboarding(): Flow<Boolean> = repo.hasSeenOnboarding()

    fun markSeen(onComplete: () -> Unit = {}) {
        viewModelScope.launch {
            repo.markSeen()
            onComplete()
        }
    }

    companion object {
        fun factory(repo: OnboardingRepository): ViewModelProvider.Factory = viewModelFactory {
            initializer { OnboardingViewModel(repo) }
        }
    }
}
