package com.company.template.onboarding

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

// --------------- Fake repository ---------------

class FakeOnboardingRepository : OnboardingRepository {
    private val _seen = MutableStateFlow(false)
    var markSeenCalled = false

    override fun hasSeenOnboarding(): Flow<Boolean> = _seen

    override suspend fun markSeen() {
        markSeenCalled = true
        _seen.value = true
    }
}

// --------------- Tests ---------------

@OptIn(ExperimentalCoroutinesApi::class)
class OnboardingViewModelTest {
    private lateinit var fakeRepo: FakeOnboardingRepository
    private lateinit var viewModel: OnboardingViewModel

    @Before
    fun setUp() {
        Dispatchers.setMain(UnconfinedTestDispatcher())
        fakeRepo = FakeOnboardingRepository()
        viewModel = OnboardingViewModel(fakeRepo)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `hasSeenOnboarding emits false initially`() =
        runTest {
            val value = viewModel.hasSeenOnboarding().first()
            assertEquals(false, value)
        }

    @Test
    fun `markSeen calls repository and flow emits true`() =
        runTest {
            var callbackInvoked = false
            viewModel.markSeen(onComplete = { callbackInvoked = true })

            assertTrue(fakeRepo.markSeenCalled)
            assertTrue(callbackInvoked)
            val value = viewModel.hasSeenOnboarding().first()
            assertEquals(true, value)
        }
}
