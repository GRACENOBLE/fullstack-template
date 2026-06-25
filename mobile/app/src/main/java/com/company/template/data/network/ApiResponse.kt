package com.company.template.data.network

import kotlinx.serialization.Serializable

@Serializable
data class ApiResponse<T>(val data: T)

@Serializable
data class ApiErrorDetail(val code: String, val message: String)

@Serializable
data class ApiErrorResponse(val error: ApiErrorDetail)
