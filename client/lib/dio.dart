import 'dart:developer';

import 'package:client/encrypter.dart';
import 'package:dio/dio.dart' as dio;

class EncryptionInterceptor extends dio.InterceptorsWrapper {
  @override
  void onRequest(
    dio.RequestOptions options,
    dio.RequestInterceptorHandler handler,
  ) {
    if (options.data is Map<String, dynamic>) {
      options.data = encryptMap(options.data as Map<String, dynamic>);
    } else if (options.data is dio.FormData) {
      final fields = <String, dynamic>{};
      final files = <String, dynamic>{};
      log('files: ${options.data.files}', name: 'Dio Encrypt');
      final formData = options.data as dio.FormData;
      for (final item in formData.fields) {
        if (item.value.isNotEmpty) {
          fields[item.key] = item.value;
        } else {
          fields.remove(item.key);
        }
      }
      for (final item in formData.files) {
        files[item.key] = item.value;
      }
      log('files: ${formData.files}', name: 'Dio Encrypt');
      final encodedFormData = dio.FormData.fromMap({
        ...encryptMap(fields),
        if (formData.files.isNotEmpty)
          formData.files.first.key: formData.files.map((e) => e.value).toList(),
      });
      options.data = encodedFormData;
    }
    handler.next(options);
  }

  @override
  void onResponse(
    dio.Response<dynamic> response,
    dio.ResponseInterceptorHandler handler,
  ) {
    try {
      if (response.data != null) {
        if (response.data['data'] is List) {
          response.data['data'] =
              decryptList(response.data['data'] as List<dynamic>);
        } else {
          if (response.data is Map<String, dynamic>) {
            if ((response.data as Map<String, dynamic>).containsKey('data')) {
              response.data['data'] = decryptMap(
                (response.data['data'] as Map<String, dynamic>?) ?? {},
              );
            }
          }
        }
      }
      handler.next(response);
    } catch (e) {
      log('Error decrypting response: $e', name: 'Dio Decrypt');
      handler.next(response);
    }
  }
}