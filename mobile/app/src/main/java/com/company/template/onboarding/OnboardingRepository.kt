package com.company.template.onboarding

import kotlinx.coroutines.flow.Flow

interface OnboardingRepository {
    fun hasSeenOnboarding(): Flow<Boolean>

    suspend fun markSeen()
}
