package com.company.template.navigation

import com.company.template.auth.AuthRepository
import com.company.template.auth.User
import com.company.template.onboarding.OnboardingRepository
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Test

// --------------- Fakes ---------------

class FakeAuthRepositoryForNav : AuthRepository {
    private val _authStateFlow = MutableStateFlow<User?>(null)
    override val authStateFlow: StateFlow<User?> = _authStateFlow

    fun setUser(user: User?) {
        _authStateFlow.value = user
    }

    override suspend fun signInWithEmail(email: String, password: String): Result<Unit> =
        Result.success(Unit)

    override suspend fun registerWithEmail(name: String, email: String, password: String): Result<Unit> =
        Result.success(Unit)

    override suspend fun signInWithGoogle(googleIdToken: String): Result<Unit> = Result.success(Unit)

    override suspend fun signOut() {
        _authStateFlow.value = null
    }
}

class FakeOnboardingRepositoryForNav : OnboardingRepository {
    private val _seen = MutableStateFlow(false)

    fun setSeen(seen: Boolean) {
        _seen.value = seen
    }

    override fun hasSeenOnboarding(): Flow<Boolean> = _seen

    override suspend fun markSeen() {
        _seen.value = true
    }
}

// --------------- Tests ---------------

@OptIn(ExperimentalCoroutinesApi::class)
class AppViewModelTest {

    private lateinit var fakeAuth: FakeAuthRepositoryForNav
    private lateinit var fakeOnboarding: FakeOnboardingRepositoryForNav
    private lateinit var viewModel: AppViewModel

    @Before
    fun setUp() {
        Dispatchers.setMain(UnconfinedTestDispatcher())
        fakeAuth = FakeAuthRepositoryForNav()
        fakeOnboarding = FakeOnboardingRepositoryForNav()
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    private fun createViewModel() {
        viewModel = AppViewModel(fakeAuth, fakeOnboarding)
    }

    @Test
    fun `startDestination is Onboarding when not seen and not signed in`() = runTest {
        fakeOnboarding.setSeen(false)
        fakeAuth.setUser(null)
        createViewModel()

        val dest = viewModel.startDestination.first { it != null }
        assertEquals(StartDestination.ONBOARDING, dest)
    }

    @Test
    fun `startDestination is Login when onboarding seen and not signed in`() = runTest {
        fakeOnboarding.setSeen(true)
        fakeAuth.setUser(null)
        createViewModel()

        val dest = viewModel.startDestination.first { it != null }
        assertEquals(StartDestination.LOGIN, dest)
    }

    @Test
    fun `startDestination is Home when user is signed in`() = runTest {
        fakeOnboarding.setSeen(true)
        fakeAuth.setUser(User(uid = "uid123", email = "a@b.com", displayName = "Alice", photoUrl = null))
        createViewModel()

        val dest = viewModel.startDestination.first { it != null }
        assertEquals(StartDestination.HOME, dest)
    }

    @Test
    fun `startDestination transitions to Login after onboarding is marked seen`() = runTest {
        fakeOnboarding.setSeen(false)
        fakeAuth.setUser(null)
        createViewModel()

        val firstDest = viewModel.startDestination.first { it != null }
        assertEquals(StartDestination.ONBOARDING, firstDest)

        fakeOnboarding.setSeen(true)
        val secondDest = viewModel.startDestination.first { it == StartDestination.LOGIN }
        assertEquals(StartDestination.LOGIN, secondDest)
    }

    @Test
    fun `startDestination transitions to Login after sign out`() = runTest {
        fakeOnboarding.setSeen(true)
        fakeAuth.setUser(User(uid = "uid123", email = "a@b.com", displayName = "Alice", photoUrl = null))
        createViewModel()

        val homeDest = viewModel.startDestination.first { it == StartDestination.HOME }
        assertEquals(StartDestination.HOME, homeDest)

        fakeAuth.signOut()
        val loginDest = viewModel.startDestination.first { it == StartDestination.LOGIN }
        assertEquals(StartDestination.LOGIN, loginDest)
    }

    @Test
    fun `markOnboardingSeen persists the flag`() = runTest {
        fakeOnboarding.setSeen(false)
        fakeAuth.setUser(null)
        createViewModel()

        viewModel.markOnboardingSeen()

        val seen = fakeOnboarding.hasSeenOnboarding().first()
        assertEquals(true, seen)
    }
}
