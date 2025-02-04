import 'dart:developer';

import 'package:client/dio.dart';
import 'package:dio/dio.dart' as dio;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';

Future<void> main() async {
  await dotenv.load(fileName: ".env");
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final _loginEmailController = TextEditingController();
  final _loginPasswordController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();

  final _dio = dio.Dio(dio.BaseOptions(baseUrl: dotenv.env['BASE_URL'] ?? ''))
    ..options.connectTimeout = const Duration(minutes: 1)
    ..options.receiveTimeout = const Duration(minutes: 1)
    ..options.sendTimeout = const Duration(minutes: 1)
    ..interceptors.add(EncryptionInterceptor())
    ..interceptors.add(
      PrettyDioLogger(
        requestBody: true,
        requestHeader: true,
        responseHeader: true,
        logPrint: (value) {
          if (kDebugMode) {
            log(value.toString(), name: 'Dio');
          }
        },
      ),
    );

  Future<void> register() async {
    setState(() {
      registerLoading = true;
    });
    try {
      final response = await _dio.post('/auth/register', data: {
        'email': _emailController.text,
        'password': _passwordController.text,
        'first_name': _firstNameController.text,
        'last_name': _lastNameController.text,
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('User registered: ${response.data['data']}'),
        ),
      );
    } on dio.DioException catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('User logged in: ${e.toString()}'),
        ),
      );
      log('Error logging in: $e', name: 'Login');
    }
    setState(() {
      registerLoading = false;
    });
  }

  Future<void> login() async {
    setState(() {
      loginLoading = true;
    });
    try {
      final response = await _dio.post('/auth/login', data: {
        'email': _loginEmailController.text,
        'password': _loginPasswordController.text,
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('User logged in: ${response.data['data']}'),
        ),
      );
    } on dio.DioException catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('User logged in: ${e.toString()}'),
        ),
      );
      log('Error logging in: $e', name: 'Login');
    }
    setState(() {
      loginLoading = false;
    });
  }

  bool loginLoading = false;
  bool registerLoading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: PageView(
        children: [
          Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              children: [
                TextFormField(
                  controller: _loginEmailController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'E-mail'),
                ),
                SizedBox(height: 20),
                TextFormField(
                  controller: _loginPasswordController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'Password'),
                ),
                SizedBox(height: 20),
                ElevatedButton(
                  onPressed: loginLoading ? null : login,
                  child: const Text('Login'),
                ),
                SizedBox(height: 40),
                TextFormField(
                  controller: _emailController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'E-mail'),
                ),
                SizedBox(height: 20),
                TextFormField(
                  controller: _passwordController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'Password'),
                ),
                SizedBox(height: 20),
                TextFormField(
                  controller: _firstNameController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'First Name'),
                ),
                SizedBox(height: 20),
                TextFormField(
                  controller: _lastNameController,
                  onTapOutside: (_) {
                    FocusScope.of(context).unfocus();
                  },
                  decoration: const InputDecoration(labelText: 'Last Name'),
                ),
                SizedBox(height: 20),
                ElevatedButton(
                  onPressed: registerLoading ? null : register,
                  child: const Text('Register'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
