package com.company.template.navigation

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.company.template.auth.AuthUiState
import com.company.template.auth.AuthViewModel
import com.google.firebase.auth.FirebaseUser
import com.company.template.auth.LoginScreen
import com.company.template.auth.RegisterScreen
import com.company.template.home.HomeScreen
import com.company.template.onboarding.OnboardingScreen

private const val ROUTE_ONBOARDING = "onboarding"
private const val ROUTE_LOGIN = "login"
private const val ROUTE_REGISTER = "register"
private const val ROUTE_HOME = "home"

@Composable
fun AppNavGraph(
    appViewModel: AppViewModel,
    authViewModel: AuthViewModel,
    navController: NavHostController = rememberNavController(),
    modifier: Modifier = Modifier
) {
    val startDestination by appViewModel.startDestination.collectAsStateWithLifecycle()
    val authUiState by authViewModel.uiState.collectAsStateWithLifecycle()
    val currentUser: FirebaseUser? by authViewModel.currentUser.collectAsStateWithLifecycle()

    // Navigate away from auth screens when the user successfully signs in
    LaunchedEffect(authUiState) {
        if (authUiState is AuthUiState.Success) {
            navController.navigate(ROUTE_HOME) {
                popUpTo(0) { inclusive = true }
            }
        }
    }

    val resolvedStart = when (startDestination) {
        StartDestination.ONBOARDING -> ROUTE_ONBOARDING
        StartDestination.LOGIN -> ROUTE_LOGIN
        StartDestination.HOME -> ROUTE_HOME
        null -> return // wait until resolved
    }

    NavHost(
        navController = navController,
        startDestination = resolvedStart,
        modifier = modifier
    ) {
        composable(ROUTE_ONBOARDING) {
            OnboardingScreen(
                onGetStarted = {
                    navController.navigate(ROUTE_LOGIN) {
                        popUpTo(ROUTE_ONBOARDING) { inclusive = true }
                    }
                }
            )
        }
        composable(ROUTE_LOGIN) {
            LoginScreen(
                uiState = authUiState,
                onSignIn = { email, password -> authViewModel.signIn(email, password) },
                onNavigateToRegister = {
                    navController.navigate(ROUTE_REGISTER)
                },
                onClearError = authViewModel::clearError
            )
        }
        composable(ROUTE_REGISTER) {
            RegisterScreen(
                uiState = authUiState,
                onRegister = { name, email, password, confirm ->
                    authViewModel.register(name, email, password, confirm)
                },
                onNavigateToLogin = {
                    navController.popBackStack()
                },
                onClearError = authViewModel::clearError
            )
        }
        composable(ROUTE_HOME) {
            HomeScreen(
                displayName = currentUser?.displayName ?: currentUser?.email ?: "",
                onSignOut = { authViewModel.signOut() }
            )
        }
    }
}
