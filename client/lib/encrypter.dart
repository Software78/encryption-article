import 'package:encrypt/encrypt.dart' as crypt;
import 'package:flutter_dotenv/flutter_dotenv.dart';

class AESCryptoSystem {
  const AESCryptoSystem({
    required this.key,
    required this.vector,
  });
  final String key;
  final String vector;

  String encrypt(String? input) {
    try {
      final encrypter = crypt.Encrypter(
        crypt.AES(crypt.Key.fromUtf8(key), mode: crypt.AESMode.cbc),
      );
      return input == null
          ? ''
          : encrypter.encrypt(input, iv: crypt.IV.fromUtf8(vector)).base64;
    } catch (e) {
      return '';
    }
  }

  String decrypt(String? input) {
    try {
      final deckey = crypt.Key.fromUtf8(key);
      final decrypter =
          crypt.Encrypter(crypt.AES(deckey, mode: crypt.AESMode.cbc));
      final deciv = crypt.IV.fromUtf8(vector);
      return input == null ? '' : decrypter.decrypt64(input, iv: deciv);
    } catch (e) {
      return '';
    }
  }
}

Map<String, dynamic> encryptMap(Map<String, dynamic> payload) {
  final decryptedPayload = <String, dynamic>{};
 final encrypter = AESCryptoSystem(
    key: dotenv.env['KEY'] ?? '' ,// 'my32lengthkey',
    vector: dotenv.env['VECTOR'] ?? '', // 'my16lengthvector',
  );
  payload.forEach((key, value) {
    if (value is String) {
      // Encrypt string values
      decryptedPayload[key] = encrypter.encrypt(value);
    } else if (value is Map<String, dynamic>) {
      // Recursively encrypt nested maps
      decryptedPayload[key] = encryptMap(value);
    } else if (value is List<dynamic>) {
      // Recursively encrypt list elements
      decryptedPayload[key] = encryptList(value);
    } else {
      // Handle other data types (if needed)
      decryptedPayload[key] = value;
    }
  });

  return decryptedPayload;
}

List<dynamic> encryptList(List<dynamic> list) {
  final encryptedList = <dynamic>[];
final encrypter = AESCryptoSystem(
    key: dotenv.env['KEY'] ?? '' ,// 'my32lengthkey',
    vector: dotenv.env['VECTOR'] ?? '', // 'my16lengthvector',
  );
  for (final element in list) {
    if (element is String) {
      // Encrypt string elements
      encryptedList.add(encrypter.encrypt(element));
    } else if (element is Map<String, dynamic>) {
      // Recursively encrypt nested maps
      encryptedList.add(encryptMap(element));
    } else if (element is List<dynamic>) {
      // Recursively encrypt nested lists
      encryptedList.add(
        encryptList(
          element,
        ),
      );
    } else {
      // Handle other data types (if needed)
      encryptedList.add(element);
    }
  }

  return encryptedList;
}

Map<String, dynamic> decryptMap(Map<String, dynamic> payload) {
  final encryptedPayload = <String, dynamic>{};
 final encrypter = AESCryptoSystem(
    key: dotenv.env['KEY'] ?? '' ,// 'my32lengthkey',
    vector: dotenv.env['VECTOR'] ?? '', // 'my16lengthvector',
  );
  payload.forEach((key, value) {
    if (value is String) {
      final decryptedValue = encrypter.decrypt(value);
      final number = num.tryParse(decryptedValue);
      if (number != null) {
        encryptedPayload[key] = number;
        return;
      }
      if (decryptedValue.toLowerCase() == 'true' ||
          decryptedValue.toLowerCase() == 'false') {
        encryptedPayload[key] = decryptedValue.toLowerCase() == 'true';
        return;
      }

      encryptedPayload[key] = encrypter.decrypt(value);
    } else if (value is Map<String, dynamic>) {
      // Recursively encrypt nested maps
      encryptedPayload[key] = decryptMap(value);
    } else if (value is List<dynamic>) {
      // Recursively encrypt list elements
      encryptedPayload[key] = decryptList(value);
    } else if (value is List<Map<String, dynamic>>) {
      // Recursively encrypt list elements
      encryptedPayload[key] = value.map(decryptMap).toList();
    } else {
      // Handle other data types (if needed)
      encryptedPayload[key] = value;
    }
  });
  return encryptedPayload;
}

List<dynamic> decryptList(List<dynamic> list) {
  final decryptedList = <dynamic>[];
final encrypter = AESCryptoSystem(
    key: dotenv.env['KEY'] ?? '' ,// 'my32lengthkey',
    vector: dotenv.env['VECTOR'] ?? '', // 'my16lengthvector',
  );

  for (final element in list) {
    if (element is String) {
      // Encrypt string elements
      decryptedList.add(encrypter.decrypt(element));
    } else if (element is Map<String, dynamic>) {
      // Recursively encrypt nested maps
      decryptedList.add(decryptMap(element));
    } else if (element is List<dynamic>) {
      // Recursively encrypt nested lists
      decryptedList.add(
        decryptList(
          element,
        ),
      );
    } else {
      // Handle other data types (if needed)
      decryptedList.add(element);
    }
  }
  return decryptedList;
}
