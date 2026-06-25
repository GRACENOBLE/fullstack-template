package com.company.template.fcm

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import androidx.core.app.NotificationCompat
import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody

class MyFirebaseMessagingService : FirebaseMessagingService() {
    private val serviceScope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private val httpClient = OkHttpClient()

    override fun onNewToken(token: String) {
        serviceScope.launch {
            registerTokenWithBackend(token)
        }
    }

    override fun onMessageReceived(message: RemoteMessage) {
        message.notification?.let { notif ->
            showNotification(notif.title.orEmpty(), notif.body.orEmpty())
        }
    }

    private fun registerTokenWithBackend(token: String) {
        val backendUrl =
            getString(
                applicationContext.resources
                    .getIdentifier(
                        "backend_base_url",
                        "string",
                        packageName,
                    ).takeIf { it != 0 } ?: return,
            )
        val payload = FcmRegistrationPayload(token = token)
        val body =
            Json
                .encodeToString(FcmRegistrationPayload.serializer(), payload)
                .toRequestBody("application/json".toMediaType())
        val request =
            Request
                .Builder()
                .url("$backendUrl/api/v1/fcm/register")
                .post(body)
                .build()
        runCatching { httpClient.newCall(request).execute().close() }
    }

    private fun showNotification(
        title: String,
        body: String,
    ) {
        val channelId = "fcm_default"
        val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel =
                NotificationChannel(
                    channelId,
                    "Push Notifications",
                    NotificationManager.IMPORTANCE_DEFAULT,
                )
            manager.createNotificationChannel(channel)
        }

        val notification =
            NotificationCompat
                .Builder(this, channelId)
                .setSmallIcon(android.R.drawable.ic_dialog_info)
                .setContentTitle(title)
                .setContentText(body)
                .setAutoCancel(true)
                .build()
        manager.notify(System.currentTimeMillis().toInt(), notification)
    }
}
