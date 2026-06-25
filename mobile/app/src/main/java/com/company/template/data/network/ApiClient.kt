package com.company.template.data.network

import com.google.firebase.auth.FirebaseAuth
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.Response
import java.util.concurrent.TimeUnit

object ApiClient {
    // Interceptor that attaches the Firebase ID token on every request.
    // intercept() runs on a background OkHttp dispatcher thread, so
    // getIdToken().result (synchronous) is safe — do NOT use await() here.
    private class AuthInterceptor : Interceptor {
        override fun intercept(chain: Interceptor.Chain): Response {
            val token =
                runCatching {
                    // getIdToken(false) returns the cached token if still valid
                    FirebaseAuth
                        .getInstance()
                        .currentUser
                        ?.getIdToken(false)
                        ?.result
                        ?.token
                }.getOrNull()

            val request =
                if (token != null) {
                    chain
                        .request()
                        .newBuilder()
                        .header("Authorization", "Bearer $token")
                        .build()
                } else {
                    chain.request()
                }
            return chain.proceed(request)
        }
    }

    val httpClient: OkHttpClient =
        OkHttpClient
            .Builder()
            .addInterceptor(AuthInterceptor())
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .build()
}
