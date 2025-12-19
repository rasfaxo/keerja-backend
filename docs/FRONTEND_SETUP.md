# Frontend Integration Guide

This guide provides configuration and integration instructions for frontend teams (Web and Mobile) to connect to the Keerja Backend API.

## Table of Contents

1. [Environment URLs](#environment-urls)
2. [Web Frontend (PHP/Next.js) Setup](#web-frontend-setup)
3. [Mobile (Flutter) Setup](#mobile-flutter-setup)
4. [Authentication Flow](#authentication-flow)
5. [CORS Configuration](#cors-configuration)
6. [API Documentation](#api-documentation)
7. [WebSocket Integration](#websocket-integration)
8. [Error Handling](#error-handling)
9. [Testing Guidelines](#testing-guidelines)

---

## Environment URLs

### API Endpoints

| Environment | Base URL                                        | Purpose                        |
| ----------- | ----------------------------------------------- | ------------------------------ |
| **Local**   | `http://localhost:8080/api/v1`                  | Backend developer testing      |
| **STAGING** | `http://staging-api.145.79.8.227.nip.io/api/v1` | Frontend development & testing |
| **DEMO**    | `http://demo-api.145.79.8.227.nip.io/api/v1`    | Client demonstrations          |

### Documentation

| Environment | Docs URL                               |
| ----------- | -------------------------------------- |
| **STAGING** | https://bump.sh/doc/keerja-api-staging |
| **DEMO**    | https://bump.sh/doc/keerja-api-demo    |

### Health Check Endpoints

```bash
# Check if API is running
GET /health/live

# Check if API is ready (includes DB connection)
GET /health/ready
```

---

## Web Frontend Setup

### Environment Configuration

Create environment files in your web project:

**`.env.development`** (for local development against STAGING):

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://staging-api.145.79.8.227.nip.io/api/v1
NEXT_PUBLIC_WS_URL=ws://staging-api.145.79.8.227.nip.io/ws

# Documentation
NEXT_PUBLIC_API_DOCS_URL=https://bump.sh/doc/keerja-api-staging

# Environment
NEXT_PUBLIC_ENV=staging
```

**`.env.production`** (for production/demo):

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://demo-api.145.79.8.227.nip.io/api/v1
NEXT_PUBLIC_WS_URL=ws://demo-api.145.79.8.227.nip.io/ws

# Documentation
NEXT_PUBLIC_API_DOCS_URL=https://bump.sh/doc/keerja-api-demo

# Environment
NEXT_PUBLIC_ENV=demo
```

### PHP Configuration

**`config/api.php`**:

```php
<?php

return [
    'environments' => [
        'local' => [
            'base_url' => 'http://localhost:8080/api/v1',
            'ws_url' => 'ws://localhost:8080/ws',
            'docs_url' => 'http://localhost:8080/docs',
        ],
        'staging' => [
            'base_url' => 'http://staging-api.145.79.8.227.nip.io/api/v1',
            'ws_url' => 'ws://staging-api.145.79.8.227.nip.io/ws',
            'docs_url' => 'https://bump.sh/doc/keerja-api-staging',
        ],
        'demo' => [
            'base_url' => 'http://demo-api.145.79.8.227.nip.io/api/v1',
            'ws_url' => 'ws://demo-api.145.79.8.227.nip.io/ws',
            'docs_url' => 'https://bump.sh/doc/keerja-api-demo',
        ],
    ],

    // Current environment
    'current' => env('API_ENV', 'staging'),
];
```

### API Client Example (JavaScript/TypeScript)

```typescript
// lib/api-client.ts

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ||
  "http://staging-api.145.79.8.227.nip.io/api/v1";

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  meta?: {
    page: number;
    per_page: number;
    total: number;
  };
}

class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  setToken(token: string) {
    this.token = token;
  }

  clearToken() {
    this.token = null;
  }

  private async request<T>(
    method: string,
    endpoint: string,
    data?: any
  ): Promise<ApiResponse<T>> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      Accept: "application/json",
    };

    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const config: RequestInit = {
      method,
      headers,
      credentials: "include", // For cookies if needed
    };

    if (data && method !== "GET") {
      config.body = JSON.stringify(data);
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, config);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "API request failed");
    }

    return response.json();
  }

  // HTTP Methods
  get<T>(endpoint: string) {
    return this.request<T>("GET", endpoint);
  }

  post<T>(endpoint: string, data: any) {
    return this.request<T>("POST", endpoint, data);
  }

  put<T>(endpoint: string, data: any) {
    return this.request<T>("PUT", endpoint, data);
  }

  patch<T>(endpoint: string, data: any) {
    return this.request<T>("PATCH", endpoint, data);
  }

  delete<T>(endpoint: string) {
    return this.request<T>("DELETE", endpoint);
  }

  // File upload
  async upload<T>(
    endpoint: string,
    formData: FormData
  ): Promise<ApiResponse<T>> {
    const headers: HeadersInit = {};

    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: "POST",
      headers,
      body: formData,
    });

    return response.json();
  }
}

export const api = new ApiClient(API_BASE_URL);
```

---

## Mobile (Flutter) Setup

### Environment Configuration

**`lib/config/environment.dart`**:

```dart
/// Keerja API Environment Configuration
///
/// This file contains all environment URLs for the Keerja Backend API.
///
/// Usage:
///   final apiUrl = Environment.staging.apiUrl;
///   final docsUrl = Environment.staging.docsUrl;

enum EnvironmentType { local, staging, demo }

class Environment {
  final String name;
  final String apiUrl;
  final String wsUrl;
  final String docsUrl;

  const Environment._({
    required this.name,
    required this.apiUrl,
    required this.wsUrl,
    required this.docsUrl,
  });

  /// Local development (for backend developers)
  /// Use with Android Emulator: 10.0.2.2:8080
  /// Use with iOS Simulator: localhost:8080
  /// Use with real device: Your computer's IP address
  static const local = Environment._(
    name: 'local',
    apiUrl: 'http://10.0.2.2:8080/api/v1', // Android Emulator
    wsUrl: 'ws://10.0.2.2:8080/ws',
    docsUrl: 'http://10.0.2.2:8080/docs',
  );

  /// Local for iOS Simulator
  static const localIOS = Environment._(
    name: 'local-ios',
    apiUrl: 'http://localhost:8080/api/v1',
    wsUrl: 'ws://localhost:8080/ws',
    docsUrl: 'http://localhost:8080/docs',
  );

  /// Staging environment (for frontend development)
  /// Semi-stable, may have frequent updates
  static const staging = Environment._(
    name: 'staging',
    apiUrl: 'http://staging-api.145.79.8.227.nip.io/api/v1',
    wsUrl: 'ws://staging-api.145.79.8.227.nip.io/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-staging',
  );

  /// Demo environment (for client demonstrations)
  /// Stable releases only
  static const demo = Environment._(
    name: 'demo',
    apiUrl: 'http://demo-api.145.79.8.227.nip.io/api/v1',
    wsUrl: 'ws://demo-api.145.79.8.227.nip.io/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-demo',
  );

  /// Direct IP access (for real device testing)
  static const directStaging = Environment._(
    name: 'direct-staging',
    apiUrl: 'http://145.79.8.227:8080/api/v1',
    wsUrl: 'ws://145.79.8.227:8080/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-staging',
  );

  static const directDemo = Environment._(
    name: 'direct-demo',
    apiUrl: 'http://145.79.8.227:8081/api/v1',
    wsUrl: 'ws://145.79.8.227:8081/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-demo',
  );

  /// Current environment (change this for different builds)
  static Environment current = staging;

  /// Helper to get environment by type
  static Environment fromType(EnvironmentType type) {
    switch (type) {
      case EnvironmentType.local:
        return local;
      case EnvironmentType.staging:
        return staging;
      case EnvironmentType.demo:
        return demo;
    }
  }
}
```

### Android Network Security Config

**`android/app/src/main/res/xml/network_security_config.xml`**:

```xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
    <!-- Allow cleartext traffic for development -->
    <domain-config cleartextTrafficPermitted="true">
        <!-- Local development -->
        <domain includeSubdomains="true">10.0.2.2</domain>
        <domain includeSubdomains="true">localhost</domain>

        <!-- VPS environments -->
        <domain includeSubdomains="true">145.79.8.227</domain>
        <domain includeSubdomains="true">staging-api.145.79.8.227.nip.io</domain>
        <domain includeSubdomains="true">demo-api.145.79.8.227.nip.io</domain>
    </domain-config>

    <!-- Production should use HTTPS only -->
    <!-- <base-config cleartextTrafficPermitted="false" /> -->
</network-security-config>
```

**Update `android/app/src/main/AndroidManifest.xml`**:

```xml
<application
    android:networkSecurityConfig="@xml/network_security_config"
    ... >
```

### iOS App Transport Security

**`ios/Runner/Info.plist`** (add inside `<dict>`):

```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsArbitraryLoads</key>
    <false/>
    <key>NSExceptionDomains</key>
    <dict>
        <!-- Local development -->
        <key>localhost</key>
        <dict>
            <key>NSExceptionAllowsInsecureHTTPLoads</key>
            <true/>
            <key>NSIncludesSubdomains</key>
            <true/>
        </dict>

        <!-- VPS environments -->
        <key>145.79.8.227</key>
        <dict>
            <key>NSExceptionAllowsInsecureHTTPLoads</key>
            <true/>
            <key>NSIncludesSubdomains</key>
            <true/>
        </dict>
        <key>nip.io</key>
        <dict>
            <key>NSExceptionAllowsInsecureHTTPLoads</key>
            <true/>
            <key>NSIncludesSubdomains</key>
            <true/>
        </dict>
    </dict>
</dict>
```

### API Service Example (Dart)

```dart
// lib/services/api_service.dart

import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config/environment.dart';

class ApiResponse<T> {
  final bool success;
  final String message;
  final T? data;

  ApiResponse({
    required this.success,
    required this.message,
    this.data,
  });
}

class ApiService {
  static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;
  ApiService._internal();

  String? _token;

  String get baseUrl => Environment.current.apiUrl;

  void setToken(String token) {
    _token = token;
  }

  void clearToken() {
    _token = null;
  }

  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    if (_token != null) 'Authorization': 'Bearer $_token',
  };

  Future<ApiResponse<T>> get<T>(
    String endpoint, {
    T Function(dynamic)? fromJson,
  }) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl$endpoint'),
        headers: _headers,
      );
      return _handleResponse<T>(response, fromJson);
    } catch (e) {
      return ApiResponse(success: false, message: e.toString());
    }
  }

  Future<ApiResponse<T>> post<T>(
    String endpoint, {
    Map<String, dynamic>? body,
    T Function(dynamic)? fromJson,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl$endpoint'),
        headers: _headers,
        body: body != null ? jsonEncode(body) : null,
      );
      return _handleResponse<T>(response, fromJson);
    } catch (e) {
      return ApiResponse(success: false, message: e.toString());
    }
  }

  ApiResponse<T> _handleResponse<T>(
    http.Response response,
    T Function(dynamic)? fromJson,
  ) {
    final json = jsonDecode(response.body);

    if (response.statusCode >= 200 && response.statusCode < 300) {
      return ApiResponse(
        success: true,
        message: json['message'] ?? 'Success',
        data: fromJson != null ? fromJson(json['data']) : json['data'],
      );
    } else {
      return ApiResponse(
        success: false,
        message: json['message'] ?? 'Request failed',
      );
    }
  }
}
```

---

## Authentication Flow

### Login Flow

```
1. User enters credentials
2. POST /api/v1/auth/login
   Body: { "email": "user@example.com", "password": "password" }
3. Receive tokens
   Response: { "access_token": "...", "refresh_token": "..." }
4. Store tokens securely
   - Web: httpOnly cookies or secure storage
   - Mobile: secure storage (flutter_secure_storage)
5. Include token in subsequent requests
   Header: Authorization: Bearer <access_token>
```

### Token Refresh Flow

```
1. Access token expires (401 response)
2. POST /api/v1/auth/refresh
   Body: { "refresh_token": "..." }
3. Receive new tokens
4. Retry original request with new token
```

### OAuth (Google) Flow

```
1. Initiate OAuth
   GET /api/v1/auth/oauth/google?redirect_uri=<your_callback>
2. User authenticates with Google
3. Receive callback with authorization code
4. Exchange code for tokens (handled by backend)
5. Receive access and refresh tokens
```

**Mobile OAuth Redirect URIs:**

```
keerja://oauth-callback
com.keerja.app://oauth-callback
```

---

## CORS Configuration

### Allowed Origins

The backend allows requests from:

| Origin                      | Environment                |
| --------------------------- | -------------------------- |
| `http://localhost:3000`     | Web dev (Next.js)          |
| `http://localhost:3001`     | Web dev (alternative port) |
| `http://localhost:5173`     | Web dev (Vite)             |
| `http://staging.keerja.com` | STAGING frontend           |
| `http://demo.keerja.com`    | DEMO frontend              |

### Mobile Apps

Mobile apps typically send `null` or no `Origin` header. The backend handles this by:

- Accepting requests without Origin header
- Not requiring credentials for mobile requests

### CORS Headers

```
Access-Control-Allow-Origin: <origin>
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type, Accept, Origin
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 86400
```

---

## API Documentation

### Accessing Documentation

- **STAGING:** https://bump.sh/doc/keerja-api-staging
- **DEMO:** https://bump.sh/doc/keerja-api-demo

### Common Endpoints

| Method | Endpoint         | Description       |
| ------ | ---------------- | ----------------- |
| POST   | `/auth/login`    | User login        |
| POST   | `/auth/register` | User registration |
| POST   | `/auth/refresh`  | Refresh token     |
| GET    | `/users/me`      | Get current user  |
| GET    | `/jobs`          | List jobs         |
| POST   | `/jobs`          | Create job        |
| GET    | `/companies`     | List companies    |

---

## WebSocket Integration

### Connection URL

```
STAGING: ws://staging-api.145.79.8.227.nip.io/ws
DEMO: ws://demo-api.145.79.8.227.nip.io/ws
```

### Authentication

Include token as query parameter:

```
ws://staging-api.145.79.8.227.nip.io/ws?token=<access_token>
```

### Flutter WebSocket Example

```dart
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocketService {
  WebSocketChannel? _channel;

  void connect(String token) {
    final wsUrl = '${Environment.current.wsUrl}?token=$token';
    _channel = WebSocketChannel.connect(Uri.parse(wsUrl));

    _channel?.stream.listen(
      (message) {
        // Handle incoming message
        print('Received: $message');
      },
      onError: (error) {
        print('WebSocket error: $error');
      },
      onDone: () {
        print('WebSocket closed');
      },
    );
  }

  void send(Map<String, dynamic> data) {
    _channel?.sink.add(jsonEncode(data));
  }

  void disconnect() {
    _channel?.sink.close();
  }
}
```

---

## Error Handling

### Standard Error Response

```json
{
  "success": false,
  "message": "Error description",
  "errors": {
    "email": ["Email is required"],
    "password": ["Password must be at least 8 characters"]
  }
}
```

### HTTP Status Codes

| Code | Meaning          | Action                    |
| ---- | ---------------- | ------------------------- |
| 200  | Success          | Process response          |
| 201  | Created          | Resource created          |
| 400  | Bad Request      | Check request body        |
| 401  | Unauthorized     | Refresh token or re-login |
| 403  | Forbidden        | User lacks permission     |
| 404  | Not Found        | Resource doesn't exist    |
| 422  | Validation Error | Check error details       |
| 500  | Server Error     | Report to backend team    |

---

## Testing Guidelines

### Testing Against STAGING

1. Use STAGING environment for all development testing
2. STAGING may have breaking changes - check documentation
3. Test user accounts are available on STAGING

### Testing Against DEMO

1. Use DEMO only for final QA and client demos
2. DEMO has stable, released code only
3. Don't use DEMO for regular development

### Mobile Testing Checklist

- [ ] Test on Android Emulator (use `10.0.2.2` for localhost)
- [ ] Test on iOS Simulator (use `localhost`)
- [ ] Test on real Android device (use VPS IP or nip.io domain)
- [ ] Test on real iOS device (use VPS IP or nip.io domain)
- [ ] Verify network security config allows cleartext (development only)

---

## Support

For API issues, contact the backend team:

- Email: support@keerja.com
- Slack: #keerja-backend

For documentation issues:

- Check the API docs first
- Report discrepancies to backend team
