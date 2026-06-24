package com.company.template.onboarding

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "onboarding_prefs")

class DataStoreOnboardingRepository(private val context: Context) : OnboardingRepository {

    private val hasSeenKey = booleanPreferencesKey("has_seen_onboarding")

    override fun hasSeenOnboarding(): Flow<Boolean> =
        context.dataStore.data.map { prefs -> prefs[hasSeenKey] ?: false }

    override suspend fun markSeen() {
        context.dataStore.edit { prefs -> prefs[hasSeenKey] = true }
    }
}
