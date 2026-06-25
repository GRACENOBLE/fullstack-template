package com.company.template.auth

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
    private val _authStateFlow = MutableStateFlow<User?>(null)
    override val authStateFlow: StateFlow<User?> = _authStateFlow

    var signInResult: Result<Unit> = Result.success(Unit)
    var registerResult: Result<Unit> = Result.success(Unit)
    var googleSignInResult: Result<Unit> = Result.success(Unit)
    var signOutCalled = false

    fun setUser(user: User?) {
        _authStateFlow.value = user
    }

    override suspend fun signInWithEmail(
        email: String,
        password: String,
    ): Result<Unit> = signInResult

    override suspend fun registerWithEmail(
        name: String,
        email: String,
        password: String,
    ): Result<Unit> = registerResult

    override suspend fun signInWithGoogle(googleIdToken: String): Result<Unit> = googleSignInResult

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
    fun `initial state is Idle`() {
        assertEquals(AuthUiState.Idle, viewModel.uiState.value)
    }

    @Test
    fun `signIn emits Success on repository success`() =
        runTest {
            fakeRepo.signInResult = Result.success(Unit)
            viewModel.updateLoginEmail("test@example.com")
            viewModel.updateLoginPassword("password")
            viewModel.signIn()
            assertEquals(AuthUiState.Success, viewModel.uiState.value)
        }

    @Test
    fun `signIn emits Error on repository failure`() =
        runTest {
            fakeRepo.signInResult = Result.failure(Exception("Invalid credentials"))
            viewModel.updateLoginEmail("test@example.com")
            viewModel.updateLoginPassword("wrong")
            viewModel.signIn()
            val state = viewModel.uiState.value
            assertTrue(state is AuthUiState.Error)
            assertEquals("Invalid credentials", (state as AuthUiState.Error).message)
        }

    @Test
    fun `register emits Success on repository success`() =
        runTest {
            fakeRepo.registerResult = Result.success(Unit)
            viewModel.updateRegisterName("Alice")
            viewModel.updateRegisterEmail("alice@example.com")
            viewModel.updateRegisterPassword("password1")
            viewModel.updateRegisterConfirmPassword("password1")
            viewModel.register()
            assertEquals(AuthUiState.Success, viewModel.uiState.value)
        }

    @Test
    fun `register emits Error when passwords do not match`() =
        runTest {
            viewModel.updateRegisterName("Alice")
            viewModel.updateRegisterEmail("alice@example.com")
            viewModel.updateRegisterPassword("password1")
            viewModel.updateRegisterConfirmPassword("password2")
            viewModel.register()
            val state = viewModel.uiState.value
            assertTrue(state is AuthUiState.Error)
            assertEquals("Passwords do not match", (state as AuthUiState.Error).message)
        }

    @Test
    fun `register emits Error on repository failure`() =
        runTest {
            fakeRepo.registerResult = Result.failure(Exception("Email already in use"))
            viewModel.updateRegisterName("Alice")
            viewModel.updateRegisterEmail("alice@example.com")
            viewModel.updateRegisterPassword("password1")
            viewModel.updateRegisterConfirmPassword("password1")
            viewModel.register()
            val state = viewModel.uiState.value
            assertTrue(state is AuthUiState.Error)
            assertEquals("Email already in use", (state as AuthUiState.Error).message)
        }

    @Test
    fun `signOut calls repository signOut and resets uiState to Idle`() =
        runTest {
            fakeRepo.signInResult = Result.success(Unit)
            viewModel.updateLoginEmail("test@example.com")
            viewModel.updateLoginPassword("pw")
            viewModel.signIn()
            assertEquals(AuthUiState.Success, viewModel.uiState.value)

            viewModel.signOut()
            assertTrue(fakeRepo.signOutCalled)
            assertEquals(AuthUiState.Idle, viewModel.uiState.value)
        }

    @Test
    fun `clearError resets Error state to Idle`() =
        runTest {
            fakeRepo.signInResult = Result.failure(Exception("Bad"))
            viewModel.updateLoginEmail("x@y.com")
            viewModel.updateLoginPassword("pw")
            viewModel.signIn()
            assertTrue(viewModel.uiState.value is AuthUiState.Error)

            viewModel.clearError()
            assertEquals(AuthUiState.Idle, viewModel.uiState.value)
        }

    @Test
    fun `updating login field clears existing error`() =
        runTest {
            fakeRepo.signInResult = Result.failure(Exception("Bad"))
            viewModel.updateLoginEmail("x@y.com")
            viewModel.updateLoginPassword("pw")
            viewModel.signIn()
            assertTrue(viewModel.uiState.value is AuthUiState.Error)

            viewModel.updateLoginEmail("new@example.com")
            assertEquals(AuthUiState.Idle, viewModel.uiState.value)
        }

    @Test
    fun `loginForm reflects field updates`() {
        viewModel.updateLoginEmail("a@b.com")
        viewModel.updateLoginPassword("secret")
        val form = viewModel.loginForm.value
        assertEquals("a@b.com", form.email)
        assertEquals("secret", form.password)
    }

    @Test
    fun `registerForm reflects field updates`() {
        viewModel.updateRegisterName("Bob")
        viewModel.updateRegisterEmail("bob@example.com")
        viewModel.updateRegisterPassword("pass")
        viewModel.updateRegisterConfirmPassword("pass")
        val form = viewModel.registerForm.value
        assertEquals("Bob", form.name)
        assertEquals("bob@example.com", form.email)
        assertEquals("pass", form.password)
        assertEquals("pass", form.confirmPassword)
    }
}
