/// Keerja API Environment Configuration for Flutter
/// 
/// Copy this file to your Flutter project at:
///   lib/config/environment.dart
/// 
/// Usage:
///   import 'config/environment.dart';
///   
///   final api = ApiService(Environment.current.apiUrl);
///   // or
///   Environment.setCurrent(EnvironmentType.demo);

// ignore_for_file: constant_identifier_names

/// Available environment types
enum EnvironmentType { 
  /// Local backend development
  local,
  /// Local for iOS Simulator (uses localhost instead of 10.0.2.2)
  localIOS,
  /// Staging environment for development testing
  staging, 
  /// Demo environment for client presentations
  demo,
  /// Direct IP access to staging (bypasses nginx)
  directStaging,
  /// Direct IP access to demo (bypasses nginx)
  directDemo,
}

/// Environment configuration class
/// Holds all API URLs and configuration for each environment
class Environment {
  final String name;
  final String apiUrl;
  final String wsUrl;
  final String docsUrl;
  final bool isProduction;
  
  const Environment._({
    required this.name,
    required this.apiUrl,
    required this.wsUrl,
    required this.docsUrl,
    this.isProduction = false,
  });
  
  // =========================================
  // LOCAL DEVELOPMENT
  // =========================================
  
  /// Local development for Android Emulator
  /// Android Emulator uses 10.0.2.2 to access host machine
  static const local = Environment._(
    name: 'local',
    apiUrl: 'http://10.0.2.2:8080/api/v1',
    wsUrl: 'ws://10.0.2.2:8080/ws',
    docsUrl: 'http://10.0.2.2:8080/docs',
    isProduction: false,
  );
  
  /// Local development for iOS Simulator
  /// iOS Simulator uses localhost to access host machine
  static const localIOS = Environment._(
    name: 'local-ios',
    apiUrl: 'http://localhost:8080/api/v1',
    wsUrl: 'ws://localhost:8080/ws',
    docsUrl: 'http://localhost:8080/docs',
    isProduction: false,
  );
  
  // =========================================
  // VPS ENVIRONMENTS (via Nginx reverse proxy)
  // =========================================
  
  /// Staging environment for development and testing
  /// - Updated frequently with latest features
  /// - May contain breaking changes
  /// - Use for development and internal testing
  static const staging = Environment._(
    name: 'staging',
    apiUrl: 'http://staging-api.145.79.8.227.nip.io/api/v1',
    wsUrl: 'ws://staging-api.145.79.8.227.nip.io/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-staging',
    isProduction: false,
  );
  
  /// Demo environment for client presentations
  /// - Stable releases only (tagged versions)
  /// - No breaking changes between demos
  /// - Use for client demos and QA testing
  static const demo = Environment._(
    name: 'demo',
    apiUrl: 'http://demo-api.145.79.8.227.nip.io/api/v1',
    wsUrl: 'ws://demo-api.145.79.8.227.nip.io/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-demo',
    isProduction: true,
  );
  
  // =========================================
  // DIRECT IP ACCESS (bypasses Nginx)
  // Use when nip.io domain doesn't resolve
  // =========================================
  
  /// Direct staging access via IP:port
  static const directStaging = Environment._(
    name: 'direct-staging',
    apiUrl: 'http://145.79.8.227:8080/api/v1',
    wsUrl: 'ws://145.79.8.227:8080/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-staging',
    isProduction: false,
  );
  
  /// Direct demo access via IP:port
  static const directDemo = Environment._(
    name: 'direct-demo',
    apiUrl: 'http://145.79.8.227:8081/api/v1',
    wsUrl: 'ws://145.79.8.227:8081/ws',
    docsUrl: 'https://bump.sh/doc/keerja-api-demo',
    isProduction: true,
  );
  
  // =========================================
  // CURRENT ENVIRONMENT
  // =========================================
  
  /// The currently active environment
  /// Change this value for different build configurations
  static Environment _current = staging;
  
  /// Get the current environment
  static Environment get current => _current;
  
  /// Set the current environment
  static void setCurrent(EnvironmentType type) {
    _current = fromType(type);
  }
  
  /// Get environment by type
  static Environment fromType(EnvironmentType type) {
    switch (type) {
      case EnvironmentType.local:
        return local;
      case EnvironmentType.localIOS:
        return localIOS;
      case EnvironmentType.staging:
        return staging;
      case EnvironmentType.demo:
        return demo;
      case EnvironmentType.directStaging:
        return directStaging;
      case EnvironmentType.directDemo:
        return directDemo;
    }
  }
  
  // =========================================
  // HELPER METHODS
  // =========================================
  
  /// Check if current environment is local
  bool get isLocal => name.startsWith('local');
  
  /// Check if current environment uses direct IP
  bool get isDirect => name.startsWith('direct');
  
  /// Get full health check URL
  String get healthUrl => '$apiUrl/../health/live';
  
  @override
  String toString() => 'Environment($name)';
}

// =========================================
// BUILD FLAVOR CONFIGURATION
// =========================================

/// Configure environment based on build flavor
/// 
/// Usage in main.dart:
///   void main() {
///     setupEnvironment();
///     runApp(MyApp());
///   }
void setupEnvironment() {
  // You can use compile-time variables for build flavors
  const flavor = String.fromEnvironment('FLAVOR', defaultValue: 'staging');
  
  switch (flavor) {
    case 'local':
      Environment.setCurrent(EnvironmentType.local);
      break;
    case 'staging':
      Environment.setCurrent(EnvironmentType.staging);
      break;
    case 'demo':
    case 'production':
      Environment.setCurrent(EnvironmentType.demo);
      break;
    default:
      Environment.setCurrent(EnvironmentType.staging);
  }
  
  // Print current environment for debugging
  assert(() {
    print('🌍 Environment: ${Environment.current.name}');
    print('📡 API URL: ${Environment.current.apiUrl}');
    return true;
  }());
}

// =========================================
// USAGE EXAMPLES
// =========================================

/*
// Example 1: Basic API call
import 'package:http/http.dart' as http;

Future<void> fetchUsers() async {
  final response = await http.get(
    Uri.parse('${Environment.current.apiUrl}/users'),
    headers: {'Authorization': 'Bearer $token'},
  );
}

// Example 2: WebSocket connection
import 'package:web_socket_channel/web_socket_channel.dart';

void connectWebSocket(String token) {
  final channel = WebSocketChannel.connect(
    Uri.parse('${Environment.current.wsUrl}?token=$token'),
  );
}

// Example 3: Environment-specific UI
Widget buildDebugBanner() {
  if (Environment.current.isProduction) {
    return SizedBox.shrink();
  }
  return Banner(
    message: Environment.current.name.toUpperCase(),
    location: BannerLocation.topStart,
  );
}

// Example 4: Switch environment at runtime (for debug menu)
void showEnvironmentSelector(BuildContext context) {
  showDialog(
    context: context,
    builder: (context) => AlertDialog(
      title: Text('Select Environment'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: EnvironmentType.values.map((type) {
          return ListTile(
            title: Text(type.name),
            selected: Environment.current == Environment.fromType(type),
            onTap: () {
              Environment.setCurrent(type);
              Navigator.pop(context);
              // Restart app or reload data
            },
          );
        }).toList(),
      ),
    ),
  );
}
*/
