package com.company.template.auth

import android.content.Context
import androidx.credentials.ClearCredentialStateRequest
import androidx.credentials.CredentialManager
import com.google.firebase.auth.FirebaseAuth
import com.google.firebase.auth.FirebaseUser
import com.google.firebase.auth.GoogleAuthProvider
import com.google.firebase.auth.UserProfileChangeRequest
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.callbackFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.tasks.await

class FirebaseAuthRepository(
    private val context: Context,
    private val auth: FirebaseAuth = FirebaseAuth.getInstance(),
) : AuthRepository {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Default)

    override val authStateFlow: StateFlow<User?> =
        callbackFlow<FirebaseUser?> {
            val listener = FirebaseAuth.AuthStateListener { trySend(it.currentUser) }
            auth.addAuthStateListener(listener)
            awaitClose { auth.removeAuthStateListener(listener) }
        }.map { it?.toDomain() }
            .stateIn(
                scope = scope,
                started = SharingStarted.WhileSubscribed(5_000),
                initialValue = auth.currentUser?.toDomain(),
            )

    override suspend fun signInWithEmail(
        email: String,
        password: String,
    ): Result<Unit> =
        runCatching {
            auth.signInWithEmailAndPassword(email, password).await()
            Unit
        }

    override suspend fun registerWithEmail(
        name: String,
        email: String,
        password: String,
    ): Result<Unit> =
        runCatching {
            val result = auth.createUserWithEmailAndPassword(email, password).await()
            result.user
                ?.updateProfile(
                    UserProfileChangeRequest.Builder().setDisplayName(name).build(),
                )?.await()
            Unit
        }

    override suspend fun signInWithGoogle(googleIdToken: String): Result<Unit> =
        runCatching {
            val credential = GoogleAuthProvider.getCredential(googleIdToken, null)
            auth.signInWithCredential(credential).await()
            Unit
        }

    override suspend fun signOut() {
        auth.signOut()
        // Clear saved Google credential so the next sign-in shows the account picker
        CredentialManager
            .create(context)
            .clearCredentialState(ClearCredentialStateRequest())
    }

    private fun FirebaseUser.toDomain() =
        User(
            uid = uid,
            email = email,
            displayName = displayName,
            photoUrl = photoUrl?.toString(),
        )
}
