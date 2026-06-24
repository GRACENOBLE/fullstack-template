package com.company.template.auth

import android.app.Activity
import com.google.firebase.auth.FirebaseUser
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

// --------------- Fake repository (no Mockito) ---------------

class FakeAuthRepository : AuthRepository {
    private val _authStateFlow = MutableStateFlow<FirebaseUser?>(null)
    override val authStateFlow: StateFlow<FirebaseUser?> = _authStateFlow

    var signInResult: Result<Unit> = Result.success(Unit)
    var registerResult: Result<Unit> = Result.success(Unit)
    var googleSignInResult: Result<Unit> = Result.success(Unit)
    var signOutCalled = false

    override suspend fun signInWithEmail(email: String, password: String): Result<Unit> =
        signInResult

    override suspend fun registerWithEmail(name: String, email: String, password: String): Result<Unit> =
        registerResult

    override suspend fun signInWithGoogle(activity: Activity): Result<Unit> = googleSignInResult

    override suspend fun signOut() {
        signOutCalled = true
        _authStateFlow.value = null
    }
}

// --------------- Tests ---------------

@OptIn(ExperimentalCoroutinesApi::class)
class AuthViewModelTest {

    private lateinit var fakeRepo: FakeAuthRepository
    private lateinit var viewModel: AuthViewModel

    @Before
    fun setUp() {
        Dispatchers.setMain(UnconfinedTestDispatcher())
        fakeRepo = FakeAuthRepository()
        viewModel = AuthViewModel(fakeRepo)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `signIn emits Loading then Success on repository success`() = runTest {
        val states = mutableListOf<AuthUiState>()
        // UnconfinedTestDispatcher runs coroutines eagerly so we can observe synchronously
        fakeRepo.signInResult = Result.success(Unit)

        viewModel.signIn("test@example.com", "password")

        // After UnconfinedTestDispatcher finishes, final state should be Success
        assertEquals(AuthUiState.Success, viewModel.uiState.value)
    }

    @Test
    fun `signIn emits Error on repository failure`() = runTest {
        fakeRepo.signInResult = Result.failure(Exception("Invalid credentials"))

        viewModel.signIn("test@example.com", "wrong")

        val state = viewModel.uiState.value
        assertTrue(state is AuthUiState.Error)
        assertEquals("Invalid credentials", (state as AuthUiState.Error).message)
    }

    @Test
    fun `register emits Success on repository success`() = runTest {
        fakeRepo.registerResult = Result.success(Unit)

        viewModel.register("Alice", "alice@example.com", "password1", "password1")

        assertEquals(AuthUiState.Success, viewModel.uiState.value)
    }

    @Test
    fun `register emits Error when passwords do not match`() = runTest {
        viewModel.register("Alice", "alice@example.com", "password1", "password2")

        val state = viewModel.uiState.value
        assertTrue(state is AuthUiState.Error)
        assertEquals("Passwords do not match", (state as AuthUiState.Error).message)
    }

    @Test
    fun `register emits Error on repository failure`() = runTest {
        fakeRepo.registerResult = Result.failure(Exception("Email already in use"))

        viewModel.register("Alice", "alice@example.com", "password1", "password1")

        val state = viewModel.uiState.value
        assertTrue(state is AuthUiState.Error)
        assertEquals("Email already in use", (state as AuthUiState.Error).message)
    }

    @Test
    fun `initial state is Idle`() {
        assertEquals(AuthUiState.Idle, viewModel.uiState.value)
    }

    @Test
    fun `signOut calls repository signOut`() = runTest {
        viewModel.signOut()
        assertTrue(fakeRepo.signOutCalled)
    }
}
