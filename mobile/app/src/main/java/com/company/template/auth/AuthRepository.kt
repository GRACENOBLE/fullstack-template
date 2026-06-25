package com.company.template.auth

import kotlinx.coroutines.flow.StateFlow

interface AuthRepository {
    val authStateFlow: StateFlow<User?>

    suspend fun signInWithEmail(
        email: String,
        password: String,
    ): Result<Unit>

    suspend fun registerWithEmail(
        name: String,
        email: String,
        password: String,
    ): Result<Unit>

    suspend fun signInWithGoogle(googleIdToken: String): Result<Unit>

    suspend fun signOut()
}
