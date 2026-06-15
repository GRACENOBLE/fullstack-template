package com.company.template

import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

/**
 * Unit tests for Sentry initialisation logic.
 *
 * The real SentryAndroid.init() is an Android framework call and cannot run on
 * the JVM. The production code therefore guards the call behind
 * [shouldInitSentry], which IS pure Kotlin and testable here.
 */
class SentryInitTest {

    @Test
    fun `shouldInitSentry returns false when DSN is empty`() {
        assertFalse(shouldInitSentry(""))
    }

    @Test
    fun `shouldInitSentry returns false when DSN is blank`() {
        assertFalse(shouldInitSentry("   "))
    }

    @Test
    fun `shouldInitSentry returns true when DSN is non-empty`() {
        assertTrue(shouldInitSentry("https://example@sentry.io/123"))
    }
}
