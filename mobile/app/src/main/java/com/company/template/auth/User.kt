package com.company.template.auth

data class User(
    val uid: String,
    val email: String?,
    val displayName: String?,
    val photoUrl: String?,
)
