package com.company.template.auth

import android.app.Activity
import com.google.firebase.auth.FirebaseUser
import kotlinx.coroutines.flow.StateFlow

interface AuthRepository {
    val authStateFlow: StateFlow<FirebaseUser?>
    suspend fun signInWithEmail(email: String, password: String): Result<Unit>
    suspend fun registerWithEmail(name: String, email: String, password: String): Result<Unit>
    suspend fun signInWithGoogle(activity: Activity): Result<Unit>
    suspend fun signOut()
}
