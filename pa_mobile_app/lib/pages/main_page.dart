import 'package:flutter/material.dart';
import 'package:pa_mobile_app/pages/map_page.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MainPage extends StatelessWidget {
  const MainPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF132555),
      body: FutureBuilder<UserInfo?>(
          future: _getDisplayName(),
          builder: (context, snapshot) {
            if (snapshot.hasData && snapshot.data != null) {
              return SafeArea(
                child: Center(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      const Text('Welcome,', style: TextStyle(color: Colors.white)),
                      Text(
                        snapshot.data!.userName,
                        style: const TextStyle(color: Colors.white),
                      ),
                      Text(snapshot.data!.eMail, style: const TextStyle(color: Colors.white)),
                      const SizedBox(height: 20),
                      const Expanded(child: MapPage()),
                      //Text(snapshot.data!.responseBody, style: const TextStyle(color: Colors.white)),
                    ],
                  ),
                ),
              );
            } else {
              return const Text('No Data');
            }
          }),
    );
  }

  Future<UserInfo> _getDisplayName() async {
    final SharedPreferences pref = await SharedPreferences.getInstance();
    return UserInfo(pref.getString('Name')!, pref.getString('IdToken')!, pref.getString('Email')!);
  }
}

class UserInfo {
  final String userName;
  final String eMail;
  final String idToken;

  UserInfo(this.userName, this.idToken, this.eMail);
}
