package com.company.template.navigation

import android.app.Activity
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.company.template.auth.AuthUiState
import com.company.template.auth.AuthViewModel
import com.company.template.auth.User
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
    modifier: Modifier = Modifier,
) {
    val startDestination by appViewModel.startDestination.collectAsStateWithLifecycle()
    val authUiState by authViewModel.uiState.collectAsStateWithLifecycle()
    val currentUser: User? by authViewModel.currentUser.collectAsStateWithLifecycle()
    val loginForm by authViewModel.loginForm.collectAsStateWithLifecycle()
    val registerForm by authViewModel.registerForm.collectAsStateWithLifecycle()
    val activity = LocalContext.current as Activity

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
        modifier = modifier,
    ) {
        composable(ROUTE_ONBOARDING) {
            OnboardingScreen(
                onGetStarted = {
                    appViewModel.markOnboardingSeen()
                    navController.navigate(ROUTE_LOGIN) {
                        popUpTo(ROUTE_ONBOARDING) { inclusive = true }
                    }
                }
            )
        }
        composable(ROUTE_LOGIN) {
            LoginScreen(
                loginForm = loginForm,
                uiState = authUiState,
                onEmailChange = authViewModel::updateLoginEmail,
                onPasswordChange = authViewModel::updateLoginPassword,
                onSignIn = authViewModel::signIn,
                onSignInWithGoogle = { authViewModel.signInWithGoogle(activity) },
                onNavigateToRegister = { navController.navigate(ROUTE_REGISTER) },
            )
        }
        composable(ROUTE_REGISTER) {
            RegisterScreen(
                registerForm = registerForm,
                uiState = authUiState,
                onNameChange = authViewModel::updateRegisterName,
                onEmailChange = authViewModel::updateRegisterEmail,
                onPasswordChange = authViewModel::updateRegisterPassword,
                onConfirmPasswordChange = authViewModel::updateRegisterConfirmPassword,
                onRegister = authViewModel::register,
                onNavigateToLogin = { navController.popBackStack() },
            )
        }
        composable(ROUTE_HOME) {
            HomeScreen(
                displayName = currentUser?.displayName ?: currentUser?.email ?: "",
                onSignOut = { authViewModel.signOut() },
            )
        }
    }
}
